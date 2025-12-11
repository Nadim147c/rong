package models

import (
	"fmt"
	"slices"
	"strings"
	"unicode"

	"github.com/Nadim147c/material/v2/color"
	"github.com/Nadim147c/rong/v3/internal/base16"
	"github.com/Nadim147c/rong/v3/internal/material"
)

// Output contains all values that will execute templates
type Output struct {
	Material `json:"material"`
	Base16   `json:"base16"`
	Image    string       `json:"image"`
	Colors   []NamedColor `json:"colors"`
}

// NewOutput create output struct for templates execution
func NewOutput(
	source string,
	base16Colors base16.Base16,
	materialColors map[string]color.ARGB,
	customColors map[string]material.CustomColor,
) Output {
	colors := make([]NamedColor, 0, len(materialColors)+16)

	b, basedSlice := NewBase16(base16Colors)
	colors = append(colors, basedSlice...)

	for key, value := range materialColors {
		colors = append(colors, NewNamedColor(key, value))
	}

	for key, value := range customColors {
		colors = append(
			colors,
			NewNamedColor(key, value.Color),
			NewNamedColor("on_"+key, value.OnColor),
			NewNamedColor(key+"_container", value.ColorContainer),
			NewNamedColor("on_"+key+"_container", value.OnColor),
		)
	}

	slices.SortFunc(colors, func(a, b NamedColor) int {
		return strings.Compare(a.Name.Snake, b.Name.Snake)
	})

	material := NewMaterial(materialColors, customColors)

	return Output{
		Material: material,
		Base16:   b,
		Image:    source,
		Colors:   colors,
	}
}

// NamedColor represents a named color with various color format
// representations.
type NamedColor struct {
	// Name represents the key of the color in various case styles.
	Name ColorName `json:"name"`
	// Color contains the different representations of the color value.
	Color FormatedColor `json:"value"`
}

// ColorName holds different case variations of a color name.
type ColorName struct {
	// Snake case format (e.g., on_primary)
	Snake string `json:"snake"`
	// Camel case format (e.g., onPrimary)
	Camel string `json:"camel"`
	// Kebab case format (e.g., on-primary)
	Kebab string `json:"kebab"`
	// Pascal case format (e.g., OnPrimary)
	Pascal string `json:"pascal"`
}

var _ fmt.Stringer = (*ColorName)(nil)

// String implements fmt.Stringer
func (c ColorName) String() string {
	return c.Snake
}

// FormatedColor stores a color in multiple string and numeric formats.
type FormatedColor struct {
	// HexRGB is the RGB hexadecimal representation with '#' prefix (e.g.,
	// #FF0000)
	HexRGB string `json:"hex_rgb"`
	// TrimmedHexRGB is the RGB hexadecimal representation without '#' (e.g.,
	// FF0000)
	TrimmedHexRGB string `json:"trimmed_hex_rgb"`
	// HexRGBA is the RGBA hexadecimal representation with '#' prefix (e.g.,
	// #FF0000FF)
	HexRGBA string `json:"hex_rgba"`
	// TrimmedHexRGBA is the RGBA hexadecimal representation without '#' (e.g.,
	// FF0000FF)
	TrimmedHexRGBA string `json:"trimmed_hex_rgba"`
	// RGB is the RGB string format with values 0–255 (e.g., rgb(0,255,0))
	RGB string `json:"rgb"`
	// TrimmedRGB is the comma-separated RGB values 0–255 without the
	// "rgb(...)" wrapper (e.g., 0,255,0)
	TrimmedRGB string `json:"trimmed_rgb"`
	// RGBA is the RGBA string format with values 0–255 (e.g.,
	// rgba(0,255,0,255))
	RGBA string `json:"rgba"`
	// TrimmedRGBA is the comma-separated RGBA values 0–255 without the
	// "rgba(...)" wrapper (e.g., 0,255,0,255)
	TrimmedRGBA string `json:"trimmed_rgba"`
	// LinearRGB is the RGB string format with values in 0–1 range (e.g.,
	// rgb(0,1,0))
	LinearRGB string `json:"linear_rgb"`
	// TrimmedLinearRGB is the comma-separated RGB values in 0–1 range (e.g.,
	// 0,1,0)
	TrimmedLinearRGB string `json:"trimmed_linear_rgb"`
	// LinearRGBA is the RGBA string format with values in 0–1 range (e.g.,
	// rgba(0,1,0,1))
	LinearRGBA string `json:"linear_rgba"`
	// TrimmedLinearRGBA is the comma-separated RGBA values in 0–1 range
	// (e.g., 0,1,0,1)
	TrimmedLinearRGBA string `json:"trimmed_linear_rgba"`
	// AnsiColor is the terminals ansi sequence without the escape character
	// (e.g., 38;2;1;2;3)
	AnsiForeground string `json:"ansi_foreground"`
	// AnsiColor is the terminals ansi sequence without the escape character
	// (e.g., 48;2;1;2;3)
	AnsiBackground string `json:"ansi_background"`

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

var _ fmt.Stringer = (*FormatedColor)(nil)

// String is useful when using .NameValue. It will be converted to "#XXXXXX".
func (cv FormatedColor) String() string {
	return cv.HexRGB
}

// NewNamedColor creates a Color
func NewNamedColor(key string, rgb color.ARGB) NamedColor {
	// Convert snake_case to other cases
	var name ColorName
	key = strings.ToLower(key)
	name.Snake = key
	name.Camel = toCamelCase(key, false)
	name.Kebab = strings.ReplaceAll(key, "_", "-")
	name.Pascal = toCamelCase(key, true)

	value := NewFormatedColor(rgb)
	return NamedColor{Name: name, Color: value}
}

func lf(c uint8) float64 { return float64(c) / 255.0 }

// NewFormatedColor create a ColorValue
func NewFormatedColor(rgb color.ARGB) FormatedColor {
	var value FormatedColor
	if rgb == 0 {
		return value
	}

	alpha, red, green, blue := rgb.Components()

	// Format color representations
	value.HexRGB = fmt.Sprintf("#%02X%02X%02X", red, green, blue)
	value.TrimmedHexRGB = fmt.Sprintf("%02X%02X%02X", red, green, blue)
	value.HexRGBA = fmt.Sprintf("#%02X%02X%02X%02X", red, green, blue, alpha)
	value.TrimmedHexRGBA = fmt.Sprintf(
		"%02X%02X%02X%02X",
		red,
		green,
		blue,
		alpha,
	)
	value.RGB = fmt.Sprintf("rgb(%d,%d,%d)", red, green, blue)
	value.TrimmedRGB = fmt.Sprintf("%d,%d,%d", red, green, blue)
	value.RGBA = fmt.Sprintf("rgba(%d,%d,%d,%d)", red, green, blue, alpha)
	value.TrimmedRGBA = fmt.Sprintf("%d,%d,%d,%d", red, green, blue, alpha)
	value.LinearRGB = fmt.Sprintf(
		"rgb(%.3f,%.3f,%.3f)",
		lf(red),
		lf(green),
		lf(blue),
	)
	value.TrimmedLinearRGB = fmt.Sprintf(
		"%.3f,%.3f,%.3f",
		lf(red),
		lf(green),
		lf(blue),
	)
	value.LinearRGBA = fmt.Sprintf(
		"rgba(%.3f,%.3f,%.3f,%.3f)",
		lf(red),
		lf(green),
		lf(blue),
		lf(alpha),
	)
	value.TrimmedLinearRGBA = fmt.Sprintf(
		"%.3f,%.3f,%.3f,%.3f",
		lf(red),
		lf(green),
		lf(blue),
		lf(alpha),
	)
	value.AnsiForeground = fmt.Sprintf("38;2;%d;%d;%d", red, green, blue)
	value.AnsiBackground = fmt.Sprintf("48;2;%d;%d;%d", red, green, blue)

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
