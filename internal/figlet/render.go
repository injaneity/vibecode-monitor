// Package figlet provides ASCII art text rendering using the "small" figlet font.
package figlet

import (
	"strings"
)

// Render converts a string into figlet ASCII art using the small font.
// Returns a multi-line string with the rendered text.
func Render(text string) string {
	height := FontHeight()
	lines := make([]strings.Builder, height)

	for _, char := range text {
		pattern, ok := SmallFont[char]
		if !ok {
			// Use space for unknown characters
			pattern = SmallFont[' ']
		}

		// Find the max width for this character
		maxWidth := 0
		for _, line := range pattern {
			if len(line) > maxWidth {
				maxWidth = len(line)
			}
		}

		// Add each line, padding to max width
		for i := 0; i < height; i++ {
			if i < len(pattern) {
				line := pattern[i]
				// Pad with spaces to consistent width
				for len(line) < maxWidth {
					line += " "
				}
				lines[i].WriteString(line)
			} else {
				// Add empty space if pattern is shorter than height
				lines[i].WriteString(strings.Repeat(" ", maxWidth))
			}
		}
	}

	result := make([]string, height)
	for i, line := range lines {
		result[i] = line.String()
	}

	return strings.Join(result, "\n")
}

// RenderColored renders text with ANSI color codes.
func RenderColored(text string, colorCode string) string {
	rendered := Render(text)
	if colorCode == "" {
		return rendered
	}

	lines := strings.Split(rendered, "\n")
	for i, line := range lines {
		lines[i] = colorCode + line + "\033[0m"
	}
	return strings.Join(lines, "\n")
}

// GetWidth returns the total width of rendered text.
func GetWidth(text string) int {
	if len(text) == 0 {
		return 0
	}

	width := 0
	for _, char := range text {
		pattern, ok := SmallFont[char]
		if !ok {
			pattern = SmallFont[' ']
		}
		if len(pattern) > 0 {
			width += len(pattern[0])
		}
	}
	return width
}
