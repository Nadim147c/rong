package base16

import (
	"github.com/Nadim147c/material/color"
	"github.com/Nadim147c/material/num"
)

// standard ansi color
var staticColors = []string{}

// SourceColors is all source colors for static generation and fallback for
// dynamic generation
type SourceColors struct {
	// Black is terminal color 0,8
	Black color.ARGB
	// Red is terminal color 1,9
	Red color.ARGB
	// Green is terminal color 2,A
	Green color.ARGB
	// Yellow is terminal color 3,B
	Yellow color.ARGB
	// Blue is terminal color 4,C
	Blue color.ARGB
	// Magenta is terminal color 5,D
	Magenta color.ARGB
	// Cyan is terminal color 6,E
	Cyan color.ARGB
	// White is terminal color 7,F
	White color.ARGB
}

var defaultSrcColors = SourceColors{
	Black:   0xFF000000, // #000000
	Red:     0xFF800000, // #800000
	Green:   0xFF008000, // #008000
	Yellow:  0xFF808000, // #808000
	Blue:    0xFF0044FF, // #0044FF
	Magenta: 0xFF800080, // #800080
	Cyan:    0xFF008080, // #008080
	White:   0xFFC0C0C0, // #C0C0C0
}

// GenerateStatic generates base16 colors from pre-defined colors
func GenerateStatic(primary color.ARGB, src SourceColors) Base16 {
	primaryLab := primary.ToLab()

	black := blend(src.Black.ToLab(), primaryLab, BlendRatio)
	red := blend(src.Red.ToLab(), primaryLab, BlendRatio)
	green := blend(src.Green.ToLab(), primaryLab, BlendRatio)
	yellow := blend(src.Yellow.ToLab(), primaryLab, BlendRatio)
	blue := blend(src.Blue.ToLab(), primaryLab, BlendRatio)
	magenta := blend(src.Magenta.ToLab(), primaryLab, BlendRatio)
	cyan := blend(src.Cyan.ToLab(), primaryLab, BlendRatio)
	white := blend(src.White.ToLab(), primaryLab, BlendRatio)

	if num.DifferenceDegrees(blue.Hue, white.Hue) < 60 {
		white.Hue = blue.Hue - 60
	}

	based := NewBase16()
	based.SetWhite(white)
	based.SetBlack(black)
	based.SetRed(red)
	based.SetGreen(green)
	based.SetYellow(yellow)
	based.SetBlue(blue)
	based.SetMagenta(magenta)
	based.SetCyan(cyan)

	return based
}

func blend(src color.Lab, dst color.Lab, ratio float64) color.Hct {
	if ratio <= 0 {
		return src.ToXYZ().ToHct()
	}
	if ratio >= 1 {
		return dst.ToXYZ().ToHct()
	}

	b := color.Lab{
		L: src.L + (dst.L-src.L)*ratio,
		A: src.A + (dst.A-src.A)*ratio,
		B: src.B + (dst.B-src.B)*ratio,
	}
	return b.ToXYZ().ToHct()
}
