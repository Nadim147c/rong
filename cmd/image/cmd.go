package image

import (
	"encoding/json"
	"fmt"
	"image"
	"log/slog"
	"os"

	"github.com/Nadim147c/rong/v5/internal/base16"
	"github.com/Nadim147c/rong/v5/internal/cache"
	"github.com/Nadim147c/rong/v5/internal/config"
	"github.com/Nadim147c/rong/v5/internal/material"
	"github.com/Nadim147c/rong/v5/internal/models"
	"github.com/Nadim147c/rong/v5/internal/pathutil"
	"github.com/Nadim147c/rong/v5/internal/templates"
	"github.com/spf13/cobra"

	_ "image/jpeg" // for jpeg encoding
	_ "image/png"  // for png encoding

	_ "golang.org/x/image/webp" // for webp encoding
)

// Command is the image command.
var Command = &cobra.Command{
	Use:   "image <image>",
	Short: "Generate colors from a image",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		imagePath := args[0]

		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		imagePath, err = pathutil.FindPath(cwd, imagePath)
		if err != nil {
			return fmt.Errorf("failed to find image path: %w", err)
		}

		slog.Info("Generating color", "from", imagePath)

		hash, err := cache.Hash(imagePath)
		if err != nil {
			return fmt.Errorf("failed to get xxh sum: %w", err)
		}

		quantized, err := cache.LoadCache(hash)
		if err != nil {
			if !os.IsNotExist(err) {
				slog.Error("Failed to load cache", "error", err)
			}

			file, err := os.Open(imagePath)
			if err != nil {
				return fmt.Errorf("failed to open image file: %w", err)
			}
			defer file.Close()

			img, _, err := image.Decode(file)
			if err != nil {
				return fmt.Errorf("failed to decode image: %w", err)
			}

			pixels := material.GetPixelsFromImage(img)
			quantized, err = material.Quantize(ctx, pixels)
			if err != nil {
				return err
			}

			if err := cache.SaveCache(hash, quantized); err != nil {
				slog.Warn("Failed to save colors to cache", "error", err)
			}
		}

		cfg := material.GetConfig()

		colorMap, wu, err := material.GenerateFromQuantized(quantized, cfg)
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

		output := models.NewOutput(imagePath, based, colorMap, customs)

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

		if err := cache.SaveState(imagePath, hash, quantized); err != nil {
			slog.Warn("Failed to save colors to cache", "error", err)
		}

		return templates.Execute(ctx, output)
	},
}
