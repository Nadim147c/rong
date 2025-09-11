package cmd

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/MatusOllah/slogcolor"
	"github.com/Nadim147c/material/dynamic"
	"github.com/Nadim147c/rong/cmd/cache"
	"github.com/Nadim147c/rong/cmd/color"
	"github.com/Nadim147c/rong/cmd/image"
	"github.com/Nadim147c/rong/cmd/regen"
	"github.com/Nadim147c/rong/cmd/video"
	"github.com/Nadim147c/rong/internal/pathutil"
	"github.com/carapace-sh/carapace"
	"github.com/carapace-sh/carapace/pkg/style"
	termcolor "github.com/fatih/color"
	"github.com/mattn/go-isatty"
	slogmulti "github.com/samber/slog-multi"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

	imageComp := carapace.Gen(image.Command)
	imageComp.FlagCompletion(actions)
	imageComp.PositionalAnyCompletion(carapace.ActionFiles())

	videoComp := carapace.Gen(video.Command)
	videoComp.FlagCompletion(actions)
	videoComp.PositionalAnyCompletion(carapace.ActionFiles())

	cacheComp := carapace.Gen(cache.Command)
	cacheComp.FlagCompletion(actions)
	cacheComp.PositionalAnyCompletion(carapace.ActionFiles())

	colorComp := carapace.Gen(color.Command)
	colorComp.FlagCompletion(actions)
	nameCompletions := make([]string, 0, len(color.Names)*2)
	for name, value := range color.Names {
		_, r, g, b := value.Values()
		nameCompletions = append(nameCompletions, name, style.TrueColor(r, g, b))
	}
	colorComp.PositionalAnyCompletion(carapace.ActionStyledValues(nameCompletions...))

	regenComp := carapace.Gen(regen.Command)
	regenComp.FlagCompletion(actions)

	rootComp := carapace.Gen(Command)
	rootComp.Standalone()
	rootComp.FlagCompletion(carapace.ActionMap{
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
			viper.SetOptions(viper.WithLogger(logger))
		} else {
			err := os.MkdirAll(filepath.Dir(logFilePath), 0755)
			if err != nil {
				slog.Error("Failed to create parent directory for log file", "error", err)
			}

			file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
			if err != nil {
				logger := slog.New(stderrHandler)
				slog.SetDefault(logger)
				viper.SetOptions(viper.WithLogger(logger))
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
				viper.SetOptions(viper.WithLogger(logger))
				logfile = file
			}
		}

		viper.AddConfigPath("/etc/rong")
		viper.AddConfigPath(pathutil.ConfigDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")

		viper.SetEnvPrefix("rong")

		viper.SetDefault("dark", true)
		viper.SetDefault("variant", dynamic.Expressive)
		viper.SetDefault("platform", dynamic.Phone)

		cfgFlag := cmd.Flags().Lookup("config")
		if cfgFlag != nil && cfgFlag.Changed {
			viper.SetConfigFile(cfgFlag.Value.String())
		}

		return viper.ReadInConfig()
	},
}
