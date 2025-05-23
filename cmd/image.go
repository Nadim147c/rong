package cmd

import (
	"fmt"
	"image"
	"os"

	_ "image/jpeg"
	_ "image/png"

	"github.com/Nadim147c/rong/internal/material"
	"github.com/Nadim147c/rong/internal/models"
	"github.com/Nadim147c/rong/templates"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	_ "golang.org/x/image/webp" // Add WebP support
)

var imageCmd = &cobra.Command{
	Use:   "image [flags] <image>",
	Short: "",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return viper.BindPFlags(cmd.Flags())
	},
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		imagePath := args[0]

		file, err := os.Open(imagePath)
		if err != nil {
			return fmt.Errorf("failed to open image file: %w", err)
		}
		defer file.Close()

		img, _, err := image.Decode(file)
		if err != nil {
			return fmt.Errorf("failed to decode image: %w", err)
		}

		colorMap, err := material.GenerateColorsFromImage(img, !viper.GetBool("light"))
		if err != nil {
			return err
		}
		materialColors := material.MaterialColorFromMap(colorMap)

		colors := make([]models.Color, 0, len(colorMap))
		for key, value := range colorMap {
			colors = append(colors, models.NewColor(key, value))
		}

		output := models.Output{
			Colors:        colors,
			MaterialColor: materialColors,
		}

		templates.Execute(output)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(imageCmd)
	imageCmd.Flags().BoolP("light", "l", false, "use light theme")
}
