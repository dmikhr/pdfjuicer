package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"

	"github.com/gen2brain/go-fitz"

	"github.com/dmikhr/pdfjuicer/internal/imageutils"
)

const thumbnailsDir = "thumbnails"

type Page struct {
	doc        *fitz.Document
	imgType    string
	savePath   string
	prefix     string
	postfix    string
	scaleDown  float64
	thumbnails Thumbnail
}

type Thumbnail struct {
	isActive  bool
	scaleDown float64
}

func extractPage(p Page, pageNum int) error {
	// extracting source image
	srcImg, err := p.doc.Image(pageNum)
	if err != nil {
		return err
	}

	// creating file for dst image
	f, err := os.Create(filepath.Join(p.savePath,
		fmt.Sprintf("%s%03d%s.%s", p.prefix, pageNum+1, p.postfix, p.imgType)))
	if err != nil {
		return err
	}

	var dstImg *image.RGBA
	if p.scaleDown != 1.0 {
		dstImg = imageutils.ScaleResize(srcImg, p.scaleDown)
	} else {
		fmt.Println("Saving img without resizing")
		dstImg = srcImg
	}

	err = saveImg(f, p.imgType, dstImg)
	if err != nil {
		return err
	}

	if p.thumbnails.isActive {
		f, err = os.Create(filepath.Join(p.savePath, thumbnailsDir,
			fmt.Sprintf("thumbnail_%03d.%s", pageNum+1, p.imgType)))
		if err != nil {
			return err
		}

		thumbnail := imageutils.ScaleResize(srcImg, p.thumbnails.scaleDown)
		err = saveImg(f, p.imgType, thumbnail)
		if err != nil {
			return err
		}
	}

	f.Close()

	return nil
}

func saveImg(f *os.File, imgType string, dstImg *image.RGBA) error {
	var err error

	switch imgType {
	case "jpg":
		err = jpeg.Encode(f, dstImg, &jpeg.Options{jpeg.DefaultQuality})
	case "png":
		err = png.Encode(f, dstImg)
	}

	if err != nil {
		return err
	}
	return nil
}
