package templates

import (
	"encoding/json"
	"fmt"
	"strings"
	"text/template"
)

var funcs = template.FuncMap{
	"upper": strings.ToUpper,
	"lower": strings.ToLower,
	"qoute": qoute,
	"json":  jsonString,
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
