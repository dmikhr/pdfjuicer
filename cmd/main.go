package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/dmikhr/pdfjuicer/internal/input"
	"github.com/gen2brain/go-fitz"
	"github.com/spf13/pflag"
)

const (
	version       = "1.0.0"
	workersNumber = 4
)

func main() {
	var sourcePath, saveDir, imgType, imgSize, pages, thumbnailsScale, prefix, postfix string
	var imgScaleDown, thumbScaleDown float64
	var createThumbnails, force, versionFlag bool
	var workersNum int
	workersNumDefault := runtime.NumCPU()

	// ImgFormatValidator ImgSizeValidator
	pflag.StringVarP(&sourcePath, "source", "s", "test.pdf", "Specify path to source file (pdf/pptx)")
	pflag.StringVarP(&saveDir, "output", "o", "", "Specify output folder path")

	// FilenameValidator
	pflag.StringVarP(&prefix, "prefix", "p", "page", "Prefix for a filename")
	pflag.StringVarP(&postfix, "postfix", "x", "", "Postfix for a filename")

	pflag.StringVarP(&imgSize, "size", "S", "", "Specify image size, example 640x480, if not specified will output default size from document")
	pflag.Float64VarP(&imgScaleDown, "scale", "C", 1.0, "Specify image scaling down factor, example 5, for example 5 means output image will be 5 times smaller than original image")
	pflag.StringVarP(&imgType, "format", "F", "png", "Specify output image format (png/jpg)")

	pflag.StringVarP(&pages, "pages", "P", "", "Use this flag to extract specific pages, example: 2,3,6-8,10")

	pflag.BoolVarP(&createThumbnails, "thumb", "t", false, "enable thumbnails generation")
	pflag.Float64VarP(&thumbScaleDown, "tscale", "c", 10.0, "Specify thumbnails scaling down factor, for example 5 means thumbnail will be 5 times smaller than original image")
	pflag.StringVarP(&thumbnailsScale, "tsize", "z", "64x64", "Specify thumbnails size e.g. 64x64")

	pflag.BoolVarP(&force, "force", "f", false, "Don't ask for rewriting is directory contains files")
	pflag.BoolVar(&versionFlag, "version", false, "Show version")

	pflag.IntVarP(&workersNum, "workers", "w", workersNumDefault, "Set number of anynchronous workers")

	pflag.Parse()

	log.Println(fmt.Sprintf("Setting image format to %s, save folder: %s", imgType, saveDir))
	if pages != "" {
		log.Println(fmt.Sprintf("Extracting selected pages: %s", pages))
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

	// todo logic with savedir if not exist and thumbnails dir
	var createDir string
	if createThumbnails {
		createDir = filepath.Join(workDir, thumbnailsDir, saveDir)
	} else {
		createDir = filepath.Join(workDir, saveDir)
	}
	// todo check if folder exists, and if not empty
	err = os.MkdirAll(createDir, 0755)
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
		savePath:   filepath.Join(workDir, saveDir),
		prefix:     prefix,
		postfix:    postfix,
		scaleDown:  imgScaleDown,
		thumbnails: thumbnails,
	}

	// Extract pages as images
	for _, pageNum := range pagesToExtract {
		err = extractPage(page, pageNum-1)
		if err != nil {
			log.Println("Error extracting page #", pageNum, err)
		}
		log.Println(fmt.Sprintf("Page %d extracted to %s", pageNum, saveDir))
	}
}
