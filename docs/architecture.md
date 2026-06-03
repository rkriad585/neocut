# Architecture

## Overview

neocut follows a standard Go project layout with `cmd/` for the entry point and `internal/` for private packages. The CLI is built with [cobra](https://github.com/spf13/cobra), audio processing is handled by [godub](https://github.com/Vernacular-ai/godub) (vendored and patched in-tree), and the interactive form uses [huh](https://github.com/charmbracelet/huh).

## Module map

```
cmd/neocut/main.go
    │
    ▼
internal/cmd/root.go         ◄── cobra root command
    │
    ├── run()                 ◄── default RunE
    │   ├── config.InitConfigFile()
    │   ├── config.ReadConfig()
    │   ├── config.PrintBanner()
    │   ├── config.EnsureConfigDir()
    │   ├── if --tui     → tui.RunConfigForm() → core.Process()
    │   ├── if --config  → tui.RunConfigEditor()
    │   └── else         → core.Process(cfg)
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
- Reading `.version` file at startup (`ReadVersion()` — checks ldflags, then embedded, then .version file)
- Printing the banner to stdout (`PrintBanner()`)
- Resolving the output directory (`GetOutputDir()` — returns `~/Downloads/neocut/` by default)
- Resolving the config directory (`ConfigDir()` → `~/.config/neostore/neocut/`)
- JSONL config file management (`jsonl.go`):
  - `InitConfigFile()` — creates config.jsonl with meta, defaults, and 3 presets (aggressive, gentle, speech)
  - `ReadConfig()` — returns presets and default entry
  - `WriteDefaults()` — replaces the default entry in config.jsonl
  - `AppendHistory()` — records each processing run
  - `FindPreset()` — case-insensitive preset lookup
- Holding ldflags-injected values: `Commit`, `PublisherName`, `PublisherEmail`

### `internal/cmd`

Contains cobra command definitions:
- `root.go` — root command with all flags (including `--format`, `--bitrate`, `--dry-run`), the `self-update` subcommand, and the `run()` function that applies defaults/presets/CLI overrides
- `uninstall.go` — the `--selfuninstall` logic with platform-specific binary removal

### `internal/core`

Audio processing pipeline. The `Process()` function:
1. Ensures ffmpeg is available (`ffmpeg.Ensure()`)
2. Loads the MP3 file via `godub.NewLoader().Load()`
3. Splits on silence via `godub.SplitOnSilence()` (vendored godub, patched normalization)
4. Recombines chunks via `chunks[0].Append(chunks[1:])`
5. Exports the result with optional `WithDstFormat()` and `WithBitRate()`

Each step is wrapped in `step()` which shows an animated spinner with panic recovery. Export uses `exportWithProgress()` for a unicode-block progress bar. Dry-run mode skips step 5 entirely.

### `internal/ffmpeg`

Manages the ffmpeg dependency:
- `Ensure()` — checks `exec.LookPath("ffmpeg")`, if missing, triggers auto-download
- `download.go` — downloads ffmpeg from well-known static build URLs with a progress bar
- On Windows, creates a `which.cmd` shim so godub's internal `which` detection works
- PATH is extended with the config bin directory to pick up the downloaded ffmpeg
- Archive extraction supports `.zip` (Windows) and `.tar.xz` (Linux/macOS)

Download sources:
| Platform | Source |
|----------|--------|
| Windows | [Gyan.dev](https://www.gyan.dev/ffmpeg/builds/) |
| Linux | [johnvansickle.com](https://johnvansickle.com/ffmpeg/) |
| macOS | [evermeet.cx](https://evermeet.cx/ffmpeg/) |

### `internal/tui`

- `form.go` — Interactive processing form using [huh](https://github.com/charmbracelet/huh). Active when `--tui` flag is passed. Returns a populated `config.Config` struct.
- `configedit.go` — Config editor TUI. Active when `--config` / `-c` flag is passed. Loads, displays, and saves `config.jsonl` defaults, presets, and history.

### `internal/update`

Self-update mechanism:
- `LatestVersion()` — HTTP GET to `raw.githubusercontent.com/rkriad585/neocut/main/.version`
- `DownloadURL()` — builds the release binary URL from `runtime.GOOS`/`runtime.GOARCH`
- `Run()` — orchestrates version check → download with progress (`io.TeeReader` + `WriteCounter`) → binary replacement
- Cross-platform rename with deferred cleanup:
  - **Windows:** `os.Rename(exe, exe.old)` → `os.Rename(tmp, exe)` — Windows permits renaming a running executable. Restores `.old` on failure.
  - **Unix:** `os.Rename(tmp, exe)` with `os.Chmod(0755)`

## Platform support

| Feature | Windows | Linux | macOS |
|---------|---------|-------|-------|
| CLI flags | ✓ | ✓ | ✓ |
| TUI mode | ✓ | ✓ | ✓ |
| ffmpeg auto-download | ✓ | ✓ | ✓ |
| self-update | ✓ (rename) | ✓ (rename) | ✓ (rename) |
| --selfuninstall | ✓ (RemoveAll) | ✓ (RemoveAll) | ✓ (RemoveAll) |
| Build script | build.ps1 / Makefile | build.sh / Makefile | build.sh / Makefile |
| Docker | — | ✓ (Dockerfile) | ✓ (Dockerfile) |
| Installer | installer.ps1 | installer.sh | installer.sh |

## Testing

Each internal package has dedicated unit tests in `*_test.go` files:

| Package | Test file | What it covers |
|---------|-----------|----------------|
| config | `config_test.go` | ReadVersion, PrintBanner, GetOutputDir, ConfigDir, directories |
| config | `jsonl_test.go` | InitConfigFile, ReadConfig, WriteDefaults, AppendHistory, FindPreset |
| core | `processor_test.go` | fmtDuration (21 sub-cases), SetQuietMode thread safety |
| ffmpeg | `ffmpeg_test.go` | BinDir, pathContains, addToPATH, downloadURL, which shim |
| ffmpeg | `download_test.go` | extractZip, downloadWithProgress (HTTP test server) |
| update | `update_test.go` | DownloadURL, filepathEval |
| cmd | `root_test.go` | Execute, flag registration, short-form uniqueness |

Run all tests with:
```bash
go test ./internal/...
```

## Automated CI/CD release

### Trigger

Push a semver tag (`git tag vx.x.x && git push --tags`) or run manually via the GitHub Actions UI.

### Pipeline (`.github/workflows/release.yml`)

```
prepare ──────────────────────────────────────────────┐
   └─ fetches .version from GitHub raw URL            │
   └─ extracts commit SHA + pre-release flag          │
                                                      ▼
build (6 parallel) ──────────────────────────────► release ──► notify-on-failure
   windows/amd64  (machine type  9)                   ▲
   windows/arm64  (machine type 12)                   ▲
   linux/amd64                                        ▲
   linux/arm64    (cross-compiled w/ aarch64-gcc)     ▲
   darwin/amd64   (macos-latest Intel runner)         ▲
   darwin/arm64   (macos-latest M2 native runner)     ▲
                                                      ▲
changelog ────────────────────────────────────────────┘
   └─ groups commits: feat / fix / perf / docs / other
```

### Output binaries

| Platform | File |
|---|---|
| Windows AMD64 | `neocut-windows-amd64.exe` |
| Windows ARM64 | `neocut-windows-arm64.exe` |
| Linux AMD64 | `neocut-linux-amd64` |
| Linux ARM64 | `neocut-linux-arm64` |
| macOS Intel | `neocut-darwin-amd64` |
| macOS Silicon | `neocut-darwin-arm64` |

### ldflags injected

| Variable | Source |
|---|---|
| `main.Version` | `.version` file (fetched from GitHub raw URL) |
| `main.Commit` | `git rev-parse --short=8 HEAD` |
| `main.PublisherName` | `rkriad585` |
| `main.PublisherEmail` | `rkriad585@gmail.com` |

### Local builds

- Version embedded via `//go:embed version.txt` as fallback
- `go:generate` directive syncs `.version` → `version.txt` via `gen.go`
- Cross-platform scripts (`build.ps1` / `build.sh`) build 6 binaries to `bin/`
- Installer scripts download from GitHub Releases
- Quiet mode (`cfg.Quiet`) propagated to `core.SetQuietMode()` skips animated spinners
