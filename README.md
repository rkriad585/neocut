# neocut

> Remove silence from MP3 audio files — automatically detect, split, and recombine.

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

> [!NOTE]
> Type: CLI Tool | Language: Go | Status: active

## Features

- Detects silent portions in MP3 files and removes them
- Configurable silence threshold, minimum length, and keep-silence margin
- Animated progress spinner and progress bars during processing
- Interactive TUI mode via `--tui` (powered by [huh](https://github.com/charmbracelet/huh))
- Output in MP3, WAV, or FLAC via `--format` with configurable `--bitrate`
- Preview stats without exporting via `--dry-run` (iterate quickly)
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

# Using Make (Linux/macOS)
make

# Or directly with Go (ldflags inject into main package)
go build -ldflags "-X main.Version=$(cat .version) -X main.Commit=$(git rev-parse --short HEAD)" -o neocut ./cmd/neocut/
```

### Cross-platform build

```bash
# Windows
.\build.ps1

# Unix
chmod +x build.sh && ./build.sh

# Or using Make
make build-all
```

Outputs 6 binaries to `bin/`:
- `neocut-windows-amd64.exe`, `neocut-windows-arm64.exe`
- `neocut-linux-amd64`, `neocut-linux-arm64`
- `neocut-darwin-amd64`, `neocut-darwin-arm64`

### Automated release (CI/CD)

Push a tag and GitHub Actions builds + publishes all 6 binaries:

```bash
git tag v1.0.3
git push --tags
```

The workflow:
1. **Prepare** — fetches `.version` from GitHub, resolves tag/commit/prerelease
2. **Build (6× parallel)** — cross-compiles for Windows/macOS/Linux × amd64/arm64
3. **Changelog** — groups commits by `feat:` / `fix:` / `perf:` / `docs:` / other
4. **Release** — creates GitHub Release with binaries + SHA-256 checksums
5. **Notify** — failure alert (extensible to Slack)

Zero manual steps. Everything in `.github/workflows/release.yml`.

### Docker

```bash
# Build the image
make docker

# Or build manually
docker build -t neocut .

# Run
docker run --rm -v "$(pwd):/workspace" neocut -i /workspace/input.mp3
```

## Usage

```
neocut -i input.mp3 [-o output.mp3] [flags]
```

### Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--input` | `-i` | `""` | Input MP3 file (required) |
| `--output` | `-o` | auto | Output filename (saved to output-dir) |
| `--output-dir` | `-d` | `~/Downloads/neostore/neocut/` | Output directory |
| `--tui` | `-t` | `false` | Interactive TUI mode |
| `--config` | `-c` | `false` | Edit project config interactively (huh TUI) |
| `--quiet` | `-q` | `false` | Suppress banner, spinners, and progress |
| `--preset` | | `""` | Load preset from config (aggressive, gentle, speech) |
| `--theme` | | `""` | Theme name (e.g. dark, sunny_beach_day, olive_garden_feast) |
| `--color-mode` | | `"auto"` | Color mode: auto, dark, light |
| `--format` | `-f` | `mp3` | Output format: mp3, wav, flac |
| `--bitrate` | `-b` | `0` | Output bitrate in kbps (e.g. 192, 320) |
| `--dry-run` | | `false` | Preview stats without exporting |
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

# Use a preset from the config file (aggressive, gentle, or speech)
neocut -i podcast.mp3 --preset speech

# WAV output — for further editing
neocut -i recording.mp3 -f wav

# Dry run — see stats without writing a file
neocut -i podcast.mp3 --dry-run

# FLAC with custom bitrate
neocut -i master.mp3 -f flac -b 320

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
 │  4. Export  │  Writes the result (MP3/WAV/FLAC) to output-dir
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

- Processed files are saved to `~/Downloads/neostore/neocut/`
- Default output name: `{input_name}_no_silence.{ext}` (mp3, wav, flac)
- Default format: MP3 (set via `--format`)
- Config directory: `~/.config/neostore/neocut/`
- ffmpeg binary: `~/.config/neostore/neocut/bin/`

## Updating

```bash
neocut self-update
```

Fetches the latest version from GitHub and replaces the current binary. Works on all platforms.

## Project config

neocut stores a TOML config file at `~/.config/neostore/neocut/config.toml` and a processing history log at `~/.config/neostore/neocut/history.log`.

### Editing interactively

```bash
neocut --config   # or neocut -c
```

Opens a TUI (powered by [huh](https://github.com/charmbracelet/huh)) to edit:
- Default processing parameters
- Theme and color mode
- View configured presets
- Browse recent processing history
- Save changes back to config.toml

### File format

```toml
[default]
min_silence_len = 1000
silence_thresh = -16.0
keep_silence = 100
seek_step = 1
output_dir = ""
theme = "sunny_beach_day"
color_mode = "auto"

[[preset]]
name = "aggressive"
min_silence_len = 500
silence_thresh = -24.0
keep_silence = 50
seek_step = 1

[[preset]]
name = "gentle"
min_silence_len = 2000
silence_thresh = -10.0
keep_silence = 200
seek_step = 5

[[preset]]
name = "speech"
min_silence_len = 800
silence_thresh = -20.0
keep_silence = 80
seek_step = 1
```

- The `[default]` section sets base parameters — override any field with CLI flags
- `[[preset]]` entries are named collections of parameters, loaded via `--preset`
- `theme` and `color_mode` in `[default]` control the UI color scheme
- Processing history is appended to `history.log` after each successful run
- CLI flags always take precedence over config values

## Theme system

neocut ships with 13 built-in color themes that can be set via `--theme`, config, or the config editor (`neocut -c`):

| Theme | Description |
|-------|-------------|
| `dark` | Classic dark mode (default fallback for dark color modes) |
| `light` | Clean light theme (default fallback for light color modes) |
| `sunny_beach_day` | Warm tropical palette (default) |
| `olive_garden_feast` | Earthy tones |
| `summer_ocean_breeze` | Cool coastal blues |
| `refreshing_summer_fun` | Bright citrus |
| `black_gold_elegance` | Premium dark + gold |
| `vibrant_color_fiesta` | Neon festival |
| `light_steel` | Industrial grays |
| `golden_twilight` | Sunset amber |
| `deep_sea` | Deep ocean blues |
| `bright_green` | Fresh nature green |
| `vivid_nightfall` | Purple twilight |

The `--color-mode` flag forces the UI into dark or light mode regardless of the selected theme:
- `auto` — use theme's native colors (default)
- `dark` — force Dark Theme (overrides the selected theme)
- `light` — force Light Theme (overrides the selected theme)

## Testing

```bash
go test ./internal/...
```

85+ unit tests across all 6 internal packages covering config I/O, audio pipeline helpers, ffmpeg management, self-update logic, and CLI flag registration. No external test dependencies.

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

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

This project is released under the [MIT License](LICENSE) and follows the [Contributor Covenant](CODE_OF_CONDUCT.md). Security vulnerabilities can be reported via [SECURITY.md](SECURITY.md).
