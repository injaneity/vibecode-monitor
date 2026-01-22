// Package display handles terminal output formatting with colors and progress bars.
package display

// ANSI color codes for terminal output
const (
	Reset = "\033[0m"

	// Claude brand colors (rust orange theme)
	ClaudeOrange     = "\033[38;2;204;85;0m"    // #CC5500 - Rust orange
	ClaudeOrangeDark = "\033[38;2;166;65;0m"    // #A64100 - Darker rust
	ClaudeCream      = "\033[38;2;254;243;199m" // #FEF3C7 - Light cream
	ClaudeRust       = "\033[38;2;183;65;14m"   // #B7410E - Deep rust

	// Status colors
	Green  = "\033[38;2;34;197;94m" // Success/low usage
	Yellow = "\033[38;2;234;179;8m" // Warning/medium usage
	Red    = "\033[38;2;239;68;68m" // Critical/high usage

	// Neutral colors
	Gray     = "\033[38;2;156;163;175m" // Muted text
	White    = "\033[38;2;255;255;255m" // Bright text
	DimWhite = "\033[38;2;209;213;219m" // Slightly dimmed

	// Box drawing (for progress bar)
	BoxColor = "\033[38;2;100;116;139m" // Slate gray for borders
)

// Bold modifier
const Bold = "\033[1m"

// GetUsageColor returns an appropriate color based on usage percentage.
// 0-50%: Green, 50-75%: Yellow, 75-100%: Red
func GetUsageColor(percentage float64) string {
	switch {
	case percentage < 50:
		return Green
	case percentage < 75:
		return Yellow
	default:
		return Red
	}
}

// Colorize wraps text with a color code and reset.
func Colorize(text, color string) string {
	return color + text + Reset
}

// BoldColorize wraps text with bold and color.
func BoldColorize(text, color string) string {
	return Bold + color + text + Reset
}
