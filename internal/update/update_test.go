package update

import (
	"fmt"
	"net/http"
	"net/http/httptest"
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

func TestReplaceUnix(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping unix-specific test on windows")
	}

	td := t.TempDir()
	exePath := filepath.Join(td, "neocut")
	oldContent := []byte("old binary")
	newContent := []byte("new binary")

	if err := os.WriteFile(exePath, oldContent, 0755); err != nil {
		t.Fatal(err)
	}
	tmpPath := exePath + ".update.tmp"
	if err := os.WriteFile(tmpPath, newContent, 0644); err != nil {
		t.Fatal(err)
	}

	if err := replaceUnix(exePath, tmpPath); err != nil {
		t.Fatalf("replaceUnix() error: %v", err)
	}

	data, err := os.ReadFile(exePath)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "new binary" {
		t.Errorf("expected new binary content, got: %s", string(data))
	}

	if _, err := os.Stat(tmpPath); !os.IsNotExist(err) {
		t.Errorf("temp file should be gone, but stat returned: %v", err)
	}
}

func TestDownloadWithProgress(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	}))
	defer ts.Close()

	td := t.TempDir()
	dest := filepath.Join(td, "output.tmp")

	err := downloadWithProgress(ts.URL, dest)
	if err != nil {
		t.Fatalf("downloadWithProgress() error: %v", err)
	}

	data, err := os.ReadFile(dest)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "hello world" {
		t.Errorf("expected 'hello world', got '%s'", string(data))
	}
}

func TestDownloadWithProgressHTTPError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	td := t.TempDir()
	dest := filepath.Join(td, "output.tmp")

	err := downloadWithProgress(ts.URL, dest)
	if err == nil {
		t.Fatal("expected error for HTTP 404")
	}
	if !strings.Contains(err.Error(), "HTTP 404") {
		t.Errorf("expected HTTP 404 error, got: %v", err)
	}
}

func TestDownloadWithProgressLargeFile(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := make([]byte, 1024*64)
		for i := range data {
			data[i] = byte(i % 256)
		}
		w.Write(data)
	}))
	defer ts.Close()

	td := t.TempDir()
	dest := filepath.Join(td, "output.tmp")

	err := downloadWithProgress(ts.URL, dest)
	if err != nil {
		t.Fatalf("downloadWithProgress() error: %v", err)
	}

	fi, err := os.Stat(dest)
	if err != nil {
		t.Fatal(err)
	}
	if fi.Size() != 1024*64 {
		t.Errorf("expected 64KB file, got %d bytes", fi.Size())
	}
}

func TestReplaceBinaryDispatchesCorrectly(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("replaceBinary on windows uses bat script, test separately")
	}

	td := t.TempDir()
	exePath := filepath.Join(td, "neocut")
	os.WriteFile(exePath, []byte("old"), 0755)
	tmpPath := exePath + ".update.tmp"
	os.WriteFile(tmpPath, []byte("new"), 0644)

	if err := replaceBinary(exePath, tmpPath); err != nil {
		t.Fatalf("replaceBinary() error: %v", err)
	}

	data, _ := os.ReadFile(exePath)
	if string(data) != "new" {
		t.Errorf("expected new binary, got: %s", string(data))
	}
}
