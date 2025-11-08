package cache

import (
	"context"
	"os"
	"path/filepath"

	"github.com/Nadim147c/rong/v3/internal/ffmpeg"
	"github.com/Nadim147c/rong/v3/internal/pathutil"
)

// GetPreview returns the preview image
func GetPreview(src string, hash string) (string, error) {
	path := filepath.Join(pathutil.CacheDir, hash+".webp")
	if _, err := os.Stat(path); err == nil {
		return path, nil
	}
	return path, ffmpeg.GeneratePreview(context.Background(), src, path)
}
