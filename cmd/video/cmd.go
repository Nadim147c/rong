package video

import (
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"

	"github.com/Nadim147c/material"
	"github.com/Nadim147c/material/color"
	"github.com/Nadim147c/material/dynamic"
	"github.com/Nadim147c/rong/internal/config"
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
}

// Command is the image command
var Command = &cobra.Command{
	Use:   "video [flags] <image>",
	Short: "Generate color from a video",
	Args:  cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, _ []string) error {
		return shared.ValidateGeneratorFlags(cmd)
	},
	RunE: func(_ *cobra.Command, args []string) error {
		videoPath := args[0]

		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		videoPath, err = config.FindPath(cwd, videoPath)
		if err != nil {
			return fmt.Errorf("failed to find image path: %w", err)
		}

		if err != nil {
			return err
		}

		ffmpeg := exec.Command("ffmpeg",
			"-i", videoPath,
			"-vframes", "5", // Limit to 1 frame
			"-f", "rawvideo",
			"-pix_fmt", "rgb24",
			"-")

		out, err := ffmpeg.Output()
		if err != nil {
			return err
		}

		totalBytes := len(out)

		pixels := make([]color.ARGB, 0, totalBytes/3)
		for i := 0; i+2 < totalBytes; i += 3 {
			c := color.ARGBFromRGB(out[i], out[i+1], out[i+2])
			pixels = append(pixels, c)
		}

		colorMap, err := material.GenerateFromPixels(pixels,
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
			Image:    videoPath,
			Colors:   colors,
			Material: material,
		}

		templates.Execute(output)
		return nil
	},
}
