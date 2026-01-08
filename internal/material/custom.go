package material

import (
	"fmt"
	"log/slog"
	"regexp"

	blendPkg "github.com/Nadim147c/material/v2/blend"
	"github.com/Nadim147c/material/v2/color"
	"github.com/Nadim147c/rong/v5/internal/config"
)

// CustomColor is colors generated from user defined colors.
type CustomColor struct {
	Color            color.ARGB
	OnColor          color.ARGB
	ColorContainer   color.ARGB
	OnColorContainer color.ARGB
}

var nameRe = regexp.MustCompile("^[A-Za-z0-9][A-Za-z0-9_]+$")

// GenerateCustomColors returns all custom colors.
func GenerateCustomColors(primary color.ARGB) (map[string]CustomColor, error) {
	defined := config.MaterialCustomColors.Value()
	if len(defined) == 0 {
		return map[string]CustomColor{}, nil
	}
	dark := config.Dark.Value()
	blend := config.MaterialCustomBlend.Value()

	m := make(map[string]CustomColor, len(defined))
	for name, col := range defined {
		if !nameRe.MatchString(name) {
			//nolint
			return nil, fmt.Errorf("custom color name should only contains alphanumeric values or underscore: name=%s", name)
		}
		slog.Debug("Making custom color", "name", name, "color", col.String())
		m[name] = createCustomColor(col, primary, dark, blend)
	}
	return m, nil
}

func createCustomColor(src, to color.ARGB, dark bool, blend float64) CustomColor {
	var hct color.Hct

	if blend != 0 {
		hct = blendPkg.HctHueDirect(src, to, blend)
	} else {
		hct = src.ToHct()
	}

	if dark {
		return CustomColor{
			Color:            tone(hct, 40),
			OnColor:          tone(hct, 100),
			ColorContainer:   tone(hct, 90),
			OnColorContainer: tone(hct, 10),
		}
	}
	return CustomColor{
		Color:            tone(hct, 80),
		OnColor:          tone(hct, 20),
		ColorContainer:   tone(hct, 30),
		OnColorContainer: tone(hct, 90),
	}
}

// tone creates an ARGB color from hct color with given tone (brightness).
func tone(hct color.Hct, tone float64) color.ARGB {
	hct.Tone = tone
	return hct.ToARGB()
}
