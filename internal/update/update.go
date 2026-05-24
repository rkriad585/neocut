package update

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

const (
	gitHubUser = "rkriad585"
	project    = "neocut"
)

var httpClient = &http.Client{Timeout: 10 * time.Second}

func LatestVersion() (string, error) {
	url := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/main/.version", gitHubUser, project)
	resp, err := httpClient.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch version: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("version fetch returned HTTP %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(data)), nil
}

func DownloadURL(version string) string {
	binary := fmt.Sprintf("%s-%s-%s", project, runtime.GOOS, runtime.GOARCH)
	if runtime.GOOS == "windows" {
		binary += ".exe"
	}
	return fmt.Sprintf("https://github.com/%s/%s/releases/download/%s/%s",
		gitHubUser, project, version, binary)
}

func Run(currentVersion string) error {
	fmt.Println()
	fmt.Println("  Checking for updates...")

	latest, err := LatestVersion()
	if err != nil {
		return fmt.Errorf("could not check for updates: %w", err)
	}

	fmt.Printf("  Current version: %s\n", currentVersion)
	fmt.Printf("  Latest version:  %s\n", latest)

	if latest == currentVersion {
		fmt.Println()
		fmt.Println("  Already up to date.")
		fmt.Println()
		return nil
	}

	fmt.Println()
	fmt.Println("  Update available! Downloading...")
	fmt.Println()

	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("cannot determine executable path: %w", err)
	}

	exePath, err = filepathEval(exePath)
	if err != nil {
		return err
	}

	url := DownloadURL(latest)
	tmpPath := exePath + ".update.tmp"

	if err := downloadWithProgress(url, tmpPath); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("download failed: %w", err)
	}

	fmt.Println()
	fmt.Println("  Installing update...")

	if err := replaceBinary(exePath, tmpPath); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("installation failed: %w", err)
	}

	fmt.Println()
	fmt.Printf("  Updated to %s successfully.\n", latest)
	fmt.Println()
	return nil
}

func replaceBinary(exePath, tmpPath string) error {
	if runtime.GOOS == "windows" {
		return replaceWindows(exePath, tmpPath)
	}
	return replaceUnix(exePath, tmpPath)
}

func replaceUnix(exePath, tmpPath string) error {
	if err := os.Chmod(tmpPath, 0755); err != nil {
		return err
	}
	if err := os.Rename(tmpPath, exePath); err != nil {
		return err
	}
	fmt.Println("  Restart the application to use the new version.")
	return nil
}

func replaceWindows(exePath, tmpPath string) error {
	script := fmt.Sprintf(`@echo off
setlocal
set "self=%%~f0"
:wait
ping -n 2 127.0.0.1 >nul 2>&1
if exist "%s" (
  del /f /q "%s" >nul 2>&1
  if exist "%s" goto wait
)
rename "%s" "%s" >nul 2>&1
if exist "%s" (
  start "" "%s"
)
del /f /q "%%self%" >nul 2>&1
`, exePath, exePath, tmpPath, exePath, exePath, exePath)

	scriptPath := exePath + ".update.bat"
	if err := os.WriteFile(scriptPath, []byte(script), 0755); err != nil {
		return err
	}

	cmd := exec.Command("cmd", "/c", scriptPath)
	if err := cmd.Start(); err != nil {
		os.Remove(scriptPath)
		return fmt.Errorf("failed to start update script: %w", err)
	}

	fmt.Println("  Update will complete after exit.")
	return nil
}

func filepathEval(path string) (string, error) {
	fi, err := os.Lstat(path)
	if err != nil {
		return "", err
	}
	if fi.Mode()&os.ModeSymlink != 0 {
		resolved, err := os.Readlink(path)
		if err != nil {
			return "", err
		}
		if !strings.Contains(resolved, "/") && !strings.Contains(resolved, "\\") {
			resolved = resolved
		}
		return resolved, nil
	}
	return path, nil
}

func downloadWithProgress(url, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	f, err := os.Create(dest)
	if err != nil {
		return err
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
				return writeErr
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
			return readErr
		}
	}

	close(done)
	time.Sleep(60 * time.Millisecond)
	fmt.Println()
	return nil
}
