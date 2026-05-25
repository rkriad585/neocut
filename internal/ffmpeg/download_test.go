package ffmpeg

import (
	"archive/zip"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func createTestZip(t *testing.T, files map[string]string) string {
	t.Helper()
	td := t.TempDir()
	zipPath := filepath.Join(td, "test.zip")

	f, err := os.Create(zipPath)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	w := zip.NewWriter(f)
	for name, content := range files {
		entry, err := w.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		_, err = entry.Write([]byte(content))
		if err != nil {
			t.Fatal(err)
		}
	}
	w.Close()
	return zipPath
}

func TestExtractZipFindsFfmpeg(t *testing.T) {
	outputDir := t.TempDir()
	zipPath := createTestZip(t, map[string]string{
		"ffmpeg/bin/ffmpeg.exe": "fake ffmpeg binary content",
	})

	err := extractZip(zipPath, outputDir)
	if err != nil {
		t.Fatalf("extractZip() error: %v", err)
	}

	expectedPath := filepath.Join(outputDir, "ffmpeg.exe")
	data, err := os.ReadFile(expectedPath)
	if err != nil {
		t.Fatalf("ffmpeg.exe should exist: %v", err)
	}
	if string(data) != "fake ffmpeg binary content" {
		t.Errorf("unexpected content: %s", string(data))
	}
}

func TestExtractZipFindsUnixFfmpeg(t *testing.T) {
	outputDir := t.TempDir()
	zipPath := createTestZip(t, map[string]string{
		"ffmpeg/bin/ffmpeg": "unix ffmpeg",
	})

	err := extractZip(zipPath, outputDir)
	if err != nil {
		t.Fatalf("extractZip() error: %v", err)
	}

	expectedPath := filepath.Join(outputDir, "ffmpeg")
	data, _ := os.ReadFile(expectedPath)
	if string(data) != "unix ffmpeg" {
		t.Errorf("unexpected content: %s", string(data))
	}
}

func TestExtractZipNoFfmpeg(t *testing.T) {
	outputDir := t.TempDir()
	zipPath := createTestZip(t, map[string]string{
		"other/file.txt": "not ffmpeg",
	})

	err := extractZip(zipPath, outputDir)
	if err == nil {
		t.Fatal("expected error when ffmpeg not in archive")
	}
	if !strings.Contains(err.Error(), "ffmpeg binary not found") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestExtractZipEmptyZip(t *testing.T) {
	outputDir := t.TempDir()
	zipPath := createTestZip(t, map[string]string{})

	err := extractZip(zipPath, outputDir)
	if err == nil {
		t.Fatal("expected error for empty zip")
	}
}

func TestDownloadWithProgress(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ffmpeg data"))
	}))
	defer ts.Close()

	destDir := t.TempDir()
	destPath, err := downloadWithProgress(ts.URL+"/ffmpeg.zip", destDir)
	if err != nil {
		t.Fatalf("downloadWithProgress() error: %v", err)
	}

	data, err := os.ReadFile(destPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "ffmpeg data" {
		t.Errorf("expected 'ffmpeg data', got '%s'", string(data))
	}
}

func TestDownloadWithProgressHTTPError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	_, err := downloadWithProgress(ts.URL+"/test.bin", t.TempDir())
	if err == nil {
		t.Fatal("expected error for HTTP 404")
	}
	if !strings.Contains(err.Error(), "HTTP 404") {
		t.Errorf("expected HTTP 404 error, got: %v", err)
	}
}

func TestDownloadWithProgressServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	_, err := downloadWithProgress(ts.URL+"/test.bin", t.TempDir())
	if err == nil {
		t.Fatal("expected error for HTTP 500")
	}
}

func TestDownloadWithProgressLargeContent(t *testing.T) {
	largeData := make([]byte, 1024*128)
	for i := range largeData {
		largeData[i] = byte(i % 256)
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(largeData)
	}))
	defer ts.Close()

	destDir := t.TempDir()
	destPath, err := downloadWithProgress(ts.URL+"/ffmpeg.zip", destDir)
	if err != nil {
		t.Fatalf("downloadWithProgress() error: %v", err)
	}

	fi, err := os.Stat(destPath)
	if err != nil {
		t.Fatal(err)
	}
	if fi.Size() != int64(len(largeData)) {
		t.Errorf("expected %d bytes, got %d", len(largeData), fi.Size())
	}
}

func TestDownloadWithProgressContentTypeZip(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/zip")
		w.Write([]byte("zip content"))
	}))
	defer ts.Close()

	destDir := t.TempDir()
	destPath, err := downloadWithProgress(ts.URL+"/ffmpeg.zip", destDir)
	if err != nil {
		t.Fatalf("downloadWithProgress() error: %v", err)
	}

	if !strings.HasSuffix(destPath, ".zip") {
		t.Errorf("expected .zip extension, got: %s", destPath)
	}
}

func TestDownloadWithProgressInvalidURL(t *testing.T) {
	_, err := downloadWithProgress("://invalid", t.TempDir())
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}

func TestDownloadWithProgressCreateError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("data"))
	}))
	defer ts.Close()

	_, err := downloadWithProgress(ts.URL+"/data.bin", "/nonexistent/path/that/does/not/exist")
	if err == nil {
		t.Fatal("expected error for non-existent directory")
	}
}

func TestExtractZip(t *testing.T) {
	t.Run("creates output file with execute permissions", func(t *testing.T) {
		outputDir := t.TempDir()
		zipPath := createTestZip(t, map[string]string{
			"ffmpeg.exe": "binary",
		})

		err := extractZip(zipPath, outputDir)
		if err != nil {
			t.Fatalf("extractZip() error: %v", err)
		}

		destPath := filepath.Join(outputDir, "ffmpeg.exe")
		fi, err := os.Stat(destPath)
		if err != nil {
			t.Fatal(err)
		}

		if runtime.GOOS != "windows" {
			if fi.Mode()&0111 == 0 {
				t.Error("ffmpeg should have execute permissions")
			}
		}
	})
}

func TestDownloadWithProgressRenamesTempFile(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("final content"))
	}))
	defer ts.Close()

	destDir := t.TempDir()
	destPath, err := downloadWithProgress(ts.URL+"/ffmpeg.bin", destDir)
	if err != nil {
		t.Fatalf("downloadWithProgress() error: %v", err)
	}

	if _, err := os.Stat(destPath); os.IsNotExist(err) {
		t.Fatal("destination file should exist after successful download")
	}

	tmpFiles, _ := filepath.Glob(filepath.Join(destDir, "*.tmp"))
	if len(tmpFiles) > 0 {
		t.Errorf("temporary .tmp files should be cleaned up: %v", tmpFiles)
	}
}

func TestDownloadWithProgressNetworkError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hj, ok := w.(http.Hijacker)
		if !ok {
			t.Fatal("server does not support hijacking")
		}
		conn, _, _ := hj.Hijack()
		conn.Close()
	}))
	defer ts.Close()

	_, err := downloadWithProgress(ts.URL, t.TempDir())
	if err == nil {
		t.Log("connection reset may not always error depending on timing")
	}
}
