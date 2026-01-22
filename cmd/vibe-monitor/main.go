// vibe-monitor is a standalone CLI tool for tracking Claude Code usage limits.
package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/injaneity/vibe-monitor/internal/claude"
	"github.com/injaneity/vibe-monitor/internal/config"
	"github.com/injaneity/vibe-monitor/internal/display"
)

var (
	version = "dev"
)

func main() {
	tierFlag := flag.String("tier", "", "Subscription tier (free, pro, max_5x, max_20x, auto)")
	compactFlag := flag.Bool("compact", false, "Single-line compact format")
	noColorFlag := flag.Bool("no-color", false, "Disable colored output")
	widthFlag := flag.Int("width", 42, "Progress bar width (20-100)")
	refreshFlag := flag.Int("refresh", 0, "Auto-refresh every N seconds (0=disabled)")
	versionFlag := flag.Bool("version", false, "Print version and exit")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("vibe-monitor %s\n", version)
		os.Exit(0)
	}

	cfg := loadConfig()

	if *tierFlag != "" {
		cfg.ClaudeTier = *tierFlag
	}
	if *noColorFlag || os.Getenv("NO_COLOR") != "" {
		cfg.NoColor = true
	}
	if *widthFlag != 42 {
		cfg.Width = *widthFlag
	}

	if cfg.ClaudeTier == "" || cfg.ClaudeTier == "auto" {
		if tier, err := claude.DetectTier(); err == nil && tier != "" {
			cfg.ClaudeTier = tier
		} else {
			cfg.ClaudeTier = "pro"
		}
	}

	if *refreshFlag > 0 {
		runWatchMode(cfg, *refreshFlag, *compactFlag)
	} else {
		displayOnce(cfg, *compactFlag)
	}
}

func runWatchMode(cfg *config.Config, interval int, compact bool) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	if !compact {
		fmt.Print("\033[2J\033[?25l")
		defer fmt.Print("\033[?25h")
	}

	displayOnce(cfg, compact)

	for {
		select {
		case <-ticker.C:
			if !compact {
				fmt.Print("\033[2J\033[H")
			}
			displayOnce(cfg, compact)
		case <-sigChan:
			if !compact {
				fmt.Print("\033[?25h")
			}
			fmt.Println("\nMonitoring stopped.")
			return
		}
	}
}

func displayOnce(cfg *config.Config, compact bool) {
	output := display.NewOutput(cfg.NoColor, cfg.Width)
	tracker := claude.NewTracker(cfg.ClaudeTier)

	usage, err := tracker.Calculate()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}

	if usage.SessionsCount == 0 {
		fmt.Println("No Claude Code usage data found.")
		fmt.Println("Session files: ~/.claude/projects/")
		return
	}

	if compact {
		fmt.Println(output.RenderCompact(usage))
	} else {
		fmt.Print(output.Render(usage))
	}
}

func loadConfig() *config.Config {
	if cfg, err := config.LoadFromWorkingDir(); err == nil && cfg != nil {
		return cfg
	}

	if home, err := os.UserHomeDir(); err == nil {
		if cfg, err := config.Load(filepath.Join(home, ".config", "vibe-monitor")); err == nil && cfg != nil {
			return cfg
		}
	}

	if execPath, err := os.Executable(); err == nil {
		if cfg, err := config.Load(filepath.Dir(execPath)); err == nil && cfg != nil {
			return cfg
		}
	}

	return config.DefaultConfig()
}
