package material

import (
	"fmt"
	"regexp"

	blendPkg "github.com/Nadim147c/material/v2/blend"
	"github.com/Nadim147c/material/v2/color"
	"github.com/spf13/viper"
)

// CustomColor is colors generated from user defined colors.
type CustomColor struct {
	Color            color.ARGB
	OnColor          color.ARGB
	ColorContainer   color.ARGB
	OnColorContainer color.ARGB
}

var nameRe = regexp.MustCompile("^([A-Za-z0-9_])+$")

func init() {
	viper.SetDefault("material.custom.blend", true)
	viper.SetDefault("material.custom.ratio", 0.35)
}

// GenerateCustomColors returns all custom colors.
func GenerateCustomColors(primary color.ARGB) (map[string]CustomColor, error) {
	defined := viper.GetStringMapString("material.custom.colors")
	if len(defined) == 0 {
		return map[string]CustomColor{}, nil
	}
	m := make(map[string]CustomColor, len(defined))
	dark := viper.GetBool("dark")
	blend := viper.GetBool("material.custom.blend")
	ratio := viper.GetFloat64("material.custom.ratio")

	for name, col := range defined {
		if !nameRe.MatchString(name) {
			return nil, fmt.Errorf( //nolint
				"custom color name should only contains alphanumeric values or underscore: name=%s",
				name,
			)
		}
		argb, err := color.ARGBFromHex(col)
		if err != nil {
			return nil, err
		}
		m[name] = createCustomColor(argb, primary, dark, blend, ratio)
	}
	return m, nil
}

func createCustomColor(
	src, to color.ARGB,
	dark, blend bool,
	ratio float64,
) CustomColor {
	var hct color.Hct

	if blend {
		hct = blendPkg.HctHueDirect(src, to, ratio)
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
