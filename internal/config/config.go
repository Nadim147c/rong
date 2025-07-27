package config

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"

	"github.com/BurntSushi/toml"
	"github.com/Nadim147c/material/dynamic"
	"github.com/goccy/go-yaml"
)

// Links is map all file the will be linked after generating color
type Links map[string][]string

// Config is the rong configuration
type Config struct {
	Light     bool
	Constrast float64
	Version   dynamic.Version
	Variant   dynamic.Variant
	Platform  dynamic.Platform
	Links     Links
}

// parsedConfig is the config used during toml parse
type parsedConfig struct {
	Light     *bool             `toml:"light" yaml:"light"`
	Constrast *float64          `toml:"constrast" yaml:"constrast"`
	Version   *dynamic.Version  `toml:"version" yaml:"version"`
	Variant   *dynamic.Variant  `toml:"variant" yaml:"variant"`
	Platform  *dynamic.Platform `toml:"platform" yaml:"platform"`
	Links     parseLinks        `toml:"links" yaml:"links"`
}

type parseLinks map[string]StringSlice

// Global config holds the
var Global = defaultConfig()

// defaultConfig returns the default config values
func defaultConfig() Config {
	return Config{
		Light:   false,
		Version: dynamic.V2021,
		Links:   Links{},
	}
}

// LoadConfig parses the TOML config and merges it with defaults
func LoadConfig(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return err
	}

	parsed := parsedConfig{}

	switch filepath.Ext(path) {
	case ".yaml", ".yml":
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		if err := yaml.NewDecoder(file).Decode(&parsed); err != nil {
			return err
		}
	case ".toml":
		if _, err := toml.DecodeFile(path, &parsed); err != nil {
			return err
		}
	default:
		return errors.New("Invalid config extention")
	}

	slog.Debug("Parsed", "config", parsed)

	mergeConfig(&Global, &parsed)
	slog.Debug("Config", "config", Global)
	return nil
}

var validVariants = map[dynamic.Variant]bool{
	dynamic.Monochrome: true, dynamic.Neutral: true, dynamic.TonalSpot: true,
	dynamic.Vibrant: true, dynamic.Expressive: true, dynamic.Fidelity: true,
	dynamic.Content: true, dynamic.Rainbow: true, dynamic.FruitSalad: true,
}

// mergeConfig merges non-zero values from `src` into `dst`
func mergeConfig(dst *Config, src *parsedConfig) {
	dstVal := reflect.ValueOf(dst).Elem()
	srcVal := reflect.ValueOf(src).Elem()

	for i := 0; i < srcVal.NumField(); i++ {
		srcField := srcVal.Field(i)
		if srcField.IsNil() {
			continue
		}

		fieldName := srcVal.Type().Field(i).Name

		if fieldName == "Links" {
			continue
		}

		dstField := dstVal.FieldByName(fieldName)
		srcValue := srcField.Elem().Interface()

		// Add field-specific validation
		switch fieldName {
		case "Version":
			if version, ok := srcValue.(dynamic.Version); ok {
				if version != dynamic.V2021 && version != dynamic.V2025 {
					slog.Error("Invalid platform", "version", version)
					continue
				}
			}
		case "Constrast":
			if contrast, ok := srcValue.(float64); ok {
				if contrast < -1 || contrast > 1 {
					slog.Error("Invalid contrast value", "value", contrast)
					continue
				}
			}
		case "Platform":
			if platform, ok := srcValue.(dynamic.Platform); ok {
				if platform != dynamic.Phone && platform != dynamic.Watch {
					slog.Error("Invalid platform", "platform", platform)
					continue
				}
			}
		case "Variant":
			if variant, ok := srcValue.(dynamic.Variant); ok {
				if !validVariants[variant] {
					slog.Error("Invalid scheme variant", "variant", variant)
					continue
				}
			}
		}
		// Set the value if validation passed
		dstField.Set(reflect.ValueOf(srcValue))
	}

	if src.Links != nil {
		if dst.Links == nil {
			dst.Links = Links{}
		}
		for n, v := range src.Links {
			dst.Links[n] = []string(v)
		}
	}
}
