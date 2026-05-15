package cgroup

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsSubpath(t *testing.T) {
	if !IsSubpath("/a/b", "/a/b/c") {
		t.Fatal("expected subpath")
	}
	if IsSubpath("/a/b", "/a/c") {
		t.Fatal("expected not subpath")
	}
	if !IsSubpath("/a/b", "/a/b") {
		t.Fatal("equal should be ok")
	}
}

func TestPathUnderOptionalPrefix(t *testing.T) {
	if !PathUnderOptionalPrefix("/x/y", "") {
		t.Fatal("empty prefix allows")
	}
	if !PathUnderOptionalPrefix("/a/b/c", "/a/b") {
		t.Fatal("expected under prefix")
	}
	if PathUnderOptionalPrefix("/a/other", "/a/b") {
		t.Fatal("expected outside prefix")
	}
}

func TestDiscoverableDirUnderPrefix(t *testing.T) {
	root := t.TempDir()
	mustMkdir(t, filepath.Join(root, "karl.slice", "cell-a"))
	ok, _, blocked := DiscoverableDir(root, filepath.Join(root, "karl.slice", "cell-a"), filepath.Join(root, "karl.slice"))
	if !ok || len(blocked) != 0 {
		t.Fatalf("ok=%v blocked=%v", ok, blocked)
	}
}

func TestDiscoverableDirOutsidePrefixBlocked(t *testing.T) {
	root := t.TempDir()
	mustMkdir(t, filepath.Join(root, "other", "x"))
	ok, _, blocked := DiscoverableDir(root, filepath.Join(root, "other", "x"), filepath.Join(root, "karl.slice"))
	if ok || len(blocked) == 0 {
		t.Fatalf("expected blocked ok=%v blocked=%v", ok, blocked)
	}
}

func TestDiscoverableDirSymlinkEscapeBlocked(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	mustMkdir(t, filepath.Join(outside, "secret"))
	link := filepath.Join(root, "trap")
	if err := os.Symlink(filepath.Join(outside, "secret"), link); err != nil {
		t.Fatal(err)
	}
	ok, _, blocked := DiscoverableDir(root, link, "")
	if ok {
		t.Fatal("expected symlink escape blocked")
	}
	if len(blocked) == 0 {
		t.Fatalf("expected blocked reasons: %#v", blocked)
	}
}

func mustMkdir(t *testing.T, p string) {
	t.Helper()
	if err := os.MkdirAll(p, 0o755); err != nil {
		t.Fatal(err)
	}
}
