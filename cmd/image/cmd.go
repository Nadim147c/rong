package image

import (
	"fmt"
	"image"
	"os"
	"slices"
	"strings"

	// for jpeg encoding
	_ "image/jpeg"
	_ "image/png"

	"github.com/Nadim147c/rong/internal/config"
	"github.com/Nadim147c/rong/internal/material"
	"github.com/Nadim147c/rong/internal/models"
	"github.com/Nadim147c/rong/templates"
	"github.com/spf13/cobra"

	// for webp encoding
	_ "golang.org/x/image/webp"
)

var light = false

func init() {
	Command.Flags().BoolVarP(&light, "light", "l", light, "use light theme")
}

// Command is the image command
var Command = &cobra.Command{
	Use:   "image [flags] <image>",
	Short: "Generate color from a image",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
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

		colorMap, err := material.GenerateColorsFromImage(img, !light)
		if err != nil {
			return err
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
