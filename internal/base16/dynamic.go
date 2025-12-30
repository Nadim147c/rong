package base16

import (
	"log/slog"
	"math"
	"math/rand"
	"sort"

	"github.com/Nadim147c/material/v2/color"
	"github.com/Nadim147c/material/v2/num"
)

// GenerateDynamic generates base16 colors from selecting quantizes color. It
// takes color with long chroma distance to ensure colors has more variety.
func GenerateDynamic(
	fg, bg color.ARGB,
	colors []color.ARGB,
) Base16 {
	hct := make([]color.Hct, len(colors))
	for i, v := range colors {
		hct[i] = v.ToHct()
	}

	selected := SelectColors(ensureHueVariety(hct), 6)

	based := NewBase16()
	based.SetWhite(fg.ToHct())
	based.SetBlack(bg.ToHct())

	based.SetRed(getColor(selected, 0))
	based.SetGreen(getColor(selected, 1))
	based.SetYellow(getColor(selected, 2))
	based.SetBlue(getColor(selected, 3))
	based.SetMagenta(getColor(selected, 4))
	based.SetCyan(getColor(selected, 5))

	return based
}

func getColor(colors []color.Hct, i int) color.Hct {
	if i < len(colors) {
		return colors[i]
	}
	return color.Hct{}
}

// distance returns the shortest distance between two angles on a circle.
func distance(a, b float64) float64 {
	diff := math.Abs(a - b)
	return math.Min(diff, 360-diff)
}

func hueSpread(colors []color.Hct) float64 {
	if len(colors) < 2 {
		return 0
	}
	minHue, maxHue := 360.0, 0.0
	for _, c := range colors {
		if c.Hue < minHue {
			minHue = c.Hue
		}
		if c.Hue > maxHue {
			maxHue = c.Hue
		}
	}
	return maxHue - minHue
}

func ensureHueVariety(colors []color.Hct) []color.Hct {
	if len(colors) == 0 {
		return []color.Hct{defaultSrcColors.Red.ToHct()}
	}

	out := make([]color.Hct, len(colors))
	copy(out, colors)

	var i int
	deg := 80.0
	for hueSpread(out) < 100 {
		base := colors[i]
		newHue := num.NormalizeDegree(base.Hue + deg)
		newColor := base
		newColor.Hue = newHue
		slog.Info("Generating random color", "hue", newHue)
		out = append(out, newColor)
		i = (i + 1) % len(colors)
		deg *= 2
	}

	return out
}

// SelectColors selects k Hct colors maximizing angular separation.
func SelectColors(colors []color.Hct, k int) []color.Hct {
	n := len(colors)
	if k >= n {
		return colors
	}

	selected := []color.Hct{colors[rand.Intn(n)]} //nolint:gosec
	remaining := make([]color.Hct, 0, n-1)
	for _, c := range colors {
		if c != selected[0] {
			remaining = append(remaining, c)
		}
	}

	for len(selected) < k {
		var best color.Hct
		bestMinDist := -1.0

		for _, candidate := range remaining {
			minDist := math.MaxFloat64
			for _, s := range selected {
				d := distance(candidate.Hue, s.Hue)
				if d < minDist {
					minDist = d
				}
			}
			if minDist > bestMinDist {
				bestMinDist = minDist
				best = candidate
			}
		}

		selected = append(selected, best)

		// Remove best from remaining
		newRemaining := remaining[:0]
		for _, c := range remaining {
			if c != best {
				newRemaining = append(newRemaining, c)
			}
		}
		remaining = newRemaining
	}

	sort.Slice(selected, func(i, j int) bool {
		return selected[i].Hue < selected[j].Hue
	})

	return selected
}
