package material

import (
	"errors"
	"image"

	"github.com/Nadim147c/goyou/color"
)

var ErrNoColorFound = errors.New("no color found")

// Colors is key and color
type Colors = map[string]color.ARGB

// GenerateColorsFromImage colors from an image.Image
func GenerateColorsFromImage(img image.Image, isDark bool) (Colors, error) {
	bounds := img.Bounds()
	pixels := make([]color.ARGB, 0, bounds.Dx()*bounds.Dy())

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.At(x, y)
			argb := color.ARGBFromInterface(c)
			pixels = append(pixels, argb)
		}
	}

	return GenerateColorsFromPixels(pixels, isDark)
}
