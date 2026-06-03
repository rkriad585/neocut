package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadVersion(t *testing.T) {
	t.Run("ldflags version takes precedence", func(t *testing.T) {
		orig := Version
		Version = "v9.9.9"
		defer func() { Version = orig }()

		v := ReadVersion()
		if v != "v9.9.9" {
			t.Errorf("expected v9.9.9, got %s", v)
		}
	})

	t.Run("embedded version fallback", func(t *testing.T) {
		orig := Version
		origEmbedded := embeddedVersion
		Version = ""
		embeddedVersion = "v2.0.0\n"
		defer func() {
			Version = orig
			embeddedVersion = origEmbedded
		}()

		v := ReadVersion()
		if v != "v2.0.0" {
			t.Errorf("expected v2.0.0, got %s", v)
		}
	})

	t.Run(".version file fallback", func(t *testing.T) {
		orig := Version
		origEmbedded := embeddedVersion
		Version = ""
		embeddedVersion = ""
		defer func() {
			Version = orig
			embeddedVersion = origEmbedded
		}()

		td := t.TempDir()
		versionPath := filepath.Join(td, ".version")
		if err := os.WriteFile(versionPath, []byte("v0.0.1\n"), 0644); err != nil {
			t.Fatal(err)
		}

		origWd, _ := os.Getwd()
		os.Chdir(td)
		defer os.Chdir(origWd)

		v := ReadVersion()
		if v != "v0.0.1" {
			t.Errorf("expected v0.0.1, got %s", v)
		}
	})

	t.Run("unknown when nothing available", func(t *testing.T) {
		orig := Version
		origEmbedded := embeddedVersion
		Version = ""
		embeddedVersion = ""
		defer func() {
			Version = orig
			embeddedVersion = origEmbedded
		}()

		td := t.TempDir()
		origWd, _ := os.Getwd()
		os.Chdir(td)
		defer os.Chdir(origWd)

		v := ReadVersion()
		if v != "unknown" {
			t.Errorf("expected unknown, got %s", v)
		}
	})
}

func TestPrintBanner(t *testing.T) {
	prev := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	PrintBanner("v1.0.0", "abc123")

	w.Close()
	os.Stdout = prev

	var buf strings.Builder
	readBuf := make([]byte, 1024)
	n, _ := r.Read(readBuf)
	buf.Write(readBuf[:n])

	output := buf.String()
	if !strings.Contains(output, "neocut") {
		t.Errorf("banner should contain 'neocut'")
	}
	if !strings.Contains(output, "v1.0.0") {
		t.Errorf("banner should contain version")
	}
	if !strings.Contains(output, "abc123") {
		t.Errorf("banner should contain commit")
	}
}

func TestGetOutputDir(t *testing.T) {
	t.Run("uses config OutputDir when set", func(t *testing.T) {
		cfg := &Config{OutputDir: "/custom/path"}
		got := GetOutputDir(cfg)
		if got != "/custom/path" {
			t.Errorf("expected /custom/path, got %s", got)
		}
	})

	t.Run("nil config returns default", func(t *testing.T) {
		home, err := os.UserHomeDir()
		if err != nil {
			t.Fatal(err)
		}
		got := GetOutputDir(nil)
		want := filepath.Join(home, "Downloads", "neostore", "neocut")
		if got != want {
			t.Errorf("expected %s, got %s", want, got)
		}
	})

	t.Run("empty OutputDir returns default", func(t *testing.T) {
		home, err := os.UserHomeDir()
		if err != nil {
			t.Fatal(err)
		}
		cfg := &Config{OutputDir: ""}
		got := GetOutputDir(cfg)
		want := filepath.Join(home, "Downloads", "neostore", "neocut")
		if got != want {
			t.Errorf("expected %s, got %s", want, got)
		}
	})
}

func TestConfigDir(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}
	got := ConfigDir()
	want := filepath.Join(home, ".config", "neostore", "neocut")
	if got != want {
		t.Errorf("expected %s, got %s", want, got)
	}
}

func TestConfigFile(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}

	t.Run("toml config path", func(t *testing.T) {
		got := ConfigFile("config.toml")
		want := filepath.Join(home, ".config", "neostore", "neocut", "config.toml")
		if got != want {
			t.Errorf("expected %s, got %s", want, got)
		}
	})

	t.Run("history log path", func(t *testing.T) {
		got := HistoryFile()
		want := filepath.Join(home, ".config", "neostore", "neocut", "history.log")
		if got != want {
			t.Errorf("expected %s, got %s", want, got)
		}
	})
}

func TestEnsureConfigDir(t *testing.T) {
	origHome := os.Getenv("USERPROFILE")
	if origHome == "" {
		origHome = os.Getenv("HOME")
	}

	td := t.TempDir()
	if _, err := os.Stat(td); os.IsNotExist(err) {
		t.Fatal("temp dir should exist")
	}

	os.Setenv("USERPROFILE", td)
	os.Setenv("HOME", td)
	defer func() {
		os.Setenv("USERPROFILE", origHome)
		os.Setenv("HOME", origHome)
	}()

	EnsureConfigDir()

	expectedDir := filepath.Join(td, ".config", "neostore", "neocut")
	if _, err := os.Stat(expectedDir); os.IsNotExist(err) {
		t.Errorf("config dir should exist at %s", expectedDir)
	}
}
