package cgroup

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateCgroupPathForTelemetryPrefix(t *testing.T) {
	root := t.TempDir()
	good := filepath.Join(root, "karl.slice", "cell")
	if err := os.MkdirAll(good, 0o755); err != nil {
		t.Fatal(err)
	}
	pfx := filepath.Join(root, "karl.slice")
	res, _, blocked := ValidateCgroupPathForTelemetry(good, pfx)
	if len(blocked) != 0 || res == "" {
		t.Fatalf("res=%q blocked=%v", res, blocked)
	}
	outside := filepath.Join(root, "other")
	if err := os.MkdirAll(outside, 0o755); err != nil {
		t.Fatal(err)
	}
	_, _, blocked2 := ValidateCgroupPathForTelemetry(outside, pfx)
	if len(blocked2) == 0 {
		t.Fatal("expected prefix block")
	}
}

func TestReadFileInResolvedDirSymlinkEscape(t *testing.T) {
	root := t.TempDir()
	cell := filepath.Join(root, "cell")
	if err := os.MkdirAll(cell, 0o755); err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(root, "secret")
	if err := os.WriteFile(out, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	trap := filepath.Join(cell, "cpu.stat")
	if err := os.Symlink(out, trap); err != nil {
		t.Fatal(err)
	}
	res, _, _ := ValidateCgroupPathForTelemetry(cell, "")
	_, _, blocked := ReadFileInResolvedDir(res, "cpu.stat")
	if len(blocked) == 0 {
		t.Fatal("expected symlink escape block")
	}
}

func TestReadFileInResolvedDirOk(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "cpu.stat"), []byte("usage_usec 1\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	res, _, blocked := ValidateCgroupPathForTelemetry(root, "")
	if len(blocked) != 0 {
		t.Fatal(blocked)
	}
	data, _, b2 := ReadFileInResolvedDir(res, "cpu.stat")
	if len(b2) != 0 || string(data) != "usage_usec 1\n" {
		t.Fatalf("data=%q b=%v", data, b2)
	}
}
