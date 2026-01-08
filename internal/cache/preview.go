package cache

import (
	"context"
	"os"
	"path/filepath"

	"github.com/Nadim147c/rong/v4/internal/config"
	"github.com/Nadim147c/rong/v4/internal/ffmpeg"
	"github.com/Nadim147c/rong/v4/internal/pathutil"
)

// GetPreview returns the preview image.
func GetPreview(src string, hash string) (string, error) {
	format := config.PreviewFormat.Value()
	path := filepath.Join(pathutil.CacheDir, hash+"."+format.String())
	if _, err := os.Stat(path); err == nil {
		return path, nil
	}
	return path, ffmpeg.GeneratePreview(context.Background(), src, path)
}
