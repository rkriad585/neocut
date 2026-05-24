# neocut

> Remove silence from MP3 audio files — automatically detect, split, and recombine.

> [!NOTE]
> Type: CLI Tool | Language: Go | Status: active

## Features

- Detects silent portions in MP3 files and removes them
- Configurable silence threshold, minimum length, and keep-silence margin
- Animated progress spinner and progress bars during processing
- Interactive TUI mode via `--tui` (powered by [huh](https://github.com/charmbracelet/huh))
- Auto-downloads ffmpeg when missing (Windows, Linux, macOS)
- `self-update` command to upgrade to the latest release
- `--selfuninstall` flag to fully remove neocut from the system
- One-liner installers for Windows (PowerShell) and Unix (bash)

## Quick Start

```bash
# Process an MP3 file
neocut -i input.mp3

# Use interactive TUI
neocut --tui
```

## Installation

### One-liner installers

**Windows (PowerShell 5/7+):**
```powershell
irm https://raw.githubusercontent.com/rkriad585/neocut/main/installer.ps1 | iex
```

**Linux / macOS:**
```sh
curl -fsSL https://raw.githubusercontent.com/rkriad585/neocut/main/installer.sh | sh
```

### Build from source

```bash
git clone https://github.com/rkriad585/neocut.git
cd neocut
go build -ldflags "-X neocut/internal/config.Commit=$(git rev-parse --short HEAD)" -o neocut ./cmd/neocut/
```

### Cross-platform build

```bash
# Windows
.\build.ps1

# Unix
chmod +x build.sh && ./build.sh
```

Outputs 6 binaries to `bin/`:
- `neocut-windows-amd64.exe`, `neocut-windows-arm64.exe`
- `neocut-linux-amd64`, `neocut-linux-arm64`
- `neocut-darwin-amd64`, `neocut-darwin-arm64`

## Usage

```
neocut -i input.mp3 [-o output.mp3] [flags]
```

### Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--input` | `-i` | `""` | Input MP3 file (required) |
| `--output` | `-o` | auto | Output filename (saved to `~/Downloads/neocut/`) |
| `--tui` | `-t` | `false` | Interactive TUI mode |
| `--min-silence-len` | `-m` | `1000` | Minimum silence length in ms |
| `--silence-thresh` | `-s` | `-16` | Silence threshold in dBFS |
| `--keep-silence` | `-k` | `100` | Silence to keep at boundaries in ms |
| `--seek-step` | `-e` | `1` | Seek step in ms |
| `--config`/`--cnf` | `-c` | `""` | Path to config file |
| `--selfuninstall` | | `false` | Remove neocut and its config directory |
| `--version` | `-v` | | Print version |
| `--help` | `-h` | | Print help |

### Commands

| Command | Description |
|---------|-------------|
| `self-update` | Update neocut to the latest version from GitHub |
| `completion` | Generate shell autocompletion scripts |
| `help` | Help about any command |

### Examples

```bash
# Basic usage
neocut -i podcast.mp3

# Custom output name
neocut -i recording.mp3 -o cleaned.mp3

# Aggressive silence removal
neocut -i lecture.mp3 -m 500 -s -24 -k 50

# Interactive TUI
neocut --tui
```

## Output

- Processed files are saved to `~/Downloads/neocut/`
- Default output name: `{input_name}_no_silence.mp3`
- Config directory: `~/.config/neostore/neocut/`
- ffmpeg binary: `~/.config/neostore/neocut/bin/`

## Updating

```bash
neocut self-update
```

Fetches the latest version from GitHub and replaces the current binary. Works on all platforms.

## Uninstalling

```bash
neocut --selfuninstall
```

Removes the config directory (`~/.config/neostore/neocut/`), deletes the binary, and prints PATH cleanup instructions.

## Related

- [Detailed overview](docs/overview.md)
- [Getting started guide](docs/getting-started.md)
- [CLI usage reference](docs/usage.md)
- [Architecture](docs/architecture.md)
- [Configuration](docs/configuration.md)
- [Troubleshooting](docs/troubleshooting.md)
