// Package config handles application configuration via .env files.
package config

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// Config holds application configuration.
type Config struct {
	ClaudeTier string // Claude subscription tier (free, pro, max_5x, max_20x)
	NoColor    bool   // Disable colors in output
	Width      int    // Progress bar width
}

// DefaultConfig returns default configuration.
func DefaultConfig() *Config {
	return &Config{
		ClaudeTier: "pro",
		NoColor:    false,
		Width:      42,
	}
}

// Load reads configuration from .env file in the given directory.
func Load(dir string) (*Config, error) {
	cfg := DefaultConfig()

	envPath := filepath.Join(dir, ".env")
	file, err := os.Open(envPath)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil // Use defaults if no .env
		}
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip comments and empty lines
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse KEY=VALUE
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove quotes if present
		value = strings.Trim(value, `"'`)

		switch key {
		case "CLAUDE_TIER":
			cfg.ClaudeTier = value
		case "NO_COLOR":
			cfg.NoColor = value == "1" || strings.ToLower(value) == "true"
		case "PROGRESS_WIDTH":
			// Simple integer parsing
			width := 0
			for _, c := range value {
				if c >= '0' && c <= '9' {
					width = width*10 + int(c-'0')
				}
			}
			if width >= 20 && width <= 100 {
				cfg.Width = width
			}
		}
	}

	return cfg, scanner.Err()
}

// LoadFromWorkingDir loads config from current working directory.
func LoadFromWorkingDir() (*Config, error) {
	dir, err := os.Getwd()
	if err != nil {
		return DefaultConfig(), nil
	}
	return Load(dir)
}
