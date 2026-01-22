# vibe-monitor

> A standalone CLI tool for tracking Claude Code usage limits with live monitoring

[![License: MIT](https://img.shields.io/badge/License-MIT-orange.svg)](LICENSE)

## ✨ Features

- 🎨 **Figlet ASCII Art** - Beautiful "small" font header with Claude orange branding
- 📊 **Progress Bars** - 3-line bars showing current usage, limits, and time until reset
- 🔄 **Watch Mode** - Auto-refresh display every N seconds for live monitoring
- 🔍 **Auto-Tier Detection** - Automatically detects your tier from `~/.claude/.credentials.json`
- 🧡 **Claude Orange Theme** - Authentic Claude branding colors throughout
- ⚡ **Fast & Efficient** - Local JSONL parsing with zero external dependencies
- 🔒 **Privacy-First** - 100% local processing, no network requests ever

## 📦 Installation

```bash
git clone https://github.com/injaneity/vibe-monitor
cd vibe-monitor

# Install to ~/.local/bin (recommended)
make install-user

# OR install system-wide (requires sudo)
make install

# Verify installation
vibe-monitor --version
```

## 🚀 Quick Start

```bash
# Show current usage once
vibe-monitor

# Watch mode - refresh every 5 seconds
vibe-monitor --refresh 5

# Compact single-line format
vibe-monitor --compact

# Disable colors
vibe-monitor --no-color
```

## 📖 Usage

### Basic Usage

```bash
vibe-monitor
```

Output:
```
  ___  _                   _             ___            _       
 / __|| |  __ _  _  _   __| |  ___      / __|  ___   __| |  ___ 
| (__ | | / _` || || | / _` | / -_)    | (__  / _ \ / _` | / -_)
 \___||_| \__,_| \_,_| \__,_| \___|     \___| \___/ \__,_| \___|

    Sonnet: 5.0 / 80.0h

┌───────── 77h 29m until reset ──────────┐
│██░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░│
└────────── 6.2% (5.0 / 80.0h) ──────────┘
```

### Watch Mode

Monitor your usage in real-time with automatic updates:

```bash
# Refresh every 5 seconds (default interval)
vibe-monitor --refresh 5

# Fast updates (every 1 second)
vibe-monitor --refresh 1

# Compact watch mode
vibe-monitor --refresh 5 --compact
```

Press `Ctrl+C` to stop monitoring.

### Command Line Options

```
Options:
  -tier string          Subscription tier (free, pro, max_5x, max_20x, auto)
  -compact              Single-line compact format
  -no-color             Disable colored output
  -width int            Progress bar width (default 42)
  -refresh int          Auto-refresh every N seconds (0=disabled)
  -version              Print version and exit
```

## ⚙️ Configuration

Create a `.env` file in the project directory for persistent settings:

```bash
cp .env.example .env
```

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `CLAUDE_TIER` | `auto` | Subscription tier: `free`, `pro`, `max_5x`, `max_20x`, or `auto` |
| `NO_COLOR` | — | Set to `1` to disable colors |
| `PROGRESS_WIDTH` | `42` | Width of the progress bar (20-100) |

## 📊 Tier Limits

| Tier | 5h Prompts | Weekly Sonnet | Weekly Opus |
|------|------------|---------------|-------------|
| Free/Pro | 10-40 | 40-80h | — |
| Max 5x | 50-200 | 140-280h | 15-35h |
| Max 20x | 200-800 | 240-480h | 24-40h |

*Note: Limits vary based on your specific subscription. vibe-monitor auto-detects your tier.*

## 🛠️ How It Works

1. **Scans** `~/.claude/projects/**/*.jsonl` session files
2. **Calculates** 5-hour prompt cycles and weekly model hours from session data
3. **Detects** your tier automatically from `~/.claude/.credentials.json`
4. **Displays** current usage with beautiful ASCII art and progress bars

All processing happens locally on your machine. No data is sent anywhere.

## 🎯 Examples

**Quick check:**
```bash
vibe-monitor --compact
# Output: Claude: 5.0/80.0h (6%) | 77h 29m
```

**Monitor during heavy usage:**
```bash
vibe-monitor --refresh 10
# Updates every 10 seconds, clears screen, shows full display
```

**Override tier detection:**
```bash
vibe-monitor --tier max_5x
# Forces Max 5x tier limits even if auto-detection differs
```

**No colors for piping:**
```bash
vibe-monitor --no-color --compact >> usage.log
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📝 License

MIT License - see [LICENSE](LICENSE) file for details

## 💙 Credits

Built with ❤️ for the Claude Code community

---

**Not affiliated with Anthropic. Claude and Claude Code are trademarks of Anthropic, PBC.**
