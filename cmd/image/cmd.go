package image

import (
	"encoding/json"
	"fmt"
	"image"
	"log/slog"
	"os"

	"github.com/Nadim147c/go-config"
	"github.com/Nadim147c/material/dynamic"
	"github.com/Nadim147c/rong/internal/base16"
	"github.com/Nadim147c/rong/internal/cache"
	"github.com/Nadim147c/rong/internal/material"
	"github.com/Nadim147c/rong/internal/models"
	"github.com/Nadim147c/rong/internal/pathutil"
	"github.com/Nadim147c/rong/templates"
	"github.com/spf13/cobra"

	_ "image/jpeg" // for jpeg encoding
	_ "image/png"  // for png encoding

	_ "golang.org/x/image/webp" // for webp encoding
)

func init() {
	Command.Flags().Bool("dark", false, "generate dark color palette")
	Command.Flags().String("variant", string(dynamic.TonalSpot), "variant to use (e.g., tonal_spot, vibrant, expressive)")
	Command.Flags().Float64("contrast", 0.0, "contrast adjustment (-1.0 to 1.0)")
	Command.Flags().String("platform", string(dynamic.Phone), "target platform (phone or watch)")
	Command.Flags().Int("version", int(dynamic.V2021), "version of the theme (2021 or 2025)")
	Command.Flags().BoolP("json", "j", false, "print generated colors as json")
	Command.Flags().Bool("dry-run", false, "generate colors without applying templates")
}

// Command is the image command
var Command = &cobra.Command{
	Use:   "image [flags] <image>",
	Short: "Generate colors from a image",
	Args:  cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, _ []string) {
		config.SetPflagSet(cmd.Flags())
	},
	RunE: func(_ *cobra.Command, args []string) error {
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

		quantized, err := cache.LoadCache(imagePath)
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
			quantized = material.Quantize(pixels)
		}

		var cfg material.GeneratorConfig
		if err := config.Bind("", &cfg); err != nil {
			return err
		}
		colorMap, wu, err := material.GenerateFromQuantized(quantized, cfg)
		if err != nil {
			return fmt.Errorf("failed to generate colors: %w", err)
		}

		fg, bg := colorMap["on_background"], colorMap["background"]
		based := base16.Generate(fg, bg, wu)

		output := models.NewOutput(imagePath, based, colorMap)

		if err := cache.SaveCache(imagePath, quantized); err != nil {
			slog.Warn("Failed to save colors to cache", "error", err)
		}

		if config.GetBool("json") {
			if err := json.NewEncoder(os.Stdout).Encode(output); err != nil {
				slog.Error("Failed to encode output", "error", err)
			}
		}

		if !config.GetBool("dry-run") {
			return templates.Execute(output)
		}

		return nil
	},
}
