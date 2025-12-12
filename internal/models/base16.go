package models

import "github.com/Nadim147c/rong/v4/internal/base16"

// Base16 is 16 colors
type Base16 struct {
	Color0 FormatedColor `json:"color_0"`
	Color1 FormatedColor `json:"color_1"`
	Color2 FormatedColor `json:"color_2"`
	Color3 FormatedColor `json:"color_3"`
	Color4 FormatedColor `json:"color_4"`
	Color5 FormatedColor `json:"color_5"`
	Color6 FormatedColor `json:"color_6"`
	Color7 FormatedColor `json:"color_7"`
	Color8 FormatedColor `json:"color_8"`
	Color9 FormatedColor `json:"color_9"`
	ColorA FormatedColor `json:"color_a"`
	ColorB FormatedColor `json:"color_b"`
	ColorC FormatedColor `json:"color_c"`
	ColorD FormatedColor `json:"color_d"`
	ColorE FormatedColor `json:"color_e"`
	ColorF FormatedColor `json:"color_f"`
}

// NewBase16 creates a Base16 from a base16.Base16 definition
func NewBase16(based base16.Base16) (Base16, []NamedColor) {
	var b Base16

	black := NewNamedColor("color_0", based.Black)
	red := NewNamedColor("color_1", based.Red)
	green := NewNamedColor("color_2", based.Green)
	yellow := NewNamedColor("color_3", based.Yellow)
	blue := NewNamedColor("color_4", based.Blue)
	magenta := NewNamedColor("color_5", based.Magenta)
	cyan := NewNamedColor("color_6", based.Cyan)
	white := NewNamedColor("color_7", based.White)

	brightBlack := NewNamedColor("color_8", based.BrightBlack)
	brightRed := NewNamedColor("color_9", based.BrightRed)
	brightGreen := NewNamedColor("color_a", based.BrightGreen)
	brightYellow := NewNamedColor("color_b", based.BrightYellow)
	brightBlue := NewNamedColor("color_c", based.BrightBlue)
	brightMagenta := NewNamedColor("color_d", based.BrightMagenta)
	brightCyan := NewNamedColor("color_e", based.BrightCyan)
	brightWhite := NewNamedColor("color_f", based.BrightWhite)

	// Assign to Base16 struct
	b.Color0 = black.Color
	b.Color1 = red.Color
	b.Color2 = green.Color
	b.Color3 = yellow.Color
	b.Color4 = blue.Color
	b.Color5 = magenta.Color
	b.Color6 = cyan.Color
	b.Color7 = white.Color
	b.Color8 = brightBlack.Color
	b.Color9 = brightRed.Color
	b.ColorA = brightGreen.Color
	b.ColorB = brightYellow.Color
	b.ColorC = brightBlue.Color
	b.ColorD = brightMagenta.Color
	b.ColorE = brightCyan.Color
	b.ColorF = brightWhite.Color

	// Return ordered slice
	s := []NamedColor{
		black,
		red,
		green,
		yellow,
		blue,
		magenta,
		cyan,
		white,
		brightBlack,
		brightRed,
		brightGreen,
		brightYellow,
		brightBlue,
		brightMagenta,
		brightCyan,
		brightWhite,
	}

	return b, s
}
