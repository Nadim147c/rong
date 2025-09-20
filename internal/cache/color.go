package cache

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/Nadim147c/rong/internal/material"
	"github.com/Nadim147c/rong/internal/pathutil"
	"github.com/zeebo/xxh3"
)

// Sum is the xxh3_128 hash result
type Sum struct {
	bytes []byte
}

// NewSum return Sum
func NewSum(s xxh3.Uint128) Sum {
	b := s.Bytes()
	return Sum{bytes: b[:]}
}

func (s Sum) String() string {
	return hex.EncodeToString(s.bytes)
}

// MarshalText implements the encoding.TextMarshaler interface.
// It encodes the Sum as a lowercase hexadecimal string.
func (s Sum) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// It decodes a lowercase hexadecimal string into the Sum.
// The input must represent exactly 16 bytes (32 hex characters).
func (s *Sum) UnmarshalText(text []byte) error {
	b, err := hex.DecodeString(string(text))
	if err != nil {
		return err
	}
	if len(b) != 16 {
		return fmt.Errorf("invalid length: expected 16 bytes, got %d", len(b))
	}
	s.bytes = b
	return nil
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

	return NewSum(h.Sum128()), nil
}

// IsCached checks if the file is colors is cached or not
func IsCached(hash Sum) bool {
	p1 := filepath.Join(pathutil.CacheDir, hash.String()+".json")
	p2 := filepath.Join(pathutil.CacheDir, hash.String()+".webp")
	_, e1 := os.Stat(p1)
	_, e2 := os.Stat(p2)
	return e1 == nil && e2 == nil
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
