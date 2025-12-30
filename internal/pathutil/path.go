package pathutil

import (
	"errors"
	"path/filepath"
	"regexp"

	"github.com/adrg/xdg"
)

const app = "rong"

var separator = regexp.MustCompile(`[/\\]+`)

var (
	// ConfigDir is the rong config directory.
	ConfigDir = filepath.Clean(filepath.Join(xdg.ConfigHome, app))
	// CacheDir is the rong cache directory.
	CacheDir = filepath.Clean(filepath.Join(xdg.CacheHome, app))
	// StateDir is the rong state directory.
	StateDir = filepath.Clean(filepath.Join(xdg.StateHome, app))
)

// ErrEmptyPath means user difined an empty string as path.
var ErrEmptyPath = errors.New("empty path")

// FindPath expands environment-aware variables and returns an absolute path
// from a given base path (not CWD).
// Path Prefix Expands:
//   - $HOME or ~      : xdg.Home
//   - $XDG_CONFIG_HOME: xdg.ConfigHome
//   - $XDG_CACHE_HOME : xdg.CacheHome
//   - $XDG_DATA_HOME  : xdg.DataHome
//   - $RONG_CONFIG    : xdg.ConfigHome + "rong"
//   - $RONG_CACHE     : xdg.CacheHome + "rong"
//   - $RONG_DATA      : xdg.DataHome + "rong"
func FindPath(basePath, inputPath string) (string, error) {
	if inputPath == "" {
		return "", ErrEmptyPath
	}

	if filepath.IsAbs(inputPath) {
		return filepath.Clean(inputPath), nil
	}

	split := separator.Split(inputPath, 2)
	if len(split) != 2 {
		joined := filepath.Join(basePath, inputPath)
		return filepath.Clean(joined), nil
	}

	base, rest := split[0], split[1]

	var path string
	// Check the top level directory
	switch base {
	case "~", "$HOME":
		path = filepath.Join(xdg.Home, rest)
	case "$XDG_CONFIG_HOME":
		path = filepath.Join(xdg.ConfigHome, rest)
	case "$RONG_CONFIG":
		path = filepath.Join(xdg.ConfigHome, app, rest)
	case "$XDG_CACHE_HOME":
		path = filepath.Join(xdg.CacheHome, rest)
	case "$RONG_CACHE":
		path = filepath.Join(xdg.CacheHome, app, rest)
	case "$XDG_DATA_HOME":
		path = filepath.Join(xdg.DataHome, rest)
	case "$RONG_DATA":
		path = filepath.Join(xdg.DataHome, app, rest)
	}

	return filepath.Clean(path), nil
}
