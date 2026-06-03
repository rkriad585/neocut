package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

type DefaultEntry struct {
	MinSilenceLen int     `toml:"min_silence_len" json:"min_silence_len"`
	SilenceThresh float64 `toml:"silence_thresh" json:"silence_thresh"`
	KeepSilence   int     `toml:"keep_silence" json:"keep_silence"`
	SeekStep      int     `toml:"seek_step" json:"seek_step"`
	OutputDir     string  `toml:"output_dir" json:"output_dir"`
}

type PresetEntry struct {
	Name          string  `toml:"name" json:"name"`
	MinSilenceLen int     `toml:"min_silence_len" json:"min_silence_len"`
	SilenceThresh float64 `toml:"silence_thresh" json:"silence_thresh"`
	KeepSilence   int     `toml:"keep_silence" json:"keep_silence"`
	SeekStep      int     `toml:"seek_step" json:"seek_step"`
}

type HistoryEntry struct {
	Type      string `json:"type"`
	Input     string `json:"input"`
	Output    string `json:"output"`
	Dir       string `json:"dir"`
	Timestamp string `json:"timestamp"`
}

type tomlConfig struct {
	Default DefaultEntry  `toml:"default"`
	Presets []PresetEntry `toml:"preset"`
}

func InitConfigFile() error {
	path := ConfigFile("config.toml")
	if path == "" {
		return fmt.Errorf("cannot determine config directory")
	}

	if _, err := os.Stat(path); err == nil {
		return nil
	}

	EnsureConfigDir()

	jsonlPath := ConfigFile("config.jsonl")
	if data, err := os.ReadFile(jsonlPath); err == nil {
		var defaults *DefaultEntry
		var presets []PresetEntry
		for _, line := range strings.Split(string(data), "\n") {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			var raw struct {
				Type string `json:"type"`
			}
			if json.Unmarshal([]byte(line), &raw) != nil {
				continue
			}
			switch raw.Type {
			case "default":
				var d DefaultEntry
				if json.Unmarshal([]byte(line), &d) == nil {
					defaults = &d
				}
			case "preset":
				var p PresetEntry
				if json.Unmarshal([]byte(line), &p) == nil {
					presets = append(presets, p)
				}
			}
		}
		if defaults != nil || len(presets) > 0 {
			cfg := tomlConfig{}
			if defaults != nil {
				cfg.Default = *defaults
			}
			if len(presets) > 0 {
				cfg.Presets = presets
			}
			return writeTomlConfig(path, cfg)
		}
	}

	cfg := tomlConfig{
		Default: DefaultEntry{
			MinSilenceLen: 1000,
			SilenceThresh: -16.0,
			KeepSilence:   100,
			SeekStep:      1,
		},
		Presets: []PresetEntry{
			{Name: "aggressive", MinSilenceLen: 500, SilenceThresh: -24.0, KeepSilence: 50, SeekStep: 1},
			{Name: "gentle", MinSilenceLen: 2000, SilenceThresh: -10.0, KeepSilence: 200, SeekStep: 5},
			{Name: "speech", MinSilenceLen: 800, SilenceThresh: -20.0, KeepSilence: 80, SeekStep: 1},
		},
	}
	return writeTomlConfig(path, cfg)
}

func ReadConfig() ([]PresetEntry, *DefaultEntry, error) {
	path := ConfigFile("config.toml")
	if path == "" {
		return nil, nil, nil
	}

	var cfg tomlConfig
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return nil, nil, nil
	}

	d := cfg.Default
	return cfg.Presets, &d, nil
}

func WriteDefaults(d DefaultEntry) error {
	path := ConfigFile("config.toml")
	if path == "" {
		return fmt.Errorf("cannot determine config directory")
	}

	var cfg tomlConfig
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		EnsureConfigDir()
		cfg = tomlConfig{}
	}
	cfg.Default = d

	return writeTomlConfig(path, cfg)
}

func AppendHistory(cfg *Config) error {
	path := HistoryFile()
	if path == "" || cfg == nil {
		return nil
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	outputDir := GetOutputDir(cfg)
	outputName := cfg.OutputName
	if outputName == "" && cfg.InputFile != "" {
		base := filepath.Base(cfg.InputFile)
		ext := filepath.Ext(base)
		outputName = strings.TrimSuffix(base, ext) + "_no_silence.mp3"
	}

	entry := HistoryEntry{
		Type:      "history",
		Input:     cfg.InputFile,
		Output:    outputName,
		Dir:       outputDir,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(append(data, '\n'))
	return err
}

func writeTomlConfig(path string, cfg tomlConfig) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	var buf strings.Builder
	if err := toml.NewEncoder(&buf).Encode(cfg); err != nil {
		return err
	}

	return os.WriteFile(path, []byte(buf.String()), 0644)
}

func FindPreset(presets []PresetEntry, name string) *PresetEntry {
	for _, p := range presets {
		if strings.EqualFold(p.Name, name) {
			return &p
		}
	}
	return nil
}
