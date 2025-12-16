package cache

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

// Supported image and video extensions
var mediaExtensions = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".bmp": true, ".webp": true,
	".mp4": true, ".mov": true, ".avi": true, ".mkv": true, ".flv": true, ".wmv": true,
	".webm": true, ".mpeg": true, ".mpg": true,
}

func isMediaFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return mediaExtensions[ext]
}

// ScanPaths scans the given paths and returns absolute paths of image/video
// files
func find(ctx context.Context, inputs []string, paths chan<- string) error {
	handler := func(path string, info os.FileInfo, err error) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err != nil {
			return err
		}

		if !info.IsDir() && isMediaFile(path) {
			abs, err := filepath.Abs(path)
			if err == nil {
				paths <- abs
			}
		}
		return nil
	}

	for _, p := range inputs {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		fileInfo, err := os.Stat(p)
		if err != nil {
			slog.Error("Failed to find file/directory", "path", p, "error", err)
			continue // Skip unreadable path
		}

		if fileInfo.IsDir() {
			filepath.Walk(p, handler)
			continue
		}

		if isMediaFile(p) {
			abs, err := filepath.Abs(p)
			if err == nil {
				paths <- abs
			}
		}
	}

	return nil
}
