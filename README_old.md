# vibe-monitor

A standalone CLI tool for tracking **Claude Code** usage limits with ASCII art display and live monitoring.

## Features

- 🎨 **Figlet ASCII Art** - Beautiful "small" font header
- 📊 **Progress Bars** - 3-line bars with time-until-reset and percentage
- � **Watch Mode** - Auto-refresh display every N seconds
- 🔍 **Auto-Tier Detection** - Reads from `~/.claude/.credentials.json`
- 🧡 **Claude Orange Theme** - Authentic Claude branding colors
- ⚡ **Fast** - Efficient JSONL parsing with no external dependencies
- 🔒 **Privacy-First** - 100% local, no network requests

## Installation

```bash
git clone https://github.com/injaneity/vibe-monitor
cd vibe-monitor

# Install to ~/.local/bin
make install-user

# OR install system-wide (requires sudo)
make install

# Verify
vibe-monitor --version
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
vibe-monitor                  # Show current usage
vibe-monitor --refresh 5      # Update every 5 seconds (watch mode)
vibe-monitor --compact        # Single-line format
vibe-monitor --no-color       # Disable colors
vibe-monitor --tier max_5x    # Override tier
```

### Command Line Options

```
-tier string          Override subscription tier (free, pro, max_5x, max_20x, auto)
-compact              Output single-line format
-no-color             Disable colored output
-width int            Progress bar width (default 42)
-refresh int          Auto-refresh every N seconds (0=disabled)
-version              Print version and exit
```

## Tier Limits

| Tier | 5h Prompts | Weekly Sonnet | Weekly Opus |
|------|------------|---------------|-------------|
| Free/Pro | 10-40 | 40-80h | — |
| Max 5x | 50-200 | 140-280h | 15-35h |
| Max 20x | 200-800 | 240-480h | 24-40h |

## Fastfetch Integration

> Requires fastfetch **2.x or later**.

### Quick Setup

1. **Install vibecode-monitor to PATH:**
   ```bash
   cd vibecode-monitor
   make install-user  # Or make install
   ```

2. **Add to fastfetch config** (`~/.config/fastfetch/config.jsonc`):
   
   ```jsonc
   {
       "modules": [
           "title",
           "separator",
           "os",
           "host",
           {
               "type": "command",
               "key": "Claude",
               "text": "vibecode-monitor"
           }
       ]
   }
   ```

3. **That's it!** 
   ```bash
   fastfetch
   ```

   > **Auto-detection:** vibecode-monitor automatically uses compact mode when piped (like in fastfetch) and full ASCII art mode when run directly in your terminal. No configuration needed!

### Alternative Configurations

**Full ASCII art in terminal, compact in fastfetch (default behavior):**
- Just run `vibecode-monitor` - it auto-detects!
- In terminal: Shows full ASCII art with colors
- In fastfetch: Automatically uses compact single-line format

**Force full ASCII art mode even in fastfetch:**
```jsonc
{
    "logo": {
        "type": "none"
    },
    "modules": [
        "title",
        "separator",
        {
            "type": "command",
            "key": " ",
            "keyWidth": 0,
            "text": "vibecode-monitor --offset 20"
        }
    ]
}
```

**Force compact mode everywhere:**
```bash
vibecode-monitor --compact
```

### Troubleshooting

**"vibecode-monitor: command not found"**
- Binary not in PATH. Run `make install-user` or use full path:
  ```jsonc
  "text": "/full/path/to/vibecode-monitor --fastfetch"
  ```
- Verify with: `which vibecode-monitor`

**Output doesn't appear in fastfetch**
- Ensure binary has execute permissions: `chmod +x /path/to/vibecode-monitor`
- Test command directly: `vibecode-monitor --fastfetch`
- Check fastfetch config syntax is valid JSON

**Want full ASCII art in fastfetch instead of compact?**
- Disable fastfetch logo and use `--offset 20` flag (see Alternative Configurations above)

**Run diagnostics:**
```bash
vibecode-monitor --test-fastfetch
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
