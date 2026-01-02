package video

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/Nadim147c/rong/v4/internal/base16"
	"github.com/Nadim147c/rong/v4/internal/cache"
	"github.com/Nadim147c/rong/v4/internal/config"
	"github.com/Nadim147c/rong/v4/internal/ffmpeg"
	"github.com/Nadim147c/rong/v4/internal/material"
	"github.com/Nadim147c/rong/v4/internal/models"
	"github.com/Nadim147c/rong/v4/internal/pathutil"
	"github.com/Nadim147c/rong/v4/internal/templates"
	"github.com/spf13/cobra"
)

// Command is the video command.
var Command = &cobra.Command{
	Use:   "video <video>",
	Short: "Generate colors from a video",
	Example: `
# Generate from a video
rong video path/to/video.mkv

# Generate from a image
rong video path/to/image.webp

# Get generate colors as json
rong video path/to/image.mp4 --dry-run --json | jq
  `,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		videoPath := args[0]

		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		videoPath, err = pathutil.FindPath(cwd, videoPath)
		if err != nil {
			return fmt.Errorf("failed to find image path: %w", err)
		}

		slog.Info("Generating color", "from", videoPath)

		hash, err := cache.Hash(videoPath)
		if err != nil {
			return fmt.Errorf("failed to get xxh sum: %w", err)
		}

		quantized, err := cache.LoadCache(hash)
		if err != nil {
			if !os.IsNotExist(err) {
				slog.Error("Failed to load cache", "error", err)
			}

			frames := config.FFmpegFrames.Value()
			duration := config.FFmpegDuration.Value().Seconds()
			pixels, err := ffmpeg.GetPixels(ctx, videoPath, frames, duration)
			if err != nil {
				return fmt.Errorf("failed to get pixels from media: %w", err)
			}
			quantized, err = material.Quantize(ctx, pixels)
			if err != nil {
				return err
			}

			if err := cache.SaveCache(hash, quantized); err != nil {
				slog.Warn("Failed to save colors to cache", "error", err)
			}
		}

		slog.Info("Generating colors from source")

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

		path, err := cache.GetPreview(videoPath, hash)
		if err != nil {
			slog.Warn("Failed to generate preview image", "error", err)
			path = videoPath
		} else {
			slog.Info("Using generated preview", "path", path)
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

		if err := cache.SaveState(videoPath, hash, quantized); err != nil {
			slog.Warn("Failed to save colors to cache", "error", err)
		}

		return templates.Execute(ctx, output)
	},
}
