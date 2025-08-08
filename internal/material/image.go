package material

import (
	"image"

	"github.com/Nadim147c/material/color"
	"github.com/Nadim147c/rong/internal/config"
)

// GetPixelsFromImage returns pixels from image.Imaget interface
func GetPixelsFromImage(img image.Image) []color.ARGB {
	bounds := img.Bounds()
	pixels := make([]color.ARGB, 0, bounds.Dx()*bounds.Dy())

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.At(x, y)
			argb := color.ARGBFromInterface(c)
			pixels = append(pixels, argb)
		}
	}

	return pixels
}

// GenerateFromImage colors from an image.Image
func GenerateFromImage(
	img image.Image,
	cfg config.GeneratorConfig,
) (Colors, []color.ARGB, error) {
	pixels := GetPixelsFromImage(img)
	return GenerateFromPixels(pixels, cfg)
}
