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
)

func main() {
	var sourcePath, saveDir, imgType, imgSize, pages, thumbnailsSize, prefix, postfix string
	var imgScaleDown, thumbScaleDown float64
	var createThumbnails, force, versionFlag bool
	var workersNum int
	var sizeX, sizeY, thumbSizeX, thumbSizeY int
	var err error

	workersNumDefault := runtime.NumCPU()

	// ImgFormatValidator ImgSizeValidator
	pflag.StringVarP(&sourcePath, "source", "s", "test.pdf",
		"Specify path to source file (pdf/pptx)")
	pflag.StringVarP(&saveDir, "output", "o", "", "Specify output folder path")

	// FilenameValidator
	pflag.StringVarP(&prefix, "prefix", "p", "page", "Prefix for a filename")
	pflag.StringVarP(&postfix, "postfix", "x", "", "Postfix for a filename")

	pflag.StringVarP(&imgSize, "size", "S", "",
		"Specify image size, example 640x480, if not specified will output default size from document")
	pflag.Float64VarP(&imgScaleDown, "scale", "C", imgScaleDownDefault,
		"Specify image scaling down factor, example 5, for example 5 means output image will be 5 times smaller than original image")
	pflag.StringVarP(&imgType, "format", "F", "png",
		"Specify output image format (png/jpg)")

	pflag.StringVarP(&pages, "pages", "P", "",
		"Use this flag to extract specific pages, example: 2,3,6-8,10")

	pflag.BoolVarP(&createThumbnails, "thumb", "t", false, "enable thumbnails generation")
	pflag.Float64VarP(&thumbScaleDown, "tscale", "c", thumbScaleDownDefault,
		"Specify thumbnails scaling down factor, for example 5 means thumbnail will be 5 times smaller than original image")
	pflag.StringVarP(&thumbnailsSize, "tsize", "z", "",
		"Specify thumbnails size e.g. 64x64")

	pflag.BoolVarP(&force, "force", "f", false,
		"Don't ask for rewriting is directory contains files")
	pflag.BoolVar(&versionFlag, "version", false, "Show version")

	pflag.IntVarP(&workersNum, "workers", "w", workersNumDefault,
		"Set number of anynchronous workers")

	pflag.Parse()

	if imgSize != "" && imgScaleDown != imgScaleDownDefault {
		log.Println(fmt.Sprintf("Choose either scaling factor (--scale) or exact image size for resizing (--size)"))
		return
	}

	if err := input.ImgFormatValidator(imgType); err != nil {
		log.Println(fmt.Sprintf("Unsupported image type: %s", imgType))
		return
	}
	if thumbnailsSize != "" && thumbScaleDown != thumbScaleDownDefault {
		log.Println(fmt.Sprintf("Choose either scaling factor (--scale) or exact image size for resizing (--size)"))
		return
	}

	if prefix != "" {
		if err = input.FilenameValidator(prefix); err != nil {
			log.Println(fmt.Sprintf("Invalid prefix: %s. Error: %s", prefix, err))
		}
	}
	if postfix != "" {
		if err = input.FilenameValidator(postfix); err != nil {
			log.Println(fmt.Sprintf("Invalid postfix: %s. Error: %s", postfix, err))
		}
	}

	if workersNum <= 0 {
		log.Println("Number of workers must be at least 1")
		return
	} else {
		log.Println("Number of workers is", workersNum)
	}

	log.Println(fmt.Sprintf("Setting image format to %s, save folder: %s", imgType, saveDir))
	if pages != "" {
		log.Println(fmt.Sprintf("Selected pages will be extracted: %s", pages))
	}
	if imgSize != "" {
		sizeX, sizeY, err = input.ImgSizeExtractor(imgSize)
		if err != nil {
			log.Println(fmt.Sprintf("Invalid image size (example: 120x256): %s", err))
			return
		}
		log.Println(fmt.Sprintf("Extracted images size will be set to: %dx%d", sizeX, sizeY))
	} else if imgScaleDown != 1.0 {
		log.Println(fmt.Sprintf("Extracted images size will be scaled down with factor %.2f ", imgScaleDown))
	}

	if thumbnailsSize != "" {
		thumbSizeX, thumbSizeY, err = input.ImgSizeExtractor(thumbnailsSize)
		if err != nil {
			log.Println(fmt.Sprintf("Invalid thumbnail size (example: 120x256): %s", err))
			return
		}
		log.Println(fmt.Sprintf("Thumbnails size will be set to: %dx%d", thumbSizeX, thumbSizeY))
	} else if imgScaleDown != 1.0 {
		log.Println(fmt.Sprintf("Thumbnails will be resized with scaling down factor %.2f ", imgScaleDown))
	}

	if versionFlag {
		fmt.Printf("pdfjuicer version %s", version)
		return
	}

	workDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	doc, err := fitz.New(filepath.Join(workDir, sourcePath))
	if err != nil {
		log.Fatal(err)
	}

	defer doc.Close()

	pageCount := doc.NumPage()
	var pagesToExtract []int
	if pages != "" {
		pagesToExtract, err = input.PagesExtractor(pages, pageCount)
		if err != nil {
			log.Println(err)
			return
		}
	} else {
		for i := 1; i <= pageCount; i++ {
			pagesToExtract = append(pagesToExtract, i)
		}
	}

	var savePath string
	if createThumbnails {
		savePath = filepath.Join(workDir, thumbnailsDir, saveDir)
	} else {
		savePath = filepath.Join(workDir, saveDir)
	}
	err = os.MkdirAll(savePath, 0755)
	if err != nil {
		log.Fatal(err)
	}

	thumbnails := Thumbnail{
		isActive:  createThumbnails,
		scaleDown: thumbScaleDown,
	}

	page := Page{
		doc:        doc,
		imgType:    imgType,
		savePath:   savePath,
		prefix:     prefix,
		postfix:    postfix,
		scaleDown:  imgScaleDown,
		sizeX:      sizeX,
		sizeY:      sizeY,
		thumbnails: thumbnails,
	}

	var wg sync.WaitGroup
	numJobs := len(pagesToExtract)
	jobs := make(chan Job, numJobs)
	jobErrors := make(chan JobErr, numJobs)
	done := make(chan struct{}, numJobs)

	bar := progressbar.Default(int64(numJobs))

	for w := 1; w <= workersNum; w++ {
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
			bar.Add(1)
		}
	}()

	wg.Wait()

	close(jobErrors)

	for jobErr := range jobErrors {
		if jobErr.err != nil {
			log.Printf("Worker %d failed: %v", jobErr.workerID, jobErr.err)
		}
	}
}
