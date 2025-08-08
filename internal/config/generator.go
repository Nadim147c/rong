package config

import (
	"fmt"

	"github.com/Nadim147c/material/dynamic"
	"github.com/spf13/viper"
)

// GeneratorConfig is configuration used to generate colors
type GeneratorConfig struct {
	Variant   dynamic.Variant
	Platform  dynamic.Platform
	Version   dynamic.Version
	Dark      bool
	Constrast float64
}

// GetGeneratorConfig returns validated color generator config
func GetGeneratorConfig() (GeneratorConfig, error) {
	var cfg GeneratorConfig
	cfg.Dark = viper.GetBool("dark")

	contrast := viper.GetFloat64("contrast")
	if contrast < -1.0 || contrast > 1.0 {
		return cfg, fmt.Errorf("contrast must be between -1.0 and 1.0, got %.2f", contrast)
	}
	cfg.Constrast = contrast

	version := dynamic.Version(viper.GetInt("version"))
	if version != dynamic.V2021 && version != dynamic.V2025 {
		return cfg, fmt.Errorf("invalid version: %d (must be 2021 or 2025)", version)
	}
	cfg.Version = version

	variant := dynamic.Variant(viper.GetString("variant"))
	validVariants := map[dynamic.Variant]bool{
		dynamic.Monochrome: true, dynamic.Neutral: true, dynamic.TonalSpot: true,
		dynamic.Vibrant: true, dynamic.Expressive: true, dynamic.Fidelity: true,
		dynamic.Content: true, dynamic.Rainbow: true, dynamic.FruitSalad: true,
	}
	if _, ok := validVariants[variant]; !ok {
		return cfg, fmt.Errorf("invalid variant: %s", variant)
	}
	cfg.Variant = variant

	platform := dynamic.Platform(viper.GetString("platform"))
	if platform != dynamic.Phone && platform != dynamic.Watch {
		return cfg, fmt.Errorf("invalid platform: %s", platform)
	}
	cfg.Platform = platform

	return cfg, nil
}
