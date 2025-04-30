package main

import "fmt"

type Color string

const (
	Reset      Color = "\033[0m"
	Bold             = "\033[1m"
	ColorRed   Color = "\033[31m"
	ColorGreen Color = "\033[32m"
	ColorBlue  Color = "\033[34m"
)

// color apply color and optional bold formatting to text
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

type formattedLabel interface {
	~string | ~float64
}

// fbg formatting text in bold and color green
func fbg[T formattedLabel](label T, noFormat bool) string {
	if s, ok := any(label).(string); ok {
		return color(s, ColorGreen, true, noFormat)
	}
	return color(fmt.Sprintf("%.2f", any(label).(float64)), ColorGreen, true, noFormat)
}
