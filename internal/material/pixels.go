package material

import (
	"errors"

	"github.com/Nadim147c/material/color"
	"github.com/Nadim147c/material/dynamic"
	"github.com/Nadim147c/material/palettes"
	"github.com/Nadim147c/material/quantizer"
	"github.com/Nadim147c/material/score"
)

// ErrNoColorFound is a error
var ErrNoColorFound = errors.New("no color found")

// Colors is key and color
type Colors = map[string]color.ARGB

// GenerateFromPixels generates color from a slice of pixels
func GenerateFromPixels(
	pixels []color.ARGB,
	variant dynamic.Variant,
	dark bool,
	constrast float64,
	platform dynamic.Platform,
	version dynamic.Version,
) (Colors, []color.ARGB, error) {
	quantizedWu := quantizer.QuantizeWu(pixels, 100)
	colors := make([]color.Lab, len(quantizedWu))
	for i, c := range quantizedWu {
		colors[i] = c.ToLab()
	}

	quantized := quantizer.QuantizeWsMeans(pixels, colors, 4)
	if len(quantized) == 0 {
		return Colors{}, quantizedWu, ErrNoColorFound
	}

	scored := score.Score(quantized, score.ScoreOptions{Desired: 4, Fallback: score.FallbackColor})
	if len(scored) == 0 {
		return Colors{}, quantizedWu, ErrNoColorFound
	}

	primary := palettes.NewFromARGB(scored[0])

	scheme := dynamic.NewDynamicScheme(
		scored[0].ToHct(), variant, constrast, dark,
		platform, version, primary,
		nil, nil, nil, nil, nil,
	)

	dcs := scheme.ToColorMap()

	colorMap := map[string]color.ARGB{}
	for key, value := range dcs {
		if value != nil {
			colorMap[key] = value.GetArgb(scheme)
		}
	}
	return colorMap, quantizedWu, nil
}
