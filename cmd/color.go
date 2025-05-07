package main

import "fmt"

// Color represents an ANSI escape code string used for terminal text formatting and coloring
type Color string

const (
	// Reset resets all terminal formatting to default
	Reset Color = "\033[0m"
	// Bold applies bold formatting to terminal text
	Bold = "\033[1m"
	// ColorGreen sets the terminal text color to green
	ColorGreen Color = "\033[32m"
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
