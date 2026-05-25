package core

import (
	"testing"
	"time"
)

func TestFmtDuration(t *testing.T) {
	tests := []struct {
		duration time.Duration
		expected string
	}{
		{0, "0.0s"},
		{100 * time.Millisecond, "0.1s"},
		{500 * time.Millisecond, "0.5s"},
		{1 * time.Second, "1.0s"},
		{1*time.Second + 500*time.Millisecond, "1.5s"},
		{59 * time.Second, "59.0s"},
		{60 * time.Second, "1m 0s"},
		{61*time.Second + 500*time.Millisecond, "1m 1s"},
		{120 * time.Second, "2m 0s"},
		{59*time.Minute + 59*time.Second, "59m 59s"},
		{60 * time.Minute, "1h 0m 0s"},
		{1*time.Hour + 30*time.Minute + 15*time.Second, "1h 30m 15s"},
		{25*time.Hour + 10*time.Minute + 5*time.Second, "25h 10m 5s"},
		{3661 * time.Second, "1h 1m 1s"},
	}

	for _, tt := range tests {
		t.Run(tt.duration.String(), func(t *testing.T) {
			got := fmtDuration(tt.duration)
			if got != tt.expected {
				t.Errorf("fmtDuration(%v) = %s, want %s", tt.duration, got, tt.expected)
			}
		})
	}
}

func TestFmtDurationRounding(t *testing.T) {
	tests := []struct {
		duration time.Duration
		expected string
	}{
		{1500 * time.Millisecond, "1.5s"},
		{1550 * time.Millisecond, "1.6s"},
		{1549 * time.Millisecond, "1.5s"},
		{100 * time.Microsecond, "0.0s"},
		{999 * time.Millisecond, "1.0s"},
	}

	for _, tt := range tests {
		t.Run(tt.duration.String(), func(t *testing.T) {
			got := fmtDuration(tt.duration)
			if got != tt.expected {
				t.Errorf("fmtDuration(%v) = %s, want %s", tt.duration, got, tt.expected)
			}
		})
	}
}

func TestSetQuietMode(t *testing.T) {
	SetQuietMode(true)
	if !isQuiet() {
		t.Error("expected quiet mode to be true")
	}

	SetQuietMode(false)
	if isQuiet() {
		t.Error("expected quiet mode to be false")
	}
}

func TestIsQuietThreadSafe(t *testing.T) {
	done := make(chan struct{})
	go func() {
		SetQuietMode(true)
		close(done)
	}()

	SetQuietMode(false)

	<-done
}

func TestFmtDurationEdgeCases(t *testing.T) {
	t.Run("negative duration", func(t *testing.T) {
		got := fmtDuration(-5 * time.Second)
		if got != "-5.0s" && got != "0.0s" {
			t.Logf("negative duration formats as: %s", got)
		}
	})

	t.Run("very large duration", func(t *testing.T) {
		d := 100*time.Hour + 30*time.Minute
		got := fmtDuration(d)
		if got != "100h 30m 0s" {
			t.Errorf("expected 100h 30m 0s, got %s", got)
		}
	})
}
