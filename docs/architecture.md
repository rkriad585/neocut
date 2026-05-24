# Architecture

## Overview

neocut follows a standard Go project layout with `cmd/` for the entry point and `internal/` for private packages. The CLI is built with [cobra](https://github.com/spf13/cobra), audio processing is handled by [godub](https://github.com/Vernacular-ai/godub), and the interactive form uses [huh](https://github.com/charmbracelet/huh).

## Module map

```
cmd/neocut/main.go
    │
    ▼
internal/cmd/root.go         ◄── cobra root command
    │
    ├── run()                 ◄── default RunE
    │   ├── config.PrintBanner()
    │   ├── config.EnsureConfigDir()
    │   ├── if --tui → tui.RunConfigForm()
    │   └── core.Process(cfg) ◄── audio pipeline
    │
    ├── selfUpdateCmd         ◄── "self-update" subcommand
    │   └── update.Run(version)
    │
    └── --selfuninstall flag
        └── runSelfUninstall()
```

## Package details

### `internal/config`

Responsible for:
- Reading `.version` file at startup
- Printing the banner to stdout
- Resolving the output directory (`~/Downloads/neocut/`)
- Resolving the config directory (`~/.config/neostore/neocut/`)
- Holding ldflags-injected values: `Commit`, `PublisherName`, `PublisherEmail`

### `internal/cmd`

Contains cobra command definitions:
- `root.go` — root command with all flags, the `self-update` subcommand, and error handling
- `uninstall.go` — the `--selfuninstall` logic with platform-specific binary removal

### `internal/core`

Audio processing pipeline. The `Process()` function:
1. Ensures ffmpeg is available (`ffmpeg.Ensure()`)
2. Loads the MP3 file via `godub.NewLoader().Load()`
3. Splits on silence via `godub.SplitOnSilence()`
4. Recombines chunks via `chunks[0].Append(chunks[1:])`
5. Exports the result as MP3

Each step is wrapped in `step()` which shows an animated spinner with panic recovery. Export uses `exportWithProgress()` for a unicode-block progress bar.

### `internal/ffmpeg`

Manages the ffmpeg dependency:
- `Ensure()` — checks `exec.LookPath("ffmpeg")`, if missing, triggers auto-download
- `download.go` — downloads ffmpeg from well-known static build URLs
- On Windows, creates a `which.cmd` shim so godub's internal `which` detection works
- PATH is extended with the config bin directory to pick up the downloaded ffmpeg

Download sources:
| Platform | Source |
|----------|--------|
| Windows | [Gyan.dev](https://www.gyan.dev/ffmpeg/builds/) |
| Linux | [johnvansickle.com](https://johnvansickle.com/ffmpeg/) |
| macOS | [evermeet.cx](https://evermeet.cx/ffmpeg/) |

### `internal/tui`

Interactive form using [huh](https://github.com/charmbracelet/huh). Only active when `--tui` flag is passed. Returns a populated `config.Config` struct.

### `internal/update`

Self-update mechanism:
- `LatestVersion()` — HTTP GET to `raw.githubusercontent.com/rkriad585/neocut/main/.version` (10s timeout)
- `DownloadURL()` — builds the release binary URL from `runtime.GOOS`/`runtime.GOARCH`
- `Run()` — orchestrates version check → download with progress → binary replacement
- `replaceBinary()` — platform-specific replacement:
  - Unix: `os.Rename(tmp, exePath)` (inode stays alive)
  - Windows: deferred `.bat` script that waits, deletes old, renames new, restarts

## Platform support

| Feature | Windows | Linux | macOS |
|---------|---------|-------|-------|
| CLI flags | ✓ | ✓ | ✓ |
| TUI mode | ✓ | ✓ | ✓ |
| ffmpeg auto-download | ✓ | ✓ | ✓ |
| self-update | ✓ (bat) | ✓ (rename) | ✓ (rename) |
| --selfuninstall | ✓ (bat) | ✓ (RemoveAll) | ✓ (RemoveAll) |
| Build script | build.ps1 | build.sh | build.sh |
| Installer | installer.ps1 | installer.sh | installer.sh |

## Build & release

- Version is embedded via `//go:embed version.txt`, always available at runtime from any directory
- Can be overridden via ldflags (`Version`) at build time
- ldflags inject `Commit`, `PublisherName`, `PublisherEmail` into the binary
- `go:generate` directive syncs `.version` → `version.txt` via `gen.go`
- Quiet mode (`cfg.Quiet`) is propagated to `core.SetQuietMode()` which skips animated spinners and progress bars, using direct function calls with panic recovery instead
- Cross-platform scripts build 6 binaries: `{os}-{arch}` for windows/linux/darwin × amd64/arm64
- Installer scripts download from `https://github.com/rkriad585/neocut/releases/download/{version}/{binary}`
