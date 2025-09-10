package material

import (
	"errors"

	"github.com/Nadim147c/material/color"
	"github.com/Nadim147c/material/dynamic"
	"github.com/Nadim147c/material/palettes"
	"github.com/Nadim147c/material/quantizer"
	"github.com/Nadim147c/material/score"
)

// Quantized is quantized colors
type Quantized struct {
	Celebi map[color.ARGB]int `json:"celebi"`
	Wu     []color.ARGB       `json:"wu"`
}

// Quantize quantizes list of pixels
func Quantize(pixels []color.ARGB) Quantized {
	wu := quantizer.QuantizeWu(pixels, 100)
	colors := make([]color.Lab, len(wu))
	for i, c := range wu {
		colors[i] = c.ToLab()
	}

	celebi := quantizer.QuantizeWsMeans(pixels, colors, 4)
	return Quantized{celebi, wu}
}

// ErrNoColorFound is a error
var ErrNoColorFound = errors.New("no color found")

// Colors is key and color
type Colors = map[string]color.ARGB

// GenerateFromPixels generates color from a slice of pixels
func GenerateFromPixels(pixels []color.ARGB, cfg Config) (Colors, []color.ARGB, error) {
	q := Quantize(pixels)
	return GenerateFromQuantized(q, cfg)
}

// GenerateFromQuantized generates color from a cached quantized
func GenerateFromQuantized(quantized Quantized, cfg Config) (Colors, []color.ARGB, error) {
	celebi, wu := quantized.Celebi, quantized.Wu

	scored := score.Score(celebi, score.ScoreOptions{
		Desired:  4,
		Fallback: score.FallbackColor,
	})

	if len(scored) == 0 {
		return Colors{}, wu, ErrNoColorFound
	}

	primary := palettes.NewFromARGB(scored[0])

	scheme := dynamic.NewDynamicScheme(
		scored[0].ToHct(),
		cfg.Variant, cfg.Constrast, cfg.Dark,
		cfg.Platform, cfg.Version, primary,
		nil, nil, nil, nil, nil,
	)

	dcs := scheme.ToColorMap()

	colorMap := map[string]color.ARGB{}
	for key, value := range dcs {
		if value != nil {
			colorMap[key] = value.GetArgb(scheme)
		}
	}
	return colorMap, wu, nil
}
