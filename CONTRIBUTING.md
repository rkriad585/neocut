# Contributing to neocut

First off, thanks for taking the time to contribute!

## How to Contribute

### Report Bugs

Open an issue at https://github.com/rkriad585/neocut/issues with:
- A clear title and description
- Steps to reproduce
- Expected vs actual behavior
- Your OS and neocut version (`neocut --version`)

### Suggest Features

Open an issue with:
- A clear title and description of the feature
- Why this would be useful
- Any implementation ideas

### Submit Code Changes

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/your-feature`)
3. Make your changes
4. Run tests (`go test ./...`)
5. Run `go vet ./...` to check for issues
6. Commit with a clear message
7. Push and open a Pull Request

### Pull Request Guidelines

- Keep changes focused — one feature/fix per PR
- Write tests for new functionality
- Ensure all existing tests pass
- Update documentation if needed
- Follow Go conventions (`gofmt` your code)

## Development Setup

```bash
git clone https://github.com/rkriad585/neocut.git
cd neocut
go generate ./...
go build -o neocut.exe ./cmd/neocut/
```

## Project Structure

```
neocut/
├── cmd/neocut/          # Entry point
├── internal/
│   ├── cmd/             # Cobra command definitions
│   ├── config/          # Config, version, JSONL
│   ├── core/            # Audio processing pipeline
│   ├── ffmpeg/          # FFmpeg detection & download
│   ├── tui/             # Interactive TUI forms
│   └── update/          # Self-update logic
├── docs/                # Documentation
├── vendor/              # Vendored dependencies
├── build.ps1            # Windows build script
├── build.sh             # Unix build script
├── installer.ps1        # Windows installer
└── installer.sh         # Unix installer
```

## Code of Conduct

This project adheres to the [Contributor Covenant](CODE_OF_CONDUCT.md).
By participating, you agree to uphold this code.
