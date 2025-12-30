package cache

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/Nadim147c/rong/v4/internal/material"
	"github.com/Nadim147c/rong/v4/internal/pathutil"
)

// State is the current generation state.
type State struct {
	Path      string             `json:"filename"`
	Hash      string             `json:"hash"`
	Quantized material.Quantized `json:"quantized"`
}

// SaveState saves state to state dir.
func SaveState(source, hash string, output material.Quantized) error {
	path := filepath.Join(pathutil.StateDir, "state.json")

	if err := os.MkdirAll(pathutil.CacheDir, 0o750); err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	state := State{source, hash, output}

	return json.NewEncoder(file).Encode(state)
}

// LoadState saves state to state dir.
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
