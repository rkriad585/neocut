package theme

import (
	"testing"
)

func TestFind(t *testing.T) {
	t.Run("finds theme by name", func(t *testing.T) {
		th, ok := Find("dark")
		if !ok {
			t.Fatal("expected to find 'dark' theme")
		}
		if th.Label != "Dark Theme" {
			t.Errorf("expected 'Dark Theme', got %s", th.Label)
		}
	})

	t.Run("case insensitive", func(t *testing.T) {
		th, ok := Find("DARK")
		if !ok {
			t.Fatal("expected case-insensitive match")
		}
		if th.Name != "dark" {
			t.Errorf("expected 'dark', got %s", th.Name)
		}
	})

	t.Run("returns false for unknown", func(t *testing.T) {
		_, ok := Find("nonexistent")
		if ok {
			t.Error("expected false for unknown theme")
		}
	})
}

func TestNames(t *testing.T) {
	names := Names()
	if len(names) != 13 {
		t.Errorf("expected 13 themes, got %d", len(names))
	}
	if names[0] != "dark" {
		t.Errorf("first theme should be 'dark', got %s", names[0])
	}
}

func TestLabels(t *testing.T) {
	labels := Labels()
	if len(labels) != 13 {
		t.Errorf("expected 13 labels, got %d", len(labels))
	}
}

func TestDefaultTheme(t *testing.T) {
	dt := DefaultTheme()
	if dt.Name != "sunny_beach_day" {
		t.Errorf("expected 'sunny_beach_day', got %s", dt.Name)
	}
}

func TestHex(t *testing.T) {
	th := Theme{Name: "test", Colors: []string{"#ff0000", "#00ff00", "#0000ff"}}

	t.Run("returns color for valid role", func(t *testing.T) {
		if h := th.Hex(RolePrimary); h != "#ff0000" {
			t.Errorf("expected #ff0000, got %s", h)
		}
	})

	t.Run("wraps around with modulo", func(t *testing.T) {
		if h := th.Hex(RoleAccent); h != "#00ff00" {
			t.Errorf("expected #00ff00 (4%%3=1), got %s", h)
		}
	})
}

func TestResolve(t *testing.T) {
	t.Run("respects theme name in auto mode", func(t *testing.T) {
		th := Resolve("dark", "auto")
		if th.Name != "dark" {
			t.Errorf("expected 'dark', got %s", th.Name)
		}
	})

	t.Run("forces dark theme in dark mode", func(t *testing.T) {
		th := Resolve("light", "dark")
		if th.Name != "dark" {
			t.Errorf("expected 'dark' when color_mode=dark, got %s", th.Name)
		}
	})

	t.Run("forces light theme in light mode", func(t *testing.T) {
		th := Resolve("dark", "light")
		if th.Name != "light" {
			t.Errorf("expected 'light' when color_mode=light, got %s", th.Name)
		}
	})

	t.Run("falls back to default for unknown name", func(t *testing.T) {
		th := Resolve("nonexistent", "auto")
		if th.Name != "sunny_beach_day" {
			t.Errorf("expected default, got %s", th.Name)
		}
	})
}

func TestResolveColors(t *testing.T) {
	rc := ResolveColors("dark", "auto")
	if rc.Primary == "" {
		t.Error("Primary should not be empty")
	}
	if rc.Success == "" {
		t.Error("Success should not be empty")
	}
	if rc.Accent == "" {
		t.Error("Accent should not be empty")
	}
}

func TestIsDark(t *testing.T) {
	thDark := Theme{Name: "dark", Colors: []string{"#000000"}}
	if !thDark.IsDark() {
		t.Error("black should be dark")
	}

	thLight := Theme{Name: "light", Colors: []string{"#ffffff"}}
	if thLight.IsDark() {
		t.Error("white should not be dark")
	}
}

func TestSetActive(t *testing.T) {
	SetActive("dark", "auto")
	a := Active()
	if a.Name != "dark" {
		t.Errorf("expected 'dark', got %s", a.Name)
	}

	if RoleColor(RolePrimary) == "" {
		t.Error("RoleColor should return non-empty string")
	}

	if h := RoleColor(RolePrimary); h != "#0f172a" {
		t.Errorf("expected '#0f172a', got %s", h)
	}
}

func TestSprintf(t *testing.T) {
	result := Sprintf("hello", "#ff0000")
	expected := "\x1b[38;2;255;0;0mhello\x1b[0m"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestSprintfBold(t *testing.T) {
	result := SprintfBold("hello", "#00ff00")
	expected := "\x1b[1;38;2;0;255;0mhello\x1b[0m"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}
