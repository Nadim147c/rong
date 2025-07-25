package cache

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/Nadim147c/rong/internal/config"
	"github.com/Nadim147c/rong/internal/models"
	"github.com/cespare/xxhash"
)

func hash(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	h := xxhash.New()
	if _, err := io.Copy(h, file); err != nil {
		return "", err
	}

	return fmt.Sprint(h.Sum64()), nil
}

// IsCached checks if the file is colors is cached or not
func IsCached(file string) bool {
	_, _, err := LoadCache(file)
	if err != nil {
		return false
	}
	return true
}

// LoadCache tries to load cached colors for this image
func LoadCache(file string) (models.Output, []byte, error) {
	var jsonb []byte
	var output models.Output

	name, err := hash(file)
	if err != nil {
		return output, jsonb, err
	}

	path := filepath.Join(config.CacheDir, name+".json")

	cache, err := os.ReadFile(path)
	if err != nil {
		return output, jsonb, err
	}

	if err := json.Unmarshal(cache, &output); err != nil {
		return output, cache, err
	}
	return output, cache, nil
}

// SaveCache saves output colors to cache dir
func SaveCache(image string, output models.Output) ([]byte, error) {
	var jsonb []byte

	name, err := hash(image)
	if err != nil {
		return jsonb, err
	}

	path := filepath.Join(config.CacheDir, name+".json")

	if err := os.MkdirAll(config.CacheDir, 0755); err != nil {
		return jsonb, err
	}

	jsonb, err = json.Marshal(output)
	if err != nil {
		return jsonb, err
	}

	err = os.WriteFile(path, jsonb, 0644)

	return jsonb, err
}
