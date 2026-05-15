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

func TestGoldenCollectEvidenceStdout(t *testing.T) {
	t.Setenv("KHR_TEST_CGROUP_VERSION", "v2")
	t.Setenv("KHR_TEST_COLLECTED_AT", "2026-05-15T16:00:00Z")
	t.Setenv("KHR_TEST_TELEMETRY_NOW", "2026-05-15T15:00:00Z")
	root := moduleRoot(t)
	cfg := filepath.Join(root, "examples", "khr", "khr-linux-agent-config.yaml")
	cmdDir := filepath.Join(root, "cmd", "khr-linux-agent")
	cell := filepath.Join(root, "examples", "khr", "evidence", "collect-evidence-input-cell.json")
	lease := filepath.Join(root, "examples", "khr", "resourcelease-linux-envelope-full.json")
	port := filepath.Join(root, "examples", "khr", "resourceport-linux-envelope-full.json")
	cellCtx := filepath.Join(root, "examples", "khr", "golden", "inputs", "cell-context-allowed.json")

	t.Run("ready", func(t *testing.T) {
		r := t.TempDir()
		if err := os.MkdirAll(filepath.Join(r, "karl.slice", "karl-shell-dev-linux-systemd-001.scope"), 0o755); err != nil {
			t.Fatal(err)
		}
		scope := filepath.Join(r, "karl.slice", "karl-shell-dev-linux-systemd-001.scope")
		writeMetricFixtures(t, scope)
		golden := filepath.Join(root, "examples", "khr", "evidence", "collect-evidence-output-ready.json")
		runCollectEvidenceGolden(t, cmdDir, cfg, r, cell, nil, golden)
	})

	t.Run("blocked_no_cgroup", func(t *testing.T) {
		r := t.TempDir()
		golden := filepath.Join(root, "examples", "khr", "evidence", "collect-evidence-output-blocked-no-cgroup.json")
		runCollectEvidenceGolden(t, cmdDir, cfg, r, cell, nil, golden)
	})

	t.Run("with_dryrun", func(t *testing.T) {
		r := t.TempDir()
		if err := os.MkdirAll(filepath.Join(r, "karl.slice", "karl-shell-dev-linux-systemd-001.scope"), 0o755); err != nil {
			t.Fatal(err)
		}
		scope := filepath.Join(r, "karl.slice", "karl-shell-dev-linux-systemd-001.scope")
		writeMetricFixtures(t, scope)
		golden := filepath.Join(root, "examples", "khr", "evidence", "collect-evidence-output-with-dryrun.json")
		extra := []string{"-lease-input", lease, "-resource-port-input", port, "-cell-context", cellCtx}
		runCollectEvidenceGolden(t, cmdDir, cfg, r, cell, extra, golden)
	})

	t.Run("dryrun_skipped", func(t *testing.T) {
		r := t.TempDir()
		if err := os.MkdirAll(filepath.Join(r, "karl.slice", "karl-shell-dev-linux-systemd-001.scope"), 0o755); err != nil {
			t.Fatal(err)
		}
		scope := filepath.Join(r, "karl.slice", "karl-shell-dev-linux-systemd-001.scope")
		writeMetricFixtures(t, scope)
		golden := filepath.Join(root, "examples", "khr", "evidence", "collect-evidence-output-dryrun-skipped.json")
		extra := []string{"-lease-input", lease}
		runCollectEvidenceGolden(t, cmdDir, cfg, r, cell, extra, golden)
	})
}

func writeMetricFixtures(t *testing.T, scope string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(scope, "cpu.stat"), []byte("usage_usec 100\nnr_periods 2\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(scope, "memory.current"), []byte("8192\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(scope, "memory.max"), []byte("max\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(scope, "memory.events"), []byte("pgfault 5\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(scope, "io.stat"), []byte("259:0 rbytes=10\n\n259:1 wbytes=20\n"), 0o644); err != nil {
		t.Fatal(err)
	}
}

func runCollectEvidenceGolden(t *testing.T, cmdDir, cfg, cgroupRoot, cellInput string, extra []string, goldenPath string) {
	t.Helper()
	args := []string{
		"run", ".",
		"-mode", "collect-evidence",
		"-config", cfg,
		"-cell-input", cellInput,
		"-cgroup-root", cgroupRoot,
	}
	args = append(args, extra...)
	cmd := exec.Command("go", args...)
	cmd.Dir = cmdDir
	cmd.Env = append(os.Environ(),
		"KHR_TEST_CGROUP_VERSION=v2",
		"KHR_TEST_COLLECTED_AT=2026-05-15T16:00:00Z",
		"KHR_TEST_TELEMETRY_NOW=2026-05-15T15:00:00Z",
	)
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
		t.Fatalf("golden mismatch %s\n--- want ---\n%s\n--- got ---\n%s", goldenPath, wantStr, stdout.String())
	}
}
