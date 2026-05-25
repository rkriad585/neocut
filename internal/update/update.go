package update

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
)

const (
	gitHubUser = "rkriad585"
	project    = "neocut"
)

func LatestVersion() (string, error) {
	url := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/main/.version", gitHubUser, project)
	resp, err := http.Get(url)
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
		fmt.Println("  ✓ Already up to date!")
		fmt.Println()
		return nil
	}

	fmt.Println()
	fmt.Printf("  Updating to %s...\n", latest)

	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to locate running executable: %w", err)
	}

	exePath, err = filepathEval(exePath)
	if err != nil {
		return err
	}

	url := DownloadURL(latest)
	tmpPath := exePath + ".tmp"
	os.Remove(tmpPath)

	fmt.Printf("  Downloading from: %s\n", url)

	client := &http.Client{Timeout: 5 * time.Minute}
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to download release binary: server returned %s", resp.Status)
	}

	tmpFile, err := os.OpenFile(tmpPath, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0755)
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}

	cleanup := true
	defer func() {
		if cleanup {
			tmpFile.Close()
			os.Remove(tmpPath)
		}
	}()

	progress := &WriteCounter{}
	_, err = io.Copy(tmpFile, io.TeeReader(resp.Body, progress))
	if err != nil {
		return fmt.Errorf("failed to write binary content: %w", err)
	}
	tmpFile.Close()
	fmt.Println()

	oldPath := exePath + ".old"
	os.Remove(oldPath)

	if runtime.GOOS == "windows" {
		err = os.Rename(exePath, oldPath)
		if err != nil {
			return fmt.Errorf("failed to rename running executable: %w", err)
		}
		err = os.Rename(tmpPath, exePath)
		if err != nil {
			os.Rename(oldPath, exePath)
			return fmt.Errorf("failed to install new executable: %w", err)
		}
		cleanup = false
		fmt.Printf("  ✓ Success! neocut has been updated to %s.\n", latest)
		fmt.Println("  Note: You can safely delete the old executable (neocut.exe.old) after closing this session.")
	} else {
		err = os.Rename(tmpPath, exePath)
		if err != nil {
			return fmt.Errorf("failed to install new executable: %w", err)
		}
		cleanup = false
		os.Chmod(exePath, 0755)
		fmt.Printf("  ✓ Success! neocut has been updated to %s.\n", latest)
	}

	fmt.Println()
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
		return resolved, nil
	}
	return path, nil
}

type WriteCounter struct {
	Total uint64
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	fmt.Printf("\r  Downloaded: %.2f MB", float64(wc.Total)/1024/1024)
	return n, nil
}
