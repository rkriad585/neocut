package ffmpeg

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestBinDir(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}
	got := BinDir()
	want := filepath.Join(home, ".config", "neostore", "neocut", "bin")
	if got != want {
		t.Errorf("BinDir() = %s, want %s", got, want)
	}
}

func TestPathContains(t *testing.T) {
	sep := string(filepath.ListSeparator)

	tests := []struct {
		name  string
		path  string
		dir   string
		want  bool
	}{
		{"exact match", "a" + sep + "b", "a", true},
		{"not present", "a" + sep + "b", "c", false},
		{"empty path", "", "a", false},
		{"empty dir", "a" + sep + "b", "", false},
		{"last entry", "a" + sep + "c", "c", true},
		{"middle entry", "a" + sep + "b" + sep + "c", "b", true},
		{"single entry", "a", "a", true},
		{"single entry no match", "a", "b", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pathContains(tt.path, tt.dir)
			if got != tt.want {
				t.Errorf("pathContains(%q, %q) = %v, want %v", tt.path, tt.dir, got, tt.want)
			}
		})
	}
}

func TestAddToPATH(t *testing.T) {
	original := os.Getenv("PATH")
	defer os.Setenv("PATH", original)

	td := t.TempDir()
	newDir := filepath.Join(td, "newbin")

	os.Setenv("PATH", td)
	addToPATH(newDir)

	current := os.Getenv("PATH")
	if !strings.Contains(current, newDir) {
		t.Errorf("PATH should contain %s: %s", newDir, current)
	}
}

func TestAddToPATHMultiple(t *testing.T) {
	original := os.Getenv("PATH")
	defer os.Setenv("PATH", original)

	td := t.TempDir()
	dir1 := filepath.Join(td, "bin1")
	dir2 := filepath.Join(td, "bin2")

	os.Setenv("PATH", "")
	addToPATH(dir1, dir2)

	current := os.Getenv("PATH")
	if !strings.Contains(current, dir1) {
		t.Errorf("PATH should contain %s: %s", dir1, current)
	}
	if !strings.Contains(current, dir2) {
		t.Errorf("PATH should contain %s: %s", dir2, current)
	}
}

func TestAddToPATHDeduplicates(t *testing.T) {
	original := os.Getenv("PATH")
	defer os.Setenv("PATH", original)

	td := t.TempDir()

	os.Setenv("PATH", td)
	addToPATH(td)

	current := os.Getenv("PATH")
	entries := strings.Split(current, string(os.PathListSeparator))
	count := 0
	for _, e := range entries {
		if e == td {
			count++
		}
	}
	if count != 1 {
		t.Errorf("expected 1 occurrence of %s in PATH, got %d", td, count)
	}
}

func TestAddToPATHEmptyDir(t *testing.T) {
	original := os.Getenv("PATH")
	defer os.Setenv("PATH", original)

	os.Setenv("PATH", "/usr/bin")
	addToPATH("")

	if os.Getenv("PATH") != "/usr/bin" {
		t.Errorf("PATH should remain unchanged")
	}
}

func TestDownloadURL(t *testing.T) {
	url := downloadURL()

	if runtime.GOOS == "windows" {
		if !strings.Contains(url, "gyan.dev") {
			t.Errorf("Windows URL should contain gyan.dev, got: %s", url)
		}
		if !strings.HasSuffix(url, ".zip") {
			t.Errorf("Windows URL should end with .zip, got: %s", url)
		}
	} else if runtime.GOOS == "linux" {
		if !strings.Contains(url, "johnvansickle.com") {
			t.Errorf("Linux URL should contain johnvansickle.com, got: %s", url)
		}
		if !strings.HasSuffix(url, ".tar.xz") {
			t.Errorf("Linux URL should end with .tar.xz, got: %s", url)
		}
		if runtime.GOARCH == "arm64" {
			if !strings.Contains(url, "aarch64") {
				t.Errorf("Linux arm64 URL should contain aarch64, got: %s", url)
			}
		}
	} else if runtime.GOOS == "darwin" {
		if !strings.Contains(url, "evermeet.cx") {
			t.Errorf("macOS URL should contain evermeet.cx, got: %s", url)
		}
	} else {
		if url != "" {
			t.Errorf("unknown OS %s should return empty URL, got: %s", runtime.GOOS, url)
		}
	}
}

func TestEnsureWhichShim(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("which shim is Windows-only")
	}

	td := t.TempDir()

	ensureWhichShim(td)

	whichPath := filepath.Join(td, "which.cmd")
	if _, err := os.Stat(whichPath); os.IsNotExist(err) {
		t.Errorf("which.cmd should exist at %s", whichPath)
	}

	data, err := os.ReadFile(whichPath)
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)
	if !strings.Contains(content, "where %1") {
		t.Errorf("which.cmd should contain 'where %%1', got: %s", content)
	}
}

func TestEnsureWhichShimIdempotent(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("which shim is Windows-only")
	}

	td := t.TempDir()
	whichPath := filepath.Join(td, "which.cmd")
	os.WriteFile(whichPath, []byte("existing"), 0644)

	ensureWhichShim(td)

	data, _ := os.ReadFile(whichPath)
	if string(data) != "existing" {
		t.Errorf("existing which.cmd should not be overwritten")
	}
}

func TestEnsureWhichShimEmptyDir(t *testing.T) {
	ensureWhichShim("")
}

func TestExtractFfmpegUnsupportedFormat(t *testing.T) {
	err := extractFfmpeg("/tmp/test.tar.gz", t.TempDir())
	if err == nil {
		t.Error("expected error for unsupported archive format")
	}
	if !strings.Contains(err.Error(), "unsupported archive format") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestPathContainsEmptyDir(t *testing.T) {
	if pathContains("/a:/b", "") {
		t.Error("empty dir should not match")
	}
}

func TestPathContainsCaseSensitivity(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("case sensitivity test is platform-dependent")
	}

	result := pathContains("C:\\Users\\test\\bin", "C:\\Users\\test\\BIN")
	if !result {
		t.Log("case-insensitive filesystem detected, PATH comparison may vary")
	}
}
