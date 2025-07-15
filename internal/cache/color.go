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

// LoadCache tries to load cached colors for this image
func LoadCache(image string) (models.Output, error) {
	var output models.Output

	file, err := os.Open(image)
	if err != nil {
		return output, err
	}
	defer file.Close()

	h := xxhash.New()
	if _, err := io.Copy(h, file); err != nil {
		return output, err
	}

	hash := fmt.Sprint(h.Sum64())
	path := filepath.Join(config.CacheDir, hash+".json")

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
	file, err := os.Open(image)
	if err != nil {
		return err
	}
	defer file.Close()

	h := xxhash.New()
	if _, err := io.Copy(h, file); err != nil {
		return err
	}
	hash := fmt.Sprint(h.Sum64())

	if err := os.MkdirAll(config.CacheDir, 0755); err != nil {
		return err
	}

	path := filepath.Join(config.CacheDir, hash+".json")

	cache, err := os.Create(path)
	if err != nil {
		return err
	}
	defer cache.Close()

	return json.NewEncoder(cache).Encode(output)
}
