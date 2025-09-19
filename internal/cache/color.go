package cache

import (
	"encoding/hex"
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"github.com/Nadim147c/rong/internal/material"
	"github.com/Nadim147c/rong/internal/pathutil"
	"github.com/zeebo/xxh3"
)

func hash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}

	h := xxh3.New128()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	sum := h.Sum128().Bytes()
	return hex.EncodeToString(sum[:]), nil
}

// IsCached checks if the file is colors is cached or not
func IsCached(file string) bool {
	_, err := LoadCache(file)
	if err != nil {
		return false
	}
	return true
}

// State is the current generation state
type State struct {
	Path      string             `json:"filename"`
	Quantized material.Quantized `json:"quantized"`
}

// SaveState saves state to state dir
func SaveState(source string, output material.Quantized) error {
	path := filepath.Join(pathutil.StateDir, "state.json")

	if err := os.MkdirAll(pathutil.CacheDir, 0755); err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	state := State{source, output}

	return json.NewEncoder(file).Encode(state)
}

// LoadState saves state to state dir
func LoadState() (State, error) {
	path := filepath.Join(pathutil.StateDir, "state.json")

	var state State
	file, err := os.Open(path)
	if err != nil {
		return state, err
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&state)
	return state, err
}

// LoadCache tries to load cached colors for this image
func LoadCache(source string) (material.Quantized, error) {
	var output material.Quantized

	name, err := hash(source)
	if err != nil {
		return output, err
	}

	path := filepath.Join(pathutil.CacheDir, name+".json")

	file, err := os.Open(path)
	if err != nil {
		return output, err
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&output)
	return output, err
}

// SaveCache saves output colors to cache dir
func SaveCache(source string, output material.Quantized) error {
	name, err := hash(source)
	if err != nil {
		return err
	}

	path := filepath.Join(pathutil.CacheDir, name+".json")

	if err := os.MkdirAll(pathutil.CacheDir, 0755); err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(output)
}
