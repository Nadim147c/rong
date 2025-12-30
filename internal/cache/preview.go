package cache

import (
	"context"
	"os"
	"path/filepath"

	"github.com/Nadim147c/rong/v4/internal/ffmpeg"
	"github.com/Nadim147c/rong/v4/internal/pathutil"
	"github.com/spf13/viper"
)

// GetPreview returns the preview image.
func GetPreview(src string, hash string) (string, error) {
	format := viper.GetString("preview-format")
	path := filepath.Join(pathutil.CacheDir, hash+"."+format)
	if _, err := os.Stat(path); err == nil {
		return path, nil
	}
	return path, ffmpeg.GeneratePreview(context.Background(), src, path)
}
