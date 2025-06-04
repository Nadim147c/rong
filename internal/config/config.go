package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

// Link is map all file the will be linked after generating color
type Link map[string]StringSlice

// Config is the rong configuration
type Config struct {
	Dark       NullBool   `toml:"dark"`    // true
	Version    NullInt    `toml:"version"` // 2021/2025
	MagickPath NullString `toml:"magick_path"`
	Link       Link       `toml:"link"`
}

// Global config holds the
var Global = defaultConfig()

// defaultConfig returns the default config values
func defaultConfig() Config {
	return Config{
		Dark:    NullBool{true, false},
		Version: NullInt{2021, false},
		Link:    Link{},
	}
}

// LoadConfig parses the TOML config and merges it with defaults
func LoadConfig(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "config file not found: %s", path)
	}

	parsed := &Config{}
	if _, err := toml.DecodeFile(path, parsed); err != nil {
		fmt.Fprintf(os.Stderr, "parse error: %v", err)
	}

	mergeConfig(&Global, parsed)
}

// mergeConfig merges non-zero values from `src` into `dst`
func mergeConfig(dst, src *Config) {
	if !src.Dark.Null {
		dst.Dark = src.Dark
	}
	if !src.Version.Null {
		dst.Version = src.Version
	}
}
