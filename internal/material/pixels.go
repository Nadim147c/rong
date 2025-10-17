package material

import (
	"context"
	"errors"

	"github.com/Nadim147c/material/color"
	"github.com/Nadim147c/material/dynamic"
	"github.com/Nadim147c/material/quantizer"
	"github.com/Nadim147c/material/score"
)

// Quantized is quantized colors
type Quantized struct {
	Celebi map[color.ARGB]int `json:"celebi"`
	Wu     []color.ARGB       `json:"wu"`
}

// Quantize quantizes list of pixels
func Quantize(ctx context.Context, pixels []color.ARGB) (Quantized, error) {
	wu, err := quantizer.QuantizeWuContext(ctx, pixels, 100)
	if err != nil {
		return Quantized{}, err
	}

	colors := make([]color.Lab, len(wu))
	for i, c := range wu {
		colors[i] = c.ToLab()
	}

	celebi, err := quantizer.QuantizeWsMeansContext(ctx, pixels, colors, 4)
	if err != nil {
		return Quantized{}, err
	}

	return Quantized{celebi, wu}, nil
}

// ErrNoColorFound is a error
var ErrNoColorFound = errors.New("no color found")

// Colors is key and color
type Colors = map[string]color.ARGB

// GenerateFromPixels generates color from a slice of pixels
func GenerateFromPixels(
	ctx context.Context,
	pixels []color.ARGB,
	cfg Config,
) (Colors, []color.ARGB, error) {
	q, err := Quantize(ctx, pixels)
	if err != nil {
		return nil, nil, err
	}

	return GenerateFromQuantized(q, cfg)
}

// GenerateFromQuantized generates color from a cached quantized
func GenerateFromQuantized(
	quantized Quantized,
	cfg Config,
) (Colors, []color.ARGB, error) {
	celebi, wu := quantized.Celebi, quantized.Wu

	scored := score.Score(celebi, score.WithFilter())

	if len(scored) == 0 {
		return nil, wu, ErrNoColorFound
	}

	scheme := dynamic.NewDynamicScheme(
		scored[0].ToHct(),
		cfg.Variant, cfg.Constrast, cfg.Dark,
		cfg.Platform, cfg.Version,
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
