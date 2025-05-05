package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/gen2brain/go-fitz"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/pflag"

	"github.com/dmikhr/pdfjuicer/internal/input"
)

const (
	version               = "1.0.0"
	imgScaleDownDefault   = 1.0
	thumbScaleDownDefault = 10.0
	defaultFilenamePrefix = "page"
	defaultImgFormat      = "png"
)

type config struct {
	sourcePath string
	saveDir    string
	prefix     string
	postfix    string
	pages      string
	image      struct {
		imgSize      string
		imgScaleDown float64
		imgType      string
	}
	thumb struct {
		createThumbnails bool
		thumbScaleDown   float64
		thumbnailsSize   string
	}
	workersNum  int
	force       bool
	versionFlag bool
	quiet       bool
}

func main() {
	var sizeX, sizeY, thumbSizeX, thumbSizeY int
	var err error
	var anyErr bool

	var cfg config

	workersNumDefault := runtime.NumCPU()

	pflag.StringVarP(&cfg.sourcePath, "source", "s", "",
		"Specify path to source file (pdf)")
	pflag.StringVarP(&cfg.saveDir, "output", "o", "", "Specify output folder path")

	pflag.StringVarP(&cfg.prefix, "prefix", "p", defaultFilenamePrefix, "Prefix for a filename")
	pflag.StringVarP(&cfg.postfix, "postfix", "x", "", "Postfix for a filename")

	pflag.StringVarP(&cfg.image.imgSize, "size", "S", "",
		"Specify image size, example 640x480, if not specified will output default size from document")
	pflag.Float64VarP(&cfg.image.imgScaleDown, "scale", "C", imgScaleDownDefault,
		"Specify image scaling down factor, example 5, for example 5 means output image will be 5 times smaller than original image")
	pflag.StringVarP(&cfg.image.imgType, "format", "F", defaultImgFormat,
		"Specify output image format (png/jpg)")

	pflag.StringVarP(&cfg.pages, "pages", "P", "",
		"Use this flag to extract specific pages, example: 2,3,6-8,10")

	pflag.BoolVarP(&cfg.thumb.createThumbnails, "thumb", "t", false, "enable thumbnails generation")
	pflag.Float64VarP(&cfg.thumb.thumbScaleDown, "tscale", "c", thumbScaleDownDefault,
		"Specify thumbnails scaling down factor, for example 5 means thumbnail will be 5 times smaller than original image")
	pflag.StringVarP(&cfg.thumb.thumbnailsSize, "tsize", "z", "",
		"Specify thumbnails size e.g. 64x64")

	pflag.BoolVarP(&cfg.force, "force", "f", false,
		"Don't ask for rewriting is directory contains files")
	pflag.BoolVarP(&cfg.versionFlag, "version", "v", false, "Show version")

	pflag.IntVarP(&cfg.workersNum, "workers", "w", workersNumDefault,
		"Set number of anynchronous workers")

	pflag.BoolVarP(&cfg.quiet, "quiet", "q", false, "Quiet mode (no progress bar, no colored output)")

	pflag.Parse()

	if cfg.image.imgSize != "" && cfg.image.imgScaleDown != imgScaleDownDefault {
		fmt.Fprintln(os.Stderr, "Choose either scaling factor (--scale) or exact image size for resizing (--size)")
		anyErr = true
	}
	if err = input.ImgFormatValidator(cfg.image.imgType); err != nil {
		fmt.Fprintf(os.Stderr, "Unsupported image type: %s\n", cfg.image.imgType)
		anyErr = true
	}
	if cfg.thumb.thumbnailsSize != "" && cfg.thumb.thumbScaleDown != thumbScaleDownDefault {
		fmt.Fprintln(os.Stderr, "Choose either scaling factor (--scale) or exact image size for resizing (--size)")
		anyErr = true
	}

	if cfg.image.imgSize != "" {
		sizeX, sizeY, err = input.ImgSizeExtractor(cfg.image.imgSize)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid image size (example: 120x256): %s\n", err)
			anyErr = true
		}
		fmt.Printf("Extracted images size will be set to: %s\n", fbg(cfg.image.imgSize, cfg.quiet))
	} else if cfg.image.imgScaleDown != imgScaleDownDefault {
		fmt.Printf("Extracted images size will be scaled down with factor %s\n", fbg(cfg.image.imgScaleDown, cfg.quiet))
	}

	if cfg.thumb.thumbnailsSize != "" {
		thumbSizeX, thumbSizeY, err = input.ImgSizeExtractor(cfg.thumb.thumbnailsSize)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid thumbnail size (example: 120x256): %s\n", err)
			anyErr = true
		}
		fmt.Printf("Thumbnails size will be set to: %s\n", fbg(cfg.thumb.thumbnailsSize, cfg.quiet))
	} else if cfg.thumb.thumbScaleDown != thumbScaleDownDefault {
		fmt.Printf("Thumbnails will be resized with scaling down factor %s\n", fbg(cfg.image.imgScaleDown, cfg.quiet))
	}

	if cfg.sourcePath == "" {
		fmt.Fprintln(os.Stderr, "No source pdf file was specified")
		anyErr = true
	}
	if cfg.saveDir == "" {
		fmt.Fprintln(os.Stderr, "No target directory for image extraction was specified")
		anyErr = true
	}

	if cfg.prefix != "" {
		if err = input.FilenameValidator(cfg.prefix); err != nil {
			fmt.Fprintf(os.Stderr, "Invalid prefix: %s. Error: %s\n", cfg.prefix, err)
			anyErr = true
		}
	}
	if cfg.postfix != "" {
		if err = input.FilenameValidator(cfg.postfix); err != nil {
			fmt.Fprintf(os.Stderr, "Invalid postfix: %s. Error: %s\n", cfg.postfix, err)
			anyErr = true
		}
	}

	if cfg.workersNum <= 0 {
		fmt.Fprintln(os.Stderr, "Number of workers must be at least 1")
		anyErr = true
	}

	if anyErr {
		os.Exit(1)
	}

	if cfg.versionFlag {
		fmt.Printf("pdfjuicer version %s\n", version)
		return
	}

	fmt.Printf("Setting image format to %s, save folder: %s\n",
		fbg(cfg.image.imgType, cfg.quiet),
		fbg(cfg.saveDir, cfg.quiet))
	if cfg.pages != "" {
		fmt.Printf("Selected pages will be extracted: %s\n",
			fbg(cfg.pages, cfg.quiet))
	}

	workDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Start processing...")

	doc, err := fitz.New(filepath.Join(workDir, cfg.sourcePath))
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err = doc.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	pageCount := doc.NumPage()
	var pagesToExtract []int
	if cfg.pages != "" {
		pagesToExtract, err = input.PagesExtractor(cfg.pages, pageCount)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		for i := 1; i <= pageCount; i++ {
			pagesToExtract = append(pagesToExtract, i)
		}
	}

	var savePath string
	if cfg.thumb.createThumbnails {
		savePath = filepath.Join(workDir, thumbnailsDir, cfg.saveDir)
	} else {
		savePath = filepath.Join(workDir, cfg.saveDir)
	}
	err = os.MkdirAll(savePath, 0755)
	if err != nil {
		log.Fatal(err)
	}

	thumbnails := Thumbnail{
		isActive:  cfg.thumb.createThumbnails,
		scaleDown: cfg.thumb.thumbScaleDown,
		sizeX:     thumbSizeX,
		sizeY:     thumbSizeY,
	}

	page := Page{
		doc:        doc,
		imgType:    cfg.image.imgType,
		savePath:   savePath,
		prefix:     cfg.prefix,
		postfix:    cfg.postfix,
		scaleDown:  cfg.image.imgScaleDown,
		sizeX:      sizeX,
		sizeY:      sizeY,
		thumbnails: thumbnails,
	}

	var wg sync.WaitGroup
	numJobs := len(pagesToExtract)
	jobs := make(chan Job, numJobs)
	jobErrors := make(chan JobErr, numJobs)
	done := make(chan struct{}, numJobs)

	var bar *progressbar.ProgressBar
	if !cfg.quiet {
		bar = progressbar.Default(int64(numJobs))
	}

	for w := 1; w <= cfg.workersNum; w++ {
		wg.Add(1)
		go worker(w, jobs, jobErrors, done, &wg)
	}
	for _, pageNum := range pagesToExtract {
		jobs <- Job{page: page, pageNum: pageNum - 1}
	}

	close(jobs)

	go func() {
		for i := 0; i < numJobs; i++ {
			<-done
			if !cfg.quiet {
				err = bar.Add(1)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Progress bar encountered problem: %s\n", err)
				}
			}
		}
	}()

	wg.Wait()

	if !cfg.quiet {
		err = bar.Finish()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Progress bar encountered problem: %s\n", err)
		}
	}

	close(jobErrors)

	for jobErr := range jobErrors {
		if jobErr.err != nil {
			log.Printf("Worker %d failed: %v", jobErr.workerID, jobErr.err)
		}
	}

	if len(jobErrors) == 0 {
		fmt.Println("Finished extraction")
	}
}
