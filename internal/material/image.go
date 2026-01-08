package material

import (
	"context"
	"image"

	"github.com/Nadim147c/material/v2/color"
	"github.com/Nadim147c/material/v2/dynamic"
	"github.com/Nadim147c/rong/v5/internal/config"
)

// Config is configuration used to generate colors.
type Config struct {
	Variant   dynamic.Variant
	Platform  dynamic.Platform
	Version   dynamic.Version
	Dark      bool
	Constrast float64
}

func GetConfig() Config {
	return Config{
		Variant:   config.MaterialVariant.Value(),
		Version:   config.MaterialVersion.Value(),
		Platform:  config.MaterialPlatformt.Value(),
		Constrast: config.MaterialContrast.Value(),
		Dark:      config.Dark.Value(),
	}
}

// GetPixelsFromImage returns pixels from image.Imaget interface.
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

// GenerateFromImage colors from an image.Image.
func GenerateFromImage(
	ctx context.Context,
	img image.Image,
	cfg Config,
) (Colors, []color.ARGB, error) {
	pixels := GetPixelsFromImage(img)
	return GenerateFromPixels(ctx, pixels, cfg)
}
