package cache

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/Nadim147c/rong/v4/internal/material"
	"github.com/Nadim147c/rong/v4/internal/pathutil"
	"github.com/spf13/viper"
)

// Hash returns xxh3_128 sum
func Hash(path string) (string, error) {
	name, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	b := md5.Sum([]byte(name))
	return hex.EncodeToString(b[:]), nil
}

// IsCached checks if the file is colors is cached or not
func IsCached(hash string, isVideo bool) bool {
	if !isVideo {
		jsonCache := filepath.Join(pathutil.CacheDir, hash+".json")
		_, err := os.Stat(jsonCache)
		return err == nil
	}
	format := viper.GetString("preview-format")
	preview := filepath.Join(pathutil.CacheDir, hash+"."+format)
	_, err := os.Stat(preview)
	return err == nil
}

// LoadCache tries to load cached colors for this image
func LoadCache(hash string) (material.Quantized, error) {
	var output material.Quantized
	path := filepath.Join(pathutil.CacheDir, hash+".json")

	file, err := os.Open(path)
	if err != nil {
		return output, err
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&output)
	return output, err
}

// SaveCache saves output colors to cache dir
func SaveCache(hash string, output material.Quantized) error {
	path := filepath.Join(pathutil.CacheDir, hash+".json")

	if err := os.MkdirAll(pathutil.CacheDir, 0o755); err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(output)
}
