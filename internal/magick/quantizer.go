package magick

import (
	"fmt"
	"log/slog"
	"os/exec"
	"strings"

	"github.com/Nadim147c/goyou/color"
)

// GenerateFromImage quantizes an image using magick
func GenerateFromImage(path string, count int) ([]color.ARGB, error) {
	cmd := exec.Command(
		"magick", path,
		"-resize", "25%",
		"-colors", fmt.Sprint(count),
		"-unique-colors",
		"-colorspace", "srgb",
		"-depth", "8",
		"txt:-",
	)

	colors := make([]color.ARGB, 0, count)

	output, err := cmd.Output()
	if err != nil {
		return colors, fmt.Errorf("Failed to get magick output: %v", err)
	}

	for line := range strings.Lines(string(output)) {
		fields := strings.Fields(line)
		if len(fields) != 4 {
			continue
		}

		var r, g, b, a uint8
		_, err := fmt.Sscanf(fields[1], "(%d,%d,%d,%d)", &r, &g, &b, &a)
		if err != nil {
			slog.Debug("Failed to parse color", "field", fields[2], "line", line)
			continue
		}
		rgb := color.NewARGB(a, r, g, b)
		colors = append(colors, rgb)
	}

	return colors, nil
}
