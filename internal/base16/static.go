package base16

import (
	"cmp"
	"math"
	"slices"

	"github.com/Nadim147c/material/v2/color"
	"github.com/Nadim147c/rong/v5/internal/config"
)

// GenerateStatic generates base16 colors from pre-defined colors.
func GenerateStatic(primary color.ARGB, wu []color.ARGB) Base16 {
	ratio := config.Base16Blend.Value()

	blend := makeBlendFunc(ratio, primary, wu)
	black := blend(config.Base16Black.Value())
	red := blend(config.Base16Red.Value())
	green := blend(config.Base16Green.Value())
	yellow := blend(config.Base16Yellow.Value())
	blue := blend(config.Base16Blue.Value())
	magenta := blend(config.Base16Magenta.Value())
	cyan := blend(config.Base16Cyan.Value())
	white := blend(config.Base16White.Value())

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

func makeBlendFunc(ratio float64, primary color.ARGB, wu []color.ARGB) func(color.ARGB) color.Hct {
	if ratio <= 0 {
		return func(c color.ARGB) color.Hct { return c.ToHct() }
	}
	if ratio >= 1 {
		return func(color.ARGB) color.Hct { return primary.ToHct() }
	}

	dst := primary.ToOkLab()
	if len(wu) == 0 {
		return func(c color.ARGB) color.Hct {
			src := c.ToOkLab()
			return blend(src, dst, ratio).ToXYZ().ToHct()
		}
	}

	colors := make([]color.OkLab, 0, len(wu)+1)
	colors = append(colors, primary.ToOkLab())
	for c := range slices.Values(wu) {
		colors = append(colors, c.ToOkLab())
	}

	return func(c color.ARGB) color.Hct {
		src := c.ToOkLab()
		lowest := slices.MinFunc(colors, func(a, b color.OkLab) int {
			return cmp.Compare(okLabDistance(src, a), okLabDistance(src, b))
		})
		return blend(src, lowest, ratio).ToXYZ().ToHct()
	}
}

func blend(a, b color.OkLab, ratio float64) color.OkLab {
	return color.OkLab{
		L: a.L + (b.L-a.L)*ratio,
		A: a.A + (b.A-a.A)*ratio,
		B: a.B + (b.B-a.B)*ratio,
	}
}

func okLabDistance(a, b color.OkLab) float64 {
	dL := a.L - b.L
	dA := a.A - b.A
	dB := a.B - b.B
	return math.Sqrt(dL*dL + dA*dA + dB*dB)
}
