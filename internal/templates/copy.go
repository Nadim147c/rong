package templates

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/google/renameio/v2"
)

// atomicCopy copies a file from src to dst atomically. It reads from src and
// writes to dst using Save.
func atomicCopy(src, dst string) error {
	dir := filepath.Dir(dst)
	// Ensure directory exists
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := renameio.TempFile("", dst)
	if err != nil {
		return err
	}
	defer dstFile.Cleanup()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	return dstFile.CloseAtomicallyReplace()
}

// hardlinkOrCopy tries to hardlink src to dst.
// - If dst already exists and is the same inode as src, it does nothing.
// - If dst exists but is different, it removes dst and recreates it.
// - If hardlink fails, it falls back to copying.
// - Parent directories for dst are created automatically.
func hardlinkOrCopy(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !srcInfo.Mode().IsRegular() {
		return errors.New("source is not a regular file")
	}

	// Fast path: dst exists and already hardlinked to src
	if dstInfo, err := os.Stat(dst); err == nil {
		if os.SameFile(srcInfo, dstInfo) {
			return nil
		}
		if err := os.Remove(dst); err != nil {
			return err
		}
	}

	// Ensure parent dirs exist
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}

	// Try hardlink first
	if err := os.Link(src, dst); err == nil {
		return nil
	}

	// Hardlink failed → copy
	return copyFile(src, dst)
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			out.Close()
			os.Remove(dst)
		}
	}()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	if err = out.Sync(); err != nil {
		return err
	}
	return out.Close()
}
