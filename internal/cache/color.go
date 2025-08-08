package cache

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"github.com/Nadim147c/rong/internal/material"
	"github.com/Nadim147c/rong/internal/pathutil"
	"github.com/cespare/xxhash"
)

func hash(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	sum := xxhash.New()
	if _, err := io.Copy(sum, file); err != nil {
		return "", err
	}

	result := make([]byte, 8)
	binary.BigEndian.PutUint64(result, sum.Sum64())

	return base64.RawURLEncoding.EncodeToString(result), nil
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
