package cmd

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"math"
	"os"
	"path/filepath"
	"syscall"

	"github.com/Nadim147c/fang"
	"github.com/Nadim147c/rong/cmd/cache"
	"github.com/Nadim147c/rong/cmd/color"
	"github.com/Nadim147c/rong/cmd/image"
	"github.com/Nadim147c/rong/cmd/regen"
	"github.com/Nadim147c/rong/cmd/video"
	"github.com/Nadim147c/rong/internal/pathutil"
	"github.com/carapace-sh/carapace"
	"github.com/carapace-sh/carapace/pkg/style"
	"github.com/charmbracelet/log"
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
		nameCompletions = append(
			nameCompletions,
			name,
			style.TrueColor(r, g, b),
		)
	}
	colorComp.PositionalAnyCompletion(
		carapace.ActionStyledValues(nameCompletions...),
	)

	regenComp := carapace.Gen(regen.Command)
	regenComp.FlagCompletion(actions)

	rootComp := carapace.Gen(Command)
	rootComp.Standalone()
	rootComp.FlagCompletion(carapace.ActionMap{
		"config":   carapace.ActionFiles(),
		"log-file": carapace.ActionFiles(),
	})

	Command.PersistentFlags().CountP("verbose", "v", "enable verbose logging")
	Command.PersistentFlags().BoolP("quiet", "q", false, "suppress all logs")
	Command.PersistentFlags().StringP("log-file", "l", "", "file to save logs")
	Command.PersistentFlags().
		StringP("config", "c", "$XDG_CONFIG_HOME/rong/config.{toml,yaml,yml}", "path to config (.toml|.yaml|.yml) file")
	Command.MarkFlagsMutuallyExclusive("verbose", "quiet")
}

func handleError(w io.Writer, styles fang.Styles, err error) {
	if errors.Is(err, context.Canceled) {
		err = errors.New("operation cancelled by user")
	}
	fang.DefaultErrorHandler(w, styles, err)
}

// Execute runs the cobra cli
func Execute(version string) error {
	return fang.Execute(
		context.Background(),
		Command,
		fang.WithErrorHandler(handleError),
		fang.WithFlagTypes(),
		fang.WithNotifySignal(syscall.SIGINT, syscall.SIGTERM),
		fang.WithShorthandPadding(),
		fang.WithVersion(version),
		fang.WithoutCompletions(),
	)
}

func should[T any](v T, _ error) T {
	return v
}

var logfile *os.File

// Command is root command of the cli
var Command = &cobra.Command{
	Use:          "rong",
	Short:        "A material you color generator from image or video.",
	SilenceUsage: true,
	PersistentPostRun: func(_ *cobra.Command, _ []string) {
		slog.Info("Exiting rong")
		if logfile != nil {
			logfile.Close()
		}
	},
	PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
		if cmd.Name() == "_carapace" {
			return nil
		}

		level := slog.LevelInfo

		if verbose, err := cmd.Flags().GetCount("verbose"); err == nil {
			if verbose == 0 {
				level = slog.LevelError
			} else if verbose == 1 {
				level = slog.LevelWarn
			} else if verbose == 2 {
				level = slog.LevelInfo
			} else if verbose >= 3 {
				level = slog.Level(math.MinInt)
			}
		}

		quiet := should(cmd.Flags().GetBool("quiet"))
		if quiet {
			level = slog.Level(math.MaxInt)
		}

		stderrHandler := log.NewWithOptions(os.Stderr, log.Options{
			ReportTimestamp: false, Level: log.Level(level),
		})

		if cmd.Flags().Changed("log-file") {
			logFilePath := should(cmd.Flags().GetString("log-file"))

			if err := os.MkdirAll(
				filepath.Dir(logFilePath),
				0o755,
			); err != nil {
				slog.Error(
					"Failed to create parent directory for log file",
					"error",
					err,
				)
			}

			file, err := os.OpenFile(
				logFilePath,
				os.O_CREATE|os.O_APPEND|os.O_WRONLY,
				0o666,
			)
			if err != nil {
				logger := slog.New(stderrHandler)
				slog.SetDefault(logger)
				viper.SetOptions(viper.WithLogger(logger))
				slog.Error("Failed to open log-file", "error", err)
			} else {
				fileHanlder := slog.NewJSONHandler(file, &slog.HandlerOptions{
					AddSource: true,
					// Manually enabling file logs indicates user is trying to
					// debug
					Level: slog.LevelDebug,
				})

				handler := slogmulti.Fanout(stderrHandler, fileHanlder)
				logger := slog.New(handler)
				slog.SetDefault(logger)
				viper.SetOptions(viper.WithLogger(logger))
				logfile = file
			}
		} else {
			logger := slog.New(stderrHandler)
			slog.SetDefault(logger)
			viper.SetOptions(viper.WithLogger(logger))
		}

		viper.AddConfigPath("/etc/rong")
		viper.AddConfigPath(pathutil.ConfigDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")

		viper.SetEnvPrefix("rong")
		viper.AutomaticEnv()

		cfgFlag := cmd.Flags().Lookup("config")
		if cfgFlag != nil && cfgFlag.Changed {
			viper.SetConfigFile(cfgFlag.Value.String())
		}

		return viper.ReadInConfig()
	},
}
