package cache

import (
	"context"
	"os"
	"path/filepath"

	"github.com/Nadim147c/rong/internal/ffmpeg"
	"github.com/Nadim147c/rong/internal/pathutil"
)

// GetPreview returns the preview image
func GetPreview(src string, hash Sum) (string, error) {
	path := filepath.Join(pathutil.CacheDir, hash.String()+".webp")
	if _, err := os.Stat(path); err == nil {
		return path, nil
	}
	return path, ffmpeg.GeneratePreview(context.Background(), src, path)
}
