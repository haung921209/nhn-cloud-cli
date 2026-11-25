package fsutil

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func mode(t *testing.T, p string) os.FileMode {
	t.Helper()
	fi, err := os.Stat(p)
	if err != nil {
		t.Fatal(err)
	}
	return fi.Mode().Perm()
}

func TestWritePrivateReplacesAsOwnerOnly(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "key")
	os.WriteFile(p, []byte("old"), 0644)
	before, _ := os.Stat(p)
	if err := WritePrivate(p, []byte("new")); err != nil {
		t.Fatal(err)
	}
	if got, _ := os.ReadFile(p); string(got) != "new" || mode(t, p) != 0600 {
		t.Errorf("content %q mode %v", got, mode(t, p))
	}
	// Replaced by rename, not rewritten in place: a failed write cannot
	// truncate the old file.
	if after, _ := os.Stat(p); os.SameFile(before, after) {
		t.Error("the old file was rewritten in place")
	}
	if left, _ := filepath.Glob(filepath.Join(dir, ".key-*")); len(left) != 0 {
		t.Errorf("temporary files left: %v", left)
	}
}

func TestWritePrivateRefusesDirectoryWithoutChangingIt(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "kube")
	os.Mkdir(dir, 0755)
	link := dir + "-link"
	os.Symlink(dir, link)
	for _, p := range []string{dir, link} {
		err := WritePrivate(p, []byte("x"))
		if err == nil || !strings.Contains(err.Error(), "not a regular file") {
			t.Errorf("%s: want refusal, got %v", p, err)
		}
	}
	if mode(t, dir) != 0755 {
		t.Errorf("directory mode changed to %v", mode(t, dir))
	}
}

func TestWritePrivateKeepsSymlink(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "real")
	os.WriteFile(target, []byte("old"), 0644)
	link := filepath.Join(dir, "link")
	os.Symlink(target, link)
	if err := WritePrivate(link, []byte("new")); err != nil {
		t.Fatal(err)
	}
	if fi, _ := os.Lstat(link); fi.Mode()&os.ModeSymlink == 0 {
		t.Error("the symlink was replaced by a file")
	}
	if got, _ := os.ReadFile(target); string(got) != "new" || mode(t, target) != 0600 {
		t.Errorf("target content %q mode %v", got, mode(t, target))
	}
}

func TestWritePrivateFailureKeepsOldFile(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "key")
	os.WriteFile(p, []byte("old"), 0600)
	os.Chmod(dir, 0500) // no new file can be created next to it
	defer os.Chmod(dir, 0700)
	if err := WritePrivate(p, []byte("new")); err == nil {
		t.Skip("directory is writable here (running as root?)")
	}
	if got, _ := os.ReadFile(p); string(got) != "old" {
		t.Errorf("old content lost: %q", got)
	}
}

func TestWritePrivateFollowsDanglingSymlink(t *testing.T) {
	dir := t.TempDir()
	os.Mkdir(filepath.Join(dir, "dotfiles"), 0700)
	link := filepath.Join(dir, "config")
	os.Symlink("dotfiles/config", link) // relative, target not created yet
	if err := WritePrivate(link, []byte("new")); err != nil {
		t.Fatal(err)
	}
	if fi, _ := os.Lstat(link); fi.Mode()&os.ModeSymlink == 0 {
		t.Error("the dangling symlink was replaced by a file")
	}
	target := filepath.Join(dir, "dotfiles", "config")
	if got, _ := os.ReadFile(target); string(got) != "new" || mode(t, target) != 0600 {
		t.Errorf("target content %q", got)
	}
}
