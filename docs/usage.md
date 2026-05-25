# Usage

## Command line

```
neocut [flags]
neocut [command]
```

## Global flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--input` | `-i` | string | `""` | Path to the input MP3 file |
| `--output` | `-o` | string | auto | Output filename (saved to output-dir) |
| `--output-dir` | `-d` | string | `~/Downloads/neocut/` | Custom output directory |
| `--tui` | `-t` | bool | `false` | Launch interactive TUI form |
| `--config` | `-c` | bool | `false` | Edit project config interactively |
| `--quiet` | `-q` | bool | `false` | Suppress banner, spinners, and progress |
| `--format` | `-f` | string | `mp3` | Output codec: `mp3`, `wav`, `flac` |
| `--bitrate` | `-b` | int | `0` | Output bitrate in kbps (e.g. `192`, `320`); 0 = codec default |
| `--dry-run` | | bool | `false` | Preview stats (segments, durations, removal %) without exporting |
| `--min-silence-len` | `-m` | int | `1000` | Minimum silence length in milliseconds |
| `--silence-thresh` | `-s` | float | `-16` | Silence threshold in dBFS |
| `--keep-silence` | `-k` | int | `100` | Silence to keep at boundaries in ms |
| `--seek-step` | `-e` | int | `1` | Seek step in milliseconds (lower = more precise) |
| `--selfuninstall` | | bool | `false` | Remove neocut and its config directory |
| `--version` | `-v` | | | Print version and exit |
| `--help` | `-h` | | | Print help and exit |

## Commands

| Command | Description |
|---------|-------------|
| `self-update` | Fetch and install the latest version from GitHub |
| `completion` | Generate shell autocompletion script (bash, zsh, fish, powershell) |
| `help [command]` | Show help for a specific command |

## Examples

### Basic silence removal

```bash
neocut -i podcast_episode.mp3
```

Loads `podcast_episode.mp3`, removes silent portions, saves to `~/Downloads/neocut/podcast_episode_no_silence.mp3`.

### Custom output name

```bash
neocut -i recording.mp3 -o cleaned_recording.mp3
```

### Aggressive removal (shorter silence, lower threshold)

```bash
neocut -i lecture.mp3 -m 500 -s -24 -k 50
```

- `-m 500`: treat any 500ms+ segment as silence (default: 1000ms)
- `-s -24`: silence threshold at -24 dBFS (default: -16 dBFS, more sensitive)
- `-k 50`: keep 50ms of silence at boundaries (default: 100ms)

### Gentle removal (only long pauses)

```bash
neocut -i interview.mp3 -m 2000 -s -12 -k 200
```

- `-m 2000`: only remove silences of 2 seconds or more
- `-s -12`: higher threshold detects fewer segments as silence
- `-k 200`: keep 200ms of natural silence at boundaries

### Custom output directory

```bash
neocut -i input.mp3 -d /tmp/processed
```

Saves output to `/tmp/processed/input_no_silence.mp3` instead of `~/Downloads/neocut/`.

### Quiet mode (scripting)

```bash
neocut -i input.mp3 -q
```

Suppresses the banner, animated spinners, and progress bar. Only prints the output path on success. Useful for scripts and pipes.

### High precision seeking

```bash
neocut -i vocals.mp3 -e 0.5
```

Lowers the seek step to 0.5ms for more precise silence boundaries (slower processing).

### FLAC output with custom bitrate

```bash
neocut -i recording.mp3 --format flac --bitrate 320
```

Exports to `~/Downloads/neocut/recording_no_silence.flac` at 320 kbps (FLAC supports lossless at 1411 kbps by default).

### WAV output

```bash
neocut -i lecture.mp3 --format wav
```

Exports to `~/Downloads/neocut/lecture_no_silence.wav`.

### Dry run

```bash
neocut -i podcast.mp3 --dry-run
```

Processes the file and shows stats (segments, input/output duration, removal %) but does **not** create an output file.

### Format + bitrate + quiet (scripting)

```bash
neocut -i input.mp3 -f wav -b 192 -q
```

Quiet mode: suppresses banner, spinners, and progress. Only prints the output path.

### Interactive TUI mode

```bash
neocut --tui
```

Launches a [huh](https://github.com/charmbracelet/huh) form where you fill in all options interactively instead of passing flags.

### Self-update

```bash
neocut self-update
```

Checks `https://raw.githubusercontent.com/rkriad585/neocut/main/.version`, downloads the newer binary, and replaces the running executable.

### Uninstall

```bash
neocut --selfuninstall
```

Removes `~/.config/neostore/neocut/` and the binary, then prints PATH cleanup instructions.

### Shell completion

```bash
# Bash
source <(neocut completion bash)

# zsh
source <(neocut completion zsh)

# fish
neocut completion fish | source

# PowerShell
neocut completion powershell | Out-String | Invoke-Expression
```

## Config editor

```bash
neocut -c
```

Opens an interactive TUI to edit `~/.config/neostore/neocut/config.jsonl`:

1. **Default parameters** — set min silence length, silence threshold, keep silence, seek step, output directory
2. **Presets overview** — shows all configured presets with their values
3. **History browser** — last 5 processing runs with timestamps
4. **Save** — confirms and writes changes

## Processing stats

After a successful run, neocut displays a summary:

```
    Segments:     42
    Input:        45m 30s
    Output:       38m 12s
    Removed:      7m 18s (16.0%)
```

- **Segments**: number of non-silent chunks found
- **Input**: original audio duration
- **Output**: duration after silence removal
- **Removed**: total silence cut, with percentage of original

## Output

- All processed files go to `~/Downloads/neocut/`
- Default naming: `{input_filename_without_ext}_no_silence.{ext}`
- Output format: MP3 (default), WAV, or FLAC — set via `--format`
- Output bitrate configurable via `--bitrate` (applies to MP3, FLAC)
- If the output directory doesn't exist, it is created automatically

## Processing pipeline

Each step shows an animated spinner:

1. **Loading** — reads the MP3 file via godub + ffmpeg
2. **Detecting and removing silence** — calls `godub.SplitOnSilence` with the configured parameters
3. **Recombining audio chunks** — appends all non-silent segments back together
4. **Exporting** — writes the result with a progress bar

If ffmpeg is not on your PATH, it is auto-downloaded before processing begins.
