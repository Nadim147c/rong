package material

import (
	"image"

	"github.com/Nadim147c/material/color"
	"github.com/Nadim147c/material/dynamic"
)

// GeneratorConfig is configuration used to generate colors
type GeneratorConfig struct {
	Variant   dynamic.Variant  `config:"variant" check:"enum='monochrome,neutral,tonal_spot,vibrant,expressive,fidelity,content,rainbow,fruit_salad'"`
	Platform  dynamic.Platform `config:"platform" check:"enum='phone,watch'"`
	Version   dynamic.Version  `config:"version" check:"enum='2021,2025'"`
	Dark      bool             `config:"dark"`
	Constrast float64          `config:"Constrast" check:"min=-1,max=1"`
}

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
	cfg GeneratorConfig,
) (Colors, []color.ARGB, error) {
	pixels := GetPixelsFromImage(img)
	return GenerateFromPixels(pixels, cfg)
}
