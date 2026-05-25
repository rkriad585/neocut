# v0.2.9 — Plan

## Summary

Polish the config system, add output format flexibility, and improve discoverability. No new large features — focused on finishing what's been started.

---

## 1. `--save` flag (persist current params)

**Problem:** Users can customize flags every run, but there's no way to save their preferences as the new default or as a named preset.

**Solution:**

```bash
# Save current CLI flags as the config default
neocut -i input.mp3 -m 500 -s -24 --save

# Save as a named preset
neocut -i input.mp3 -m 500 -s -24 --save mypreset
```

- Writes a `default` or `preset` entry to `config.jsonl`
- Only saves flags that differ from built-in defaults

**Files:** `internal/config/jsonl.go`, `internal/cmd/root.go`

---

## 2. `--preset list` and `list` command

**Problem:** Users can't see what presets are available without reading the JSONL file manually.

**Solution:**

```bash
# List available presets with their values
neocut --preset list

# Or as a subcommand
neocut list presets
neocut list history
```

Output:

```
Available presets:
  aggressive    m=500   s=-24.0  k=50   e=1
  gentle        m=2000  s=-10.0  k=200  e=5
  speech        m=800   s=-20.0  k=80   e=1

Recent history (last 5):
  2026-05-25  podcast.mp3 → podcast_no_silence.mp3  (45m→38m)
```

**Files:** `internal/config/jsonl.go`, `internal/cmd/root.go`, `internal/cmd/list.go`

---

## 3. `--format` flag (output codec)

**Problem:** Output is hardcoded to MP3. Users may want WAV for editing or FLAC for archival.

**Solution:**

```bash
neocut -i input.mp3 --format wav
neocut -i input.mp3 --format flac
neocut -i input.mp3 --format mp3   # default
```

- Passes format to `godub.NewExporter().WithDstFormat()`
- Auto-updates output file extension
- Updates export progress label

**Files:** `internal/core/processor.go`, `internal/config/config.go`

---

## 4. `--dry-run` flag (preview)

**Problem:** Users run a long processing job only to find the silence detection was too aggressive or too weak.

**Solution:**

```bash
neocut -i input.mp3 --dry-run
```

Prints stats (segments, input/output duration, removal %) without actually exporting. Lets users tune parameters iteratively without waiting for export.

**Files:** `internal/core/processor.go`, `internal/cmd/root.go`

---

## Effort estimate

| Item | Complexity | Files |
|------|-----------|-------|
| `--save` | medium | 2 |
| `--preset list` / `list` cmd | low | 3 |
| `--format` | low | 2 |
| `--dry-run` | low | 2 |

Total: ~1-2 hours implementation, ~30 min docs.
