package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type MetaEntry struct {
	Type    string `json:"type"`
	Project string `json:"project"`
	Version string `json:"version"`
	Created string `json:"created"`
}

type DefaultEntry struct {
	Type          string  `json:"type"`
	MinSilenceLen int     `json:"min_silence_len"`
	SilenceThresh float64 `json:"silence_thresh"`
	KeepSilence   int     `json:"keep_silence"`
	SeekStep      int     `json:"seek_step"`
	OutputDir     string  `json:"output_dir"`
}

type PresetEntry struct {
	Type          string  `json:"type"`
	Name          string  `json:"name"`
	MinSilenceLen int     `json:"min_silence_len"`
	SilenceThresh float64 `json:"silence_thresh"`
	KeepSilence   int     `json:"keep_silence"`
	SeekStep      int     `json:"seek_step"`
}

type HistoryEntry struct {
	Type      string `json:"type"`
	Input     string `json:"input"`
	Output    string `json:"output"`
	Dir       string `json:"dir"`
	Timestamp string `json:"timestamp"`
}

func ConfigFile() string {
	dir := ConfigDir()
	if dir == "" {
		return ""
	}
	return filepath.Join(dir, "config.jsonl")
}

func InitConfigFile() error {
	path := ConfigFile()
	if path == "" {
		return fmt.Errorf("cannot determine config directory")
	}

	if _, err := os.Stat(path); err == nil {
		return nil
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	entries := []string{
		`{"type":"meta","project":"neocut","version":"1","created":"` + time.Now().UTC().Format(time.RFC3339) + `"}`,
		`{"type":"default","min_silence_len":1000,"silence_thresh":-16.0,"keep_silence":100,"seek_step":1,"output_dir":""}`,
		`{"type":"preset","name":"aggressive","min_silence_len":500,"silence_thresh":-24.0,"keep_silence":50,"seek_step":1}`,
		`{"type":"preset","name":"gentle","min_silence_len":2000,"silence_thresh":-10.0,"keep_silence":200,"seek_step":5}`,
		`{"type":"preset","name":"speech","min_silence_len":800,"silence_thresh":-20.0,"keep_silence":80,"seek_step":1}`,
	}

	return os.WriteFile(path, []byte(strings.Join(entries, "\n")+"\n"), 0644)
}

func ReadConfig() ([]PresetEntry, *DefaultEntry, error) {
	path := ConfigFile()
	if path == "" {
		return nil, nil, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, nil
	}

	var presets []PresetEntry
	var defaults *DefaultEntry

	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		var raw struct {
			Type string `json:"type"`
		}
		if err := json.Unmarshal([]byte(line), &raw); err != nil {
			continue
		}

		switch raw.Type {
		case "default":
			var d DefaultEntry
			if err := json.Unmarshal([]byte(line), &d); err == nil {
				defaults = &d
			}
		case "preset":
			var p PresetEntry
			if err := json.Unmarshal([]byte(line), &p); err == nil {
				presets = append(presets, p)
			}
		}
	}

	return presets, defaults, nil
}

func AppendHistory(cfg *Config) error {
	path := ConfigFile()
	if path == "" || cfg == nil {
		return nil
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil
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

func WriteDefaults(d DefaultEntry) error {
	path := ConfigFile()
	if path == "" {
		return fmt.Errorf("cannot determine config directory")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var newLines []string
	replaced := false
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var raw struct {
			Type string `json:"type"`
		}
		if err := json.Unmarshal([]byte(line), &raw); err == nil && raw.Type == "default" {
			jsonData, _ := json.Marshal(d)
			newLines = append(newLines, string(jsonData))
			replaced = true
			continue
		}
		newLines = append(newLines, line)
	}

	if !replaced {
		jsonData, _ := json.Marshal(d)
		newLines = append(newLines, string(jsonData))
	}

	return os.WriteFile(path, []byte(strings.Join(newLines, "\n")+"\n"), 0644)
}

func FindPreset(presets []PresetEntry, name string) *PresetEntry {
	for _, p := range presets {
		if strings.EqualFold(p.Name, name) {
			return &p
		}
	}
	return nil
}
