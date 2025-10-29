package base16

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/Nadim147c/material/color"
	"github.com/charmbracelet/lipgloss"
	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// BlendRatio is the default blend ratio from static color
const BlendRatio float64 = 0.5

// Flags are the flags used for generating colors
var Flags = pflag.NewFlagSet("base16", pflag.ContinueOnError)

func formatDefault(c color.ARGB) string {
	hex := c.HexRGB()
	return lipgloss.NewStyle().Foreground(lipgloss.Color(hex)).Render(hex)
}

func init() {
	Flags.Float64(
		"base16.blend",
		BlendRatio,
		"blend ratio toward the primary color",
	)
	viper.SetDefault("base16.blend", BlendRatio)

	Flags.String(
		"base16.method",
		"static",
		"color generation method (static or dynamic)",
	)
	viper.SetDefault("base16.method", "static")

	Flags.String(
		"base16.colors.black",
		formatDefault(defaultSrcColors.Black),
		"black source color for base16 colors",
	)
	viper.SetDefault("base16.colors.black", defaultSrcColors.Black.HexRGB())

	Flags.String(
		"base16.colors.red",
		formatDefault(defaultSrcColors.Red),
		"red source color for base16 colors",
	)
	viper.SetDefault("base16.colors.red", defaultSrcColors.Red.HexRGB())

	Flags.String(
		"base16.colors.green",
		formatDefault(defaultSrcColors.Green),
		"green source color for base16 colors",
	)
	viper.SetDefault("base16.colors.green", defaultSrcColors.Green.HexRGB())

	Flags.String(
		"base16.colors.yellow",
		formatDefault(defaultSrcColors.Yellow),
		"yellow source color for base16 colors",
	)
	viper.SetDefault("base16.colors.yellow", defaultSrcColors.Yellow.HexRGB())

	Flags.String(
		"base16.colors.blue",
		formatDefault(defaultSrcColors.Blue),
		"blue source color for base16 colors",
	)
	viper.SetDefault("base16.colors.blue", defaultSrcColors.Blue.HexRGB())

	Flags.String(
		"base16.colors.magenta",
		formatDefault(defaultSrcColors.Magenta),
		"magenta source color for base16 colors",
	)
	viper.SetDefault("base16.colors.magenta", defaultSrcColors.Magenta.HexRGB())

	Flags.String(
		"base16.colors.cyan",
		formatDefault(defaultSrcColors.Cyan),
		"cyan source color for base16 colors",
	)
	viper.SetDefault("base16.colors.cyan", defaultSrcColors.Cyan.HexRGB())

	Flags.String(
		"base16.colors.white",
		formatDefault(defaultSrcColors.White),
		"white source color for base16 colors",
	)
	viper.SetDefault("base16.colors.white", defaultSrcColors.White.HexRGB())
}

var opt = viper.DecodeHook(mapstructure.DecodeHookFuncValue(
	func(from, _ reflect.Value) (any, error) {
		if from.Kind() == reflect.String {
			return color.ARGBFromHex(from.String())
		}
		return from.Interface(), nil
	},
))

// Generate generates colors from material color name and quantized colors
func Generate(
	material map[string]color.ARGB,
	quantized []color.ARGB,
) (Base16, error) {
	var static SourceColors
	if err := viper.UnmarshalKey("base16.colors", &static, opt); err != nil {
		return Base16{}, err
	}

	switch method := strings.ToLower(viper.GetString("base16.method")); method {
	case "static":
		return GenerateStatic(material["primary"], static), nil
	case "dynamic":
		fg, bg := material["on_background"], material["background"]
		return GenerateDynamic(fg, bg, quantized), nil
	default:
		return Base16{}, fmt.Errorf(
			"invalid base16 color generating method: %v",
			method,
		)
	}
}

// Base16 is the generated output
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

// NewBase16 creates a new Base16
func NewBase16() Base16 {
	b := Base16{}
	b.dark = viper.GetBool("dark")
	return b
}

// SetBlack sets the Black and Bright Black color
func (b *Base16) SetBlack(c color.Hct) {
	b.Black, b.BrightBlack = fixBlack(b.dark, c)
}

// SetRed sets the Red and Bright Red color
func (b *Base16) SetRed(c color.Hct) {
	b.Red, b.BrightRed = fix(b.dark, c)
}

// SetGreen sets the Green and Bright Green color
func (b *Base16) SetGreen(c color.Hct) {
	b.Green, b.BrightGreen = fix(b.dark, c)
}

// SetYellow sets the Yellow and Bright Yellow color
func (b *Base16) SetYellow(c color.Hct) {
	b.Yellow, b.BrightYellow = fix(b.dark, c)
}

// SetBlue sets the Blue and Bright Blue color
func (b *Base16) SetBlue(c color.Hct) {
	b.Blue, b.BrightBlue = fix(b.dark, c)
}

// SetMagenta sets the Magenta and Bright Magenta color
func (b *Base16) SetMagenta(c color.Hct) {
	b.Magenta, b.BrightMagenta = fix(b.dark, c)
}

// SetCyan sets the Cyan and Bright Cyan color
func (b *Base16) SetCyan(c color.Hct) {
	b.Cyan, b.BrightCyan = fix(b.dark, c)
}

// SetWhite sets the White and Bright White color
func (b *Base16) SetWhite(c color.Hct) {
	b.White, b.BrightWhite = fixWhite(b.dark, c)
}

const (
	darkChroma        = 80.0
	lightChroma       = 100.0
	darkChromaScaled  = 10.0
	lightChromaScaled = 15.0

	toneDark    = 95.0
	toneDarkHi  = 100.0
	toneLight   = 5.0
	toneLightHi = 15.0

	fixToneDark    = 50.0
	fixToneDarkHi  = 70.0
	fixToneLight   = 35.0
	fixToneLightHi = 25.0
)

func setToneChroma(c color.Hct, tone float64, chroma float64) color.ARGB {
	c.Tone = tone
	c.Chroma = chroma
	return c.ToARGB()
}

// shared func for fg/bg
func fixBlackWhite(dark, invert bool, c color.Hct) (color.ARGB, color.ARGB) {
	if dark != invert {
		return setToneChroma(c, toneDark, darkChromaScaled),
			setToneChroma(c, toneDarkHi, darkChromaScaled)
	}
	return setToneChroma(c, toneLight, lightChromaScaled),
		setToneChroma(c, toneLightHi, lightChromaScaled)
}

func fixWhite(dark bool, c color.Hct) (color.ARGB, color.ARGB) {
	return fixBlackWhite(dark, false, c)
}

func fixBlack(dark bool, c color.Hct) (color.ARGB, color.ARGB) {
	return fixBlackWhite(dark, true, c)
}

func fix(dark bool, c color.Hct) (color.ARGB, color.ARGB) {
	if dark {
		return setToneChroma(c, fixToneDark, darkChroma),
			setToneChroma(c, fixToneDarkHi, darkChroma)
	}
	return setToneChroma(c, fixToneLight, lightChroma),
		setToneChroma(c, fixToneLightHi, lightChroma)
}
