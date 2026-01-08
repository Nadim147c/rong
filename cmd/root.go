package cmd

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"math"
	"os"
	"path/filepath"
	"slices"
	"syscall"

	"github.com/Nadim147c/fang"
	"github.com/Nadim147c/rong/v4/cmd/cache"
	"github.com/Nadim147c/rong/v4/cmd/color"
	"github.com/Nadim147c/rong/v4/cmd/image"
	"github.com/Nadim147c/rong/v4/cmd/regen"
	"github.com/Nadim147c/rong/v4/cmd/video"
	"github.com/Nadim147c/rong/v4/internal/config"
	"github.com/Nadim147c/rong/v4/internal/pathutil"
	"github.com/carapace-sh/carapace"
	"github.com/carapace-sh/carapace/pkg/style"
	"github.com/charmbracelet/log"
	slogmulti "github.com/samber/slog-multi"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	Command.AddCommand(color.Command)
	Command.AddCommand(image.Command)
	Command.AddCommand(video.Command)
	Command.AddCommand(cache.Command)
	Command.AddCommand(regen.Command)

	commonFlags := pflag.NewFlagSet("generate", pflag.ContinueOnError)
	config.Dark.RegisterFlag(commonFlags)
	config.JSON.RegisterFlag(commonFlags)
	config.SimpleJSON.RegisterFlag(commonFlags)
	config.DryRun.RegisterFlag(commonFlags)

	config.MaterialContrast.RegisterFlag(commonFlags)
	config.MaterialCustomBlend.RegisterFlag(commonFlags)
	config.MaterialCustomColors.RegisterFlag(commonFlags)
	config.MaterialPlatformt.RegisterFlag(commonFlags)
	config.MaterialVariant.RegisterFlag(commonFlags)
	config.MaterialVersion.RegisterFlag(commonFlags)

	config.Base16Blend.RegisterFlag(commonFlags)
	config.Base16Method.RegisterFlag(commonFlags)

	config.Base16Black.RegisterFlag(commonFlags)
	config.Base16Blue.RegisterFlag(commonFlags)
	config.Base16Cyan.RegisterFlag(commonFlags)
	config.Base16Green.RegisterFlag(commonFlags)
	config.Base16Magenta.RegisterFlag(commonFlags)
	config.Base16Red.RegisterFlag(commonFlags)
	config.Base16White.RegisterFlag(commonFlags)
	config.Base16Yellow.RegisterFlag(commonFlags)

	generateCmds := []*cobra.Command{
		color.Command,
		image.Command,
		regen.Command,
		video.Command,
	}
	for cmd := range slices.Values(generateCmds) {
		cmd.Flags().AddFlagSet(commonFlags)
		carapace.Gen(cmd).FlagCompletion(config.CarapaceAction)
	}

	videoFlagSet := pflag.NewFlagSet("preview", pflag.ContinueOnError)
	config.PreviewFormat.RegisterFlag(videoFlagSet)
	config.FFmpegDuration.RegisterFlag(videoFlagSet)
	config.FFmpegFrames.RegisterFlag(videoFlagSet)

	carapace.Gen(image.Command).PositionalAnyCompletion(carapace.ActionFiles())

	video.Command.Flags().AddFlagSet(videoFlagSet)
	carapace.Gen(video.Command).PositionalAnyCompletion(carapace.ActionFiles())

	cache.Command.Flags().AddFlagSet(videoFlagSet)
	carapace.Gen(cache.Command).PositionalAnyCompletion(carapace.ActionFiles())

	colorComp := carapace.Gen(color.Command)
	colorComp.FlagCompletion(config.CarapaceAction)
	nameCompletions := make([]string, 0, len(color.Names)*2)
	for name, value := range color.Names {
		r, g, b := value.Red(), value.Green(), value.Blue()
		nameCompletions = append(
			nameCompletions,
			name,
			style.TrueColor(r, g, b),
		)
	}
	colorComp.PositionalAnyCompletion(carapace.ActionStyledValues(nameCompletions...))

	rootComp := carapace.Gen(Command)
	rootComp.Standalone()
	rootComp.FlagCompletion(config.CarapaceAction)

	persFlags := Command.PersistentFlags()
	config.Verbose.RegisterFlag(persFlags)
	config.Quiet.RegisterFlag(persFlags)
	config.LogFile.RegisterFlag(persFlags)
	config.Config.RegisterFlag(persFlags)
}

func handleError(w io.Writer, styles fang.Styles, err error) {
	if errors.Is(err, context.Canceled) {
		err = errors.New("operation cancelled by user") //nolint
	}
	fang.DefaultErrorHandler(w, styles, err)
}

// Execute runs the cobra cli.
func Execute(version string) {
	err := fang.Execute(
		context.Background(),
		Command,
		fang.WithErrorHandler(handleError),
		fang.WithFlagTypes(),
		fang.WithNotifySignal(syscall.SIGINT, syscall.SIGTERM),
		fang.WithShorthandPadding(),
		fang.WithVersion(version),
		fang.WithoutCompletions(),
	)
	if err != nil {
		os.Exit(1)
	}
}

func should[T any](v T, _ error) T {
	return v
}

var logfile *os.File

// Command is root command of the cli.
var Command = &cobra.Command{
	Use:          "rong",
	Short:        "A material you color generator from image or video.",
	SilenceUsage: true,
	PersistentPostRunE: func(*cobra.Command, []string) error {
		slog.Info("Exiting rong")
		if logfile != nil {
			return logfile.Close()
		}
		return nil
	},
	PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
		if cmd.Name() == "_carapace" {
			return nil
		}

		level := slog.LevelWarn - slog.Level(config.Verbose.Value()*4)

		quiet := should(cmd.Flags().GetBool("quiet"))
		if quiet {
			level = slog.Level(math.MaxInt)
		}

		stderrHandler := log.NewWithOptions(os.Stderr, log.Options{
			ReportTimestamp: false, Level: log.Level(level),
		})

		if cmd.Flags().Changed("log-file") {
			logFilePath := should(cmd.Flags().GetString("log-file"))

			err := os.MkdirAll(filepath.Dir(logFilePath), 0o750)
			if err != nil {
				slog.Error("Failed to create parent directory for log file", "error", err)
			}

			mod := os.O_CREATE | os.O_APPEND | os.O_WRONLY
			file, err := os.OpenFile(logFilePath, mod, 0o600)
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
		} else {
			logger := slog.New(stderrHandler)
			slog.SetDefault(logger)
			viper.SetOptions(viper.WithLogger(logger))
		}

		viper.SetEnvPrefix("rong")
		viper.AutomaticEnv()

		if value := config.Config.Value(); value != "" {
			slog.Info("Cofniguration path has been set", "value", value)
			if slices.Contains([]string{"no", "0", "false"}, value) {
				return nil
			}
			viper.SetConfigFile(value)
		} else {
			viper.AddConfigPath("/etc/rong")
			viper.AddConfigPath(pathutil.ConfigDir)
			viper.SetConfigName("config")
		}

		return viper.ReadInConfig()
	},
}
