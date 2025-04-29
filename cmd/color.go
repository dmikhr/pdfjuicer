package main

type Color string

const (
	Reset      Color = "\033[0m"
	Bold             = "\033[1m"
	ColorRed   Color = "\033[31m"
	ColorGreen Color = "\033[32m"
	ColorBlue  Color = "\033[34m"
)

func color(text string, color Color, bold, noFormat bool) string {
	if noFormat {
		return text
	}
	mode := string(color)
	if bold {
		mode += Bold
	}
	return mode + text + string(Reset)
}
