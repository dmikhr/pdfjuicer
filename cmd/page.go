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

// Page contains settings for page extraction as image and pointer to source doc
type Page struct {
	doc        *fitz.Document
	imgType    string
	savePath   string
	prefix     string
	postfix    string
	scaleDown  float64
	sizeX      int
	sizeY      int
	thumbnails Thumbnail
}

// Thumbnail contains settings for thimbnails
type Thumbnail struct {
	isActive  bool
	scaleDown float64
	sizeX     int
	sizeY     int
}

// extract page from pdf document as image
func (ps *Page) extract(pageNum int) error {
	srcImg, err := ps.doc.Image(pageNum)
	if err != nil {
		return err
	}

	imageFName := fmt.Sprintf("%s%03d%s.%s", ps.prefix, pageNum+1, ps.postfix, ps.imgType)
	imagePath := filepath.Join(ps.savePath, imageFName)
	f, err := os.Create(imagePath)
	if err != nil {
		return err
	}

	var dstImg, thumbnail *image.RGBA
	if ps.scaleDown != imgScaleDownDefault {
		dstImg = imageutils.ScaleResize(srcImg, ps.scaleDown)
	} else if ps.sizeX > 0 && ps.sizeY > 0 {
		dstImg = imageutils.Resize(srcImg, ps.sizeX, ps.sizeY)
	} else {
		dstImg = srcImg
	}

	err = saveImg(f, ps.imgType, dstImg)
	if err != nil {
		return err
	}

	if ps.thumbnails.isActive {
		f, err = os.Create(filepath.Join(ps.savePath, thumbnailsDir,
			fmt.Sprintf("thumbnail_%03d.%s", pageNum+1, ps.imgType)))
		if err != nil {
			return err
		}

		if ps.thumbnails.sizeX > 0 && ps.thumbnails.sizeY > 0 {
			thumbnail = imageutils.Resize(srcImg, ps.thumbnails.sizeX, ps.thumbnails.sizeY)
		} else {
			thumbnail = imageutils.ScaleResize(srcImg, ps.thumbnails.scaleDown)
		}
		err = saveImg(f, ps.imgType, thumbnail)
		if err != nil {
			return err
		}
	}

	f.Close()

	return nil
}

// saveImg saves image in a given image format
func saveImg(f *os.File, imgType string, dstImg *image.RGBA) error {
	var err error

	switch imgType {
	case "jpg":
		err = jpeg.Encode(f, dstImg, &jpeg.Options{Quality: jpeg.DefaultQuality})
	case "png":
		err = png.Encode(f, dstImg)
	}

	if err != nil {
		return err
	}
	return nil
}
