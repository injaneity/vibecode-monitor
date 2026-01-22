// Package display handles terminal output formatting with colors and progress bars.
package display

import (
	"fmt"
	"strings"

	"github.com/injaneity/vibe-monitor/internal/claude"
	"github.com/injaneity/vibe-monitor/internal/figlet"
)

// Output combines all display components into final terminal output.
type Output struct {
	NoColor bool
	Width   int
	Offset  int // Left padding for logo alignment
}

// NewOutput creates a new output renderer.
func NewOutput(noColor bool, width int) *Output {
	if width < 20 {
		width = 42
	}
	return &Output{
		NoColor: noColor,
		Width:   width,
	}
}

// SetOffset sets the left padding offset for logo alignment.
func (o *Output) SetOffset(offset int) {
	o.Offset = offset
}

// Render produces the complete display output for Claude Code usage.
func (o *Output) Render(usage *claude.UsageData) string {
	var sb strings.Builder

	// 1. Figlet ASCII art header
	header := o.renderHeader("Claude Code")
	if o.Offset > 0 {
		header = o.addOffset(header)
	}
	sb.WriteString(header)
	sb.WriteString("\n\n")

	// 2. Model-specific usage stats
	stats := o.renderModelStats(usage)
	if o.Offset > 0 {
		stats = o.addOffset(stats)
	}
	sb.WriteString(stats)
	sb.WriteString("\n\n")

	// 3. Progress bar
	bar := o.renderProgressBar(usage)
	if o.Offset > 0 {
		bar = o.addOffset(bar)
	}
	sb.WriteString(bar)
	sb.WriteString("\n")

	return sb.String()
}

// renderHeader creates the figlet ASCII art header.
func (o *Output) renderHeader(text string) string {
	if o.NoColor {
		return figlet.Render(text)
	}
	return figlet.RenderColored(text, ClaudeOrange)
}

// renderModelStats formats the model usage breakdown.
func (o *Output) renderModelStats(usage *claude.UsageData) string {
	var sb strings.Builder
	indent := "    "

	// Sonnet hours: "Sonnet: " white, current orange, " / Xh" white
	if o.NoColor {
		sonnetLine := fmt.Sprintf("%sSonnet: %.1f / %.1fh",
			indent, usage.WeeklySonnetHours, usage.Tier.WeeklySonnetMax)
		sb.WriteString(sonnetLine)
	} else {
		sb.WriteString(indent)
		sb.WriteString(White + "Sonnet: " + Reset)
		sb.WriteString(ClaudeOrange + fmt.Sprintf("%.1f", usage.WeeklySonnetHours) + Reset)
		sb.WriteString(White + fmt.Sprintf(" / %.1fh", usage.Tier.WeeklySonnetMax) + Reset)
	}

	// Opus hours (if available)
	if usage.Tier.HasOpus() {
		sb.WriteString("\n")
		if o.NoColor {
			opusLine := fmt.Sprintf("%sOpus:   %.1f / %.1fh",
				indent, usage.WeeklyOpusHours, usage.Tier.WeeklyOpusMax)
			sb.WriteString(opusLine)
		} else {
			sb.WriteString(indent)
			sb.WriteString(White + "Opus:   " + Reset)
			sb.WriteString(ClaudeOrange + fmt.Sprintf("%.1f", usage.WeeklyOpusHours) + Reset)
			sb.WriteString(White + fmt.Sprintf(" / %.1fh", usage.Tier.WeeklyOpusMax) + Reset)
		}
	}

	return sb.String()
}

// renderProgressBar creates the 3-line progress bar.
func (o *Output) renderProgressBar(usage *claude.UsageData) string {
	totalHours := usage.TotalWeeklyHours()
	maxHours := usage.Tier.GetTotalWeeklyMax()
	timeLeft := claude.FormatResetTime(usage.WeeklyResetIn)

	bar := &ProgressBar{
		Width:    o.Width,
		Current:  totalHours,
		Total:    maxHours,
		TimeLeft: timeLeft,
		Unit:     "h",
		NoColor:  o.NoColor,
	}

	return bar.Render()
}

// RenderCompact produces a single-line compact output for status bars.
func (o *Output) RenderCompact(usage *claude.UsageData) string {
	totalHours := usage.TotalWeeklyHours()
	maxHours := usage.Tier.GetTotalWeeklyMax()
	percentage := usage.WeeklyPercentage()
	resetTime := claude.FormatResetTime(usage.WeeklyResetIn)

	line := fmt.Sprintf("Claude: %.1f/%.1fh (%.0f%%) | %s",
		totalHours, maxHours, percentage, resetTime)

	if o.NoColor {
		return line
	}
	return Colorize(line, GetUsageColor(percentage))
}

// addOffset adds left padding to multi-line text.
func (o *Output) addOffset(text string) string {
	if o.Offset <= 0 {
		return text
	}
	padding := strings.Repeat(" ", o.Offset)
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		lines[i] = padding + line
	}
	return strings.Join(lines, "\n")
}
