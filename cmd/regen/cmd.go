package regen

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/Nadim147c/rong/v4/internal/base16"
	"github.com/Nadim147c/rong/v4/internal/cache"
	"github.com/Nadim147c/rong/v4/internal/config"
	"github.com/Nadim147c/rong/v4/internal/material"
	"github.com/Nadim147c/rong/v4/internal/models"
	"github.com/Nadim147c/rong/v4/internal/templates"
	"github.com/gabriel-vasile/mimetype"
	"github.com/spf13/cobra"
)

// Command is the image command.
var Command = &cobra.Command{
	Use:   "regen [flags]",
	Short: "Regenerate colors from previous generation",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		ctx := cmd.Context()
		state, err := cache.LoadState()
		if err != nil {
			return fmt.Errorf("failed load current state: %v", err) //nolint
		}

		slog.Info("Generating color from cached state", "path", state.Path)

		cfg := material.GetConfig()

		colorMap, wu, err := material.GenerateFromQuantized(
			state.Quantized,
			cfg,
		)
		if err != nil {
			return fmt.Errorf("failed to generate colors: %w", err)
		}

		customs, err := material.GenerateCustomColors(colorMap["primary"])
		if err != nil {
			return err
		}

		based, err := base16.Generate(colorMap, wu)
		if err != nil {
			return err
		}

		path := state.Path
		mtype, err := mimetype.DetectFile(state.Path)
		if err == nil && strings.HasPrefix(mtype.String(), "video") {
			if preview, err := cache.GetPreview(path, state.Hash); err == nil {
				path = preview
			}
		}

		output := models.NewOutput(path, based, colorMap, customs)

		if config.JSON.Value() {
			err := json.NewEncoder(cmd.OutOrStdout()).Encode(output)
			if err != nil {
				slog.Error("Failed to encode output", "error", err)
			}
		}

		if config.SimpleJSON.Value() {
			err := models.WriteSimpleJSON(cmd.OutOrStdout(), output)
			if err != nil {
				slog.Error("Failed to encode output", "error", err)
			}
		}

		if config.DryRun.Value() {
			return nil
		}

		return templates.Execute(ctx, output)
	},
}
