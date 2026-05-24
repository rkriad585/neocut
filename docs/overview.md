# Overview

> neocut is a cross-platform CLI tool that removes silence from MP3 audio files. It detects silent regions, splits them out, and recombines the non-silent segments into a single tight output file.

## What it does

Given an MP3 file with pauses, gaps, or silent sections, neocut:

1. **Loads** the audio via [godub](https://github.com/Vernacular-ai/godub) (Go port of [pydub](https://github.com/jiaaro/pydub))
2. **Detects silence** using configurable threshold, minimum length, and seek precision
3. **Splits** on silence boundaries, discarding the silent chunks
4. **Recombines** the remaining audio segments into one continuous file
5. **Exports** the result as MP3 to `~/Downloads/neocut/`

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
│       └── main.go          # Entry point
├── internal/
│   ├── cmd/
│   │   ├── root.go          # Cobra root command + self-update
│   │   └── uninstall.go     # --selfuninstall logic
│   ├── config/
│   │   └── config.go        # Config struct, banner, version, directories
│   ├── core/
│   │   └── processor.go     # Audio pipeline (load → split → combine → export)
│   ├── ffmpeg/
│   │   ├── ffmpeg.go        # ffmpeg detection, PATH setup, which.cmd shim
│   │   └── download.go      # HTTP download with progress bar, archive extraction
│   ├── tui/
│   │   └── form.go          # huh interactive form
│   └── update/
│       └── update.go        # self-update: version fetch, download, binary replace
├── docs/                    # Documentation
├── build.ps1                # Windows cross-platform build script
├── build.sh                 # Unix cross-platform build script
├── installer.ps1            # Windows one-liner installer
├── installer.sh             # Unix one-liner installer
├── .version                 # Current version (v0.1.0)
├── go.mod / go.sum          # Go module
└── README.md
```

## Key technologies

| Component | Library |
|-----------|---------|
| CLI framework | [cobra](https://github.com/spf13/cobra) |
| Audio processing | [godub](https://github.com/Vernacular-ai/godub) |
| TUI forms | [huh](https://github.com/charmbracelet/huh) |
| Audio codec | ffmpeg (auto-downloaded) |
