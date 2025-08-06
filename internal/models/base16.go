package models

import "github.com/Nadim147c/material/color"

// Base16 is 16 colors
type Base16 struct {
	Color0 ColorValue `json:"color_0"`
	Color1 ColorValue `json:"color_1"`
	Color2 ColorValue `json:"color_2"`
	Color3 ColorValue `json:"color_3"`
	Color4 ColorValue `json:"color_4"`
	Color5 ColorValue `json:"color_5"`
	Color6 ColorValue `json:"color_6"`
	Color7 ColorValue `json:"color_7"`
	Color8 ColorValue `json:"color_8"`
	Color9 ColorValue `json:"color_9"`
	ColorA ColorValue `json:"color_a"`
	ColorB ColorValue `json:"color_b"`
	ColorC ColorValue `json:"color_c"`
	ColorD ColorValue `json:"color_d"`
	ColorE ColorValue `json:"color_e"`
	ColorF ColorValue `json:"color_f"`
}

// NewBase16 create Base16 from map of base16 color
func NewBase16(colorMap map[string]color.ARGB) Base16 {
	return Base16{
		Color0: NewColorValue(colorMap["color_0"]),
		Color1: NewColorValue(colorMap["color_1"]),
		Color2: NewColorValue(colorMap["color_2"]),
		Color3: NewColorValue(colorMap["color_3"]),
		Color4: NewColorValue(colorMap["color_4"]),
		Color5: NewColorValue(colorMap["color_5"]),
		Color6: NewColorValue(colorMap["color_6"]),
		Color7: NewColorValue(colorMap["color_7"]),
		Color8: NewColorValue(colorMap["color_8"]),
		Color9: NewColorValue(colorMap["color_9"]),
		ColorA: NewColorValue(colorMap["color_a"]),
		ColorB: NewColorValue(colorMap["color_b"]),
		ColorC: NewColorValue(colorMap["color_c"]),
		ColorD: NewColorValue(colorMap["color_d"]),
		ColorE: NewColorValue(colorMap["color_e"]),
		ColorF: NewColorValue(colorMap["color_f"]),
	}
}
