package regen

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/Nadim147c/rong/internal/base16"
	"github.com/Nadim147c/rong/internal/cache"
	"github.com/Nadim147c/rong/internal/material"
	"github.com/Nadim147c/rong/internal/models"
	"github.com/Nadim147c/rong/templates"
	"github.com/gabriel-vasile/mimetype"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	Command.Flags().AddFlagSet(material.GeneratorFlags)
}

// Command is the image command
var Command = &cobra.Command{
	Use:   "regen [flags]",
	Short: "Regenerate colors from previous generation",
	Args:  cobra.NoArgs,
	PreRun: func(cmd *cobra.Command, _ []string) {
		viper.BindPFlags(cmd.Flags())
	},
	RunE: func(_ *cobra.Command, _ []string) error {
		state, err := cache.LoadState()
		if err != nil {
			return fmt.Errorf("failed load current state: %v", err)
		}

		slog.Info("Generating color from cached state", "path", state.Path)

		cfg, err := material.GetConfig()
		if err != nil {
			return err
		}

		colorMap, wu, err := material.GenerateFromQuantized(state.Quantized, cfg)
		if err != nil {
			return fmt.Errorf("failed to generate colors: %w", err)
		}

		fg, bg := colorMap["on_background"], colorMap["background"]
		based := base16.Generate(fg, bg, wu)

		path := state.Path
		mtype, err := mimetype.DetectFile(state.Path)
		if err == nil && strings.HasPrefix(mtype.String(), "video") {
			if preview, err := cache.GetPreview(path, state.Hash); err == nil {
				path = preview
			}
		}

		output := models.NewOutput(path, based, colorMap)

		if viper.GetBool("json") {
			if err := json.NewEncoder(os.Stdout).Encode(output); err != nil {
				slog.Error("Failed to encode output", "error", err)
			}
		}

		if !viper.GetBool("dry-run") {
			return templates.Execute(output)
		}

		return nil
	},
}
