// Package fsutil writes files that hold credentials.
package fsutil

import (
	"fmt"
	"os"
	"path/filepath"
)

// WritePrivate replaces path with data as an owner-only (0600) file.
//
// The data goes to a new temporary file in the same directory, which is
// renamed over path, so a failed write leaves the old file as it was and the
// content is never in a file another user can open. A symlink keeps pointing
// to its target, which is replaced. A directory or other non-regular file is
// refused without changing it.
func WritePrivate(path string, data []byte) error {
	path, err := resolveLinks(path)
	if err != nil {
		return err
	}
	if fi, err := os.Stat(path); err == nil && !fi.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", path)
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), "."+filepath.Base(path)+"-*")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name()) // no-op after the rename
	if err := tmp.Chmod(0600); err != nil {
		tmp.Close()
		return err
	}
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Sync(); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), path)
}

// resolveLinks follows symlinks one hop at a time, so a link whose target
// does not exist yet still leads to that target (EvalSymlinks fails there).
func resolveLinks(path string) (string, error) {
	for i := 0; i < 40; i++ {
		fi, err := os.Lstat(path)
		if err != nil || fi.Mode()&os.ModeSymlink == 0 {
			return path, nil
		}
		target, err := os.Readlink(path)
		if err != nil {
			return "", err
		}
		if !filepath.IsAbs(target) {
			target = filepath.Join(filepath.Dir(path), target)
		}
		path = target
	}
	return "", fmt.Errorf("%s: too many levels of symbolic links", path)
}
