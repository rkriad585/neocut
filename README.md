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
go build -ldflags "-X neocut/internal/config.Commit=$(git rev-parse --short HEAD) -X neocut/internal/config.Version=$(cat .version)" -o neocut ./cmd/neocut/
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
| `--output` | `-o` | auto | Output filename (saved to output-dir) |
| `--output-dir` | `-d` | `~/Downloads/neocut/` | Output directory |
| `--tui` | `-t` | `false` | Interactive TUI mode |
| `--quiet` | `-q` | `false` | Suppress banner, spinners, and progress |
| `--min-silence-len` | `-m` | `1000` | Minimum silence length in ms |
| `--silence-thresh` | `-s` | `-16` | Silence threshold in dBFS |
| `--keep-silence` | `-k` | `100` | Silence to keep at boundaries in ms |
| `--seek-step` | `-e` | `1` | Seek step in ms |
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
# Basic usage — removes all pauses from a podcast
neocut -i podcast_episode.mp3

# Custom output filename
neocut -i recording.mp3 -o cleaned_recording.mp3

# Custom output directory (for organizing projects)
neocut -i interview.mp3 -d "D:\Projects\audio\cleaned"

# Aggressive silence removal — catches short gaps, lower threshold
neocut -i lecture.mp3 -m 500 -s -24 -k 50

# Gentle removal — only long pauses, keep more natural silence
neocut -i audiobook.mp3 -m 2000 -s -10 -k 200

# High precision — slower but more accurate boundaries
neocut -i vocals.mp3 -e 0.5

# Quiet mode — suppress banners, only show output path (for scripts)
neocut -i batch_input.mp3 -q

# Interactive TUI mode — fill in options visually
neocut --tui
```

## How it works

neocut processes audio in four steps:

```
 Input MP3
    │
    ▼
 ┌─────────────┐
 │  1. Load    │  Reads the MP3 via ffmpeg + godub
 └──────┬──────┘
        ▼
 ┌─────────────┐
 │  2. Detect  │  Scans for silent regions using:
 │   Silence   │    • silence-thresh: volume below this = silent
 │             │    • min-silence-len: shortest silence to cut
 │             │    • seek-step: how finely to scan
 └──────┬──────┘
        ▼
 ┌─────────────┐
 │  3. Split   │  Cuts out silent chunks, keeps non-silent
 │   & Rejoin  │  segments, re-appends them with no gap
 └──────┬──────┘
        ▼
 ┌─────────────┐
 │  4. Export  │  Writes the result as MP3 to output-dir
 └─────────────┘
```

The algorithm uses [godub.SplitOnSilence](https://github.com/Vernacular-ai/godub), which walks through the audio frame-by-frame (at `seek-step` granularity), marks frames below `silence-thresh` as silent, groups consecutive silent frames into regions, discards regions longer than `min-silence-len`, and keeps `keep-silence` ms of the boundary to avoid abrupt cuts.

After processing, neocut shows a summary:

```
    Segments:     42              ← non-silent chunks found
    Input:        45m 30s         ← original duration
    Output:       38m 12s         ← duration after removal
    Removed:      7m 18s (16.0%)  ← silence cut
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
