# Troubleshooting

## ffmpeg not found

neocut requires ffmpeg to load and export audio. If ffmpeg is not on your PATH, neocut attempts to auto-download it.

**Auto-download fails:**
- Check your internet connection
- Ensure you have write permission to `~/.config/neostore/neocut/bin/`
- Manually install ffmpeg from [ffmpeg.org](https://ffmpeg.org/download.html)
- Set a custom PATH that includes ffmpeg before running neocut

**Manual installation (Windows):**
```
# Download from https://www.gyan.dev/ffmpeg/builds/
# Extract ffmpeg.exe and place it in:
%USERPROFILE%\.config\neostore\neocut\bin\
```

## "no audio remaining after silence removal"

This error means all audio segments were classified as silence and removed.

**Possible causes:**
- `--silence-thresh` is too sensitive (too high, e.g. `0` dBFS) — lower it to `-24` or `-30`
- `--min-silence-len` is too low — increase from `1000` to `2000` or higher
- The input file is already silent or contains only noise

**Fix:** Adjust your detection parameters:
```bash
neocut -i input.mp3 -s -30 -m 500 -k 100
```

## "no non-silent audio detected" on all files

This error in v0.2.3 meant every file was being classified as all-silence. This was caused by a godub bug where the audio normalization target (`-20 dBFS`) was below the silence threshold (`-16 dBFS`), making the entire normalized signal appear silent.

**Fix:** Upgrade to v1.0.1+ which corrects the normalization target to `-10 dBFS`.

## Split on silence panics

godub's `SplitOnSilence` can panic with "index out of range [-1]" on certain audio files with unusual silence patterns.

**Symptom:** The spinner shows "Detecting and removing silence" then the tool crashes or recovers.

**Workarounds:**
- Try a different `--seek-step` value (e.g. `-e 2` instead of `1`)
- Adjust `--silence-thresh` slightly
- Pre-process the audio with a different tool first

The `step()` function in the processor recovers panics and returns them as errors, so processing will fail gracefully rather than crash.

## `--selfuninstall` does not delete the binary on Windows

On Windows, a running executable cannot be deleted. neocut handles this by:
1. Creating a batch script (`%TEMP%\neocut-uninstall.bat`)
2. The script waits 1 second, then deletes the binary

If the batch script fails to run:
- Manually delete `%USERPROFILE%\.config\neostore\neocut\`
- Manually delete the binary
- Clean up your PATH

## `self-update` reports HTTP 404

This means the `.version` file on the `main` branch on GitHub does not match or the release does not exist yet.

- Make sure a GitHub release exists for the version returned by `.version`
- Check that the release binary follows the naming convention `neocut-{os}-{arch}[.exe]`

## `self-update` deletes the binary but doesn't replace it

This was a bug in v0.2.1–v1.0.1 where the Windows update script:
1. Had a format-string arg mismatch (7 verbs, 6 args)
2. Checked for the temp file instead of the exe in the wait-loop (infinite loop)
3. Renamed the exe to itself (no-op) instead of renaming the temp file

**Fix:** Upgrade to v1.0.2+ which replaces the fragile batch script with a direct rename (rename running exe → `.old`, then rename temp → exe). Windows allows renaming a running executable.
```powershell
irm https://raw.githubusercontent.com/rkriad585/neocut/main/installer.ps1 | iex
```

## Binary not found after installation

**Windows:** Restart your terminal. The installer updates your user PATH, but existing sessions don't pick it up.

**Unix:** Either restart your terminal or run:
```bash
source ~/.bashrc   # or source ~/.zshrc
```

## "ffmpeg setup failed"

This wraps an internal error from ffmpeg detection or download. Common causes:

- No internet connection for the auto-download
- The downloaded archive failed to extract
- Disk space or permission issues in the config directory

**Solution:** Manually place ffmpeg in `~/.config/neostore/neocut/bin/` and ensure it's executable.
