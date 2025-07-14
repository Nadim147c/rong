package image

import (
	"fmt"
	"image"
	"os"
	"slices"
	"strings"

	"github.com/Nadim147c/material"
	"github.com/Nadim147c/material/dynamic"
	"github.com/Nadim147c/rong/internal/config"
	"github.com/Nadim147c/rong/internal/models"
	"github.com/Nadim147c/rong/internal/shared"
	"github.com/Nadim147c/rong/templates"
	"github.com/spf13/cobra"

	_ "golang.org/x/image/webp" // for webp encoding
	_ "image/jpeg"              // for jpeg encoding
	_ "image/png"               // for png encoding
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
	RunE: func(cmd *cobra.Command, args []string) error {
		imagePath := args[0]

		cwd, _ := os.Getwd()
		imagePath, err := config.FindPath(cwd, imagePath)
		if err != nil {
			return fmt.Errorf("failed to find image path: %w", err)
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

		variant, _ := cmd.Flags().GetString("variant")
		light, _ := cmd.Flags().GetBool("light")
		contrast, _ := cmd.Flags().GetFloat64("contrast")
		platform, _ := cmd.Flags().GetString("platform")
		version, _ := cmd.Flags().GetInt("version")

		colorMap, err := material.GenerateFromImage(img,
			dynamic.Variant(variant), !light, contrast,
			dynamic.Platform(platform), dynamic.Version(version),
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

		templates.Execute(output)
		return nil
	},
}
