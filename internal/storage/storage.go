package storage

import (
	"errors"
	"os"
)

type PathType int

const (
	IsFile PathType = iota
	IsDir
	IsOther
)

var (
	PathNotExistErr = errors.New("path doesn't exist")
	NotDirErr       = errors.New("not a directory")
)

func CheckPath(path string) (PathType, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return -1, PathNotExistErr
		}
		return -1, err
	} else if info.IsDir() {
		return IsDir, nil
	} else if info.Mode().IsRegular() {
		return IsFile, nil
	}
	return IsOther, nil
}
