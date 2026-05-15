package main

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestGoldenDiscoverCgroupsStdout(t *testing.T) {
	t.Setenv("KHR_TEST_CGROUP_VERSION", "v2")
	root := moduleRoot(t)
	cfg := filepath.Join(root, "examples", "khr", "khr-linux-agent-config.yaml")
	cmdDir := filepath.Join(root, "cmd", "khr-linux-agent")
	cellInput := filepath.Join(root, "examples", "khr", "discovery", "cgroup-discovery-input-cell.json")

	t.Run("found", func(t *testing.T) {
		r := t.TempDir()
		if err := os.MkdirAll(filepath.Join(r, "karl.slice", "karl-shell-dev-linux-systemd-001.scope"), 0o755); err != nil {
			t.Fatal(err)
		}
		goldenPath := filepath.Join(root, "examples", "khr", "discovery", "cgroup-discovery-output-found.json")
		runDiscoverCompare(t, cmdDir, cfg, r, "", cellInput, goldenPath)
	})

	t.Run("not_found", func(t *testing.T) {
		r := t.TempDir()
		goldenPath := filepath.Join(root, "examples", "khr", "discovery", "cgroup-discovery-output-not-found.json")
		runDiscoverCompare(t, cmdDir, cfg, r, "", "", goldenPath)
	})

	t.Run("blocked_prefix", func(t *testing.T) {
		r := t.TempDir()
		if err := os.MkdirAll(filepath.Join(r, "outside", "bad"), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.MkdirAll(filepath.Join(r, "karl.slice", "karl-shell-dev-linux-systemd-001.scope"), 0o755); err != nil {
			t.Fatal(err)
		}
		cellTemplate := filepath.Join(root, "examples", "khr", "discovery", "cgroup-discovery-input-cell-explicit-outside.json")
		rawTpl, err := os.ReadFile(cellTemplate)
		if err != nil {
			t.Fatal(err)
		}
		cellFile := filepath.Join(t.TempDir(), "cell.json")
		if err := os.WriteFile(cellFile, []byte(strings.ReplaceAll(string(rawTpl), "__CGROUP_ROOT__", r)), 0o600); err != nil {
			t.Fatal(err)
		}
		goldenPath := filepath.Join(root, "examples", "khr", "discovery", "cgroup-discovery-output-blocked-prefix.json")
		runDiscoverCompare(t, cmdDir, cfg, r, filepath.Join(r, "karl.slice"), cellFile, goldenPath)
	})
}

func runDiscoverCompare(t *testing.T, cmdDir, cfg, cgroupRoot, allowPrefix, cellInput, goldenPath string) {
	t.Helper()
	args := []string{
		"run", ".",
		"-mode", "discover-cgroups",
		"-config", cfg,
		"-cgroup-root", cgroupRoot,
	}
	if allowPrefix != "" {
		args = append(args, "-allow-path-prefix", allowPrefix)
	}
	if cellInput != "" {
		args = append(args, "-cell-input", cellInput)
	}
	cmd := exec.Command("go", args...)
	cmd.Dir = cmdDir
	cmd.Env = append(os.Environ(), "KHR_TEST_CGROUP_VERSION=v2")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("go run: %v stderr=%s", err, stderr.String())
	}
	wantRaw, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatal(err)
	}
	wantStr := strings.ReplaceAll(string(wantRaw), "__CGROUP_ROOT__", cgroupRoot)
	var want, got any
	if err := json.Unmarshal([]byte(wantStr), &want); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(stdout.Bytes(), &got); err != nil {
		t.Fatalf("stdout json: %v body=%q", err, stdout.String())
	}
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("golden mismatch for %s\n--- want ---\n%s\n--- got ---\n%s", goldenPath, wantStr, stdout.String())
	}
}
