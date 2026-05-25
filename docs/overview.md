# Overview

> neocut is a cross-platform CLI tool that removes silence from MP3 audio files. It detects silent regions, splits them out, and recombines the non-silent segments into a single tight output file.

## What it does

Given an MP3 file with pauses, gaps, or silent sections, neocut:

1. **Loads** the audio via [godub](https://github.com/Vernacular-ai/godub) (Go port of [pydub](https://github.com/jiaaro/pydub))
2. **Detects silence** using configurable threshold, minimum length, and seek precision
3. **Splits** on silence boundaries, discarding the silent chunks
4. **Recombines** the remaining audio segments into one continuous file
5. **Exports** the result as MP3 (or WAV/FLAC via `--format`) to `~/Downloads/neocut/`

Every step shows an animated spinner or progress bar. Panics during silence detection are recovered gracefully.

## When to use it

| Use case | Example |
|----------|---------|
| Podcast editing | Remove gaps between sentences and breaths |
| Lecture recordings | Eliminate pauses between slides |
| Voice notes | Tighten up recordings with long silences |
| Pre-processing | Clean audio before further processing |

## Project structure

```
neocut/
├── cmd/
│   └── neocut/
│       └── main.go              # Entry point
├── internal/
│   ├── cmd/
│   │   ├── root.go              # Cobra root command + self-update
│   │   ├── root_test.go         # Cobra command and flag tests
│   │   └── uninstall.go         # --selfuninstall logic
│   ├── config/
│   │   ├── config.go            # Config struct, banner, version, directories
│   │   ├── config_test.go       # ReadVersion, GetOutputDir, ConfigDir tests
│   │   ├── jsonl.go             # JSONL config I/O (defaults, presets, history)
│   │   ├── jsonl_test.go        # InitConfigFile, ReadConfig, WriteDefaults tests
│   │   ├── gen.go               # go:generate helper (.version → version.txt)
│   │   └── version.txt          # Embedded version file
│   ├── core/
│   │   ├── processor.go         # Audio pipeline (load → split → combine → export)
│   │   └── processor_test.go    # fmtDuration, SetQuietMode tests
│   ├── ffmpeg/
│   │   ├── ffmpeg.go            # ffmpeg detection, PATH setup, which.cmd shim
│   │   ├── ffmpeg_test.go       # BinDir, pathContains, addToPATH tests
│   │   ├── download.go          # HTTP download with progress bar, archive extraction
│   │   └── download_test.go     # extractZip, downloadWithProgress tests
│   ├── tui/
│   │   ├── form.go              # huh interactive processing form
│   │   └── configedit.go        # huh config editor form
│   └── update/
│       ├── update.go            # self-update: version fetch, download, binary replace
│       └── update_test.go       # DownloadURL, filepathEval, replaceUnix tests
├── docs/                        # Documentation
├── vendor/                      # Vendored dependencies (godub patched in-tree)
├── build.ps1                    # Windows cross-platform build script
├── build.sh                     # Unix cross-platform build script
├── installer.ps1                # Windows one-liner installer
├── installer.sh                 # Unix one-liner installer
├── .version                     # Current version (v1.0.2)
├── go.mod / go.sum              # Go module
└── README.md
```

## Key technologies

| Component | Library |
|-----------|---------|
| CLI framework | [cobra](https://github.com/spf13/cobra) |
| Audio processing | [godub](https://github.com/Vernacular-ai/godub) (vendored + patched) |
| TUI forms | [huh](https://github.com/charmbracelet/huh) |
| Config format | JSONL (JSON Lines) |
| Audio codec | ffmpeg (auto-downloaded) |

## Version

- Embedded into the binary via `//go:embed` — always available, from any directory
- Can be overridden via ldflags at build time (`Version`)
- Also injected/embedded: `Commit` (git hash), `PublisherName`, `PublisherEmail`

## Processing stats

After processing, neocut shows a summary with:
- Input duration
- Output duration
- Silence removed (absolute and percentage)

## Unit tests

All internal packages have unit tests covering core logic, file I/O, HTTP downloads, CLI flag registration, and duration formatting:

```bash
go test ./internal/...
```

85+ tests across 6 packages — no external test dependencies required.
