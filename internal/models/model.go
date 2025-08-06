package models

import (
	"encoding"
	"fmt"
	"slices"
	"strings"
	"unicode"

	"github.com/Nadim147c/material/color"
)

// Output contains all values that will execute templates
type Output struct {
	Material `json:"material"`
	Base16   `json:"base16"`
	Image    string  `json:"image"`
	Colors   []Color `json:"colors"`
}

// NewOutput create output struct for templates execution
func NewOutput(source string, base16 map[string]color.ARGB, colorMap map[string]color.ARGB) Output {
	colors := make([]Color, 0, len(colorMap)+16)

	for key, value := range base16 {
		colors = append(colors, NewColor(key, value))
	}
	for key, value := range colorMap {
		colors = append(colors, NewColor(key, value))
	}

	slices.SortFunc(colors, func(a, b Color) int {
		return strings.Compare(a.Name.Snake, b.Name.Snake)
	})

	material := NewMaterial(colorMap)
	based := NewBase16(base16)

	return Output{
		Material: material,
		Base16:   based,
		Image:    source,
		Colors:   colors,
	}
}

// Color represents a named color with various color format representations.
type Color struct {
	// Name represents the key of the color in various case styles.
	Name Case `json:"name"`
	// Color contains the different representations of the color value.
	Color ColorValue `json:"value"`
}

// Case holds different case variations of a color name.
type Case struct {
	// Snake case format (e.g., on_primary)
	Snake string `json:"snake"`
	// Camel case format (e.g., onPrimary)
	Camel string `json:"camel"`
	// Kebab case format (e.g., on-primary)
	Kebab string `json:"kebab"`
	// Pascal case format (e.g., OnPrimary)
	Pascal string `json:"pascal"`
}

var (
	_ encoding.TextMarshaler   = (*Case)(nil)
	_ encoding.TextUnmarshaler = (*Case)(nil)
	_ fmt.Stringer             = (*Case)(nil)
)

// MarshalText implements encoding.TextMarshaler
func (c Case) MarshalText() ([]byte, error) {
	return []byte(c.Snake), nil
}

// UnmarshalText implements encoding.TextUnmarshaler
func (c *Case) UnmarshalText(b []byte) error {
	key := string(b)
	c.Snake = key
	c.Camel = toCamelCase(key, false)
	c.Kebab = strings.ReplaceAll(key, "_", "-")
	c.Pascal = toCamelCase(key, true)
	return nil
}

// String implements fmt.Stringer
func (c Case) String() string {
	return c.Snake
}

// ColorValue stores a color in multiple string and numeric formats.
type ColorValue struct {
	// HexRGB is the RGB hexadecimal representation with '#' prefix (e.g., #FF0000)
	HexRGB string `json:"hex_rgb"`
	// TrimmedHexRGB is the RGB hexadecimal representation without '#' (e.g., FF0000)
	TrimmedHexRGB string `json:"trimmed_hex_rgb"`
	// HexRGBA is the RGBA hexadecimal representation with '#' prefix (e.g., #FF0000FF)
	HexRGBA string `json:"hex_rgba"`
	// TrimmedHexRGBA is the RGBA hexadecimal representation without '#' (e.g., FF0000FF)
	TrimmedHexRGBA string `json:"trimmed_hex_rgba"`
	// RGB is the RGB string format with values 0–255 (e.g., rgb(0,255,0))
	RGB string `json:"rgb"`
	// TrimmedRGB is the comma-separated RGB values 0–255 without the "rgb(...)" wrapper (e.g., 0,255,0)
	TrimmedRGB string `json:"trimmed_rgb"`
	// RGBA is the RGBA string format with values 0–255 (e.g., rgba(0,255,0,255))
	RGBA string `json:"rgba"`
	// TrimmedRGBA is the comma-separated RGBA values 0–255 without the "rgba(...)" wrapper (e.g., 0,255,0,255)
	TrimmedRGBA string `json:"trimmed_rgba"`
	// LinearRGB is the RGB string format with values in 0–1 range (e.g., rgb(0,1,0))
	LinearRGB string `json:"linear_rgb"`
	// TrimmedLinearRGB is the comma-separated RGB values in 0–1 range (e.g., 0,1,0)
	TrimmedLinearRGB string `json:"trimmed_linear_rgb"`
	// LinearRGBA is the RGBA string format with values in 0–1 range (e.g., rgba(0,1,0,1))
	LinearRGBA string `json:"linear_rgba"`
	// TrimmedLinearRGBA is the comma-separated RGBA values in 0–1 range (e.g., 0,1,0,1)
	TrimmedLinearRGBA string `json:"trimmed_linear_rgba"`

	// Red channel value (0–255)
	Red uint8 `json:"red"`
	// Green channel value (0–255)
	Green uint8 `json:"green"`
	// Blue channel value (0–255)
	Blue uint8 `json:"blue"`
	// Alpha channel value (0–255)
	Alpha uint8 `json:"alpha"`

	Int color.ARGB
}

var _ fmt.Stringer = (*ColorValue)(nil)

// String is useful when using .NameValue. It will be converted to "#XXXXXX".
func (cv ColorValue) String() string {
	return cv.HexRGB
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
	value.RGB = fmt.Sprintf("rgb(%d,%d,%d)", red, green, blue)
	value.TrimmedRGB = fmt.Sprintf("%d,%d,%d", red, green, blue)
	value.RGBA = fmt.Sprintf("rgba(%d,%d,%d,%d)", red, green, blue, alpha)
	value.TrimmedRGBA = fmt.Sprintf("%d,%d,%d,%d", red, green, blue, alpha)
	value.LinearRGB = fmt.Sprintf("rgb(%.3f,%.3f,%.3f)", lf(red), lf(green), lf(blue))
	value.TrimmedLinearRGB = fmt.Sprintf("%.3f,%.3f,%.3f", lf(red), lf(green), lf(blue))
	value.LinearRGBA = fmt.Sprintf("rgba(%.3f,%.3f,%.3f,%.3f)", lf(red), lf(green), lf(blue), lf(alpha))
	value.TrimmedLinearRGBA = fmt.Sprintf("%.3f,%.3f,%.3f,%.3f", lf(red), lf(green), lf(blue), lf(alpha))

	value.Alpha = alpha
	value.Red = red
	value.Green = green
	value.Blue = blue
	value.Int = rgb

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
