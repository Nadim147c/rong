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
	_, err := LoadCache(file)
	if err != nil {
		return false
	}
	return true
}

// LoadCache tries to load cached colors for this image
func LoadCache(file string) (models.Output, error) {
	var output models.Output

	name, err := hash(file)
	if err != nil {
		return output, err
	}

	path := filepath.Join(config.CacheDir, name+".json")

	cache, err := os.Open(path)
	if err != nil {
		return output, err
	}
	defer cache.Close()

	if err := json.NewDecoder(cache).Decode(&output); err != nil {
		return output, err
	}
	return output, nil
}

// SaveCache saves output colors to cache dir
func SaveCache(image string, output models.Output) error {
	name, err := hash(image)
	if err != nil {
		return err
	}

	path := filepath.Join(config.CacheDir, name+".json")

	if err := os.MkdirAll(config.CacheDir, 0755); err != nil {
		return err
	}

	cache, err := os.Create(path)
	if err != nil {
		return err
	}
	defer cache.Close()

	return json.NewEncoder(cache).Encode(output)
}
