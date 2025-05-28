package config

import (
	"fmt"
	"strconv"

	"github.com/BurntSushi/toml"
)

// Distance includes Hue, Chroma, Tone
type Distance [3]float64

// NewDistance creates Distance
func NewDistance(h, c, t float64) Distance {
	return Distance{h, c, t}
}

// NullValue is a generic wrapper
type NullValue[T any] struct {
	Value T
	Null  bool
}

type (
	// StringSlice is slice of string that parses string or []string from toml
	StringSlice []string
	// NullDistance is nullable distance
	NullDistance NullValue[Distance]
	// NullString is nullable string
	NullString NullValue[string]
	// NullFloat is nullable float64
	NullFloat NullValue[float64]
	// NullInt is nullable int64
	NullInt NullValue[int64]
	// NullBool is nullable bool
	NullBool NullValue[bool]
)

var null = []byte("null")

// MarshalTOML used from json/toml/yaml marshaling
func (ns NullString) MarshalTOML() ([]byte, error) {
	if ns.Null {
		return []byte("null"), nil
	}
	return fmt.Appendf([]byte{}, "%q", ns.Value), nil
}

// UnmarshalTOML used from json/toml/yaml marshaling
func (ns *NullString) UnmarshalTOML(data any) error {
	if data == nil {
		ns.Null = true
		ns.Value = ""
		return nil
	}
	str, ok := data.(string)
	if !ok {
		return fmt.Errorf("expected string, got %T", data)
	}
	ns.Null = false
	ns.Value = str
	return nil
}

// MarshalTOML used from json/toml/yaml marshaling
func (ni NullInt) MarshalTOML() ([]byte, error) {
	if ni.Null {
		return []byte("null"), nil
	}
	return []byte(strconv.FormatInt(ni.Value, 10)), nil
}

// UnmarshalTOML used from json/toml/yaml marshaling
func (ni *NullInt) UnmarshalTOML(data any) error {
	if data == nil {
		ni.Null = true
		ni.Value = 0
		return nil
	}
	v, ok := data.(int64)
	if !ok {
		return fmt.Errorf("expected int64, got %T", data)
	}
	ni.Null = false
	ni.Value = v
	return nil
}

// MarshalTOML used from json/toml/yaml marshaling
func (nf NullFloat) MarshalTOML() ([]byte, error) {
	if nf.Null {
		return []byte("null"), nil
	}
	return []byte(strconv.FormatFloat(nf.Value, 'f', -1, 64)), nil
}

// UnmarshalTOML used from json/toml/yaml marshaling
func (nf *NullFloat) UnmarshalTOML(data any) error {
	if data == nil {
		nf.Null = true
		nf.Value = 0
		return nil
	}
	v, ok := data.(float64)
	if !ok {
		return fmt.Errorf("expected float64, got %T", data)
	}
	nf.Null = false
	nf.Value = v
	return nil
}

// MarshalTOML used from json/toml/yaml marshaling
func (nb NullBool) MarshalTOML() ([]byte, error) {
	if nb.Null {
		return []byte("null"), nil
	}
	if nb.Value {
		return []byte("true"), nil
	}
	return []byte("false"), nil
}

// UnmarshalTOML used from json/toml/yaml marshaling
func (nb *NullBool) UnmarshalTOML(data any) error {
	if data == nil {
		nb.Null = true
		nb.Value = false
		return nil
	}
	v, ok := data.(bool)
	if !ok {
		return fmt.Errorf("expected bool, got %T", data)
	}
	nb.Null = false
	nb.Value = v
	return nil
}

// MarshalTOML always marshals as an array of strings
func (s StringSlice) MarshalTOML() ([]byte, error) {
	if len(s) == 1 {
		return toml.Marshal(s[0])
	}
	return toml.Marshal(s)
}

// UnmarshalTOML handles both string and []string TOML inputs
func (s *StringSlice) UnmarshalTOML(data any) error {
	switch v := data.(type) {
	case string:
		*s = StringSlice{v}
	case []any:
		result := make([]string, len(v))
		for i, item := range v {
			str, ok := item.(string)
			if !ok {
				return fmt.Errorf("expected string in array, got %T", item)
			}
			result[i] = str
		}
		*s = StringSlice(result)
	default:
		return fmt.Errorf("unsupported type for StringSlice: %T", data)
	}
	return nil
}
