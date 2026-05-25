# Configuration

## CLI flags

All configuration is passed via CLI flags. Persistent defaults and presets are stored in `~/.config/neostore/neocut/config.jsonl`.

| Flag | Default | Description |
|------|---------|-------------|
| `-i / --input` | `""` | Input MP3 file path (required unless using `--tui`) |
| `-o / --output` | auto | Output filename; defaults to `{input}_no_silence.{format}` |
| `-d / --output-dir` | `~/Downloads/neocut/` | Custom output directory |
| `-q / --quiet` | `false` | Suppress banner, spinners, and progress |
| `-f / --format` | `mp3` | Output codec: `mp3`, `wav`, `flac` |
| `-b / --bitrate` | `0` | Output bitrate in kbps (`0` = codec default) |
| `--dry-run` | `false` | Preview stats without exporting |
| `-m / --min-silence-len` | `1000` ms | Minimum duration of silence to detect |
| `-s / --silence-thresh` | `-16` dBFS | Volume threshold; segments below this are considered silence |
| `-k / --keep-silence` | `100` ms | Silence to preserve at the boundaries of kept segments |
| `-e / --seek-step` | `1` ms | Precision of silence detection (lower = more accurate but slower) |
| `--preset` | `""` | Load a named preset from config.jsonl (`aggressive`, `gentle`, `speech`) |
| `-c / --config` | `false` | Edit project config interactively via TUI |
| `-t / --tui` | `false` | Launch interactive form instead of parsing flags |

## Config file: config.jsonl

Located at `~/.config/neostore/neocut/config.jsonl`, this JSONL file stores:

| Entry type | Purpose |
|------------|---------|
| `meta` | Project metadata (version, creation time) |
| `default` | Persistent default values for all flags |
| `preset` | Named presets (e.g. `aggressive`, `gentle`, `speech`) |
| `history` | Processing history with timestamps |

Defaults and presets are only applied to flags NOT explicitly set by the user on the command line. CLI flags always take precedence.

Built-in presets:
| Preset | Min Silence | Threshold | Keep Silence | Seek Step |
|--------|-------------|-----------|--------------|-----------|
| aggressive | 500ms | -24 dBFS | 50ms | 1ms |
| gentle | 2000ms | -10 dBFS | 200ms | 5ms |
| speech | 800ms | -20 dBFS | 80ms | 1ms |

## Silence detection parameters

### `--min-silence-len` (`-m`)

Minimum length of silence in milliseconds. Segments of silence shorter than this are ignored.

- Higher values (e.g. `2000`) → only removes long pauses
- Lower values (e.g. `500`) → removes short gaps too, more aggressive compression

### `--silence-thresh` (`-s`)

The volume threshold in dBFS. Any audio below this level is considered silence.

- `-16` dBFS (default) — moderate sensitivity; catches typical room noise
- `-24` dBFS — more sensitive; catches quieter sections
- `-8` dBFS — less sensitive; only removes very quiet sections
- `0` dBFS — maximum volume; nothing would be detected as silence

### `--keep-silence` (`-k`)

Amount of silence in milliseconds to retain at the start and end of each kept segment. Prevents abrupt cuts.

- `100` (default) — smooth transitions
- `0` — aggressive cuts (may sound choppy)
- `200` — softer transitions

### `--seek-step` (`-e`)

The step size in milliseconds for detecting silence. Lower values give more precise boundaries but take longer to process.

- `1` (default) — precise but slower
- `10` — faster but less accurate boundaries
- `0.5` — highest precision, slowest

## Directories

| Directory | Purpose | Configurable |
|-----------|---------|-------------|
| `~/Downloads/neocut/` | Default output directory for processed files | Via `--output-dir` / `-d` |
| `~/.config/neostore/neocut/` | Config and runtime directory (ffmpeg binary, which.cmd shim, config.jsonl) | No |
| `~/.config/neostore/neocut/bin/` | Downloaded ffmpeg and installed neocut binary | No |

## TUI mode

When `--tui` is passed, neocut launches an interactive form powered by [huh](https://github.com/charmbracelet/huh). The form provides the same configuration options as CLI flags but in a visual interface:

- **Input file** — file picker / manual path
- **Silence threshold (dBFS)** — number input
- **Minimum silence length (ms)** — number input
- **Silence to keep (ms)** — number input
- **Seek step (ms)** — number input
- **Output file name** — text input

The TUI is triggered with:
```bash
neocut --tui
```
