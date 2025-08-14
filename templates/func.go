package templates

import (
	"encoding/json"
	"fmt"
	"strings"
	"text/template"

	"github.com/Nadim147c/material/color"
	"github.com/Nadim147c/rong/internal/models"
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

func parse(s string) models.ColorValue {
	c := color.ARGBFromHexMust(s)
	return models.NewColorValue(c)
}

func chroma(c models.ColorValue, chroma float64) models.ColorValue {
	hct := c.Int.ToHct()
	hct.Chroma = chroma
	return models.NewColorValue(hct.ToARGB())
}

func tone(c models.ColorValue, t float64) models.ColorValue {
	hct := c.Int.ToHct()
	hct.Tone = t
	return models.NewColorValue(hct.ToARGB())
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
