package material

import (
	"context"
	"fmt"
	"image"

	"github.com/Nadim147c/material/color"
	"github.com/Nadim147c/material/dynamic"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Config is configuration used to generate colors
type Config struct {
	Variant   dynamic.Variant
	Platform  dynamic.Platform
	Version   dynamic.Version
	Dark      bool
	Constrast float64
}

// GetConfig return Config from viper
func GetConfig() (Config, error) {
	config := Config{}

	variant, err := dynamic.ParseVariant(viper.GetString("variant"))
	if err != nil {
		return config, err
	}
	config.Variant = variant

	version, err := dynamic.ParseVersion(viper.GetString("version"))
	if err != nil {
		return config, err
	}
	config.Version = version

	platform, err := dynamic.ParsePlatform(viper.GetString("platform"))
	if err != nil {
		return config, err
	}
	config.Platform = platform

	config.Constrast = viper.GetFloat64("contrast")
	if config.Constrast < -1 || config.Constrast > 1 {
		return config, fmt.Errorf(
			"contrast must between -1 to 1 but got %.2f",
			config.Constrast,
		)
	}

	config.Dark = viper.GetBool("dark")

	return config, nil
}

// GeneratorFlags are the flags used for generating colors
var GeneratorFlags = pflag.NewFlagSet("generate", pflag.ContinueOnError)

func init() {
	GeneratorFlags.Bool("dark", false, "generate dark color palette")
	GeneratorFlags.Bool(
		"dry-run",
		false,
		"generate colors without applying templates",
	)
	GeneratorFlags.Bool("json", false, "print generated colors as json")
	GeneratorFlags.Float64("contrast", 0.0, "contrast adjustment (-1.0 to 1.0)")
	GeneratorFlags.String(
		"version",
		dynamic.Version2025.String(),
		"version of the theme (2021 or 2025)",
	)
	GeneratorFlags.String(
		"platform",
		dynamic.PlatformPhone.String(),
		"target platform (phone or watch)",
	)
	GeneratorFlags.String(
		"variant",
		dynamic.VariantTonalSpot.String(),
		"variant to use (e.g., tonal_spot, vibrant, expressive)",
	)
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
	ctx context.Context,
	img image.Image,
	cfg Config,
) (Colors, []color.ARGB, error) {
	pixels := GetPixelsFromImage(img)
	return GenerateFromPixels(ctx, pixels, cfg)
}
