package cmd

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/MatusOllah/slogcolor"
	"github.com/Nadim147c/go-config"
	"github.com/Nadim147c/material/dynamic"
	"github.com/Nadim147c/rong/cmd/cache"
	"github.com/Nadim147c/rong/cmd/color"
	"github.com/Nadim147c/rong/cmd/image"
	"github.com/Nadim147c/rong/cmd/regen"
	"github.com/Nadim147c/rong/cmd/video"
	"github.com/Nadim147c/rong/internal/pathutil"
	"github.com/carapace-sh/carapace"
	termcolor "github.com/fatih/color"
	"github.com/mattn/go-isatty"
	slogmulti "github.com/samber/slog-multi"
	"github.com/spf13/cobra"
)

func init() {
	Command.AddCommand(color.Command)
	Command.AddCommand(image.Command)
	Command.AddCommand(video.Command)
	Command.AddCommand(cache.Command)
	Command.AddCommand(regen.Command)

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

	regenCara := carapace.Gen(regen.Command)
	regenCara.FlagCompletion(actions)

	rootCara := carapace.Gen(Command)
	rootCara.Standalone()
	rootCara.FlagCompletion(carapace.ActionMap{
		"config":   carapace.ActionFiles(),
		"log-file": carapace.ActionFiles(),
	})

	Command.PersistentFlags().BoolP("verbose", "v", false, "enable verbose logging")
	Command.PersistentFlags().BoolP("quiet", "q", false, "suppress all logs")
	Command.PersistentFlags().String("log-file", "", "file to save logs")
	Command.PersistentFlags().StringP("config", "c", "$XDG_CONFIG_HOME/rong/config.{toml,yaml,yml}", "path to config (.toml|.yaml|.yml) file")
	Command.MarkFlagsMutuallyExclusive("verbose", "quiet")
}

var logfile *os.File

// Command is root command of the cli
var Command = &cobra.Command{
	Use:          "rong",
	Short:        "A material you color generator from image or video.",
	SilenceUsage: true,
	PersistentPostRun: func(_ *cobra.Command, _ []string) {
		if logfile != nil {
			slog.Info("Exiting rong")
			logfile.Close()
		}
	},
	PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
		if cmd.Name() == "_carapace" {
			return nil
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

		stderrHandler := slogcolor.NewHandler(os.Stderr, opts)

		logFilePath, err := cmd.Flags().GetString("log-file")
		if err != nil || logFilePath == "" {
			logger := slog.New(stderrHandler)
			slog.SetDefault(logger)
			config.SetLogger(logger)
		} else {
			err := os.MkdirAll(filepath.Dir(logFilePath), 0755)
			if err != nil {
				slog.Error("Failed to create parent directory for log file", "error", err)
			}

			file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
			if err != nil {
				logger := slog.New(stderrHandler)
				slog.SetDefault(logger)
				config.SetLogger(logger)
				slog.Error("Failed to open log-file", "error", err)
			} else {
				fileHanlder := slog.NewJSONHandler(file, &slog.HandlerOptions{
					AddSource: true,
					// Manually enabling file logs indicates user is trying to debug
					Level: slog.LevelDebug,
				})

				handler := slogmulti.Fanout(stderrHandler, fileHanlder)
				logger := slog.New(handler)
				slog.SetDefault(logger)
				config.SetLogger(logger)
				logfile = file
			}
		}

		config.AddPath("/etc/rong")
		config.AddPath(pathutil.ConfigDir)
		config.SetFormat("yaml")
		config.SetEnvPrefix("rong")

		config.SetDefault("dark", true)
		config.SetDefault("variant", dynamic.Expressive)
		config.SetDefault("platform", dynamic.Phone)

		cfgFile, err := cmd.Flags().GetString("config")
		if err == nil {
			config.AddFile(cfgFile)
		}

		return config.ReadConfig()
	},
}
