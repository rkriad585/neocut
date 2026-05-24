package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/Vernacular-ai/godub"
	"github.com/charmbracelet/huh"
)

var commit = "unknown"

func readVersion() string {
	data, err := os.ReadFile(".version")
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(data))
}

func printBanner(version, commit string) {
	fmt.Println()
	fmt.Println("╭──────────────── neocut ───────────────────╮")
	fmt.Printf("│      Author : RK Riad Khan                │\n")
	fmt.Printf("│      Version: %-32s│\n", version)
	fmt.Printf("│      Commit : %-32s│\n", commit)
	fmt.Println("│      GitHub : rkriad585/neocut            │")
	fmt.Println("╰───────────────────────────────────────────╯")
	fmt.Println()
}

func getOutputDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return filepath.Join(home, "Downloads", "neocut")
}

func ensureConfigDir() {
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}
	configDir := filepath.Join(home, ".config", "neostore", "neocut")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		fmt.Printf("Warning: could not create config dir: %v\n", err)
	}
}

func runStep(label string, fn func() error) error {
	done := make(chan error, 1)
	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

	go func() {
		done <- fn()
	}()

	i := 0
	for {
		select {
		case err := <-done:
			if err != nil {
				fmt.Printf("\r  \u2717 %s failed: %v\n", label, err)
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
	frames := []string{"▏", "▎", "▍", "▌", "▋", "▊", "▉", "█"}
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
					bar += frames[part]
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

func main() {
	version := readVersion()
	printBanner(version, commit)

	ensureConfigDir()

	var (
		inputFile        string
		outputName       string
		minSilenceLen    = 1000
		silenceThreshStr = "-16.0"
		keepSilence      = 100
	)

	outputDir := getOutputDir()

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Input File").
				Description("Path to the MP3 file to process").
				Placeholder("/path/to/input.mp3").
				Value(&inputFile),
			huh.NewInput().
				Title("Output Filename").
				Description(fmt.Sprintf("Saved to %s", outputDir)).
				Placeholder("output.mp3").
				Value(&outputName),
		),
		huh.NewGroup(
			huh.NewNote().Title("Silence Detection Settings"),
			huh.NewSelect[int]().
				Title("Min Silence Length (ms)").
				Description("Minimum silence duration to split on").
				Options(
					huh.NewOption("500 ms", 500),
					huh.NewOption("1000 ms (default)", 1000),
					huh.NewOption("1500 ms", 1500),
					huh.NewOption("2000 ms", 2000),
				).
				Value(&minSilenceLen),
			huh.NewInput().
				Title("Silence Threshold (dBFS)").
				Description("Volume below which is considered silence").
				Placeholder("-16.0").
				Value(&silenceThreshStr),
			huh.NewSelect[int]().
				Title("Keep Silence (ms)").
				Description("Silence to keep at chunk boundaries").
				Options(
					huh.NewOption("0 ms", 0),
					huh.NewOption("100 ms (default)", 100),
					huh.NewOption("250 ms", 250),
					huh.NewOption("500 ms", 500),
				).
				Value(&keepSilence),
		),
	)

	if err := form.Run(); err != nil {
		fmt.Println("Cancelled.")
		os.Exit(1)
	}

	if inputFile == "" {
		fmt.Println("Input file is required.")
		os.Exit(1)
	}

	silenceThresh, err := strconv.ParseFloat(silenceThreshStr, 64)
	if err != nil {
		fmt.Println("Invalid silence threshold.")
		os.Exit(1)
	}

	fmt.Println()

	var segment *godub.AudioSegment
	if err := runStep("Loading audio", func() error {
		var loadErr error
		segment, loadErr = godub.NewLoader().Load(inputFile)
		return loadErr
	}); err != nil {
		os.Exit(1)
	}

	var chunks []*godub.AudioSegment
	if err := runStep("Detecting and removing silence", func() error {
		var splitErr error
		chunks, _, splitErr = godub.SplitOnSilence(
			segment,
			minSilenceLen,
			godub.Volume(silenceThresh),
			keepSilence,
			1,
		)
		return splitErr
	}); err != nil {
		os.Exit(1)
	}

	if len(chunks) == 0 {
		fmt.Println("  No audio remaining after silence removal.")
		os.Exit(1)
	}
	fmt.Printf("    Found %d non-silent segment(s)\n", len(chunks))

	var combined *godub.AudioSegment
	if err := runStep("Recombining audio chunks", func() error {
		if len(chunks) == 1 {
			combined = chunks[0]
			return nil
		}
		var combineErr error
		combined, combineErr = chunks[0].Append(chunks[1:]...)
		return combineErr
	}); err != nil {
		os.Exit(1)
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Println("Failed to create output directory:", err)
		os.Exit(1)
	}

	if outputName == "" {
		base := filepath.Base(inputFile)
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
		fmt.Println("  Export failed:", err)
		os.Exit(1)
	}

	fmt.Printf("  \u2713 Exported to %s\n", outputPath)
	fmt.Println()
	fmt.Printf("  Done! (OS: %s/%s)\n", runtime.GOOS, runtime.GOARCH)
	fmt.Println()
}
