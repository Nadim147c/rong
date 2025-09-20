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

// Sum is the xxh3_128 hash result
type Sum struct {
	xxh3.Uint128
}

func (s Sum) String() string {
	b := s.Uint128.Bytes()
	return hex.EncodeToString(b[:])
}

// Hash returns xxh3_128 sum
func Hash(path string) (Sum, error) {
	f, err := os.Open(path)
	if err != nil {
		return Sum{}, err
	}

	h := xxh3.New128()
	if _, err := io.Copy(h, f); err != nil {
		return Sum{}, err
	}

	return Sum{h.Sum128()}, nil
}

// IsCached checks if the file is colors is cached or not
func IsCached(hash Sum) bool {
	path := filepath.Join(pathutil.CacheDir, hash.String()+".json")
	_, err := os.Stat(path)
	return err == nil
}

// LoadCache tries to load cached colors for this image
func LoadCache(hash Sum) (material.Quantized, error) {
	var output material.Quantized
	path := filepath.Join(pathutil.CacheDir, hash.String()+".json")

	file, err := os.Open(path)
	if err != nil {
		return output, err
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&output)
	return output, err
}

// SaveCache saves output colors to cache dir
func SaveCache(hash Sum, output material.Quantized) error {
	path := filepath.Join(pathutil.CacheDir, hash.String()+".json")

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
