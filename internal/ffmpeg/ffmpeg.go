package ffmpeg

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

const binSubDir = "bin"

func BinDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "neostore", "neocut", binSubDir)
}

func Ensure() error {
	binDir := BinDir()
	if binDir != "" {
		os.MkdirAll(binDir, 0755)
	}
	addToPATH(binDir)

	if _, err := exec.LookPath("ffmpeg"); err == nil {
		ensureWhichShim(binDir)
		return nil
	}

	if binDir == "" {
		return fmt.Errorf("cannot determine home directory for ffmpeg installation")
	}

	fmt.Println()
	fmt.Println("  ffmpeg not found on system — downloading automatically...")
	fmt.Println()

	if err := install(binDir); err != nil {
		return fmt.Errorf("failed to install ffmpeg: %w", err)
	}

	ensureWhichShim(binDir)

	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return fmt.Errorf("ffmpeg installed but still not found on PATH: %w", err)
	}

	return nil
}

func ensureWhichShim(binDir string) {
	if runtime.GOOS != "windows" || binDir == "" {
		return
	}
	whichPath := filepath.Join(binDir, "which.cmd")
	if _, err := os.Stat(whichPath); err == nil {
		return
	}
	content := `@echo off
where %1
exit /b %ERRORLEVEL%
`
	if err := os.WriteFile(whichPath, []byte(content), 0755); err != nil {
		fmt.Printf("  Warning: could not create which shim: %v\n", err)
	}
}

func addToPATH(dirs ...string) {
	current := os.Getenv("PATH")
	for _, d := range dirs {
		if d == "" {
			continue
		}
		if !pathContains(current, d) {
			if current == "" {
				current = d
			} else {
				current = d + string(os.PathListSeparator) + current
			}
		}
	}
	os.Setenv("PATH", current)
}

func pathContains(path, dir string) bool {
	for _, p := range filepath.SplitList(path) {
		if p == dir {
			return true
		}
	}
	return false
}

func install(binDir string) error {
	url := downloadURL()

	fmt.Printf("  Downloading from %s\n", url)

	archivePath, err := downloadWithProgress(url, binDir)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	fmt.Println()

	if err := extractFfmpeg(archivePath, binDir); err != nil {
		return fmt.Errorf("extraction failed: %w", err)
	}

	if err := os.Remove(archivePath); err != nil {
		fmt.Printf("  Warning: could not remove archive: %v\n", err)
	}

	return nil
}

func downloadURL() string {
	switch runtime.GOOS {
	case "windows":
		return "https://www.gyan.dev/ffmpeg/builds/ffmpeg-release-essentials.zip"
	case "linux":
		arch := runtime.GOARCH
		if arch == "arm64" {
			arch = "aarch64"
		}
		return fmt.Sprintf("https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-%s-static.tar.xz", arch)
	case "darwin":
		return "https://evermeet.cx/ffmpeg/get/zip"
	default:
		return ""
	}
}

func extractFfmpeg(archivePath, binDir string) error {
	switch {
	case filepath.Ext(archivePath) == ".zip":
		return extractZip(archivePath, binDir)
	case filepath.Ext(archivePath) == ".xz" || filepath.Ext(archivePath[:len(archivePath)-3]) == ".tar":
		return extractTarXz(archivePath, binDir)
	default:
		return fmt.Errorf("unsupported archive format: %s", filepath.Ext(archivePath))
	}
}
