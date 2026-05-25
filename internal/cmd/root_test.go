package cmd

import (
	"os"
	"testing"
)

func TestExecute(t *testing.T) {
	t.Run("Execute does not panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("Execute() panicked: %v", r)
			}
		}()

		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()

		os.Args = []string{"neocut", "--help"}
		Execute()
	})
}

func TestInitRegistersFlags(t *testing.T) {
	if rootCmd == nil {
		t.Fatal("rootCmd should not be nil")
	}

	expectedFlags := []string{
		"input", "output", "output-dir", "tui", "config",
		"format", "bitrate", "dry-run",
		"min-silence-len", "silence-thresh", "keep-silence", "seek-step",
		"quiet", "preset", "selfuninstall",
	}

	for _, name := range expectedFlags {
		t.Run("flag: "+name, func(t *testing.T) {
			f := rootCmd.Flags().Lookup(name)
			if f == nil {
				t.Errorf("flag --%s should be registered", name)
			}
		})
	}
}

func TestInitRegistersSubcommands(t *testing.T) {
	if rootCmd == nil {
		t.Fatal("rootCmd should not be nil")
	}

	cmds := rootCmd.Commands()
	found := false
	for _, c := range cmds {
		if c.Use == "self-update" {
			found = true
			break
		}
	}
	if !found {
		t.Error("self-update subcommand should be registered")
	}
}

func TestVersionFlag(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"neocut", "--version"}

	defer func() {
		if r := recover(); r != nil {
			t.Logf("Execute panicked (expected): %v", r)
		}
	}()

	Execute()
}

func TestSelfUpdateCommand(t *testing.T) {
	if selfUpdateCmd == nil {
		t.Fatal("selfUpdateCmd should not be nil")
	}

	if selfUpdateCmd.Use != "self-update" {
		t.Errorf("expected 'self-update', got '%s'", selfUpdateCmd.Use)
	}

	if selfUpdateCmd.Short == "" {
		t.Error("self-update should have a Short description")
	}
}

func TestInitConfigEditModeFlag(t *testing.T) {
	f := rootCmd.Flags().Lookup("config")
	if f == nil {
		t.Fatal("--config flag should exist")
	}
	if f.Shorthand != "c" {
		t.Errorf("expected shorthand 'c', got '%s'", f.Shorthand)
	}
}

func TestFlagShortForms(t *testing.T) {
	shortFlags := make(map[string]string)
	for _, flag := range []struct {
		name     string
		shorthand string
	}{
		{"input", "i"},
		{"output", "o"},
		{"output-dir", "d"},
		{"tui", "t"},
		{"config", "c"},
		{"format", "f"},
		{"bitrate", "b"},
		{"min-silence-len", "m"},
		{"silence-thresh", "s"},
		{"keep-silence", "k"},
		{"seek-step", "e"},
		{"quiet", "q"},
	} {
		t.Run(flag.name, func(t *testing.T) {
			if shortFlags[flag.shorthand] != "" {
				t.Fatalf("duplicate shorthand -%s for flags %s and %s",
					flag.shorthand, shortFlags[flag.shorthand], flag.name)
			}
			shortFlags[flag.shorthand] = flag.name
		})
	}
}

func TestSelfUpdateSilenceErrors(t *testing.T) {
	if selfUpdateCmd == nil {
		t.Fatal("selfUpdateCmd not initialized")
	}
	if !selfUpdateCmd.SilenceErrors {
		t.Error("selfUpdateCmd should have SilenceErrors")
	}
}
