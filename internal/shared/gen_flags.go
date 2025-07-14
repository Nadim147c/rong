package shared

import (
	"fmt"

	"github.com/Nadim147c/material/dynamic"
	"github.com/Nadim147c/rong/internal/config"
	"github.com/spf13/cobra"
)

// ValidateGeneratorFlags validates color generation flags
func ValidateGeneratorFlags(cmd *cobra.Command) error {
	flags := cmd.Flags()

	// Validate version
	if flags.Changed("version") {
		version, _ := flags.GetInt("version")
		switch dynamic.Version(version) {
		case dynamic.V2021, dynamic.V2025:
			// valid
		default:
			return fmt.Errorf("invalid version: %d (must be 2021 or 2025)", version)
		}
		config.Global.Version = dynamic.Version(version)
	}

	// Validate contrast
	if flags.Changed("contrast") {
		contrast, _ := flags.GetFloat64("contrast")
		if contrast < -1.0 || contrast > 1.0 {
			return fmt.Errorf("contrast must be between -1.0 and 1.0, got %.2f", contrast)
		}
		config.Global.Constrast = contrast
	}

	// Validate variant
	if flags.Changed("variant") {
		variant, _ := flags.GetString("variant")
		validVariants := map[string]bool{
			string(dynamic.Monochrome): true,
			string(dynamic.Neutral):    true,
			string(dynamic.TonalSpot):  true,
			string(dynamic.Vibrant):    true,
			string(dynamic.Expressive): true,
			string(dynamic.Fidelity):   true,
			string(dynamic.Content):    true,
			string(dynamic.Rainbow):    true,
			string(dynamic.FruitSalad): true,
		}
		if !validVariants[variant] {
			return fmt.Errorf("invalid variant: %s", variant)
		}
		config.Global.Variant = dynamic.Variant(variant)
	}

	// Validate platform
	if flags.Changed("platform") {
		platform, _ := flags.GetString("platform")
		validPlatforms := map[string]bool{
			string(dynamic.Phone): true, string(dynamic.Watch): true,
		}
		if !validPlatforms[platform] {
			return fmt.Errorf("invalid platform: %s", platform)
		}
		config.Global.Platform = dynamic.Platform(platform)
	}

	return nil
}
