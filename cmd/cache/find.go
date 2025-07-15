package cache

import (
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

// ScanPaths scans the given paths and sends absolute paths of image/video files to the provided channel
func ScanPaths(paths []string, out chan<- string) {
	defer close(out)
	for _, p := range paths {
		fileInfo, err := os.Stat(p)
		if err != nil {
			slog.Error("Failed to find file/directory", "path", p, "error", err)
			continue // Skip unreadable path
		}

		if fileInfo.IsDir() {
			filepath.Walk(p, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return nil
				}
				if !info.IsDir() && isMediaFile(path) {
					abs, err := filepath.Abs(path)
					if err == nil {
						out <- abs
					}
				}
				return nil
			})
			continue
		}

		if isMediaFile(p) {
			abs, err := filepath.Abs(p)
			if err == nil {
				out <- abs
			}
		}
	}
}
