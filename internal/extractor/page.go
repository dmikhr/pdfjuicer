package extractor

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"

	config "github.com/dmikhr/pdfjuicer/configs"
	"github.com/gen2brain/go-fitz"

	"github.com/dmikhr/pdfjuicer/internal/imageutils"
)

// Page contains settings for page extraction as image and pointer to source doc
type Page struct {
	Doc        *fitz.Document
	ImgType    string
	SavePath   string
	Prefix     string
	Postfix    string
	ScaleDown  float64
	SizeX      int
	SizeY      int
	Thumbnails Thumbnail
}

// Thumbnail contains settings for thumbnails
type Thumbnail struct {
	IsActive  bool
	ScaleDown float64
	SizeX     int
	SizeY     int
}

// Extract page from pdf document as image
func (ps *Page) Extract(pageNum int) error {
	srcImg, err := ps.Doc.Image(pageNum)
	if err != nil {
		return err
	}

	imageFName := fmt.Sprintf("%s%03d%s.%s", ps.Prefix, pageNum+1, ps.Postfix, ps.ImgType)
	imagePath := filepath.Join(ps.SavePath, imageFName)
	f, err := os.Create(imagePath)
	if err != nil {
		return err
	}

	var dstImg, thumbnail *image.RGBA
	if ps.ScaleDown != config.ImgScaleDownDefault {
		dstImg = imageutils.ScaleResize(srcImg, ps.ScaleDown)
	} else if ps.SizeX > 0 && ps.SizeY > 0 {
		dstImg = imageutils.Resize(srcImg, ps.SizeX, ps.SizeY)
	} else {
		dstImg = srcImg
	}

	err = saveImg(f, ps.ImgType, dstImg)
	if err != nil {
		return err
	}

	if ps.Thumbnails.IsActive {
		f, err = os.Create(filepath.Join(ps.SavePath, config.ThumbnailsDir,
			fmt.Sprintf("thumbnail_%03d.%s", pageNum+1, ps.ImgType)))
		if err != nil {
			return err
		}

		if ps.Thumbnails.SizeX > 0 && ps.Thumbnails.SizeY > 0 {
			thumbnail = imageutils.Resize(srcImg, ps.Thumbnails.SizeX, ps.Thumbnails.SizeY)
		} else {
			thumbnail = imageutils.ScaleResize(srcImg, ps.Thumbnails.ScaleDown)
		}
		err = saveImg(f, ps.ImgType, thumbnail)
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
	case "jpg", "jpeg":
		err = jpeg.Encode(f, dstImg, &jpeg.Options{Quality: jpeg.DefaultQuality})
	case "png":
		err = png.Encode(f, dstImg)
	}

	if err != nil {
		return err
	}
	return nil
}
