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

func TestGoldenReadTelemetryStdout(t *testing.T) {
	t.Setenv("KHR_TEST_CGROUP_VERSION", "v2")
	t.Setenv("KHR_TEST_TELEMETRY_NOW", "2026-05-15T15:00:00Z")
	root := moduleRoot(t)
	cfg := filepath.Join(root, "examples", "khr", "khr-linux-agent-config.yaml")
	cmdDir := filepath.Join(root, "cmd", "khr-linux-agent")
	cell := filepath.Join(root, "examples", "khr", "telemetry", "telemetry-input-cell.json")

	t.Run("full", func(t *testing.T) {
		r := t.TempDir()
		if err := os.WriteFile(filepath.Join(r, "cpu.stat"), []byte("usage_usec 100\nnr_periods 2\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(r, "memory.current"), []byte("8192\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(r, "memory.max"), []byte("max\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(r, "memory.events"), []byte("pgfault 5\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(r, "io.stat"), []byte("259:0 rbytes=10\n\n259:1 wbytes=20\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		golden := filepath.Join(root, "examples", "khr", "telemetry", "telemetry-output-full.json")
		runTelemetryGolden(t, cmdDir, cfg, r, r, "", cell, golden)
	})

	t.Run("missing_files", func(t *testing.T) {
		r := t.TempDir()
		if err := os.WriteFile(filepath.Join(r, "memory.max"), []byte("max\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		golden := filepath.Join(root, "examples", "khr", "telemetry", "telemetry-output-missing-files.json")
		runTelemetryGolden(t, cmdDir, cfg, r, r, "", "", golden)
	})

	t.Run("blocked_prefix", func(t *testing.T) {
		r := t.TempDir()
		if err := os.MkdirAll(filepath.Join(r, "outside", "c"), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(r, "outside", "c", "cpu.stat"), []byte("x\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		golden := filepath.Join(root, "examples", "khr", "telemetry", "telemetry-output-blocked-prefix.json")
		runTelemetryGolden(t, cmdDir, cfg, r, filepath.Join(r, "outside", "c"), filepath.Join(r, "karl.slice"), "", golden)
	})
}

func runTelemetryGolden(t *testing.T, cmdDir, cfg, templateRoot, cgroupPath, allowPrefix, cellInput, goldenPath string) {
	t.Helper()
	args := []string{
		"run", ".",
		"-mode", "read-telemetry",
		"-config", cfg,
		"-cgroup-path", cgroupPath,
	}
	if allowPrefix != "" {
		args = append(args, "-allow-path-prefix", allowPrefix)
	}
	if cellInput != "" {
		args = append(args, "-cell-input", cellInput)
	}
	cmd := exec.Command("go", args...)
	cmd.Dir = cmdDir
	cmd.Env = append(os.Environ(), "KHR_TEST_CGROUP_VERSION=v2", "KHR_TEST_TELEMETRY_NOW=2026-05-15T15:00:00Z")
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
	wantStr := strings.ReplaceAll(string(wantRaw), "__ROOT__", templateRoot)
	var want, got any
	if err := json.Unmarshal([]byte(wantStr), &want); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(stdout.Bytes(), &got); err != nil {
		t.Fatalf("stdout json: %v body=%q", err, stdout.String())
	}
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("golden mismatch %s\n--- want ---\n%s\n--- got ---\n%s", goldenPath, wantStr, stdout.String())
	}
}
