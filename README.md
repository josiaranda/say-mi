# say-mi

A CLI audio player that plays random sounds from nested categories. Self-contained binary with no external dependencies.

## Installation

### One-liner (Linux/macOS)

```bash
curl -fsSL https://raw.githubusercontent.com/josiaranda/say-mi/main/install.sh | bash
```

### Homebrew

```bash
brew tap josiaranda/say-mi
brew install say-mi
```

### Manual Download

Download from [Releases](https://github.com/josiaranda/say-mi/releases/latest):

| Platform | Architecture | Download |
|----------|--------------|----------|
| Linux | amd64 | `say-mi-linux-amd64` |
| Linux | arm64 | `say-mi-linux-arm64` |
| macOS | amd64 (Intel) | `say-mi-darwin-amd64` |
| macOS | arm64 (Apple Silicon) | `say-mi-darwin-arm64` |
| Windows | amd64 | `say-mi-windows-amd64.exe` |

```bash
# Example for Linux amd64
curl -L -o say-mi https://github.com/josiaranda/say-mi/releases/latest/download/say-mi-linux-amd64
chmod +x say-mi
sudo mv say-mi /usr/local/bin/
```

## Usage

```bash
# Basic usage - plays random audio from category
say-mi greeting

# Nested categories with dot notation (voice.say.category)
say-mi kyoko.say.greeting
say-mi jamie.say.permission

# Custom config file
say-mi -c ./voices.yaml kyoko.say.greeting
say-mi --config ./voices.yaml jamie.say.goodbye

# Verbose mode for debugging
say-mi --verbose kyoko.say.greeting

# Exit with error code when category not found
say-mi -e nonexistent      # exits 1
say-mi --error-exit missing  # exits 1

# Without -e flag, missing categories exit 0 silently
say-mi nonexistent         # exits 0, silent

# Show version
say-mi --version
```

## Updating

### One-liner / Manual install
```bash
# Update to latest
curl -fsSL https://raw.githubusercontent.com/josiaranda/say-mi/main/install.sh | bash

# Install specific version
curl -fsSL https://raw.githubusercontent.com/josiaranda/say-mi/main/install.sh | bash -s -- v1.0.0
```

### Homebrew
```bash
# Update to latest
brew upgrade josiaranda/say-mi/say-mi

# Reinstall specific version (edit formula version first)
brew reinstall josiaranda/say-mi/say-mi
```

## Flags

| Flag | Description | Default |
|------|-------------|---------|
| `-c, --config` | Path to YAML config file | `./config.yaml` |
| `-e, --error-exit` | Exit with code 1 on category/audio not found | `false` |
| `--verbose` | Enable verbose/debug output | `false` |

## Exit Codes

| Code | Description |
|------|-------------|
| `0` | Success (audio played) OR not-found (without `-e` flag) |
| `1` | Error (with `-e` flag AND category/audio not found) |

## Config Format

Create a `voices.yaml` file:

```yaml
# Optional: played when category not found
fallback: "audio/error/not-found.mp3"

voices:
  kyoko:
    voice_id: "elevenlabs-voice-id"  # optional, ignored by say-mi
    say:
      greeting:
        - text: "What are we working on today?"
          audio: "audio/kyoko/greeting/0.mp3"
        - text: "Good to see you again."
          audio: "audio/kyoko/greeting/1.mp3"
      permission:
        - text: "Is this okay?"
          audio: "audio/kyoko/permission/0.mp3"
      goodbye:
        - text: "See you later!"
          audio: "audio/kyoko/goodbye/0.mp3"

  jamie:
    voice_id: "another-voice-id"
    say:
      greeting:
        - text: "Hey, what's up?"
          audio: "audio/jamie/greeting/0.mp3"
      permission:
        - text: "Should I go ahead?"
          audio: "audio/jamie/permission/0.mp3"
```

Access categories with dot notation (`voice.say.category`):
```bash
say-mi kyoko.say.greeting      # Plays random file from kyoko's greetings
say-mi jamie.say.permission    # Plays from jamie's permissions
```

You can also use a simpler format without `text` fields:
```yaml
voices:
  kyoko:
    say:
      greeting:
        - "audio/kyoko/greeting/0.mp3"
        - "audio/kyoko/greeting/1.mp3"
```

## Supported Audio Formats

- MP3
- WAV
- OGG (Vorbis)
- FLAC

## Building from Source

Requirements:
- Go 1.21+
- On Linux: `libasound2-dev` (ALSA headers)

```bash
# Install dependencies (Linux)
sudo apt install -y libasound2-dev

# Build
go build -o say-mi .

# Cross-compile
GOOS=darwin GOARCH=amd64 go build -o say-mi-macos-amd64 .
GOOS=darwin GOARCH=arm64 go build -o say-mi-macos-arm64 .
GOOS=windows GOARCH=amd64 go build -o say-mi.exe .
```

## License

MIT
