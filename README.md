# vibecode-monitor

A CLI tool for tracking **Claude Code** usage limits with ASCII art display, designed for fastfetch integration.

## Features

- 🎨 **Figlet ASCII Art** - Beautiful "small" font header
- 📊 **Progress Bars** - 3-line bars with time-until-reset and percentage
- 🔍 **Auto-Tier Detection** - Reads from `~/.claude/.credentials.json`
- � **Claude Orange Theme** - Authentic Claude branding colors
- ⚡ **Fast** - Efficient JSONL parsing with no external dependencies
- 🔒 **Privacy-First** - 100% local, no network requests

## Installation

```bash
git clone https://github.com/injaneity/vibecode-monitor
cd vibecode-monitor
go build -o vibecode-monitor ./cmd/vibecode-monitor

# Optional: Install to PATH
sudo mv vibecode-monitor /usr/local/bin/
```

## Configuration

```bash
cp .env.example .env
```

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `CLAUDE_TIER` | `pro` | Subscription tier: `free`, `pro`, `max_5x`, `max_20x`, or `auto` |
| `NO_COLOR` | — | Set to `1` to disable colors |
| `PROGRESS_WIDTH` | `42` | Width of the progress bar (20-100) |

## Usage

```bash
vibecode-monitor           # Full display with ASCII art
vibecode-monitor --compact # Single-line format for status bars
vibecode-monitor --tier max_5x  # Override tier
```

### Command Line Options

```
-tier string      Override subscription tier
-compact          Output single-line format
-no-color         Disable colored output
-width int        Progress bar width (default 42)
-version          Print version and exit
```

## Tier Limits

| Tier | 5h Prompts | Weekly Sonnet | Weekly Opus |
|------|------------|---------------|-------------|
| Free/Pro | 10-40 | 40-80h | — |
| Max 5x | 50-200 | 140-280h | 15-35h |
| Max 20x | 200-800 | 240-480h | 24-40h |

## Fastfetch Integration

> Requires fastfetch **2.x or later**.

Add to `~/.config/fastfetch/config.jsonc`:

```jsonc
{
    "modules": [
        {
            "type": "command",
            "key": "Claude",
            "text": "vibecode-monitor --compact --no-color"
        }
    ]
}
```

For full ASCII art display:

```jsonc
{
    "type": "command",
    "key": " ",
    "keyWidth": 0,
    "text": "vibecode-monitor",
    "multithreading": true
}
```

## How It Works

1. Parses `~/.claude/projects/**/*.jsonl` session files
2. Calculates 5-hour prompt cycles and weekly model hours
3. Auto-detects tier from `~/.claude/.credentials.json`
4. Displays with figlet header and progress bar

## License

MIT License

## Credits

Inspired by [claude-code-limit-tracker](https://github.com/TylerGallenbeck/claude-code-limit-tracker)
