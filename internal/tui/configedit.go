package tui

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"

	"neocut/internal/config"
	"neocut/internal/theme"
)

func RunConfigEditor() error {
	presets, defaults, _ := config.ReadConfig()

	var (
		minSilenceLen     = 1000
		silenceThreshStr  = "-16.0"
		keepSilence       = 100
		seekStep          = 1
		outputDir         string
		themeName         = "sunny_beach_day"
		colorMode         = "auto"
		historyNote       = "  No history yet."
		save              = false
		previewNote       = ""
	)

	if defaults != nil {
		minSilenceLen = defaults.MinSilenceLen
		silenceThreshStr = strconv.FormatFloat(defaults.SilenceThresh, 'f', -1, 64)
		keepSilence = defaults.KeepSilence
		seekStep = defaults.SeekStep
		outputDir = defaults.OutputDir
		if defaults.Theme != "" {
			themeName = defaults.Theme
		}
		if defaults.ColorMode != "" {
			colorMode = defaults.ColorMode
		}
	}

	var presetLines []string
	if len(presets) == 0 {
		presetLines = append(presetLines, "  No presets configured.")
	}
	for _, p := range presets {
		presetLines = append(presetLines,
			fmt.Sprintf("  %s  m=%d  s=%.1f  k=%d  e=%d",
				p.Name, p.MinSilenceLen, p.SilenceThresh, p.KeepSilence, p.SeekStep))
	}

	historyLines, _ := readHistory()
	if len(historyLines) > 0 {
		historyNote = strings.Join(historyLines, "\n")
	}

	seekOpts := []huh.Option[int]{
		huh.NewOption("1 ms (default, most precise)", 1),
		huh.NewOption("2 ms", 2),
		huh.NewOption("5 ms", 5),
		huh.NewOption("10 ms (fast)", 10),
		huh.NewOption("20 ms (fastest)", 20),
	}

	themeNames := theme.Names()
	var themeOpts = make([]huh.Option[string], 0, len(themeNames))
	for _, n := range themeNames {
		t, _ := theme.Find(n)
		themeOpts = append(themeOpts, huh.NewOption(t.Label, n))
	}

	previewColors := func(name, mode string) string {
		rc := theme.ResolveColors(name, mode)
		return fmt.Sprintf(
			"  %s%s  %s%s  %s%s  %s%s  %s%s",
			theme.Sprintf("██", rc.Primary), rc.Primary,
			theme.Sprintf("██", rc.Success), rc.Success,
			theme.Sprintf("██", rc.Warning), rc.Warning,
			theme.Sprintf("██", rc.Error), rc.Error,
			theme.Sprintf("██", rc.Accent), rc.Accent,
		)
	}

	activeTheme := theme.Resolve(themeName, colorMode)
	previewNote = previewColors(activeTheme.Name, colorMode)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().Title("Default Processing Parameters").
				Description("These values are loaded each time neocut runs.\nCLI flags override them."),
			huh.NewSelect[int]().
				Title("Min Silence Length").
				Description("Milliseconds — shorter catches more gaps").
				Options(
					huh.NewOption("300 ms", 300),
					huh.NewOption("500 ms (aggressive)", 500),
					huh.NewOption("800 ms (speech)", 800),
					huh.NewOption("1000 ms (default)", 1000),
					huh.NewOption("1500 ms", 1500),
					huh.NewOption("2000 ms (gentle)", 2000),
				).
				Value(&minSilenceLen),
			huh.NewInput().
				Title("Silence Threshold").
				Description("dBFS — lower = more sensitive (e.g. -24 catches more)").
				Placeholder("-16.0").
				Value(&silenceThreshStr),
			huh.NewSelect[int]().
				Title("Keep Silence at Boundaries").
				Description("Milliseconds of silence to retain").
				Options(
					huh.NewOption("0 ms (tight cuts)", 0),
					huh.NewOption("50 ms", 50),
					huh.NewOption("100 ms (default)", 100),
					huh.NewOption("250 ms", 250),
					huh.NewOption("500 ms (soft)", 500),
				).
				Value(&keepSilence),
			huh.NewSelect[int]().
				Title("Seek Step").
				Description("Precision — lower is slower but more accurate").
				Options(seekOpts...).
				Value(&seekStep),
			huh.NewInput().
				Title("Output Directory").
				Description("Leave empty for ~/Downloads/neostore/neocut/").
				Value(&outputDir),
		),
		huh.NewGroup(
			huh.NewNote().Title("Appearance").
				Description("Choose theme and color mode"),
			huh.NewSelect[string]().
				Title("Theme").
				Description("Pick a color theme for the UI").
				Options(themeOpts...).
				Value(&themeName),
			huh.NewSelect[string]().
				Title("Color Mode").
				Description("auto = use theme colors; dark = force dark; light = force light").
				Options(
					huh.NewOption("Auto (use theme)", "auto"),
					huh.NewOption("Dark (force dark theme)", "dark"),
					huh.NewOption("Light (force light theme)", "light"),
				).
				Value(&colorMode),
			huh.NewNote().Title("Theme Preview").
				Description(previewNote),
		),
		huh.NewGroup(
			huh.NewNote().Title("Configured Presets").
				Description(strings.Join(presetLines, "\n")),
			huh.NewNote().Title("Recent History").
				Description(historyNote),
		huh.NewConfirm().
			Title("Save these defaults?").
			Affirmative("Save").
			Negative("Cancel").
			Value(&save),
		),
	).WithTheme(buildHuhTheme(themeName, colorMode))

	if err := form.Run(); err != nil {
		os.Exit(0)
		return nil
	}

	if !save {
		fmt.Println("  Cancelled — no changes saved.")
		return nil
	}

	silenceThresh, err := strconv.ParseFloat(silenceThreshStr, 64)
	if err != nil {
		return fmt.Errorf("invalid silence threshold: %s", silenceThreshStr)
	}

	if err := config.WriteDefaults(config.DefaultEntry{
		MinSilenceLen: minSilenceLen,
		SilenceThresh: silenceThresh,
		KeepSilence:   keepSilence,
		SeekStep:      seekStep,
		OutputDir:     outputDir,
		Theme:         themeName,
		ColorMode:     colorMode,
	}); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	theme.SetActive(themeName, colorMode)
	fmt.Println("  Defaults saved to config.toml")
	return nil
}

func readHistory() ([]string, error) {
	path := config.HistoryFile()
	if path == "" {
		return nil, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil
	}

	allLines := strings.Split(strings.TrimSpace(string(data)), "\n")
	var lines []string
	start := 0
	if len(allLines) > 5 {
		start = len(allLines) - 5
	}
	for i := start; i < len(allLines); i++ {
		line := strings.TrimSpace(allLines[i])
		if line == "" {
			continue
		}
		var h config.HistoryEntry
		if err := json.Unmarshal([]byte(line), &h); err != nil {
			continue
		}
		ts := h.Timestamp
		if len(ts) > 10 {
			ts = ts[:10]
		}
		lines = append(lines, fmt.Sprintf("  %s  %s → %s  (%s)", ts, h.Input, h.Output, h.Dir))
	}
	if len(lines) == 0 {
		lines = append(lines, "  No history yet.")
	}
	return lines, nil
}
