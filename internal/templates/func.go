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

func parse(s string) models.FormatedColor {
	c := color.ARGBFromHexMust(s)
	return models.NewFormatedColor(c)
}

func chroma(c models.FormatedColor, chroma float64) models.FormatedColor {
	hct := c.Int.ToHct()
	hct.Chroma = chroma
	return models.NewFormatedColor(hct.ToARGB())
}

func tone(c models.FormatedColor, t float64) models.FormatedColor {
	hct := c.Int.ToHct()
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
