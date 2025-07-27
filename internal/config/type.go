package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/goccy/go-yaml"
)

// StringSlice if slice of string that support encoding string or []string
type StringSlice []string

var (
	_ toml.Marshaler          = (*StringSlice)(nil)
	_ toml.Unmarshaler        = (*StringSlice)(nil)
	_ yaml.InterfaceMarshaler = (*StringSlice)(nil)
	_ yaml.BytesUnmarshaler   = (*StringSlice)(nil)
)

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

// MarshalYAML always marshals as an array of strings
func (s StringSlice) MarshalYAML() (any, error) {
	if len(s) == 1 {
		return s[0], nil
	}
	return s, nil
}

// UnmarshalYAML implements yaml.InterfaceUnmarshaler
func (s *StringSlice) UnmarshalYAML(b []byte) error {
	var data any
	err := yaml.Unmarshal(b, &data)
	if err != nil {
		return err
	}

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
