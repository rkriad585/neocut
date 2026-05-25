# v0.2.2 — Plan (completed)

## Summary

Polish the config system, add output format flexibility, and improve discoverability. No new large features — focused on finishing what's been started.

---

## ✅ 1. `--format` flag (output codec)

**Status:** Done in `v0.2.2`

Users can now choose the output codec:

```bash
neocut -i input.mp3 --format wav
neocut -i input.mp3 --format flac
neocut -i input.mp3 --format mp3   # default
```

- Passes format to `godub.NewExporter().WithDstFormat()`
- Auto-updates output file extension

---

## ✅ 2. `--bitrate` flag

**Status:** Done in `v0.2.2`

```bash
neocut -i input.mp3 --format flac --bitrate 320
```

- Passes bitrate to `godub.NewExporter().WithBitRate()`
- Shows format + bitrate in the completion summary
- `0` = codec default

---

## ✅ 3. `--dry-run` flag (preview)

**Status:** Done in `v0.2.2`

```bash
neocut -i input.mp3 --dry-run
```

Prints stats (segments, input/output duration, removal %) without actually exporting. Lets users tune parameters iteratively without waiting for export.

---

## 4. `--save` flag (persist current params)

**TODO:** Save current CLI flags as config default or named preset.

---

## 5. `--preset list` / `list` command

**TODO:** List available presets and history from config.jsonl.

---

## Effort estimate

| Item | Complexity | Status |
|------|-----------|--------|
| `--format` | low | ✅ Done |
| `--bitrate` | low | ✅ Done |
| `--dry-run` | low | ✅ Done |
| `--save` | medium | ❌ Pending |
| `--preset list` / `list` cmd | low | ❌ Pending |
