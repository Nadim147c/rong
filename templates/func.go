package templates

import (
	"encoding/json"
	"fmt"
	"strings"
	"text/template"

	"github.com/Nadim147c/rong/internal/models"
)

var funcs = template.FuncMap{
	"upper":   strings.ToUpper,
	"lower":   strings.ToLower,
	"replace": strings.ReplaceAll,
	"chroma":  chroma,
	"tone":    tone,
	"qoute":   qoute,
	"json":    jsonString,
}

func chroma(color models.ColorValue, c float64) models.ColorValue {
	hct := color.Int.ToHct()
	hct.Chroma = c
	return models.NewColorValue(hct.ToARGB())
}

func tone(color models.ColorValue, t float64) models.ColorValue {
	hct := color.Int.ToHct()
	hct.Tone = t
	return models.NewColorValue(hct.ToARGB())
}

func qoute(s any) string {
	return fmt.Sprintf("%q", s)
}

func jsonString(s any) string {
	bytes, err := json.Marshal(s)
	if err != nil {
		return "null"
	}
	return string(bytes)
}
