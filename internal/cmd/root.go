package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"neocut/internal/config"
	"neocut/internal/core"
	"neocut/internal/tui"
)

var (
	cfg     config.Config
	tuiMode bool
)

var rootCmd = &cobra.Command{
	Use:   "neocut",
	Short: "Remove silence from audio files",
	Long: `neocut removes silence from MP3 audio files.

It detects silent portions, splits them out, and recombines the
non-silent segments into a new, tighter audio file.

Output is saved to ~/Downloads/neocut/ by default.
Project config is stored in ~/.config/neostore/neocut/`,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE:          run,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Version = config.ReadVersion()

	rootCmd.Flags().StringVarP(&cfg.InputFile, "input", "i", "", "Input MP3 file")
	rootCmd.Flags().StringVarP(&cfg.OutputName, "output", "o", "", "Output filename (saved to ~/Downloads/neocut/)")
	rootCmd.Flags().StringVar(&cfg.ConfigFile, "cnf", "", "Path to config file")
	rootCmd.Flags().StringVarP(&cfg.ConfigFile, "config", "c", "", "Path to config file")
	rootCmd.Flags().BoolVarP(&tuiMode, "tui", "t", false, "Use interactive TUI mode")
	rootCmd.Flags().IntVarP(&cfg.MinSilenceLen, "min-silence-len", "m", 1000, "Minimum silence length in ms")
	rootCmd.Flags().Float64VarP(&cfg.SilenceThresh, "silence-thresh", "s", -16.0, "Silence threshold in dBFS")
	rootCmd.Flags().IntVarP(&cfg.KeepSilence, "keep-silence", "k", 100, "Silence to keep at boundaries in ms")
	rootCmd.Flags().IntVarP(&cfg.SeekStep, "seek-step", "e", 1, "Seek step in ms")
}

func run(cmd *cobra.Command, args []string) error {
	version := config.ReadVersion()
	config.PrintBanner(version, config.Commit)
	config.EnsureConfigDir()

	if tuiMode {
		tuiCfg, err := tui.RunConfigForm()
		if err != nil {
			return err
		}
		return core.Process(tuiCfg)
	}

	if cfg.InputFile == "" {
		cmd.Usage()
		return fmt.Errorf("--input / -i is required (or use --tui for interactive mode)")
	}

	return core.Process(&cfg)
}
