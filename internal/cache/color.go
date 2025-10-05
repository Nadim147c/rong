package cache

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/Nadim147c/rong/internal/material"
	"github.com/Nadim147c/rong/internal/pathutil"
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
func IsCached(hash string) bool {
	p1 := filepath.Join(pathutil.CacheDir, hash+".json")
	p2 := filepath.Join(pathutil.CacheDir, hash+".webp")
	_, e1 := os.Stat(p1)
	_, e2 := os.Stat(p2)
	return e1 == nil && e2 == nil
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
