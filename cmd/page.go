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
	sizeX      int
	sizeY      int
	thumbnails Thumbnail
}

type Thumbnail struct {
	isActive  bool
	scaleDown float64
	sizeX     int
	sizeY     int
}

func (ps *Page) extract(pageNum int) error {
	// extracting source image
	srcImg, err := ps.doc.Image(pageNum)
	if err != nil {
		return err
	}

	// creating file for dst image
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
		// fmt.Println("Saving img without resizing")
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

		if ps.thumbnails.scaleDown != thumbScaleDownDefault {
			thumbnail = imageutils.ScaleResize(srcImg, ps.thumbnails.scaleDown)
		} else if ps.thumbnails.sizeX > 0 && ps.thumbnails.sizeY > 0 {
			thumbnail = imageutils.Resize(srcImg, ps.sizeX, ps.sizeY)
		} else {
			fmt.Println("Saving thumbnail without resizing")
			thumbnail = srcImg
		}
		err = saveImg(f, ps.imgType, thumbnail)
		if err != nil {
			return err
		}
	}

	f.Close()

	// log.Printf("Page %d extracted to %s", pageNum+1, imageFName)

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
