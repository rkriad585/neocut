package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	Commit         = "unknown"
	PublisherName  = "unknown"
	PublisherEmail = "unknown"
)

type Config struct {
	InputFile     string
	OutputName    string
	ConfigFile    string
	MinSilenceLen int
	SilenceThresh float64
	KeepSilence   int
	SeekStep      int
}

func ReadVersion() string {
	data, err := os.ReadFile(".version")
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(data))
}

func PrintBanner(version, commit string) {
	fmt.Println()
	fmt.Println("╭──────────────── neocut ───────────────────╮")
	fmt.Printf("│      Author : RK Riad Khan                │\n")
	fmt.Printf("│      Version: %-28s│\n", version)
	fmt.Printf("│      Commit : %-28s│\n", commit)
	fmt.Println("│      GitHub : rkriad585/neocut             │")
	fmt.Println("╰───────────────────────────────────────────╯")
	fmt.Println()
}

func GetOutputDir() string {
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
