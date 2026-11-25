// Package fsutil writes files that hold credentials.
package fsutil

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
	return rename(tmp.Name(), path)
}

var rename = os.Rename // replaced in tests

// resolveLinks follows symlinks one hop at a time, so a link whose target
// does not exist yet still leads to that target (EvalSymlinks fails there).
// Each step resolves the parent directory on disk, so "hop/.." in a target
// means what the kernel means when hop is itself a link.
func resolveLinks(path string) (string, error) {
	path, err := realParent(path)
	if err != nil {
		return "", err
	}
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
			target = filepath.Dir(path) + string(filepath.Separator) + target
		}
		if path, err = realParent(target); err != nil {
			return "", err
		}
	}
	return "", fmt.Errorf("%s: too many levels of symbolic links", path)
}

// realParent resolves the directory part of an unclean path on disk and keeps
// the last element. A directory that cannot be resolved is an error, as it is
// for the kernel: "missing/../x" must not become "x".
func realParent(p string) (string, error) {
	i := strings.LastIndex(p, string(filepath.Separator))
	if i < 0 {
		return p, nil
	}
	dir, base := p[:i], p[i+1:]
	if dir == "" {
		dir = string(filepath.Separator)
	}
	real, err := filepath.EvalSymlinks(dir)
	if err != nil {
		return "", err
	}
	return filepath.Join(real, base), nil
}
