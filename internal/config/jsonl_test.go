package config

import (
	"os"
	"strings"
	"testing"
)

func setTempHomeDir(t *testing.T) string {
	t.Helper()
	origHome := os.Getenv("USERPROFILE")
	if origHome == "" {
		origHome = os.Getenv("HOME")
	}
	td := t.TempDir()
	os.Setenv("USERPROFILE", td)
	os.Setenv("HOME", td)
	t.Cleanup(func() {
		os.Setenv("USERPROFILE", origHome)
		os.Setenv("HOME", origHome)
	})
	return td
}

func TestInitConfigFile(t *testing.T) {
	setTempHomeDir(t)

	if err := InitConfigFile(); err != nil {
		t.Fatalf("InitConfigFile() error: %v", err)
	}

	path := ConfigFile("config.toml")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("config.toml should exist: %v", err)
	}

	content := string(data)
	checks := []struct {
		name   string
		match  string
		expect bool
	}{
		{"default section", "[default]", true},
		{"preset section", "[[preset]]", true},
		{"aggressive preset", `name = "aggressive"`, true},
		{"gentle preset", `name = "gentle"`, true},
		{"speech preset", `name = "speech"`, true},
		{"min_silence_len", "min_silence_len", true},
		{"silence_thresh", "silence_thresh", true},
	}
	for _, c := range checks {
		t.Run(c.name, func(t *testing.T) {
			got := strings.Contains(content, c.match)
			if got != c.expect {
				t.Errorf("contains %s: got %v, want %v", c.match, got, c.expect)
			}
		})
	}

	t.Run("idempotent — already exists", func(t *testing.T) {
		if err := InitConfigFile(); err != nil {
			t.Errorf("second call should succeed: %v", err)
		}
	})
}

func TestReadConfig(t *testing.T) {
	t.Run("reads presets and defaults", func(t *testing.T) {
		setTempHomeDir(t)
		InitConfigFile()

		presets, defaults, err := ReadConfig()
		if err != nil {
			t.Fatalf("ReadConfig() error: %v", err)
		}

		if defaults == nil {
			t.Fatal("defaults should not be nil")
		}
		if defaults.MinSilenceLen != 1000 {
			t.Errorf("default MinSilenceLen: expected 1000, got %d", defaults.MinSilenceLen)
		}
		if defaults.SilenceThresh != -16.0 {
			t.Errorf("default SilenceThresh: expected -16.0, got %f", defaults.SilenceThresh)
		}
		if defaults.KeepSilence != 100 {
			t.Errorf("default KeepSilence: expected 100, got %d", defaults.KeepSilence)
		}
		if defaults.SeekStep != 1 {
			t.Errorf("default SeekStep: expected 1, got %d", defaults.SeekStep)
		}

		if len(presets) != 3 {
			t.Fatalf("expected 3 presets, got %d", len(presets))
		}

		presetMap := make(map[string]PresetEntry)
		for _, p := range presets {
			presetMap[p.Name] = p
		}

		agg, ok := presetMap["aggressive"]
		if !ok {
			t.Fatal("aggressive preset not found")
		}
		if agg.MinSilenceLen != 500 {
			t.Errorf("aggressive MinSilenceLen: expected 500, got %d", agg.MinSilenceLen)
		}
		if agg.SilenceThresh != -24.0 {
			t.Errorf("aggressive SilenceThresh: expected -24.0, got %f", agg.SilenceThresh)
		}
		if agg.KeepSilence != 50 {
			t.Errorf("aggressive KeepSilence: expected 50, got %d", agg.KeepSilence)
		}
	})

	t.Run("no config file returns nil values", func(t *testing.T) {
		setTempHomeDir(t)

		presets, defaults, err := ReadConfig()
		if err != nil {
			t.Fatalf("ReadConfig() error: %v", err)
		}
		if presets != nil {
			t.Error("presets should be nil when no config file")
		}
		if defaults != nil {
			t.Error("defaults should be nil when no config file")
		}
	})
}

func TestWriteDefaults(t *testing.T) {
	t.Run("replaces existing default entry", func(t *testing.T) {
		setTempHomeDir(t)
		InitConfigFile()

		newDefault := DefaultEntry{
			MinSilenceLen: 999,
			SilenceThresh: -30.0,
			KeepSilence:   50,
			SeekStep:      10,
			OutputDir:     "/custom",
		}
		if err := WriteDefaults(newDefault); err != nil {
			t.Fatalf("WriteDefaults() error: %v", err)
		}

		_, defaults, _ := ReadConfig()
		if defaults == nil {
			t.Fatal("defaults should not be nil")
		}
		if defaults.MinSilenceLen != 999 {
			t.Errorf("expected 999, got %d", defaults.MinSilenceLen)
		}
		if defaults.SilenceThresh != -30.0 {
			t.Errorf("expected -30.0, got %f", defaults.SilenceThresh)
		}
		if defaults.OutputDir != "/custom" {
			t.Errorf("expected /custom, got %s", defaults.OutputDir)
		}
	})

	t.Run("creates file if none exists", func(t *testing.T) {
		setTempHomeDir(t)

		newDefault := DefaultEntry{
			MinSilenceLen: 500,
			SilenceThresh: -16.0,
			KeepSilence:   100,
			SeekStep:      1,
		}
		if err := WriteDefaults(newDefault); err != nil {
			t.Fatalf("WriteDefaults() error: %v", err)
		}

		_, defaults, _ := ReadConfig()
		if defaults == nil {
			t.Fatal("defaults should not be nil")
		}
	})
}

func TestAppendHistory(t *testing.T) {
	t.Run("appends history entry", func(t *testing.T) {
		setTempHomeDir(t)
		InitConfigFile()

		cfg := &Config{
			InputFile:  "/path/to/song.mp3",
			OutputName: "song_no_silence.mp3",
			OutputDir:  "/output",
		}
		if err := AppendHistory(cfg); err != nil {
			t.Fatalf("AppendHistory() error: %v", err)
		}

		data, _ := os.ReadFile(HistoryFile())
		lines := strings.Split(strings.TrimSpace(string(data)), "\n")
		lastLine := lines[len(lines)-1]

		if !strings.Contains(lastLine, `"type":"history"`) {
			t.Errorf("last line should be history entry, got: %s", lastLine)
		}
		if !strings.Contains(lastLine, "/path/to/song.mp3") {
			t.Errorf("history should contain input path")
		}
	})

	t.Run("nil config returns nil", func(t *testing.T) {
		setTempHomeDir(t)
		if err := AppendHistory(nil); err != nil {
			t.Errorf("AppendHistory(nil) should return nil, got: %v", err)
		}
	})

	t.Run("generates output name from input when empty", func(t *testing.T) {
		setTempHomeDir(t)
		InitConfigFile()

		cfg := &Config{
			InputFile:  "/path/to/audio.mp3",
			OutputName: "",
			OutputDir:  "",
		}
		AppendHistory(cfg)

		data, _ := os.ReadFile(HistoryFile())
		if !strings.Contains(string(data), "audio_no_silence.mp3") {
			t.Errorf("history should contain generated output name")
		}
	})
}

func TestFindPreset(t *testing.T) {
	presets := []PresetEntry{
		{Name: "aggressive", MinSilenceLen: 500, SilenceThresh: -24},
		{Name: "gentle", MinSilenceLen: 2000, SilenceThresh: -10},
		{Name: "speech", MinSilenceLen: 800, SilenceThresh: -20},
	}

	t.Run("finds by exact name", func(t *testing.T) {
		p := FindPreset(presets, "aggressive")
		if p == nil {
			t.Fatal("expected to find aggressive preset")
		}
		if p.MinSilenceLen != 500 {
			t.Errorf("expected 500, got %d", p.MinSilenceLen)
		}
	})

	t.Run("case insensitive match", func(t *testing.T) {
		p := FindPreset(presets, "AGGRESSIVE")
		if p == nil {
			t.Fatal("expected case-insensitive match")
		}
		if p.Name != "aggressive" {
			t.Errorf("expected 'aggressive', got %s", p.Name)
		}
	})

	t.Run("returns nil for unknown preset", func(t *testing.T) {
		p := FindPreset(presets, "nonexistent")
		if p != nil {
			t.Errorf("expected nil, got %+v", p)
		}
	})

	t.Run("empty presets returns nil", func(t *testing.T) {
		p := FindPreset(nil, "aggressive")
		if p != nil {
			t.Errorf("expected nil for empty presets")
		}
	})
}
