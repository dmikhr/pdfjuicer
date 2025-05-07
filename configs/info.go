package main

import "fmt"

const (
	version    = "1.0.0"
	author     = "Dmitrii Khramtsov"
	repository = "https://github.com/dmikhr/pdfjuicer"
	license    = "AGPL-3.0"
)

func about() string {
	return fmt.Sprintf(`pdfjuicer v%s
A command-line tool for extracting pages from a PDF file as images.
Author: %s
Repository: %s
License: %s
`, version, author, repository, license)
}
