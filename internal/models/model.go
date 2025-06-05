package models

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/Nadim147c/goyou/color"
)

// Output contains all values that will execute templates
type Output struct {
	Material
	Image  string
	Colors []Color
}

// Color represents a named color with various color format representations.
type Color struct {
	// Name represents the key of the color in various case styles.
	Name Case
	// Color contains the different representations of the color value.
	Color ColorValue
}

// Case holds different case variations of a color name.
type Case struct {
	// Snake case format (e.g., on_primary)
	Snake string
	// Camel case format (e.g., onPrimary)
	Camel string
	// Kebab case format (e.g., on-primary)
	Kebab string
	// Pascal case format (e.g., OnPrimary)
	Pascal string
}

// ColorValue stores a color in multiple string and numeric formats.
type ColorValue struct {
	// HexRGB is the RGB hexadecimal representation with '#' prefix (e.g., #FF0000)
	HexRGB string
	// TrimmedHexRGB is the RGB hexadecimal representation without '#' (e.g., FF0000)
	TrimmedHexRGB string
	// HexRGBA is the RGBA hexadecimal representation with '#' prefix (e.g., #FF0000FF)
	HexRGBA string
	// TrimmedHexRGBA is the RGBA hexadecimal representation without '#' (e.g., FF0000FF)
	TrimmedHexRGBA string
	// RGB is the RGB string format with values 0–255 (e.g., rgb(0, 255, 0))
	RGB string
	// TrimmedRGB is the comma-separated RGB values 0–255 without the "rgb(...)" wrapper (e.g., 0, 255, 0)
	TrimmedRGB string
	// RGBA is the RGBA string format with values 0–255 (e.g., rgba(0, 255, 0, 255))
	RGBA string
	// TrimmedRGBA is the comma-separated RGBA values 0–255 without the "rgba(...)" wrapper (e.g., 0, 255, 0, 255)
	TrimmedRGBA string
	// LinearRGB is the RGB string format with values in 0–1 range (e.g., rgb(0, 1, 0))
	LinearRGB string
	// TrimmedLinearRGB is the comma-separated RGB values in 0–1 range (e.g., 0, 1, 0)
	TrimmedLinearRGB string
	// LinearRGBA is the RGBA string format with values in 0–1 range (e.g., rgba(0, 1, 0, 1))
	LinearRGBA string
	// TrimmedLinearRGBA is the comma-separated RGBA values in 0–1 range (e.g., 0, 1, 0, 1)
	TrimmedLinearRGBA string

	// Red channel value (0–255)
	Red uint8
	// Green channel value (0–255)
	Green uint8
	// Blue channel value (0–255)
	Blue uint8
	// Alpha channel value (0–255)
	Alpha uint8
}

// NewColor creates a Color
func NewColor(key string, rgb color.ARGB) Color {
	// Convert snake_case to other cases
	var name Case
	name.Snake = key
	name.Camel = toCamelCase(key, false)
	name.Kebab = strings.ReplaceAll(key, "_", "-")
	name.Pascal = toCamelCase(key, true)

	value := NewColorValue(rgb)
	return Color{Name: name, Color: value}
}

func lf(c uint8) float64 { return float64(c) / 255.0 }

// NewColorValue create a ColorValue
func NewColorValue(rgb color.ARGB) ColorValue {
	var value ColorValue
	if rgb == 0 {
		return value
	}

	alpha, red, green, blue := rgb.Values()

	// Format color representations
	value.HexRGB = fmt.Sprintf("#%02X%02X%02X", red, green, blue)
	value.TrimmedHexRGB = fmt.Sprintf("%02X%02X%02X", red, green, blue)
	value.HexRGBA = fmt.Sprintf("#%02X%02X%02X%02X", red, green, blue, alpha)
	value.TrimmedHexRGBA = fmt.Sprintf("%02X%02X%02X%02X", red, green, blue, alpha)
	value.RGB = fmt.Sprintf("rgb(%d, %d, %d)", red, green, blue)
	value.TrimmedRGB = fmt.Sprintf("%d, %d, %d", red, green, blue)
	value.RGBA = fmt.Sprintf("rgba(%d, %d, %d, %d)", red, green, blue, alpha)
	value.TrimmedRGBA = fmt.Sprintf("%d, %d, %d, %d", red, green, blue, alpha)
	value.LinearRGB = fmt.Sprintf("rgb(%.3f, %.3f, %.3f)", lf(red), lf(green), lf(blue))
	value.TrimmedLinearRGB = fmt.Sprintf("%.3f, %.3f, %.3f", lf(red), lf(green), lf(blue))
	value.LinearRGBA = fmt.Sprintf("rgba(%.3f, %.3f, %.3f, %.3f)", lf(red), lf(green), lf(blue), lf(alpha))
	value.TrimmedLinearRGBA = fmt.Sprintf("%.3f, %.3f, %.3f, %.3f", lf(red), lf(green), lf(blue), lf(alpha))

	value.Alpha = alpha
	value.Red = red
	value.Green = green
	value.Blue = blue

	return value
}

func toCamelCase(s string, pascal bool) string {
	parts := strings.Split(s, "_")
	for i, part := range parts {
		if i == 0 && !pascal {
			parts[i] = strings.ToLower(part)
		} else {
			parts[i] = capitalize(part)
		}
	}
	return strings.Join(parts, "")
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}
