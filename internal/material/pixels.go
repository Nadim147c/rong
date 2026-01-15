package material

import (
	"context"
	"errors"

	"github.com/Nadim147c/material/v2/color"
	"github.com/Nadim147c/material/v2/dynamic"
	"github.com/Nadim147c/material/v2/quantizer"
	"github.com/Nadim147c/material/v2/score"
)

// Quantized is quantized colors.
type Quantized struct {
	Celebi map[color.ARGB]int `json:"celebi"`
	Wu     []color.ARGB       `json:"wu"`
}

// Quantize quantizes list of pixels.
func Quantize(ctx context.Context, pixels []color.ARGB) (Quantized, error) {
	wu, err := quantizer.QuantizeWuContext(ctx, pixels, 100)
	if err != nil {
		return Quantized{}, err
	}

	colors := make([]color.Lab, len(wu))
	for i, c := range wu {
		colors[i] = c.ToLab()
	}

	celebi, err := quantizer.QuantizeWsMeansContext(ctx, pixels, colors, 10)
	if err != nil {
		return Quantized{}, err
	}

	return Quantized{celebi, wu}, nil
}

// ErrNoColorFound means no color found from imput image.
var ErrNoColorFound = errors.New("no color found")

// Colors is key and color.
type Colors = map[string]color.ARGB

// GenerateFromPixels generates color from a slice of pixels.
func GenerateFromPixels(
	ctx context.Context,
	pixels []color.ARGB,
	cfg Config,
) (Colors, error) {
	q, err := Quantize(ctx, pixels)
	if err != nil {
		return nil, err
	}

	return GenerateFromQuantized(q, cfg)
}

// GenerateFromQuantized generates color from a cached quantized.
func GenerateFromQuantized(
	quantized Quantized,
	cfg Config,
) (Colors, error) {
	celebi := quantized.Celebi

	scored := score.Score(celebi, score.WithFilter())

	if len(scored) == 0 {
		return nil, ErrNoColorFound
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
	return colorMap, nil
}
