package cache

import (
	"context"
	"errors"
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
func ScanPaths(ctx context.Context, paths []string) ([]string, error) {
	var results []string

	handler := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		// Check context
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if !info.IsDir() && isMediaFile(path) {
			abs, err := filepath.Abs(path)
			if err == nil {
				results = append(results, abs)
			}
		}
		return nil
	}

	for _, p := range paths {
		select {
		case <-ctx.Done():
			return results, ctx.Err()
		default:
		}

		fileInfo, err := os.Stat(p)
		if err != nil {
			slog.Error("Failed to find file/directory", "path", p, "error", err)
			continue // Skip unreadable path
		}

		if fileInfo.IsDir() {
			err := filepath.Walk(p, handler)
			if err != nil {
				if errors.Is(err, ctx.Err()) {
					return results, ctx.Err()
				}
				return results, err
			}

			continue
		}

		if isMediaFile(p) {
			abs, err := filepath.Abs(p)
			if err == nil {
				results = append(results, abs)
			}
		}
	}

	return results, nil
}
