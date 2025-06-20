package video

import (
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"

	"github.com/Nadim147c/goyou/color"
	"github.com/Nadim147c/rong/internal/config"
	"github.com/Nadim147c/rong/internal/material"
	"github.com/Nadim147c/rong/internal/models"
	"github.com/Nadim147c/rong/templates"
	"github.com/spf13/cobra"
)

var light = false

func init() {
	Command.Flags().BoolVarP(&light, "light", "l", light, "use light theme")
}

// Command is the image command
var Command = &cobra.Command{
	Use:   "video [flags] <image>",
	Short: "Generate color from a video",
	Args:  cobra.ExactArgs(1),
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

		cmd := exec.Command("ffmpeg",
			"-i", videoPath,
			"-vframes", "1", // Limit to 1 frame
			"-f", "rawvideo",
			"-pix_fmt", "rgb24",
			"-")

		out, err := cmd.Output()
		if err != nil {
			return err
		}

		totalBytes := len(out)

		pixels := make([]color.ARGB, 0, totalBytes/3)
		for i := 0; i+2 < totalBytes; i += 3 {
			c := color.ARGBFromRGB(out[i], out[i+1], out[i+2])
			pixels = append(pixels, c)
		}

		colorMap, err := material.GenerateColorsFromPixels(pixels, !light)
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
			Image:    videoPath,
			Colors:   colors,
			Material: material,
		}

		templates.Execute(output)
		return nil
	},
}
