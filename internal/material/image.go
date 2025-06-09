package material

import (
	"errors"
	"image"

	"github.com/Nadim147c/goyou/color"
	"github.com/Nadim147c/goyou/dynamic"
	"github.com/Nadim147c/goyou/palettes"
	"github.com/Nadim147c/goyou/quantizer"
	"github.com/Nadim147c/goyou/score"
)

var ErrNoColorFound = errors.New("no color found")

// Colors is key and color
type Colors = map[string]color.ARGB

// GenerateColorsFromImage colors from an image.Image
func GenerateColorsFromImage(img image.Image, isDark bool) (Colors, error) {
	bounds := img.Bounds()
	pixels := make([]color.ARGB, 0, bounds.Dx()*bounds.Dy())

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.At(x, y)
			argb := color.ARGBFromInterface(c)
			pixels = append(pixels, argb)
		}
	}

	quantized := quantizer.QuantizeCelebi(pixels, 100)
	if len(quantized) == 0 {
		return Colors{}, ErrNoColorFound
	}

	scored := score.Score(quantized, score.ScoreOptions{Desired: 4, Fallback: score.FallbackColor})
	if len(scored) == 0 {
		return Colors{}, ErrNoColorFound
	}

	var primary, secondary, ternary *palettes.TonalPalette
	primary = palettes.NewFromARGB(scored[0])

	if len(scored) > 1 {
		secondary = palettes.NewFromARGB(scored[1])
	}
	if len(scored) > 2 {
		ternary = palettes.NewFromARGB(scored[2])
	}

	scheme := dynamic.NewDynamicScheme(
		scored[0].ToHct(), dynamic.Expressive, 0, isDark,
		dynamic.Phone, dynamic.V2021,
		primary, secondary, ternary,
		nil, nil, nil,
	)

	dcs := scheme.ToColorMap()

	colorMap := map[string]color.ARGB{}
	for key, value := range dcs {
		if value != nil {
			colorMap[key] = value.GetArgb(scheme)
		}
	}
	return colorMap, nil
}
