//go:generate go run gen.go

package config

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

//go:embed version.txt
var embeddedVersion string

var (
	Commit         = "unknown"
	Version        = "unknown"
	PublisherName  = "unknown"
	PublisherEmail = "unknown"
)

type Config struct {
	InputFile     string
	OutputName    string
	MinSilenceLen int
	SilenceThresh float64
	KeepSilence   int
	SeekStep      int
	OutputDir     string
	Quiet         bool
}

func ReadVersion() string {
	if Version != "unknown" && Version != "" {
		return Version
	}
	if embeddedVersion != "" {
		return strings.TrimSpace(embeddedVersion)
	}
	data, err := os.ReadFile(".version")
	if err == nil {
		return strings.TrimSpace(string(data))
	}
	return "unknown"
}

func PrintBanner(version, commit string) {
	fmt.Println()
	fmt.Println("╭──────────────── neocut ───────────────────╮")
	fmt.Printf("│      Author : RK Riad Khan                │\n")
	fmt.Printf("│      Version: %-28s│\n", version)
	fmt.Printf("│      Commit : %-28s│\n", commit)
	fmt.Printf("│      %-38s│\n", PublisherName)
	fmt.Printf("│      %-38s│\n", PublisherEmail)
	fmt.Println("│      GitHub : rkriad585/neocut             │")
	fmt.Println("╰───────────────────────────────────────────╯")
	fmt.Println()
}

func GetOutputDir(cfg *Config) string {
	if cfg != nil && cfg.OutputDir != "" {
		return cfg.OutputDir
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return filepath.Join(home, "Downloads", "neocut")
}

func ConfigDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "neostore", "neocut")
}

func EnsureConfigDir() {
	configDir := ConfigDir()
	if configDir == "" {
		return
	}
	if err := os.MkdirAll(configDir, 0755); err != nil {
		fmt.Printf("Warning: could not create config dir: %v\n", err)
	}
}
