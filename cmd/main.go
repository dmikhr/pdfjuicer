// Copyright (c) 2025 Dmitrii Khramtsov
// License: AGPL-3.0

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

	config "github.com/dmikhr/pdfjuicer/configs"
	dsp "github.com/dmikhr/pdfjuicer/internal/display"
	"github.com/dmikhr/pdfjuicer/internal/extractor"
	"github.com/dmikhr/pdfjuicer/internal/input"
)

func main() {
	var sizeX, sizeY, thumbSizeX, thumbSizeY int
	var err error
	var anyErr bool

	var cfg config.Config

	workersNumDefault := runtime.NumCPU()

	pflag.StringVarP(&cfg.SourcePath, "source", "s", "",
		"Specify path to source file (pdf)")
	pflag.StringVarP(&cfg.SaveDir, "output", "o", "", "Specify output folder path")

	pflag.StringVarP(&cfg.Prefix, "prefix", "p", config.DefaultFilenamePrefix, "Prefix for a filename")
	pflag.StringVarP(&cfg.Postfix, "postfix", "x", "", "Postfix for a filename")

	pflag.StringVarP(&cfg.Image.ImgSize, "size", "S", "",
		"Specify image size, example 640x480, if not specified will output default size from document")
	pflag.Float64VarP(&cfg.Image.ImgScaleDown, "scale", "C", config.ImgScaleDownDefault,
		"Specify image scaling down factor, example 5, for example 5 means output image will be 5 times smaller than original image")
	pflag.StringVarP(&cfg.Image.ImgType, "format", "F", config.DefaultImgFormat,
		"Specify output image format (png/jpg)")

	pflag.StringVarP(&cfg.Pages, "pages", "P", "",
		"Use this flag to extract specific pages, example: 2,3,6-8,10")

	pflag.BoolVarP(&cfg.Thumb.CreateThumbnails, "thumb", "t", false, "enable thumbnails generation")
	pflag.Float64VarP(&cfg.Thumb.ThumbScaleDown, "tscale", "c", config.ThumbScaleDownDefault,
		"Specify thumbnails scaling down factor, for example 5 means thumbnail will be 5 times smaller than original image")
	pflag.StringVarP(&cfg.Thumb.ThumbnailsSize, "tsize", "z", "",
		"Specify thumbnails size e.g. 64x64")

	pflag.BoolVarP(&cfg.VersionFlag, "version", "v", false, "Show version")

	pflag.IntVarP(&cfg.WorkersNum, "workers", "w", workersNumDefault,
		"Set number of anynchronous workers")

	pflag.BoolVarP(&cfg.Quiet, "quiet", "q", false, "Quiet mode (no progress bar, no colored output)")

	pflag.Parse()

	// show help if called with no params
	if pflag.NFlag() == 0 && pflag.NArg() == 0 {
		fmt.Println(config.About())
		pflag.Usage()
		os.Exit(0)
	}

	// if called with unsupported arguments
	if pflag.NFlag() == 0 && pflag.NArg() > 0 {
		fmt.Printf("Unsupported arguments: %s\n", pflag.Args())
		os.Exit(1)
	}

	if cfg.VersionFlag {
		fmt.Printf("pdfjuicer version %s\n", config.Version)
		return
	}

	if cfg.Image.ImgSize != "" && cfg.Image.ImgScaleDown != config.ImgScaleDownDefault {
		fmt.Fprintln(os.Stderr, "Choose either scaling factor (--scale) or exact image size for resizing (--size)")
		anyErr = true
	}
	if err = input.ImgFormatValidator(cfg.Image.ImgType); err != nil {
		fmt.Fprintf(os.Stderr, "Unsupported image type: %s\n", cfg.Image.ImgType)
		anyErr = true
	}
	if cfg.Thumb.ThumbnailsSize != "" && cfg.Thumb.ThumbScaleDown != config.ThumbScaleDownDefault {
		fmt.Fprintln(os.Stderr, "Choose either scaling factor (--scale) or exact image size for resizing (--size)")
		anyErr = true
	}

	if cfg.SourcePath == "" {
		fmt.Fprintln(os.Stderr, "No source pdf file was specified")
		anyErr = true
	}
	if cfg.SaveDir == "" {
		fmt.Fprintln(os.Stderr, "No target directory for image extraction was specified")
		anyErr = true
	}

	if cfg.Prefix != "" {
		if err = input.FilenameValidator(cfg.Prefix); err != nil {
			fmt.Fprintf(os.Stderr, "Invalid prefix: %s. Error: %s\n", cfg.Prefix, err)
			anyErr = true
		}
	}
	if cfg.Postfix != "" {
		if err = input.FilenameValidator(cfg.Postfix); err != nil {
			fmt.Fprintf(os.Stderr, "Invalid postfix: %s. Error: %s\n", cfg.Postfix, err)
			anyErr = true
		}
	}

	if cfg.WorkersNum <= 0 {
		fmt.Fprintln(os.Stderr, "Number of workers must be at least 1")
		anyErr = true
	}

	if anyErr {
		os.Exit(1)
	}

	if cfg.Image.ImgSize != "" {
		sizeX, sizeY, err = input.ImgSizeExtractor(cfg.Image.ImgSize)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid image size (example: 120x256): %s\n", err)
			return
		}
		fmt.Printf("Extracted images size will be set to: %s\n", dsp.Fbg(cfg.Image.ImgSize, cfg.Quiet))
	} else if cfg.Image.ImgScaleDown != config.ImgScaleDownDefault {
		fmt.Printf("Extracted images size will be scaled down with factor %s\n", dsp.Fbg(cfg.Image.ImgScaleDown, cfg.Quiet))
	}

	if cfg.Thumb.ThumbnailsSize != "" {
		thumbSizeX, thumbSizeY, err = input.ImgSizeExtractor(cfg.Thumb.ThumbnailsSize)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid thumbnail size (example: 120x256): %s\n", err)
			return
		}
		fmt.Printf("Thumbnails size will be set to: %s\n", dsp.Fbg(cfg.Thumb.ThumbnailsSize, cfg.Quiet))
	} else if cfg.Thumb.ThumbScaleDown != config.ThumbScaleDownDefault {
		fmt.Printf("Thumbnails will be resized with scaling down factor %s\n", dsp.Fbg(cfg.Image.ImgScaleDown, cfg.Quiet))
	}

	fmt.Printf("Setting image format to %s, save folder: %s\n",
		dsp.Fbg(cfg.Image.ImgType, cfg.Quiet),
		dsp.Fbg(cfg.SaveDir, cfg.Quiet))
	if cfg.Pages != "" {
		fmt.Printf("Selected pages will be extracted: %s\n",
			dsp.Fbg(cfg.Pages, cfg.Quiet))
	}

	workDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Start processing...")

	doc, err := fitz.New(filepath.Join(workDir, cfg.SourcePath))
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
	if cfg.Pages != "" {
		pagesToExtract, err = input.PagesExtractor(cfg.Pages, pageCount)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		for i := 1; i <= pageCount; i++ {
			pagesToExtract = append(pagesToExtract, i)
		}
	}

	savePath := filepath.Join(workDir, cfg.SaveDir)
	createPath := savePath
	if cfg.Thumb.CreateThumbnails {
		createPath = filepath.Join(createPath, config.ThumbnailsDir)
	}

	err = os.MkdirAll(createPath, 0755)
	if err != nil {
		log.Fatal(err)
	}

	thumbnails := extractor.Thumbnail{
		IsActive:  cfg.Thumb.CreateThumbnails,
		ScaleDown: cfg.Thumb.ThumbScaleDown,
		SizeX:     thumbSizeX,
		SizeY:     thumbSizeY,
	}

	page := extractor.Page{
		Doc:        doc,
		ImgType:    cfg.Image.ImgType,
		SavePath:   savePath,
		Prefix:     cfg.Prefix,
		Postfix:    cfg.Postfix,
		ScaleDown:  cfg.Image.ImgScaleDown,
		SizeX:      sizeX,
		SizeY:      sizeY,
		Thumbnails: thumbnails,
	}

	var wg sync.WaitGroup
	numJobs := len(pagesToExtract)
	jobs := make(chan extractor.Job, numJobs)
	jobErrors := make(chan extractor.JobErr, numJobs)
	done := make(chan struct{}, numJobs)

	var bar *progressbar.ProgressBar
	if !cfg.Quiet {
		bar = progressbar.Default(int64(numJobs))
	}

	for w := 1; w <= cfg.WorkersNum; w++ {
		wg.Add(1)
		go extractor.Worker(w, jobs, jobErrors, done, &wg)
	}
	for _, pageNum := range pagesToExtract {
		jobs <- extractor.Job{Page: page, PageNum: pageNum - 1}
	}

	close(jobs)

	go func() {
		for i := 0; i < numJobs; i++ {
			<-done
			if !cfg.Quiet {
				err = bar.Add(1)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Progress bar encountered problem: %s\n", err)
				}
			}
		}
	}()

	wg.Wait()

	if !cfg.Quiet {
		err = bar.Finish()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Progress bar encountered problem: %s\n", err)
		}
	}

	close(jobErrors)

	for jobErr := range jobErrors {
		if jobErr.Err != nil {
			log.Printf("Worker %d failed: %v", jobErr.WorkerID, jobErr.Err)
		}
	}

	if len(jobErrors) == 0 {
		fmt.Println("Finished extraction")
	}
}
