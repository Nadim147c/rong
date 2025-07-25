package video

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/Nadim147c/material"
	"github.com/Nadim147c/material/dynamic"
	"github.com/Nadim147c/rong/internal/cache"
	"github.com/Nadim147c/rong/internal/config"
	"github.com/Nadim147c/rong/internal/ffmpeg"
	"github.com/Nadim147c/rong/internal/models"
	"github.com/Nadim147c/rong/internal/shared"
	"github.com/Nadim147c/rong/templates"
	"github.com/spf13/cobra"
)

func init() {
	Command.Flags().Bool("light", false, "generate light color palette")
	Command.Flags().String("variant", string(dynamic.TonalSpot), "variant to use (e.g., tonal_spot, vibrant, expressive)")
	Command.Flags().Float64("contrast", 0.0, "contrast adjustment (-1.0 to 1.0)")
	Command.Flags().String("platform", string(dynamic.Phone), "target platform (phone or watch)")
	Command.Flags().Int("version", int(dynamic.V2021), "version of the theme (2021 or 2025)")
	Command.Flags().Int("frames", 5, "number of frames of vidoe to process")
	Command.Flags().BoolP("json", "j", false, "print generated colors as json")
}

// Command is the image command
var Command = &cobra.Command{
	Use:   "video [flags] <image>",
	Short: "Generate colors from a video",
	Args:  cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, _ []string) error {
		return shared.ValidateGeneratorFlags(cmd)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		videoPath := args[0]

		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		videoPath, err = config.FindPath(cwd, videoPath)
		if err != nil {
			return fmt.Errorf("failed to find image path: %w", err)
		}

		if cached, jsonb, err := cache.LoadCache(videoPath); err == nil {
			slog.Info("Loading color from cache")

			if jsonFlag, _ := cmd.Flags().GetBool("json"); jsonFlag {
				os.Stdout.Write(jsonb)
			}

			templates.Execute(cached)
			return nil
		}

		slog.Info("Couldn't load colors from cache", "error", err)
		slog.Info("Generating colors from source")

		frames, _ := cmd.Flags().GetInt("frames")
		pixels, err := ffmpeg.GetPixels(videoPath, frames)
		if err != nil {
			return fmt.Errorf("Failed to get pixels from media: %w", err)
		}

		colorMap, err := material.GenerateFromPixels(pixels,
			config.Global.Variant, !config.Global.Light,
			config.Global.Constrast, config.Global.Platform,
			config.Global.Version,
		)
		if err != nil {
			return fmt.Errorf("failed to generate colors: %w", err)
		}

		output := models.NewOutput(videoPath, colorMap)

		jsonb, err := cache.SaveCache(videoPath, output)
		if err != nil {
			slog.Warn("Failed to save colors to cache", "error", err)
		}

		if jsonFlag, _ := cmd.Flags().GetBool("json"); jsonFlag {
			os.Stdout.Write(jsonb)
		}

		templates.Execute(output)
		return nil
	},
}
