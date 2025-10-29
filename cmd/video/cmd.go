package video

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/Nadim147c/rong/internal/base16"
	"github.com/Nadim147c/rong/internal/cache"
	"github.com/Nadim147c/rong/internal/ffmpeg"
	"github.com/Nadim147c/rong/internal/material"
	"github.com/Nadim147c/rong/internal/models"
	"github.com/Nadim147c/rong/internal/pathutil"
	"github.com/Nadim147c/rong/templates"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	Command.Flags().AddFlagSet(material.GeneratorFlags)
	Command.Flags().AddFlagSet(base16.Flags)
	Command.Flags().Int("frames", 5, "number of frames of vidoe to process")
}

// Command is the video command
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
	PreRun: func(cmd *cobra.Command, _ []string) {
		viper.BindPFlags(cmd.Flags())
	},
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
			return fmt.Errorf("failed to get xxh sum: %v", err)
		}

		quantized, err := cache.LoadCache(hash)
		if err != nil {
			if !os.IsNotExist(err) {
				slog.Error("Failed to load cache", "error", err)
			}

			frames := viper.GetInt("frames")
			pixels, err := ffmpeg.GetPixels(ctx, videoPath, frames)
			if err != nil {
				return fmt.Errorf("Failed to get pixels from media: %w", err)
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

		cfg, err := material.GetConfig()
		if err != nil {
			return err
		}

		colorMap, wu, err := material.GenerateFromQuantized(quantized, cfg)
		if err != nil {
			return fmt.Errorf("failed to generate colors: %w", err)
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

		output := models.NewOutput(path, based, colorMap)

		if err := cache.SaveState(videoPath, hash, quantized); err != nil {
			slog.Warn("Failed to save colors to cache", "error", err)
		}

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
