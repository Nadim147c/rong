package models

import (
	"encoding/json"
	"io"
	"slices"
)

// WriteSimpleJSON generate json of format for out in simple key-value format.
func WriteSimpleJSON(w io.Writer, out Output) error {
	m := make(map[string]string, len(out.Colors)+1)
	m["image"] = out.Image
	for v := range slices.Values(out.Colors) {
		m[v.Name.Snake] = v.Color.HexRGB
	}
	return json.NewEncoder(w).Encode(m)
}
