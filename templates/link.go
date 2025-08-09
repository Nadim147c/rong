package templates

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

// hardlinkOrCopy tries to hardlink src to dst, or falls back to copying. If dst
// already exists and is hardlinked to src, it updates the mod time. If dst is a
// directory, it is removed before proceeding.
func hardlinkOrCopy(src, dst string) error {
	// Check if destination exists
	dstInfo, err := os.Lstat(dst)
	if err == nil {
		// If it's a directory, remove it
		if dstInfo.IsDir() {
			return fmt.Errorf("can't copy/link a directory")
		}

		srcInfo, err2 := os.Lstat(src)
		if err2 != nil {
			return fmt.Errorf("stat src: %w", err2)
		}

		// If already hardlinked, just update mod time
		if os.SameFile(srcInfo, dstInfo) {
			now := time.Now()
			err := os.Chtimes(dst, now, now)
			if err != nil {
				slog.Error("Failed to update", "file", dst)
			}
			slog.Info("Updated", "source", src, "destination", dst)
			return err
		}

		// Remove existing file
		if err := os.Remove(dst); err != nil {
			return fmt.Errorf("remove existing dst: %w", err)
		}
	}

	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		slog.Error("Failed to create parent directory", "for", dst)
		return err
	}

	// Try to create a hardlink
	if err := os.Link(src, dst); err == nil {
		slog.Info("Linked", "source", src, "destination", dst)
		return nil
	}

	// Fall back to copying
	err = copyFile(src, dst)
	if err != nil {
		slog.Error("Failed to copy", "error", err, "source", src, "destination", dst)
		return err
	}
	slog.Info("Copied", "source", src, "destination", dst)
	return err
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open src: %w", err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("create dst: %w", err)
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("copy: %w", err)
	}

	srcInfo, err := srcFile.Stat()
	if err != nil {
		return fmt.Errorf("stat src after copy: %w", err)
	}
	return os.Chmod(dst, srcInfo.Mode())
}
