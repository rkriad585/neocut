# Getting Started

## Prerequisites

- **Go 1.21+** (only required to build from source)
- **ffmpeg** — auto-downloaded if missing (Windows, Linux, macOS)

## Installation

### Option 1: One-liner (recommended)

**Windows (PowerShell 5/7+):**
```powershell
irm https://raw.githubusercontent.com/rkriad585/neocut/main/installer.ps1 | iex
```

**Linux / macOS:**
```sh
curl -fsSL https://raw.githubusercontent.com/rkriad585/neocut/main/installer.sh | sh
```

The installer:
- Detects your OS and architecture
- Downloads the latest release binary
- Installs to `~/.config/neostore/neocut/bin/`
- Adds the binary to your PATH
- On Windows: updates user PATH via registry
- On Unix: appends to `.bashrc` or `.zshrc`

### Option 2: Build from source

```bash
git clone https://github.com/rkriad585/neocut.git
cd neocut
go build -ldflags "-X main.Version=$(cat .version) -X main.Commit=$(git rev-parse --short HEAD)" -o neocut ./cmd/neocut/
```

> ldflags now target `main.Version`, `main.Commit`, `main.PublisherName`, and `main.PublisherEmail` — the same vars used by the automated CI/CD workflow.

### Option 3: Make (Unix)

```sh
make            # build for current platform
make build-all  # cross-compile all 6 platforms
make test       # run all tests
make clean      # remove build artifacts
```

### Option 4: Docker

```sh
make docker          # build image using Make
docker build -t neocut .  # or build manually
docker run --rm -v "$(pwd):/workspace" neocut -i /workspace/input.mp3
```

### Option 5: Cross-platform build script

**Windows:**
```powershell
.\build.ps1
```

**Unix:**
```sh
chmod +x build.sh && ./build.sh
```

Builds 6 binaries into `bin/`:
- `neocut-windows-amd64.exe`, `neocut-windows-arm64.exe`
- `neocut-linux-amd64`, `neocut-linux-arm64`
- `neocut-darwin-amd64`, `neocut-darwin-arm64`

## Verifying the install

```bash
neocut --version
```

Expected output:
```
neocut version v1.5.0
```

## First run

```bash
# Process an MP3 file
neocut -i path/to/your/file.mp3
```

On first run, if ffmpeg is not found on your system, neocut automatically downloads it to `~/.config/neostore/neocut/bin/`.

Output is saved to `~/Downloads/neostore/neocut/yourfile_no_silence.mp3`.

## Customize the look

```bash
# Apply a dark theme
neocut -i file.mp3 --theme dark

# Force light colors
neocut -i file.mp3 --color-mode light

# Browse all 13 themes interactively
neocut --config
```

Themes are saved to `~/.config/neostore/neocut/config.toml` and persist across runs.

## Next steps

- Read the [CLI usage reference](usage.md) for all flags
- See [configuration](configuration.md) for tuning silence detection and themes
- Check [troubleshooting](troubleshooting.md) if something goes wrong
