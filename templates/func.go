package templates

import (
	"fmt"
	"strings"
	"text/template"
)

var funcs = template.FuncMap{
	"upper": strings.ToUpper,
	"lower": strings.ToLower,
	"qoute": qoute,
}

func qoute(s any) string {
	return fmt.Sprintf("%q", s)
}
