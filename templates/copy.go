package templates

import (
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
