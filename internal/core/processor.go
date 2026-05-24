package core

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/Vernacular-ai/godub"

	"neocut/internal/config"
)

func step(label string, fn func() error) error {
	done := make(chan error, 1)
	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				done <- fmt.Errorf("panic: %v", r)
			}
		}()
		done <- fn()
	}()

	i := 0
	for {
		select {
		case err := <-done:
			if err != nil {
				fmt.Printf("\r  \u2717 %s\n", label)
				fmt.Printf("    Error: %v\n", err)
			} else {
				fmt.Printf("\r  \u2713 %s\n", label)
			}
			return err
		default:
			fmt.Printf("\r  %s %s", frames[i%len(frames)], label)
			i++
			time.Sleep(80 * time.Millisecond)
		}
	}
}

func exportWithProgress(exporter *godub.Exporter, segment *godub.AudioSegment) error {
	done := make(chan struct{})
	blocks := []string{"▏", "▎", "▍", "▌", "▋", "▊", "▉", "█"}
	width := 30

	go func() {
		i := 0
		for {
			select {
			case <-done:
				bar := strings.Repeat("█", width)
				fmt.Printf("\r  \u2713 Exporting audio [%s] 100%%", bar)
				return
			default:
				pos := i % (width * 8)
				full := pos / 8
				part := pos % 8
				bar := strings.Repeat("█", full)
				if part > 0 {
					bar += blocks[part]
				}
				bar += strings.Repeat(" ", width-full-1)
				fmt.Printf("\r    Exporting [%s]", bar)
				i++
				time.Sleep(40 * time.Millisecond)
			}
		}
	}()

	err := exporter.Export(segment)
	close(done)
	time.Sleep(60 * time.Millisecond)
	fmt.Println()
	return err
}

func Process(cfg *config.Config) error {
	var segment *godub.AudioSegment

	if err := step("Loading audio", func() error {
		var loadErr error
		segment, loadErr = godub.NewLoader().Load(cfg.InputFile)
		return loadErr
	}); err != nil {
		return err
	}

	var chunks []*godub.AudioSegment
	if err := step("Detecting and removing silence", func() error {
		var splitErr error
		chunks, _, splitErr = godub.SplitOnSilence(
			segment,
			cfg.MinSilenceLen,
			godub.Volume(cfg.SilenceThresh),
			cfg.KeepSilence,
			cfg.SeekStep,
		)
		return splitErr
	}); err != nil {
		return err
	}

	if len(chunks) == 0 {
		return fmt.Errorf("no audio remaining after silence removal")
	}
	fmt.Printf("    Found %d non-silent segment(s)\n", len(chunks))

	var combined *godub.AudioSegment
	if err := step("Recombining audio chunks", func() error {
		if len(chunks) == 1 {
			combined = chunks[0]
			return nil
		}
		var combineErr error
		combined, combineErr = chunks[0].Append(chunks[1:]...)
		return combineErr
	}); err != nil {
		return err
	}

	outputDir := config.GetOutputDir()
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	outputName := cfg.OutputName
	if outputName == "" {
		base := filepath.Base(cfg.InputFile)
		ext := filepath.Ext(base)
		outputName = strings.TrimSuffix(base, ext) + "_no_silence.mp3"
	}
	if !strings.HasSuffix(strings.ToLower(outputName), ".mp3") {
		outputName += ".mp3"
	}
	outputPath := filepath.Join(outputDir, outputName)

	exporter := godub.NewExporter(outputPath).
		WithDstFormat("mp3")

	if err := exportWithProgress(exporter, combined); err != nil {
		return fmt.Errorf("export failed: %w", err)
	}

	fmt.Printf("  \u2713 Exported to %s\n", outputPath)
	fmt.Println()
	fmt.Printf("  Done! (OS: %s/%s)\n", runtime.GOOS, runtime.GOARCH)
	fmt.Println()
	return nil
}
