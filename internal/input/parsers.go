package input

import (
	"errors"
	"sort"
	"strconv"
	"strings"
)

var (
	DoubleDashErr       = errors.New("only one dash allowed in the range")
	DashOnBoundariesErr = errors.New("range can't begin or end with a dash")
	PageNotIntErr       = errors.New("page number must be integer")
	PageStartGreaterErr = errors.New("start page can't be greater than the final")
	PageOutofRangeErr   = errors.New("page out of range")
)

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
				return []int{}, PageNotIntErr
			} else if isOutOfRange(pageNum, pageCount) {
				return []int{}, PageOutofRangeErr
			} else {
				pagesList = append(pagesList, pageNum)
			}
		} else if strings.Count(pageData, "-") > 1 {
			return []int{}, DoubleDashErr
		} else if dashPos == 0 || dashPos == len(pageData)-1 {
			return []int{}, DashOnBoundariesErr
		} else if strings.Count(pageData, "-") == 1 {
			pages = strings.Split(pageData, "-")
			pageStart, err1 := strconv.Atoi(pages[0])
			pageEnd, err2 := strconv.Atoi(pages[1])

			if err1 != nil || err2 != nil {
				return []int{}, PageNotIntErr
			} else if pageStart > pageEnd {
				return []int{}, PageStartGreaterErr
			} else if isOutOfRange(pageStart, pageCount) || isOutOfRange(pageEnd, pageCount) {
				return []int{}, PageOutofRangeErr
			} else {
				expandRange(pageStart, pageEnd, &pagesList)
			}
		}
	}
	sort.Ints(pagesList)

	return uniqueFromSortedAsc(pagesList), nil
}

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

func expandRange(pageStart, pageEnd int, pagesNum *[]int) {
	for i := pageStart; i <= pageEnd; i++ {
		*pagesNum = append(*pagesNum, i)
	}
}

func isOutOfRange(pageNum, pageCount int) bool {
	return pageNum <= 0 || pageNum > pageCount
}
