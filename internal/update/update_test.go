package update

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestDownloadURL(t *testing.T) {
	version := "v1.2.3"

	url := DownloadURL(version)

	expectedPrefix := fmt.Sprintf("https://github.com/%s/%s/releases/download/%s/%s-%s-%s",
		gitHubUser, project, version, project, runtime.GOOS, runtime.GOARCH)
	if runtime.GOOS == "windows" {
		expectedPrefix += ".exe"
	}

	if url != expectedPrefix {
		t.Errorf("DownloadURL() = %s, want %s", url, expectedPrefix)
	}

	if !strings.HasPrefix(url, "https://github.com/") {
		t.Errorf("URL should start with https://github.com/")
	}
	if !strings.Contains(url, "/releases/download/") {
		t.Errorf("URL should contain /releases/download/")
	}
	if !strings.HasSuffix(url, version+"/"+project+"-"+runtime.GOOS+"-"+runtime.GOARCH) &&
		!strings.HasSuffix(url, version+"/"+project+"-"+runtime.GOOS+"-"+runtime.GOARCH+".exe") {
		t.Errorf("URL has unexpected suffix: %s", url)
	}
}

func TestDownloadURLVersionFormat(t *testing.T) {
	tests := []struct {
		version string
	}{
		{"v0.2.1"},
		{"v1.0.0"},
		{"v10.20.30"},
	}
	for _, tt := range tests {
		t.Run(tt.version, func(t *testing.T) {
			url := DownloadURL(tt.version)
			if !strings.Contains(url, tt.version) {
				t.Errorf("URL should contain version %s: %s", tt.version, url)
			}
		})
	}
}

func TestFilepathEval(t *testing.T) {
	t.Run("regular file returns itself", func(t *testing.T) {
		td := t.TempDir()
		f := filepath.Join(td, "test.txt")
		if err := os.WriteFile(f, []byte("hello"), 0644); err != nil {
			t.Fatal(err)
		}

		got, err := filepathEval(f)
		if err != nil {
			t.Fatalf("filepathEval() error: %v", err)
		}
		if got != f {
			t.Errorf("expected %s, got %s", f, got)
		}
	})

	t.Run("non-existent file returns error", func(t *testing.T) {
		_, err := filepathEval("/nonexistent/path")
		if err == nil {
			t.Error("expected error for non-existent file")
		}
	})

	t.Run("empty string returns error", func(t *testing.T) {
		_, err := filepathEval("")
		if err == nil {
			t.Error("expected error for empty path")
		}
	})
}

func TestDownloadURLInvalidVersion(t *testing.T) {
	url := DownloadURL("")
	if url == "" {
		t.Error("DownloadURL should not return empty even for empty version")
	}
}

func TestDownloadURLNoTrailingSlash(t *testing.T) {
	url := DownloadURL("v0.2.4")
	if strings.HasSuffix(url, "/") {
		t.Error("DownloadURL should not have trailing slash")
	}
}


