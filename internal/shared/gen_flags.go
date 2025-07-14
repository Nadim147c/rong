package shared

import (
	"fmt"

	"github.com/Nadim147c/material/dynamic"
	"github.com/spf13/cobra"
)

// ValidateGeneratorFlags validates color generation flags
func ValidateGeneratorFlags(cmd *cobra.Command) error {
	// Validate contrast
	contrast, err := cmd.Flags().GetFloat64("contrast")
	if err != nil {
		return err
	}
	if contrast < -1.0 || contrast > 1.0 {
		return fmt.Errorf("contrast must be between -1.0 and 1.0, got %.2f", contrast)
	}

	// Validate variant
	variant, err := cmd.Flags().GetString("variant")
	if err != nil {
		return err
	}
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

	// Validate platform
	platform, err := cmd.Flags().GetString("platform")
	if err != nil {
		return err
	}
	validPlatforms := map[string]bool{
		string(dynamic.Phone): true,
		string(dynamic.Watch): true,
	}
	if !validPlatforms[platform] {
		return fmt.Errorf("invalid platform: %s", platform)
	}

	// Validate version
	version, err := cmd.Flags().GetInt("version")
	if err != nil {
		return err
	}
	switch dynamic.Version(version) {
	case dynamic.V2021, dynamic.V2025:
		// valid
	default:
		return fmt.Errorf("invalid version: %d (must be 2021 or 2025)", version)
	}

	return nil
}
