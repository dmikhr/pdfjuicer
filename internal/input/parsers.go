package input

import (
	"errors"
	"sort"
	"strconv"
	"strings"
)

var (
	// ErrDoubleDash is returned when more than one dash is found in the page range specification.
	// example: incorrect 5--10, correct: 5-10
	ErrDoubleDash = errors.New("only one dash allowed in the range")
	// ErrDashOnBoundaries is returned when the page range begins or ends with a dash
	ErrDashOnBoundaries = errors.New("range can't begin or end with a dash")
	// ErrPageNotInt is returned when the page number provided cannot be parsed as an integer
	ErrPageNotInt = errors.New("page number must be integer")
	// ErrPageStartGreater is returned when the start page number is greater than the end page number in a document
	ErrPageStartGreater = errors.New("start page can't be greater than the final")
	// ErrPageOutofRange is returned when the specified page number falls outside the page range of a document
	ErrPageOutofRange = errors.New("page out of range")
)

// PagesExtractor parses user input of custom pages to extract
// example: 1,4,5-8,10
func PagesExtractor(s string, pageCount int) ([]int, error) {
	sNoSpaces := strings.ReplaceAll(s, " ", "")
	pagesChunks := strings.Split(sNoSpaces, ",")

	var dashPos, pageNum int
	var err error
	var pagesList []int
	var pages []string

	for _, pageData := range pagesChunks {
		dashPos = strings.Index(pageData, "-")
		// if no dash then this is not a range - just a page number
		if dashPos == -1 {
			pageNum, err = strconv.Atoi(pageData)
			if err != nil {
				return []int{}, ErrPageNotInt
			} else if isOutOfRange(pageNum, pageCount) {
				return []int{}, ErrPageOutofRange
			}
			pagesList = append(pagesList, pageNum)
		} else if strings.Count(pageData, "-") > 1 {
			return []int{}, ErrDoubleDash
		} else if dashPos == 0 || dashPos == len(pageData)-1 {
			return []int{}, ErrDashOnBoundaries
		} else if strings.Count(pageData, "-") == 1 {
			pages = strings.Split(pageData, "-")
			pageStart, err1 := strconv.Atoi(pages[0])
			pageEnd, err2 := strconv.Atoi(pages[1])

			if err1 != nil || err2 != nil {
				return []int{}, ErrPageNotInt
			} else if pageStart > pageEnd {
				return []int{}, ErrPageStartGreater
			} else if isOutOfRange(pageStart, pageCount) || isOutOfRange(pageEnd, pageCount) {
				return []int{}, ErrPageOutofRange
			}
			expandRange(pageStart, pageEnd, &pagesList)
		}
	}
	sort.Ints(pagesList)

	return uniqueFromSortedAsc(pagesList), nil
}

// uniqueFromSortedAsc leaves only unique pages in case if user submitted duplicate pages
func uniqueFromSortedAsc(s []int) []int {
	var uniqueItems []int
	uniqueItems = append(uniqueItems, s[0])

	for i := 1; i < len(s); i++ {
		if s[i] != s[i-1] {
			uniqueItems = append(uniqueItems, s[i])
		}
	}

	return uniqueItems
}

// expandRange transforms ranges into consecutive page numbers
// example: 2-5 -> 2,3,4,5
func expandRange(pageStart, pageEnd int, pagesNum *[]int) {
	for i := pageStart; i <= pageEnd; i++ {
		*pagesNum = append(*pagesNum, i)
	}
}

// isOutOfRange checks if submitted page is out of range
func isOutOfRange(pageNum, pageCount int) bool {
	return pageNum <= 0 || pageNum > pageCount
}

var (
	// ErrNoX is returned when the image size is expected to have an 'x' (e.g., "800x600"),
	// but it is missing from the provided string.
	ErrNoX = errors.New("no x in image size")
	// ErrSizeMustBeInt is returned when the image size value cannot be parsed as an integer
	ErrSizeMustBeInt = errors.New("image size must be int")
	// ErrSizeMustBePositive is returned when the image size is zero or a negative number,
	// but a positive value is required.
	ErrSizeMustBePositive = errors.New("image size cannot be negative")
)

// ImgSizeExtractor parses submitted image size like 256x480 into x,y integers
func ImgSizeExtractor(s string) (int, int, error) {
	imgSize := strings.Split(s, "x")
	if len(imgSize) != 2 {
		return -1, -1, ErrNoX
	}
	x, errX := strconv.Atoi(imgSize[0])
	y, errY := strconv.Atoi(imgSize[1])
	if errX != nil || errY != nil {
		return -1, -1, ErrSizeMustBeInt
	}
	if x <= 0 || y <= 0 {
		return -1, -1, ErrSizeMustBePositive
	}
	return x, y, nil
}
