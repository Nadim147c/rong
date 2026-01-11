package templates

import (
	"encoding/json"
	"fmt"
	"maps"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/Nadim147c/material/v2/color"
	"github.com/Nadim147c/rong/v5/internal/models"
)

var funcs = template.FuncMap{
	"upper":   strings.ToUpper,
	"lower":   strings.ToLower,
	"replace": strings.ReplaceAll,
	"parse":   parse,
	"chroma":  chroma,
	"tone":    tone,
	"quote":   quote,
	"json":    jsonString,
}

func init() {
	maps.Copy(funcs, sprig.TxtFuncMap())
}

func parse(a any) models.FormatedColor {
	switch a := a.(type) {
	case models.FormatedColor:
		return a
	case models.NamedColor:
		return a.Color
	case color.Hct:
		return models.NewFormatedColor(a.ToARGB())
	case color.Lab:
		return models.NewFormatedColor(a.ToARGB())
	case color.OkLab:
		return models.NewFormatedColor(a.ToARGB())
	case color.XYZ:
		return models.NewFormatedColor(a.ToARGB())
	case string:
		c := color.ARGBFromHexMust(a)
		return models.NewFormatedColor(c)
	default:
		panic(fmt.Sprintf("invalid color format: %v", a))
	}
}

func chroma(c any, chroma float64) models.FormatedColor {
	hct := parse(c).Int.ToHct()
	hct.Chroma = chroma
	return models.NewFormatedColor(hct.ToARGB())
}

func tone(c any, t float64) models.FormatedColor {
	hct := parse(c).Int.ToHct()
	hct.Tone = t
	return models.NewFormatedColor(hct.ToARGB())
}

func quote(s any) string {
	return fmt.Sprintf("%q", s)
}

func jsonString(s any) string {
	bytes, err := json.Marshal(s)
	if err != nil {
		return "null"
	}
	return string(bytes)
}
