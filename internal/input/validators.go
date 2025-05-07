package input

import (
	"errors"
	"strings"
)

var allowedImgFormats = []string{"png", "jpg", "jpeg"}

// ErrUnsupportedImgFormat validates image format
var ErrUnsupportedImgFormat = errors.New("unsupported image format")

// ImgFormatValidator validates if submitted image format (e.g. png, jpg) is supported
func ImgFormatValidator(imgFormat string) error {
	imgFormatLow := strings.ToLower(imgFormat)
	for _, allowedFormat := range allowedImgFormats {
		if imgFormatLow == allowedFormat {
			return nil
		}
	}
	return ErrUnsupportedImgFormat
}

var (
	// ErrInputLong is returned when the provided input exceeds the allowed maximum length
	ErrInputLong = errors.New("input too long")
	// ErrInvalidChar is returned when the input contains invalid characters
	ErrInvalidChar = errors.New("invalid character")
)

// prefix and postfix sizes limited to:
const inputLimit = 100

// FilenameValidator validate postfix, prefix of a file
func FilenameValidator(s string) error {
	if len(s) > inputLimit {
		return ErrInputLong
	}
	allowedChars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_."
	for _, ch := range s {
		if !strings.Contains(allowedChars, string(ch)) {
			return ErrInvalidChar
		}
	}
	return nil
}
