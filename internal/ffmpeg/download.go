package ffmpeg

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func downloadWithProgress(url, destDir string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	fileName := filepath.Base(url)
	if fileName == "" || fileName == "." || fileName == "zip" {
		fileName = "ffmpeg-download"
		switch {
		case strings.Contains(resp.Header.Get("Content-Type"), "zip"):
			fileName += ".zip"
		default:
			if idx := strings.Index(url, ".zip"); idx > 0 {
				fileName = "ffmpeg.zip"
			} else if idx := strings.Index(url, ".xz"); idx > 0 {
				fileName = "ffmpeg.tar.xz"
			} else {
				fileName += ".dat"
			}
		}
	}

	destPath := filepath.Join(destDir, fileName)
	tmpPath := destPath + ".tmp"

	f, err := os.Create(tmpPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	total := resp.ContentLength
	var written int64

	done := make(chan struct{})
	go func() {
		blocks := []string{" ", "▏", "▎", "▍", "▌", "▋", "▊", "▉", "█"}
		width := 30
		i := 0
		for {
			select {
			case <-done:
				bar := strings.Repeat("█", width)
				fmt.Printf("\r  Downloading [%s] 100%%", bar)
				return
			default:
				if total > 0 {
					pct := float64(written) / float64(total)
					full := int(pct * float64(width))
					part := int((pct*float64(width) - float64(full)) * 8)
					bar := strings.Repeat("█", full)
					if part > 0 && full < width {
						bar += blocks[part]
						full++
					}
					bar += strings.Repeat(" ", width-full)
					fmt.Printf("\r  Downloading [%s] %3d%%", bar, int(pct*100))
				} else {
					pos := i % (width * 8)
					full := pos / 8
					part := pos % 8
					bar := strings.Repeat("█", full)
					if part > 0 {
						bar += blocks[part]
					}
					bar += strings.Repeat(" ", width-full-1)
					fmt.Printf("\r  Downloading [%s]", bar)
				}
				i++
				time.Sleep(40 * time.Millisecond)
			}
		}
	}()

	buf := make([]byte, 32*1024)
	for {
		n, readErr := resp.Body.Read(buf)
		if n > 0 {
			wn, writeErr := f.Write(buf[:n])
			if writeErr != nil {
				close(done)
				time.Sleep(60 * time.Millisecond)
				fmt.Println()
				return "", writeErr
			}
			written += int64(wn)
		}
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			close(done)
			time.Sleep(60 * time.Millisecond)
			fmt.Println()
			return "", readErr
		}
	}

	close(done)
	time.Sleep(60 * time.Millisecond)
	fmt.Println()

	f.Close()

	if err := os.Rename(tmpPath, destPath); err != nil {
		return "", err
	}

	return destPath, nil
}

func extractZip(src, destDir string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	var found string
	for _, f := range r.File {
		if f.FileInfo().IsDir() {
			continue
		}
		if filepath.Base(f.Name) == "ffmpeg.exe" || filepath.Base(f.Name) == "ffmpeg" {
			rc, err := f.Open()
			if err != nil {
				return err
			}
			destName := filepath.Base(f.Name)
			destPath := filepath.Join(destDir, destName)
			out, err := os.Create(destPath)
			if err != nil {
				rc.Close()
				return err
			}
			_, err = io.Copy(out, rc)
			out.Close()
			rc.Close()
			if err != nil {
				return err
			}
			os.Chmod(destPath, 0755)
			found = destPath
			break
		}
	}
	if found == "" {
		return fmt.Errorf("ffmpeg binary not found in archive")
	}
	return nil
}

func extractTarXz(src, destDir string) error {
	cmd := exec.Command("tar", "-xf", src, "-C", destDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("tar extract failed: %s: %s", err, string(output))
	}

	entries, err := os.ReadDir(destDir)
	if err != nil {
		return fmt.Errorf("reading extracted dir: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() && strings.Contains(entry.Name(), "ffmpeg") {
			binaryPath := filepath.Join(destDir, entry.Name(), "ffmpeg")
			if _, statErr := os.Stat(binaryPath); statErr == nil {
				input, readErr := os.ReadFile(binaryPath)
				if readErr != nil {
					return readErr
				}
				if writeErr := os.WriteFile(filepath.Join(destDir, "ffmpeg"), input, 0755); writeErr != nil {
					return writeErr
				}
				os.RemoveAll(filepath.Join(destDir, entry.Name()))
				return nil
			}
		}
	}

	return fmt.Errorf("ffmpeg binary not found in extracted archive")
}
