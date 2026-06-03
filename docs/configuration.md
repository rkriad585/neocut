# Configuration

## CLI flags

All configuration is passed via CLI flags. Persistent defaults and presets are stored in `~/.config/neostore/neocut/config.toml`.

| Flag | Default | Description |
|------|---------|-------------|
| `-i / --input` | `""` | Input MP3 file path (required unless using `--tui`) |
| `-o / --output` | auto | Output filename; defaults to `{input}_no_silence.{format}` |
| `-d / --output-dir` | `~/Downloads/neostore/neocut/` | Custom output directory |
| `-q / --quiet` | `false` | Suppress banner, spinners, and progress |
| `-f / --format` | `mp3` | Output codec: `mp3`, `wav`, `flac` |
| `-b / --bitrate` | `0` | Output bitrate in kbps (`0` = codec default) |
| `--dry-run` | `false` | Preview stats without exporting |
| `-m / --min-silence-len` | `1000` ms | Minimum duration of silence to detect |
| `-s / --silence-thresh` | `-16` dBFS | Volume threshold; segments below this are considered silence |
| `-k / --keep-silence` | `100` ms | Silence to preserve at the boundaries of kept segments |
| `-e / --seek-step` | `1` ms | Precision of silence detection (lower = more accurate but slower) |
| `--preset` | `""` | Load a named preset from config.toml (`aggressive`, `gentle`, `speech`) |
| `--theme` | `""` | Theme name (e.g. `dark`, `sunny_beach_day`, `olive_garden_feast`) |
| `--color-mode` | `"auto"` | Color mode: `auto`, `dark`, `light` |
| `-c / --config` | `false` | Edit project config interactively via TUI |
| `-t / --tui` | `false` | Launch interactive form instead of parsing flags |

## Config file: config.toml

Located at `~/.config/neostore/neocut/config.toml`, this TOML file stores:

| Section | Purpose |
|---------|---------|
| `[default]` | Persistent default values for all flags + UI settings |
| `[[preset]]` | Named presets (e.g. `aggressive`, `gentle`, `speech`) |

Processing history is stored in a separate `history.log` file in the same directory.

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
| `~/Downloads/neostore/neocut/` | Default output directory for processed files | Via `--output-dir` / `-d` |
| `~/.config/neostore/neocut/` | Config and runtime directory (ffmpeg binary, which.cmd shim, config.toml, history.log) | No |
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

## Theme system

neocut has 13 built-in color themes. Set via `--theme`, config, or `neocut -c` (config editor).

### Themes

| Config name | Label | Description |
|-------------|-------|-------------|
| `dark` | Dark | Classic dark mode |
| `light` | Light | Clean light theme |
| `sunny_beach_day` | Sunny Beach Day | Warm tropical palette (default) |
| `olive_garden_feast` | Olive Garden Feast | Earthy tones |
| `summer_ocean_breeze` | Summer Ocean Breeze | Cool coastal blues |
| `refreshing_summer_fun` | Refreshing Summer Fun | Bright citrus |
| `black_gold_elegance` | Black Gold Elegance | Premium dark + gold |
| `vibrant_color_fiesta` | Vibrant Color Fiesta | Neon festival |
| `light_steel` | Light Steel | Industrial grays |
| `golden_twilight` | Golden Twilight | Sunset amber |
| `deep_sea` | Deep Sea | Deep ocean blues |
| `bright_green` | Bright Green | Fresh nature green |
| `vivid_nightfall` | Vivid Nightfall | Purple twilight |

### Color modes

| Mode | Effect |
|------|--------|
| `auto` (default) | Use the selected theme's native colors |
| `dark` | Force Dark Theme (overrides the selected theme) |
| `light` | Force Light Theme (overrides the selected theme) |

### Config file fields

In `~/.config/neostore/neocut/config.toml`:

```toml
[default]
theme = "sunny_beach_day"
color_mode = "auto"
```

### CLI flags

```bash
# Use a specific theme
neocut --theme dark -i input.mp3

# Force dark mode
neocut --color-mode dark -i input.mp3

# Override both
neocut --theme olive_garden_feast --color-mode auto -i input.mp3
```
