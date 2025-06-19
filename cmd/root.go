package cmd

import (
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/MatusOllah/slogcolor"
	"github.com/Nadim147c/rong/cmd/image"
	"github.com/Nadim147c/rong/cmd/video"
	"github.com/Nadim147c/rong/internal/config"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	Command.AddCommand(image.Command)
	Command.AddCommand(video.Command)

	Command.PersistentFlags().BoolP("verbose", "v", false, "enable verbose logging")
	Command.PersistentFlags().String("log-file", "", "file to save logs")
	Command.PersistentFlags().StringP("config", "c", "$XDG_CONFIG_HOME/rong/config.toml", "path to config (.toml) file")
}

// Command is root command of the cli
var Command = &cobra.Command{
	Use:   "rong",
	Short: "A material you color generator from image.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		opts := slogcolor.DefaultOptions
		opts.TimeFormat = time.Kitchen
		opts.SrcFileMode = 0
		opts.LevelTags = map[slog.Level]string{
			slog.LevelDebug: color.New(color.FgGreen).Sprint("DBG"),
			slog.LevelInfo:  color.New(color.FgCyan).Sprint("INF"),
			slog.LevelWarn:  color.New(color.FgYellow).Sprint("WRN"),
			slog.LevelError: color.New(color.FgRed).Sprint("ERR"),
		}

		verbose, err := cmd.Flags().GetBool("verbose")
		if err == nil && verbose {
			opts.Level = slog.LevelDebug
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

		cfgFile, err := cmd.Flags().GetString("config")
		if err != nil {
			cfgFile = "$XDG_CONFIG_HOME/rong/config.toml"
		}
		config.LoadConfig(cfgFile)
	},
}
