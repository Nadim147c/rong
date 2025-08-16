package base16

import (
	"math"
	"math/rand"
	"sort"

	"github.com/Nadim147c/material/color"
	"github.com/spf13/viper"
)

// distance returns the shortest distance between two angles on a circle
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

func ensureHueVariety(colors []color.Hct, needed int) []color.Hct {
	if len(colors) == 0 {
		return []color.Hct{randomHct()}
	}
	out := make([]color.Hct, len(colors))
	copy(out, colors)

	spread := hueSpread(colors)
	if spread < 60 {
		base := colors[0]
		// Inject complementary & triadic hues
		for _, offset := range []float64{120, 180, 240} {
			if len(out) >= needed {
				break
			}
			rotated := base
			rotated.Hue = math.Mod(base.Hue+offset, 360.0)
			out = append(out, rotated)
		}
	}
	return out
}

// SelectColors selects k Hct colors maximizing angular separation
func SelectColors(colors []color.Hct, k int) []color.Hct {
	n := len(colors)
	if k >= n {
		return colors
	}

	selected := []color.Hct{colors[rand.Intn(n)]}
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

// Generate generates base16 colors from selecting quantizes color. It takes
// color with long chroma distance to ensure colors has more variety
func Generate(fg, bg color.ARGB, colors []color.ARGB) map[string]color.ARGB {
	hct := make([]color.Hct, len(colors))
	for i, v := range colors {
		hct[i] = v.ToHct()
	}

	selected := SelectColors(ensureHueVariety(hct, 10), 10)

	dark := viper.GetBool("dark")
	b := map[string]color.ARGB{}
	b["color_0"], b["color_8"] = fixbg(dark, bg.ToHct())
	b["color_1"], b["color_9"] = fix(dark, getColorWithRandFallback(selected, 0))
	b["color_2"], b["color_a"] = fix(dark, getColorWithRandFallback(selected, 1))
	b["color_3"], b["color_b"] = fix(dark, getColorWithRandFallback(selected, 2))
	b["color_4"], b["color_c"] = fix(dark, getColorWithRandFallback(selected, 3))
	b["color_5"], b["color_d"] = fix(dark, getColorWithRandFallback(selected, 4))
	b["color_6"], b["color_e"] = fix(dark, getColorWithRandFallback(selected, 5))
	b["color_7"], b["color_f"] = fixfg(dark, fg.ToHct())

	return b
}

func getColorWithRandFallback(colors []color.Hct, i int) color.Hct {
	if i < len(colors) {
		return colors[i]
	}
	return randomHct()
}

func randomHct() color.Hct {
	return color.Hct{
		Hue:    rand.Float64() * 360.0,
		Chroma: 40.0 + rand.Float64()*40.0, // Moderate to high chroma
		Tone:   30.0 + rand.Float64()*40.0, // Middle tone range
	}
}

func tc(c color.Hct, tone float64, chroma float64) color.ARGB {
	c.Tone = tone
	c.Chroma = chroma
	return c.ToARGB()
}

const (
	darkChroma  float64 = 80
	lightChroma float64 = 100
)

func fixfg(dark bool, c color.Hct) (color.ARGB, color.ARGB) {
	if dark {
		return tc(c, 100, darkChroma), tc(c, 90, darkChroma)
	}
	return tc(c, 20, lightChroma), tc(c, 35, lightChroma)
}

func fixbg(dark bool, c color.Hct) (color.ARGB, color.ARGB) {
	if dark {
		return tc(c, 20, darkChroma), tc(c, 35, darkChroma)
	}
	return tc(c, 100, lightChroma), tc(c, 90, lightChroma)
}

func fix(dark bool, c color.Hct) (color.ARGB, color.ARGB) {
	if dark {
		return tc(c, 50, darkChroma), tc(c, 70, darkChroma)
	}
	return tc(c, 35, lightChroma), tc(c, 25, lightChroma)
}

// GenerateRandom generate random base16 colors with given fg,bg and dark
func GenerateRandom(fg, bg color.ARGB) map[string]color.ARGB {
	dark := viper.GetBool("viper")

	b := map[string]color.ARGB{}
	b["color_0"], b["color_8"] = fixbg(dark, bg.ToHct())

	for _, code := range [6][2]string{
		{"color_1", "color_9"},
		{"color_2", "color_a"},
		{"color_3", "color_b"},
		{"color_4", "color_c"},
		{"color_5", "color_d"},
		{"color_6", "color_e"},
	} {
		h := randomHct()
		b[code[0]], b[code[1]] = fix(dark, h)
	}

	b["color_7"], b["color_f"] = fixfg(dark, fg.ToHct())

	return b
}
