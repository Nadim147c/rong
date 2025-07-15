package image

import (
	"fmt"
	"image"
	"log/slog"
	"os"
	"slices"
	"strings"

	"github.com/Nadim147c/material"
	"github.com/Nadim147c/material/dynamic"
	"github.com/Nadim147c/rong/internal/cache"
	"github.com/Nadim147c/rong/internal/config"
	"github.com/Nadim147c/rong/internal/models"
	"github.com/Nadim147c/rong/internal/shared"
	"github.com/Nadim147c/rong/templates"
	"github.com/spf13/cobra"

	_ "image/jpeg" // for jpeg encoding
	_ "image/png"  // for png encoding

	_ "golang.org/x/image/webp" // for webp encoding
)

func init() {
	Command.Flags().Bool("light", false, "generate light color palette")
	Command.Flags().String("variant", string(dynamic.TonalSpot), "variant to use (e.g., tonal_spot, vibrant, expressive)")
	Command.Flags().Float64("contrast", 0.0, "contrast adjustment (-1.0 to 1.0)")
	Command.Flags().String("platform", string(dynamic.Phone), "target platform (phone or watch)")
	Command.Flags().Int("version", int(dynamic.V2021), "version of the theme (2021 or 2025)")
}

// Command is the image command
var Command = &cobra.Command{
	Use:   "image [flags] <image>",
	Short: "Generate color from a image",
	Args:  cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, _ []string) error {
		return shared.ValidateGeneratorFlags(cmd)
	},
	RunE: func(_ *cobra.Command, args []string) error {
		imagePath := args[0]

		cwd, _ := os.Getwd()
		imagePath, err := config.FindPath(cwd, imagePath)
		if err != nil {
			return fmt.Errorf("failed to find image path: %w", err)
		}

		cached, err := cache.LoadCache(imagePath)
		if err == nil {
			templates.Execute(cached)
			return nil
		}
		slog.Info("Couldn't load colors from cache", "error", err)
		slog.Info("Generating colors from source")

		file, err := os.Open(imagePath)
		if err != nil {
			return fmt.Errorf("failed to open image file: %w", err)
		}
		defer file.Close()

		img, _, err := image.Decode(file)
		if err != nil {
			return fmt.Errorf("failed to decode image: %w", err)
		}

		colorMap, err := material.GenerateFromImage(img,
			config.Global.Variant, !config.Global.Light,
			config.Global.Constrast, config.Global.Platform,
			config.Global.Version,
		)
		if err != nil {
			return fmt.Errorf("failed to generate colors: %w", err)
		}

		material := models.MaterialFromMap(colorMap)

		colors := make([]models.Color, 0, len(colorMap))
		for key, value := range colorMap {
			colors = append(colors, models.NewColor(key, value))
		}

		slices.SortFunc(colors, func(a, b models.Color) int {
			return strings.Compare(a.Name.Snake, b.Name.Snake)
		})

		output := models.Output{
			Image:    imagePath,
			Colors:   colors,
			Material: material,
		}

		if err := cache.SaveCache(imagePath, output); err != nil {
			slog.Warn("Failed to save colors to cache", "error", err)
		}

		templates.Execute(output)
		return nil
	},
}
