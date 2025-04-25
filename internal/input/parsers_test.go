package input

import (
	"errors"
	"reflect"
	"testing"
)

type validPageListTestCase struct {
	comment     string
	inputValue  string
	pageCount   int
	expectedVal []int
	expectError error
}

var PageListTestCase = []validPageListTestCase{
	{
		comment:     "Consecutive",
		inputValue:  "1,2,3,4",
		pageCount:   100,
		expectedVal: []int{1, 2, 3, 4},
		expectError: nil,
	},
	{
		comment:     "Page range",
		inputValue:  "17-20",
		pageCount:   100,
		expectedVal: []int{17, 18, 19, 20},
		expectError: nil,
	},
	{
		comment:     "Page range one page",
		inputValue:  "25-25",
		pageCount:   100,
		expectedVal: []int{25},
		expectError: nil,
	},
	{
		comment:     "Page list and single pages mixed",
		inputValue:  "2, 4, 10, 20-22, 50",
		pageCount:   100,
		expectedVal: []int{2, 4, 10, 20, 21, 22, 50},
		expectError: nil,
	},
	{
		comment:     "Page list and single pages mixed two ranges",
		inputValue:  "2, 4, 10, 20-22, 50, 65-67",
		pageCount:   100,
		expectedVal: []int{2, 4, 10, 20, 21, 22, 50, 65, 66, 67},
		expectError: nil,
	},
	{
		comment:     "Only unique pages",
		inputValue:  "2,3,4,4,5,7-10, 9-10, 15, 15",
		pageCount:   100,
		expectedVal: []int{2, 3, 4, 5, 7, 8, 9, 10, 15},
		expectError: PageOutofRangeErr,
	},
	{
		comment:     "Page range incorrect",
		inputValue:  "10-5",
		pageCount:   100,
		expectedVal: []int{},
		expectError: PageStartGreaterErr,
	},
	{
		comment:     "Page out of range zero",
		inputValue:  "0-5",
		pageCount:   100,
		expectedVal: []int{},
		expectError: PageOutofRangeErr,
	},
	{
		comment:     "Page negative",
		inputValue:  "-1",
		pageCount:   100,
		expectedVal: []int{},
		expectError: DashOnBoundariesErr,
	},
	{
		comment:     "Page number must be integer",
		inputValue:  "1, qw, 24",
		pageCount:   100,
		expectedVal: []int{},
		expectError: PageNotIntErr,
	},
	{
		comment:     "Dash on boundary 1",
		inputValue:  "2, 4, 10, -20, 50",
		pageCount:   100,
		expectedVal: []int{},
		expectError: DashOnBoundariesErr,
	},
	{
		comment:     "Dash on boundary 2",
		inputValue:  "2, 4, 10, 22-, 50",
		pageCount:   100,
		expectedVal: []int{},
		expectError: DashOnBoundariesErr,
	},
	{
		comment:     "Dash on boundary 3",
		inputValue:  "2, 4, 10-, 20-22, 50",
		pageCount:   100,
		expectedVal: []int{},
		expectError: DashOnBoundariesErr,
	},
	{
		comment:     "Only one dash allowed in the range",
		inputValue:  "2, 4, 10, 20--22, 50",
		pageCount:   100,
		expectedVal: []int{},
		expectError: DoubleDashErr,
	},
}

var ImgSizeTestCase = []validParamTestCase{
	{
		comment:     "Size is valid",
		inputValue:  "320x240",
		expectError: nil,
	},
	{
		comment:     "Wrong delimiter",
		inputValue:  "320c240",
		expectError: NoXErr,
	},
	{
		comment:     "Size not int",
		inputValue:  "320x240p",
		expectError: SizeMustBeIntErr,
	},
	{
		comment:     "Only text",
		inputValue:  "qwerty",
		expectError: NoXErr,
	},
	{
		comment:     "Only digits",
		inputValue:  "123",
		expectError: NoXErr,
	},
	{
		comment:     "Negative size Y",
		inputValue:  "640x-480",
		expectError: SizeMustBePositive,
	},
	{
		comment:     "Negative size X",
		inputValue:  "-640x480",
		expectError: SizeMustBePositive,
	},
	{
		comment:     "Negative size both",
		inputValue:  "-640x-480",
		expectError: SizeMustBePositive,
	},
	{
		comment:     "Zero size X",
		inputValue:  "0x480",
		expectError: SizeMustBePositive,
	},
	{
		comment:     "Zero size Y",
		inputValue:  "640x0",
		expectError: SizeMustBePositive,
	},
	{
		comment:     "Zero size both",
		inputValue:  "0x0",
		expectError: SizeMustBePositive,
	},
}

func TestPagesExtractor(t *testing.T) {
	for _, tc := range PageListTestCase {
		t.Run(tc.comment, func(t *testing.T) {
			got, err := PagesExtractor(tc.inputValue, tc.pageCount)
			if err != nil {
				if !errors.Is(err, tc.expectError) {
					t.Errorf("%s test. want: %v, got: %v", tc.comment, tc.expectError, err)
				}
			} else {
				if !reflect.DeepEqual(got, tc.expectedVal) {
					t.Errorf("%s test. want: %v, got: %v", tc.comment, tc.expectedVal, got)
				}
			}
		})
	}
}

func TestImgSizeExtractor(t *testing.T) {
	for _, tc := range ImgSizeTestCase {
		t.Run(tc.comment, func(t *testing.T) {
			_, _, err := ImgSizeExtractor(tc.inputValue)
			if !errors.Is(err, tc.expectError) {
				t.Errorf("%s test. want: %v, got: %v", tc.comment, tc.expectError, err)
			}
		})
	}
}
