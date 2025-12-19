package material

import (
	"context"
	"fmt"
	"image"

	"github.com/Nadim147c/material/v2/color"
	"github.com/Nadim147c/material/v2/dynamic"
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

	variant, err := dynamic.ParseVariant(viper.GetString("material.variant"))
	if err != nil {
		return config, err
	}
	config.Variant = variant

	version, err := dynamic.ParseVersion(viper.GetString("material.version"))
	if err != nil {
		return config, err
	}
	config.Version = version

	platform, err := dynamic.ParsePlatform(viper.GetString("material.platform"))
	if err != nil {
		return config, err
	}
	config.Platform = platform

	config.Constrast = viper.GetFloat64("material.contrast")
	if config.Constrast < -1 || config.Constrast > 1 {
		return config, fmt.Errorf(
			"contrast must between -1 to 1 but got %.2f",
			config.Constrast,
		)
	}

	config.Dark = viper.GetBool("dark")

	return config, nil
}

// Flags are the flags used for generating colors
var Flags = pflag.NewFlagSet("material", pflag.ContinueOnError)

func init() {
	// TODO: These flags should be here
	Flags.BoolP("dark", "D", false, "generate dark color palette")
	Flags.BoolP("json", "j", false, "print generated colors as json")
	Flags.BoolP(
		"dry-run", "d", false,
		"generate colors without applying templates",
	)

	Flags.Float64("material.contrast", 0.0, "contrast adjustment (-1.0 to 1.0)")

	Flags.String(
		"material.version", dynamic.Version2025.String(),
		"version of the theme (2021 or 2025)",
	)
	viper.SetDefault("material.version", dynamic.Version2025.String())

	Flags.String(
		"material.platform", dynamic.PlatformPhone.String(),
		"target platform (phone or watch)",
	)
	viper.SetDefault("material.platform", dynamic.PlatformPhone.String())

	Flags.String(
		"material.variant", dynamic.VariantTonalSpot.String(),
		"variant to use (e.g., tonal_spot, vibrant, expressive)",
	)
	viper.SetDefault("material.variant", dynamic.VariantTonalSpot.String())

	Flags.Bool(
		"material.custom.blend", true,
		"blend custom colors with primary",
	)
	viper.SetDefault("material.custom.blend", true)

	Flags.Float64(
		"material.custom.ratio", 0.3,
		"blend custom colors with primary",
	)
	viper.SetDefault("material.custom.blend", 0.3)
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
