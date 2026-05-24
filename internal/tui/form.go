package tui

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/huh"

	"neocut/internal/config"
)

func RunConfigForm() (*config.Config, error) {
	var (
		inputFile        string
		outputName       string
		minSilenceLen    = 1000
		silenceThreshStr = "-16.0"
		keepSilence      = 100
	)

	outputDir := config.GetOutputDir()

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
		return nil, fmt.Errorf("cancelled")
	}

	if inputFile == "" {
		return nil, fmt.Errorf("input file is required")
	}

	silenceThresh, err := strconv.ParseFloat(silenceThreshStr, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid silence threshold")
	}

	return &config.Config{
		InputFile:     inputFile,
		OutputName:    outputName,
		MinSilenceLen: minSilenceLen,
		SilenceThresh: silenceThresh,
		KeepSilence:   keepSilence,
		SeekStep:      1,
	}, nil
}
