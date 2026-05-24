package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"neocut/internal/config"
)

func runSelfUninstall() int {
	configDir := config.ConfigDir()

	fmt.Println()
	fmt.Println("  Uninstalling neocut...")
	fmt.Println()

	exePath, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "  Error resolving binary path: %v\n", err)
		return 1
	}
	realPath, err := filepath.EvalSymlinks(exePath)
	if err == nil {
		exePath = realPath
	}

	binaryInside := strings.HasPrefix(exePath, configDir+string(os.PathSeparator))

	if runtime.GOOS == "windows" && binaryInside {
		filepath.WalkDir(configDir, func(path string, d os.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return nil
			}
			if !strings.EqualFold(path, exePath) {
				os.Remove(path)
			}
			return nil
		})
		binDir := filepath.Dir(exePath)
		if entries, _ := os.ReadDir(binDir); len(entries) == 1 {
			os.Remove(binDir)
		}
		fmt.Println("  ✓  Removed config files.")

		batContent := fmt.Sprintf("@echo off\r\ntimeout /t 1 /nobreak >nul\r\nrmdir /s /q \"%s\" 2>nul\r\necho   neocut has been uninstalled.\r\ndel /f /q \"%%~f0\"\r\n", configDir)
		batPath := filepath.Join(os.TempDir(), "neocut-uninstall.bat")
		if err := os.WriteFile(batPath, []byte(batContent), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "  Error creating uninstall script: %v\n", err)
			fmt.Println("  Please manually delete:", exePath)
			return 1
		}
		fmt.Println("  ✓  Uninstall script created. neocut will be fully removed shortly.")
		exec.Command("cmd", "/C", "start", "/B", batPath).Start()
	} else {
		if _, err := os.Stat(configDir); err == nil {
			if err := os.RemoveAll(configDir); err != nil {
				fmt.Fprintf(os.Stderr, "  Error removing config directory: %v\n", err)
				return 1
			}
			fmt.Println("  ✓  Removed config directory.")
		} else {
			fmt.Println("  ✓  No config directory found.")
		}

		if runtime.GOOS == "windows" {
			batContent := fmt.Sprintf("@echo off\r\ntimeout /t 1 /nobreak >nul\r\ndel /f /q \"%s\" 2>nul\r\necho   neocut has been uninstalled.\r\ndel /f /q \"%%~f0\"\r\n", exePath)
			batPath := filepath.Join(os.TempDir(), "neocut-uninstall.bat")
			if err := os.WriteFile(batPath, []byte(batContent), 0644); err != nil {
				fmt.Fprintf(os.Stderr, "  Error creating uninstall script: %v\n", err)
				fmt.Println("  Please manually delete:", exePath)
				return 1
			}
			fmt.Println("  ✓  Uninstall script created. Binary will be deleted shortly.")
			exec.Command("cmd", "/C", "start", "/B", batPath).Start()
		} else {
			if err := os.Remove(exePath); err != nil {
				fmt.Fprintf(os.Stderr, "  Error deleting binary: %v\n", err)
				fmt.Println("  Please manually delete:", exePath)
				return 1
			}
			fmt.Println("  ✓  Deleted binary:", exePath)
		}
	}

	fmt.Println()
	fmt.Println("  To remove neocut from your PATH, edit your shell rc file")
	fmt.Println("  and delete the line containing 'neostore/neocut/bin'.")
	fmt.Println()

	if runtime.GOOS == "windows" {
		fmt.Println("  Or run: irm https://raw.githubusercontent.com/rkriad585/neocut/main/installer.ps1 | iex -- --selfuninstall")
	} else {
		fmt.Println("  Or run: curl -fsSL https://raw.githubusercontent.com/rkriad585/neocut/main/installer.sh | sh -s -- --selfuninstall")
	}

	fmt.Println()
	fmt.Println("  Restart your terminal for PATH changes to take effect.")
	fmt.Println()
	return 0
}
