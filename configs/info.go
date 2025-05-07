package config

import "fmt"

const (
	Version    = "1.0.0"
	Author     = "Dmitrii Khramtsov"
	Repository = "https://github.com/dmikhr/pdfjuicer"
	License    = "AGPL-3.0"
)

func About() string {
	return fmt.Sprintf(`pdfjuicer v%s
A command-line tool for extracting pages from a PDF file as images.
Author: %s
Repository: %s
License: %s
`, Version, Author, Repository, License)
}
