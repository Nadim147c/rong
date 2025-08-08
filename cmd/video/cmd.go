package video

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/Nadim147c/material/dynamic"
	"github.com/Nadim147c/rong/internal/base16"
	"github.com/Nadim147c/rong/internal/cache"
	"github.com/Nadim147c/rong/internal/config"
	"github.com/Nadim147c/rong/internal/ffmpeg"
	"github.com/Nadim147c/rong/internal/material"
	"github.com/Nadim147c/rong/internal/models"
	"github.com/Nadim147c/rong/internal/pathutil"
	"github.com/Nadim147c/rong/templates"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	Command.Flags().Bool("dark", false, "generate dark color palette")
	Command.Flags().String("variant", string(dynamic.TonalSpot), "variant to use (e.g., tonal_spot, vibrant, expressive)")
	Command.Flags().Float64("contrast", 0.0, "contrast adjustment (-1.0 to 1.0)")
	Command.Flags().String("platform", string(dynamic.Phone), "target platform (phone or watch)")
	Command.Flags().Int("version", int(dynamic.V2021), "version of the theme (2021 or 2025)")
	Command.Flags().Int("frames", 5, "number of frames of vidoe to process")
	Command.Flags().BoolP("json", "j", false, "print generated colors as json")
	Command.Flags().Bool("dry-run", false, "generate colors without applying templates")
}

// Command is the image command
var Command = &cobra.Command{
	Use:   "video [flags] <image>",
	Short: "Generate colors from a video",
	Args:  cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, _ []string) {
		viper.BindPFlags(cmd.Flags())
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		videoPath := args[0]

		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		videoPath, err = pathutil.FindPath(cwd, videoPath)
		if err != nil {
			return fmt.Errorf("failed to find image path: %w", err)
		}

		quantized, err := cache.LoadCache(videoPath)
		if err != nil {
			if !os.IsNotExist(err) {
				slog.Error("Failed to load cache", "error", err)
			}

			frames, _ := cmd.Flags().GetInt("frames")
			pixels, err := ffmpeg.GetPixels(videoPath, frames)
			if err != nil {
				return fmt.Errorf("Failed to get pixels from media: %w", err)
			}
			quantized = material.Quantize(pixels)
		}

		slog.Info("Couldn't load colors from cache", "error", err)
		slog.Info("Generating colors from source")

		cfg, err := config.GetGeneratorConfig()
		if err != nil {
			return err
		}
		colorMap, wu, err := material.GenerateFromQuantized(quantized, cfg)
		if err != nil {
			return fmt.Errorf("failed to generate colors: %w", err)
		}

		fg, bg := colorMap["on_background"], colorMap["background"]
		based := base16.Generate(fg, bg, wu)

		output := models.NewOutput(videoPath, based, colorMap)

		if err := cache.SaveCache(videoPath, quantized); err != nil {
			slog.Warn("Failed to save colors to cache", "error", err)
		}

		if jsonFlag, _ := cmd.Flags().GetBool("json"); jsonFlag {
			if err := json.NewEncoder(os.Stdout).Encode(output); err != nil {
				slog.Error("Failed to encode output", "error", err)
			}
		}

		if dry, _ := cmd.Flags().GetBool("dry-run"); !dry {
			templates.Execute(output)
		}
		return nil
	},
}
