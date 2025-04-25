package imageutils

import (
	"image"

	"golang.org/x/image/draw"
)

func ScaleResize(srcImg *image.RGBA, scaleFactor float64) *image.RGBA {
	dstWidth := int(float64(srcImg.Bounds().Dx()) / scaleFactor)
	dstHeight := int(float64(srcImg.Bounds().Dy()) / scaleFactor)

	return Resize(srcImg, dstWidth, dstHeight)
}

func Resize(srcImg *image.RGBA, dstWidth, dstHeight int) *image.RGBA {
	dstImg := image.NewRGBA(image.Rect(0, 0, dstWidth, dstHeight))

	draw.BiLinear.Scale(
		dstImg,
		dstImg.Bounds(),
		srcImg,
		srcImg.Bounds(),
		draw.Over,
		nil,
	)
	return dstImg

}
