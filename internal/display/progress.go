// Package display handles terminal output formatting with colors and progress bars.
package display

import (
	"fmt"
	"strings"
)

// ProgressBar represents a 3-line progress bar with borders.
type ProgressBar struct {
	Width      int     // Total width including borders
	Current    float64 // Current value
	Total      float64 // Maximum value
	TimeLeft   string  // e.g., "2h 15m"
	Unit       string  // e.g., "h" for hours
	ShowCost   bool    // If true, show as currency
	CostPrefix string  // e.g., "$"
	NoColor    bool    // Disable colors
	HideTimer  bool    // Hide the timer in top border
	BarColor   string  // Custom bar color (ANSI code)
}

// Box drawing characters
const (
	TopLeft     = "┌"
	TopRight    = "┐"
	BottomLeft  = "└"
	BottomRight = "┘"
	Horizontal  = "─"
	Vertical    = "│"
	FillBlock   = "█"
	EmptyBlock  = "░"
)

// Render creates a 3-line progress bar string.
func (p *ProgressBar) Render() string {
	if p.Width < 10 {
		p.Width = 42 // Default width
	}

	innerWidth := p.Width - 2 // Subtract border characters

	// Calculate progress
	percentage := 0.0
	if p.Total > 0 {
		percentage = (p.Current / p.Total) * 100
	}
	if percentage > 100 {
		percentage = 100
	}

	filledCount := int(float64(innerWidth) * (percentage / 100))
	if filledCount > innerWidth {
		filledCount = innerWidth
	}
	emptyCount := innerWidth - filledCount

	// Build top line with orange time value
	topLine := p.buildTopBorderLine(innerWidth)

	// Build the progress bar middle line
	filled := strings.Repeat(FillBlock, filledCount)
	empty := strings.Repeat(EmptyBlock, emptyCount)
	var middleLine string
	barColor := p.BarColor
	if barColor == "" {
		barColor = ClaudeOrange // Default
	}
	if p.NoColor {
		middleLine = Vertical + filled + empty + Vertical
	} else {
		middleLine = BoxColor + Vertical + Reset + barColor + filled + Reset + Gray + empty + Reset + BoxColor + Vertical + Reset
	}

	// Build bottom line with orange current value
	bottomLine := p.buildBottomBorderLine(percentage, innerWidth)

	return topLine + "\n" + middleLine + "\n" + bottomLine
}

// buildTopBorderLine creates top border with timer or plain border.
func (p *ProgressBar) buildTopBorderLine(innerWidth int) string {
	// If HideTimer is set, just draw a plain border
	if p.HideTimer || p.TimeLeft == "" {
		dashes := strings.Repeat(Horizontal, innerWidth)
		if p.NoColor {
			return TopLeft + dashes + TopRight
		}
		return BoxColor + TopLeft + dashes + TopRight + Reset
	}

	// Format: " Xh Ym until reset "
	timeStr := p.TimeLeft
	suffix := " until reset"
	fullLabel := " " + timeStr + suffix + " "
	labelLen := len([]rune(fullLabel))

	if labelLen >= innerWidth {
		fullLabel = fullLabel[:innerWidth-3] + "..."
		labelLen = innerWidth
	}

	remaining := innerWidth - labelLen
	leftPad := remaining / 2
	rightPad := remaining - leftPad

	leftDashes := strings.Repeat(Horizontal, leftPad)
	rightDashes := strings.Repeat(Horizontal, rightPad)

	if p.NoColor {
		return TopLeft + leftDashes + fullLabel + rightDashes + TopRight
	}

	// Colored: orange time, default suffix
	return BoxColor + TopLeft + leftDashes + Reset +
		" " + ClaudeOrange + timeStr + Reset + DimWhite + suffix + Reset + " " +
		BoxColor + rightDashes + TopRight + Reset
}

// buildBottomBorderLine creates bottom border with styled percentage and values.
func (p *ProgressBar) buildBottomBorderLine(percentage float64, innerWidth int) string {
	var fullLabel string
	if p.ShowCost {
		fullLabel = fmt.Sprintf(" %.1f%% (%s%.2f / %s%.2f) ",
			percentage, p.CostPrefix, p.Current, p.CostPrefix, p.Total)
	} else {
		fullLabel = fmt.Sprintf(" %.1f%% (%.1f / %.1f%s) ",
			percentage, p.Current, p.Total, p.Unit)
	}
	labelLen := len([]rune(fullLabel))

	if labelLen >= innerWidth {
		fullLabel = fullLabel[:innerWidth-3] + "..."
		labelLen = innerWidth
	}

	remaining := innerWidth - labelLen
	leftPad := remaining / 2
	rightPad := remaining - leftPad

	leftDashes := strings.Repeat(Horizontal, leftPad)
	rightDashes := strings.Repeat(Horizontal, rightPad)

	if p.NoColor {
		return BottomLeft + leftDashes + fullLabel + rightDashes + BottomRight
	}

	// Use BarColor for the current value, fallback to orange
	valColor := p.BarColor
	if valColor == "" {
		valColor = ClaudeOrange
	}

	// Colored version with highlighted current value
	var coloredLabel string
	if p.ShowCost {
		coloredLabel = fmt.Sprintf(" %.1f%% ("+valColor+"%s%.2f"+Reset+DimWhite+" / %s%.2f) "+Reset,
			percentage, p.CostPrefix, p.Current, p.CostPrefix, p.Total)
	} else {
		coloredLabel = fmt.Sprintf(DimWhite+" %.1f%% ("+Reset+valColor+"%.1f"+Reset+DimWhite+" / %.1f%s) "+Reset,
			percentage, p.Current, p.Total, p.Unit)
	}

	return BoxColor + BottomLeft + leftDashes + Reset +
		coloredLabel +
		BoxColor + rightDashes + BottomRight + Reset
}

// NewProgressBar creates a progress bar with sensible defaults.
func NewProgressBar(current, total float64, timeLeft string) *ProgressBar {
	return &ProgressBar{
		Width:    42,
		Current:  current,
		Total:    total,
		TimeLeft: timeLeft,
		Unit:     "h",
	}
}
