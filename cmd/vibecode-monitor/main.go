// vibecode-monitor is a CLI tool for tracking Claude Code usage limits.
// It displays usage statistics with ASCII art headers and progress bars,
// designed for integration with fastfetch.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/injaneity/vibecode-monitor/internal/claude"
	"github.com/injaneity/vibecode-monitor/internal/config"
	"github.com/injaneity/vibecode-monitor/internal/display"
)

var (
	version = "dev"
)

func main() {
	// Parse flags
	tierFlag := flag.String("tier", "", "Override Claude subscription tier (free, pro, max_5x, max_20x)")
	compactFlag := flag.Bool("compact", false, "Output single-line compact format")
	noColorFlag := flag.Bool("no-color", false, "Disable colored output")
	widthFlag := flag.Int("width", 42, "Progress bar width")
	versionFlag := flag.Bool("version", false, "Print version and exit")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("vibecode-monitor %s\n", version)
		os.Exit(0)
	}

	// Load configuration
	cfg := loadConfig()

	// Apply flag overrides
	if *tierFlag != "" {
		cfg.ClaudeTier = *tierFlag
	}
	if *noColorFlag {
		cfg.NoColor = true
	}
	if *widthFlag != 42 {
		cfg.Width = *widthFlag
	}

	// Check for NO_COLOR environment variable
	if os.Getenv("NO_COLOR") != "" {
		cfg.NoColor = true
	}

	// Try to detect tier from credentials if not set
	if cfg.ClaudeTier == "" || cfg.ClaudeTier == "auto" {
		if detectedTier, err := claude.DetectTier(); err == nil && detectedTier != "" {
			cfg.ClaudeTier = detectedTier
		} else {
			cfg.ClaudeTier = "pro" // Default fallback
		}
	}

	// Create output renderer
	output := display.NewOutput(cfg.NoColor, cfg.Width)

	// Calculate Claude Code usage
	tracker := claude.NewTracker(cfg.ClaudeTier)
	usage, err := tracker.Calculate()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if usage.SessionsCount == 0 {
		fmt.Println("No Claude Code usage data found.")
		fmt.Println("")
		fmt.Println("Usage data is created when you use Claude Code CLI.")
		fmt.Println("Session files are stored in ~/.claude/projects/")
		os.Exit(0)
	}

	// Render output
	if *compactFlag {
		fmt.Println(output.RenderCompact(usage))
	} else {
		fmt.Print(output.Render(usage))
	}
}

// loadConfig loads configuration from multiple locations.
func loadConfig() *config.Config {
	// Try current directory first
	cfg, err := config.LoadFromWorkingDir()
	if err == nil && cfg != nil {
		return cfg
	}

	// Try home directory
	home, err := os.UserHomeDir()
	if err == nil {
		cfg, err = config.Load(filepath.Join(home, ".config", "vibecode-monitor"))
		if err == nil && cfg != nil {
			return cfg
		}
	}

	// Try executable directory
	execPath, err := os.Executable()
	if err == nil {
		cfg, err = config.Load(filepath.Dir(execPath))
		if err == nil && cfg != nil {
			return cfg
		}
	}

	return config.DefaultConfig()
}
