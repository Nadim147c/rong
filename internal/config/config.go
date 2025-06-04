package config

import (
	"log/slog"
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

	parsed := &Config{}
	if _, err := toml.DecodeFile(path, parsed); err != nil {
		slog.Error("Failed to load config file", "error", err)
	}

	mergeConfig(&Global, parsed)
	slog.Debug("Config", "config", Global)
}

// mergeConfig merges non-zero values from `src` into `dst`
func mergeConfig(dst, src *Config) {
	if !src.Dark.Null {
		dst.Dark = src.Dark
	}
	if !src.Version.Null {
		dst.Version = src.Version
	}

	if src.Link != nil {
		dst.Link = src.Link
	}
}
