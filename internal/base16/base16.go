package base16

import (
	"github.com/Nadim147c/material/v2/color"
	"github.com/Nadim147c/rong/v5/internal/config"
	"github.com/Nadim147c/rong/v5/internal/config/enums"
)

// Generate generates colors from material color name and quantized colors.
func Generate(
	material map[string]color.ARGB,
	quantized []color.ARGB,
) (Base16, error) {
	switch config.Base16Method.Value() {
	case enums.Base16MethodStatic:
		return GenerateStatic(material["primary"]), nil
	case enums.Base16MethodDynamic:
		fg, bg := material["on_background"], material["background"]
		return GenerateDynamic(fg, bg, quantized), nil
	default:
		panic("unreachable")
	}
}

// Base16 is the generated output.
type Base16 struct {
	dark                   bool
	Black, BrightBlack     color.ARGB
	Red, BrightRed         color.ARGB
	Green, BrightGreen     color.ARGB
	Yellow, BrightYellow   color.ARGB
	Blue, BrightBlue       color.ARGB
	Magenta, BrightMagenta color.ARGB
	Cyan, BrightCyan       color.ARGB
	White, BrightWhite     color.ARGB
}

// NewBase16 creates a new Base16.
func NewBase16() Base16 {
	return Base16{dark: config.Dark.Value()}
}

// SetBlack sets the Black and Bright Black color.
func (b *Base16) SetBlack(c color.Hct) {
	b.Black, b.BrightBlack = fixBlack(b.dark, c)
}

// SetRed sets the Red and Bright Red color.
func (b *Base16) SetRed(c color.Hct) {
	b.Red, b.BrightRed = fix(b.dark, c)
}

// SetGreen sets the Green and Bright Green color.
func (b *Base16) SetGreen(c color.Hct) {
	b.Green, b.BrightGreen = fix(b.dark, c)
}

// SetYellow sets the Yellow and Bright Yellow color.
func (b *Base16) SetYellow(c color.Hct) {
	b.Yellow, b.BrightYellow = fix(b.dark, c)
}

// SetBlue sets the Blue and Bright Blue color.
func (b *Base16) SetBlue(c color.Hct) {
	b.Blue, b.BrightBlue = fix(b.dark, c)
}

// SetMagenta sets the Magenta and Bright Magenta color.
func (b *Base16) SetMagenta(c color.Hct) {
	b.Magenta, b.BrightMagenta = fix(b.dark, c)
}

// SetCyan sets the Cyan and Bright Cyan color.
func (b *Base16) SetCyan(c color.Hct) {
	b.Cyan, b.BrightCyan = fix(b.dark, c)
}

// SetWhite sets the White and Bright White color.
func (b *Base16) SetWhite(c color.Hct) {
	b.White, b.BrightWhite = fixWhite(b.dark, c)
}

// Color adjustment parameters for theme variants.
const (
	// Chroma value for dark theme colors.
	chromaDark = 80.0
	// Chroma value for light theme colors.
	chromaLight = 100.0

	// Reduced chroma for dark theme black/white fixes.
	chromaDarkMuted = 10.0
	// Reduced chroma for light theme black/white fixes.
	chromaLightMuted = 15.0

	// Dark theme black adjustments.
	toneNearWhiteDark = 15.0
	tonePureWhiteDark = 25.0

	// Dark theme white adjustments.
	toneNearBlackDark = 95.0
	tonePureBlackDark = 100.0

	// Light theme white adjustments.
	toneNearWhiteLight = 75.0
	tonePureWhiteLight = 85.0

	// Light theme black adjustments.
	toneNearBlackLight = 25.0
	tonePureBlackLight = 35.0

	// Base tone for dark theme general fixes.
	toneDarkBase = 65.0
	// Bright tone for dark theme general fixes.
	toneDarkBright = 80.0

	// Base tone for light theme general fixes.
	toneLightBase = 50.0
	// Bright tone for light theme general fixes.
	toneLightBright = 65.0
)

// setToneChroma returns a new ARGB color with the specified tone and chroma
// values.
func setToneChroma(c color.Hct, tone float64, chroma float64) color.ARGB {
	c.Tone = tone
	c.Chroma = chroma
	if c.IsBlue() {
		c.Tone *= 0.8
		c.Chroma *= 1.1
	}
	return c.ToARGB()
}

// fixBlackWhite returns a pair of colors adjusted for black or white elements
// based on the theme and inversion settings.
func fixBlackWhite(dark, white bool, c color.Hct) (color.ARGB, color.ARGB) {
	if dark {
		if white {
			return setToneChroma(c, toneNearWhiteDark, chromaDarkMuted),
				setToneChroma(c, tonePureWhiteDark, chromaDarkMuted)
		}
		return setToneChroma(c, toneNearBlackDark, chromaDarkMuted),
			setToneChroma(c, tonePureBlackDark, chromaDarkMuted)
	}
	if white {
		return setToneChroma(c, toneNearWhiteLight, chromaLightMuted),
			setToneChroma(c, tonePureWhiteLight, chromaLightMuted)
	}
	return setToneChroma(c, toneNearBlackLight, chromaLightMuted),
		setToneChroma(c, tonePureBlackLight, chromaLightMuted)
}

// fixWhite returns a pair of colors adjusted for white elements in the given
// theme.
func fixWhite(dark bool, c color.Hct) (color.ARGB, color.ARGB) {
	return fixBlackWhite(dark, false, c)
}

// fixBlack returns a pair of colors adjusted for black elements in the given
// theme.
func fixBlack(dark bool, c color.Hct) (color.ARGB, color.ARGB) {
	return fixBlackWhite(dark, true, c)
}

// fix returns a pair of colors with general theme-appropriate adjustments.
func fix(dark bool, c color.Hct) (color.ARGB, color.ARGB) {
	if dark {
		return setToneChroma(c, toneDarkBase, chromaDark),
			setToneChroma(c, toneDarkBright, chromaDark)
	}
	return setToneChroma(c, toneLightBase, chromaLight),
		setToneChroma(c, toneLightBright, chromaLight)
}
