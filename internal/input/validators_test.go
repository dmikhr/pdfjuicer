package input

import (
	"errors"
	"testing"
)

type validParamTestCase struct {
	comment     string
	inputValue  string
	expectError error
}

var FileSubstringTestCase = []validParamTestCase{
	{
		comment:     "valid input",
		inputValue:  "some_part_1",
		expectError: nil,
	},
	{
		comment:     "valid input 2",
		inputValue:  "1.some.part",
		expectError: nil,
	},
	{
		comment:     "ivalid input 2",
		inputValue:  "1 some part",
		expectError: InvalidCharErr,
	},
	{
		comment:     "ivalid input 2",
		inputValue:  "some/part",
		expectError: InvalidCharErr,
	},
	{
		comment:     "ivalid input 3",
		inputValue:  "some?part",
		expectError: InvalidCharErr,
	},
}

var ImgFormatTestCase = []validParamTestCase{
	{
		comment:     "Supports png",
		inputValue:  "png",
		expectError: nil,
	},
	{
		comment:     "Supports jpg",
		inputValue:  "jpg",
		expectError: nil,
	},
	{
		comment:     "Supports jpeg",
		inputValue:  "jpg",
		expectError: nil,
	},
	{
		comment:     "Uppercase letters",
		inputValue:  "JPG",
		expectError: nil,
	},
	{
		comment:     "Unsupported tiff",
		inputValue:  "tiff",
		expectError: UnsupportedImgFormatErr,
	},
	{
		comment:     "Unsupported 0011",
		inputValue:  "0011",
		expectError: UnsupportedImgFormatErr,
	},
}

func TestImgFormatValidator(t *testing.T) {
	for _, tc := range ImgFormatTestCase {
		t.Run(tc.comment, func(t *testing.T) {
			got := ImgFormatValidator(tc.inputValue)
			if !errors.Is(got, tc.expectError) {
				t.Errorf("%s test. want: %v, got: %v", tc.comment, tc.expectError, got)
			}
		})
	}
}

func TestFilenameValidator(t *testing.T) {
	for _, tc := range FileSubstringTestCase {
		t.Run(tc.comment, func(t *testing.T) {
			got := FilenameValidator(tc.inputValue)
			if !errors.Is(got, tc.expectError) {
				t.Errorf("%s test. want: %v, got: %v", tc.comment, tc.expectError, got)
			}
		})
	}
}
