package config

import (
	"log/slog"
	"os"
	"reflect"

	"github.com/BurntSushi/toml"
	"github.com/Nadim147c/material/dynamic"
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

// tomlConfig is the config used during toml parse
type tomlConfig struct {
	Light     *bool             `toml:"light"`
	Constrast *float64          `toml:"constrast"`
	Version   *dynamic.Version  `toml:"version"`
	Variant   *dynamic.Variant  `toml:"variant"`
	Platform  *dynamic.Platform `toml:"platform"`
	Links     tomlLinks         `toml:"links"`
}

type tomlLinks map[string]StringSlice

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
func LoadConfig(path string) {
	cwd, err := os.Getwd()
	if err != nil {
		return
	}

	path, err = FindPath(cwd, path)
	if err != nil {
		slog.Error("Failed to find config file", "error", err)
		return
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		slog.Error("Config file doesn't exists", "error", err)
		return
	}

	parsed := &tomlConfig{}
	if _, err := toml.DecodeFile(path, parsed); err != nil {
		slog.Error("Failed to load config file", "error", err)
	}

	mergeConfig(&Global, parsed)
	slog.Debug("Config", "config", Global)
}

var validVariants = map[dynamic.Variant]bool{
	dynamic.Monochrome: true, dynamic.Neutral: true, dynamic.TonalSpot: true,
	dynamic.Vibrant: true, dynamic.Expressive: true, dynamic.Fidelity: true,
	dynamic.Content: true, dynamic.Rainbow: true, dynamic.FruitSalad: true,
}

// mergeConfig merges non-zero values from `src` into `dst`
func mergeConfig(dst *Config, src *tomlConfig) {
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
