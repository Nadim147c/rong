package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/MatusOllah/slogcolor"
	"github.com/Nadim147c/rong/cmd/image"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// Command is root command of the cli
var Command = &cobra.Command{
	Use:   "rong",
	Short: "A material you color generator from image.",
}

func init() {
	Command.AddCommand(image.Command)

	Command.PersistentFlags().BoolP("verbose", "v", false, "enable verbose logging")
	Command.PersistentFlags().String("log-file", "", "file to save logs")
}

func Execute() {
	opts := slogcolor.DefaultOptions
	opts.SrcFileMode = 0
	opts.LevelTags = map[slog.Level]string{
		slog.LevelDebug: color.New(color.FgGreen).Sprint("DEBUG"),
		slog.LevelInfo:  color.New(color.FgCyan).Sprint("INFO "),
		slog.LevelWarn:  color.New(color.FgYellow).Sprint("WARN "),
		slog.LevelError: color.New(color.FgRed).Sprint("ERROR"),
	}

	verbose, err := Command.Flags().GetBool("verbose")
	if err == nil && verbose {
		opts.Level = slog.LevelDebug
	}

	logFilePath, err := Command.Flags().GetString("log-file")
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

	if err := Command.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
