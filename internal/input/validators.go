package input

import (
	"errors"
	"strconv"
	"strings"
)

var allowedImgFormats = []string{"png", "jpg", "jpeg"}

// UnsupportedImgFormatErr validates image format
var UnsupportedImgFormatErr = errors.New("unsupported image format")

func ImgFormatValidator(imgFormat string) error {
	imgFormatLow := strings.ToLower(imgFormat)
	for _, allowedFormat := range allowedImgFormats {
		if imgFormatLow == allowedFormat {
			return nil
		}
	}
	return UnsupportedImgFormatErr
}

var (
	NoXErr             = errors.New("no x in image size")
	SizeMustBeIntErr   = errors.New("image size must be int")
	SizeMustBePositive = errors.New("image size cannot be negative")
)

// todo change to parser
func ImgSizeValidator(s string) error {
	imgSize := strings.Split(s, "x")
	if len(imgSize) != 2 {
		return NoXErr
	}
	x, errX := strconv.Atoi(imgSize[0])
	y, errY := strconv.Atoi(imgSize[1])
	if errX != nil || errY != nil {
		return SizeMustBeIntErr
	}
	if x <= 0 || y <= 0 {
		return SizeMustBePositive
	}
	return nil
}

var (
	InputLongErr   = errors.New("input too long")
	InvalidCharErr = errors.New("invalid character")
)

// prefix and postfix sizes limited to:
const inputLimit = 100

// FilenameValidator validate postfix, prefix of a file
func FilenameValidator(s string) error {
	if len(s) > inputLimit {
		return InputLongErr
	}
	allowedChars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_."
	for _, ch := range s {
		if !strings.Contains(allowedChars, string(ch)) {
			return InvalidCharErr
		}
	}
	return nil
}
