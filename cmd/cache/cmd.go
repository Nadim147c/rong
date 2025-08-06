package cache

import (
	"log/slog"

	"github.com/Nadim147c/material/dynamic"
	"github.com/Nadim147c/rong/internal/base16"
	"github.com/Nadim147c/rong/internal/cache"
	"github.com/Nadim147c/rong/internal/config"
	"github.com/Nadim147c/rong/internal/ffmpeg"
	"github.com/Nadim147c/rong/internal/material"
	"github.com/Nadim147c/rong/internal/models"
	"github.com/Nadim147c/rong/internal/shared"
	"github.com/spf13/cobra"
)

func init() {
	Command.Flags().Bool("light", false, "generate light color palette")
	Command.Flags().String("variant", string(dynamic.TonalSpot), "variant to use (e.g., tonal_spot, vibrant, expressive)")
	Command.Flags().Float64("contrast", 0.0, "contrast adjustment (-1.0 to 1.0)")
	Command.Flags().String("platform", string(dynamic.Phone), "target platform (phone or watch)")
	Command.Flags().Int("version", int(dynamic.V2021), "version of the theme (2021 or 2025)")
	Command.Flags().Int("frames", 5, "number of frames of vidoe to process")
}

// Command is cache command
var Command = &cobra.Command{
	Use:   "cache [flags] ...<image|video|directory>",
	Short: "Generate color cache from a image/video",
	Args:  cobra.MinimumNArgs(1),
	PreRunE: func(cmd *cobra.Command, _ []string) error {
		return shared.ValidateGeneratorFlags(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		paths := make(chan string)

		go ScanPaths(args, paths)

		for path := range paths {
			if cache.IsCached(path) {
				slog.Info("Skipping", "path", path, "reason", "already cached")
				continue
			}

			frames, _ := cmd.Flags().GetInt("frames")
			pixels, err := ffmpeg.GetPixels(path, frames)
			if err != nil {
				slog.Error("Failed to get pixels from media", "path", path, "error", err)
				continue
			}

			colorMap, wu, err := material.GenerateFromPixels(pixels,
				config.Global.Variant, !config.Global.Light,
				config.Global.Constrast, config.Global.Platform,
				config.Global.Version,
			)
			if err != nil {
				slog.Error("Failed to generate colors", "error", err)
				continue
			}

			fg, bg := colorMap["on_background"], colorMap["background"]
			based := base16.Generate(fg, bg, !config.Global.Light, wu)

			output := models.NewOutput(path, based, colorMap)

			if err := cache.SaveCache(output); err != nil {
				slog.Error("Failed to save cache", "path", path, "error", err)
				continue
			}

			slog.Info("Successfully cached media", "path", path)
		}
	},
}
