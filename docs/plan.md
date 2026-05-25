# Plan

## v0.2.2 — Output format flexibility

### ✅ 1. `--format` flag (output codec)

Users can now choose the output codec:

```bash
neocut -i input.mp3 --format wav
neocut -i input.mp3 --format flac
neocut -i input.mp3 --format mp3   # default
```

- Passes format to `godub.NewExporter().WithDstFormat()`
- Auto-updates output file extension

---

### ✅ 2. `--bitrate` flag

```bash
neocut -i input.mp3 --format flac --bitrate 320
```

- Passes bitrate to `godub.NewExporter().WithBitRate()`
- Shows format + bitrate in the completion summary
- `0` = codec default

---

### ✅ 3. `--dry-run` flag (preview)

```bash
neocut -i input.mp3 --dry-run
```

Prints stats (segments, input/output duration, removal %) without actually exporting. Lets users tune parameters iteratively without waiting for export.

---

## v0.2.3 — Vendored godub + clean error handling

### ✅ 4. Vendor dependencies

- `go mod vendor` creates `vendor/` directory with all dependencies
- godub source vendored in-tree at `vendor/github.com/Vernacular-ai/godub/`
- Allows patching godub without maintaining a separate fork

### ✅ 5. Patch godub SplitOnSilence

- Added `len(notSilenceRanges) == 0` guard to return clean error instead of index-out-of-range panic
- Added `len(nonsilentRanges) > 0` guard preventing potential OOB in `cmp.Equal` dead code

---

## v0.2.4 — Bug fixes

### ✅ 6. Fix normalization target in godub

- Root cause: `matchTargetAmp(seg, -20.0)` normalizes audio to -20 dBFS, but silence threshold is -16 dBFS. After normalization, the entire signal appears below threshold.
- Fix: changed normalization target from `-20.0` to `-10.0` so normalized signal sits above the -16 dBFS threshold.

### ✅ 7. Fix self-update on Windows

Three bugs in `replaceWindows`:
1. Format-string arg mismatch (7 verbs, 6 args) — last `%s` rendered as garbage
2. Wait-loop checked `tmpPath` instead of `exePath` — infinite loop
3. `rename` source was `exePath` instead of `tmpPath` — no-op

Result: old binary deleted, new one never moved, leaving no binary at all.

### ✅ 8. Unit tests

Added 7 test files with 85+ tests across all 6 internal packages:
- config (config + jsonl): 10 tests
- core: 3 tests (21 sub-cases for fmtDuration)
- ffmpeg (ffmpeg + download): 21 tests
- update: 7 tests
- cmd: 6 tests

---

## v1.0.1 — Stable release

### ✅ 9. Version bump to v1.0.1

First stable release. All known bugs from v0.2.x are resolved:
- godub normalization bug (all audio detected as silence)
- self-update Windows bat script (binary deleted but not replaced)
- Complete unit test coverage across all internal packages

---

## v1.0.2 — Fix self-update for real

### ✅ 10. Replace batch script with direct rename

The batch script approach was fundamentally unreliable:
- Format-string bugs caused args to map to wrong placeholders
- Wait-loop checking was incorrect (infinite loop)
- rename source/dest were swapped

Replaced with the proven approach from neodlp: **rename running exe → .old**, then **rename temp → exe**. Windows allows renaming a running executable (just not deleting it). No batch script needed.

---

## Future

### 10. `--save` flag (persist current params)

Save current CLI flags as config default or named preset.

### 11. `--preset list` / `list` command

List available presets and history from config.jsonl.

## Effort estimate

| Item | Complexity | Version | Status |
|------|-----------|---------|--------|
| `--format` | low | v0.2.2 | ✅ Done |
| `--bitrate` | low | v0.2.2 | ✅ Done |
| `--dry-run` | low | v0.2.2 | ✅ Done |
| Vendor godub | low | v0.2.3 | ✅ Done |
| Patch godub panic | low | v0.2.3 | ✅ Done |
| Fix normalization | medium | v0.2.4 | ✅ Done |
| Fix self-update | medium | v0.2.4 | ✅ Done |
| Unit tests | medium | v0.2.4 | ✅ Done |
| Stable release | low | v1.0.1 | ✅ Done |
| `--save` | medium | — | ❌ Pending |
| `--preset list` / `list` cmd | low | — | ❌ Pending |
