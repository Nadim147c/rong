package cmd

import (
	"log/slog"
	"os"
	"path/filepath"
	"slices"

	"github.com/MatusOllah/slogcolor"
	"github.com/Nadim147c/rong/cmd/cache"
	"github.com/Nadim147c/rong/cmd/color"
	"github.com/Nadim147c/rong/cmd/image"
	"github.com/Nadim147c/rong/cmd/video"
	"github.com/Nadim147c/rong/internal/config"
	"github.com/carapace-sh/carapace"
	termcolor "github.com/fatih/color"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
)

func init() {
	Command.AddCommand(color.Command)
	Command.AddCommand(image.Command)
	Command.AddCommand(video.Command)
	Command.AddCommand(cache.Command)

	actions := carapace.ActionMap{
		"variant": carapace.ActionValues(
			"monochrome", "neutral", "tonal_spot",
			"vibrant", "expressive", "fidelity",
			"content", "rainbow", "fruit_salad",
		),
		"version":  carapace.ActionValues("2021", "2025"),
		"platform": carapace.ActionValues("phone", "watch"),
	}

	carapace.Gen(color.Command).FlagCompletion(actions)

	imageCara := carapace.Gen(image.Command)
	imageCara.FlagCompletion(actions)
	imageCara.PositionalAnyCompletion(carapace.ActionFiles())

	videoCara := carapace.Gen(video.Command)
	videoCara.FlagCompletion(actions)
	videoCara.PositionalAnyCompletion(carapace.ActionFiles())

	cacheCara := carapace.Gen(cache.Command)
	cacheCara.FlagCompletion(actions)
	cacheCara.PositionalAnyCompletion(carapace.ActionFiles())

	rootCara := carapace.Gen(Command)
	rootCara.Standalone()
	rootCara.FlagCompletion(carapace.ActionMap{
		"config":   carapace.ActionFiles(),
		"log-file": carapace.ActionFiles(),
	})

	Command.PersistentFlags().BoolP("verbose", "v", false, "enable verbose logging")
	Command.PersistentFlags().BoolP("quiet", "q", false, "suppress all logs")
	Command.PersistentFlags().String("log-file", "", "file to save logs")
	Command.PersistentFlags().StringP("config", "c", "", "path to config (.toml|.yaml|.yml) file")
	Command.MarkFlagsMutuallyExclusive("verbose", "quiet")
}

// Command is root command of the cli
var Command = &cobra.Command{
	Use:          "rong",
	Short:        "A material you color generator from image or video.",
	SilenceUsage: true,
	PersistentPreRun: func(cmd *cobra.Command, _ []string) {
		if cmd.Name() == "_carapace" {
			return
		}

		tty := os.Getenv("TERM") != "dumb" &&
			(isatty.IsTerminal(os.Stderr.Fd()) ||
				isatty.IsCygwinTerminal(os.Stderr.Fd()))
		termcolor.NoColor = !tty

		opts := slogcolor.DefaultOptions
		opts.NoTime = true
		opts.SrcFileMode = 0
		opts.LevelTags = map[slog.Level]string{
			slog.LevelDebug: termcolor.New(termcolor.FgGreen).Sprint("DBG"),
			slog.LevelInfo:  termcolor.New(termcolor.FgCyan).Sprint("INF"),
			slog.LevelWarn:  termcolor.New(termcolor.FgYellow).Sprint("WRN"),
			slog.LevelError: termcolor.New(termcolor.FgRed).Sprint("ERR"),
		}

		verbose, _ := cmd.Flags().GetBool("verbose")
		if verbose {
			opts.Level = slog.LevelDebug
		}

		quiet, _ := cmd.Flags().GetBool("quiet")
		if quiet {
			opts.Level = slog.Level(100)
		}

		logFilePath, err := cmd.Flags().GetString("log-file")
		if err != nil || logFilePath == "" {
			slog.SetDefault(slog.New(slogcolor.NewHandler(os.Stderr, opts)))
		} else {
			err := os.MkdirAll(filepath.Dir(logFilePath), 0755)
			if err != nil {
				slog.Error("Failed to create parent directory for log file", "error", err)
			}

			file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
			if err != nil {
				slog.SetDefault(slog.New(slogcolor.NewHandler(os.Stderr, opts)))
				slog.Error("Failed to open log-file", "error", err)
			} else {
				opts.NoColor = true
				slog.SetDefault(slog.New(slogcolor.NewHandler(file, opts)))
				defer file.Close() // Close the file when done
			}
		}

		cfgFiles := []string{
			"/etc/rong/config.toml",
			"/etc/rong/config.yaml",
			"/etc/rong/config.yml",
			"$HOME/.rong.toml",
			"$HOME/.rong.yaml",
			"$HOME/.rong.yml",
			"$XDG_CONFIG_HOME/rong/config.toml",
			"$XDG_CONFIG_HOME/rong/config.yaml",
			"$XDG_CONFIG_HOME/rong/config.yml",
		}

		cfgFile, err := cmd.Flags().GetString("config")
		if err == nil {
			cfgFiles = append(cfgFiles, cfgFile)
		}

		cwd, _ := os.Getwd()
		for _, cfg := range slices.Backward(cfgFiles) {
			path, err := config.FindPath(cwd, cfg)
			if err != nil {
				continue
			}

			if err := config.LoadConfig(path); err != nil {
				if !os.IsNotExist(err) {
					slog.Error("Failed to load config", "error", err)
				}
				continue
			}

			slog.Info("Loaded config", "path", path)
			break
		}
	},
}
