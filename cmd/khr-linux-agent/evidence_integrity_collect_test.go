package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestCollectEvidenceWritesManifestAndDigest(t *testing.T) {
	t.Setenv("KHR_TEST_CGROUP_VERSION", "v2")
	t.Setenv("KHR_TEST_COLLECTED_AT", "2026-05-15T16:00:00Z")
	t.Setenv("KHR_TEST_TELEMETRY_NOW", "2026-05-15T15:00:00Z")
	t.Setenv("KHR_TEST_INTEGRITY_NOW", "2026-07-01T12:00:00Z")
	t.Setenv("KHR_TEST_INTEGRITY_CHAIN_STUB", "1")
	root := moduleRoot(t)
	cfg := filepath.Join(root, "examples", "khr", "khr-linux-agent-config.yaml")
	cell := filepath.Join(root, "examples", "khr", "evidence", "collect-evidence-input-cell.json")
	cmdDir := filepath.Join(root, "cmd", "khr-linux-agent")
	r := t.TempDir()
	if err := os.MkdirAll(filepath.Join(r, "karl.slice", "karl-shell-dev-linux-systemd-001.scope"), 0o755); err != nil {
		t.Fatal(err)
	}
	scope := filepath.Join(r, "karl.slice", "karl-shell-dev-linux-systemd-001.scope")
	writeMetricFixtures(t, scope)
	man := filepath.Join(r, "manifest.json")
	dig := filepath.Join(r, "digest.txt")
	args := []string{
		"run", ".",
		"-mode", "collect-evidence",
		"-config", cfg,
		"-cell-input", cell,
		"-cgroup-root", r,
		"-evidence-manifest-output", man,
		"-evidence-digest-output", dig,
	}
	cmd := exec.Command("go", args...)
	cmd.Dir = cmdDir
	cmd.Env = append(os.Environ(),
		"KHR_TEST_CGROUP_VERSION=v2",
		"KHR_TEST_COLLECTED_AT=2026-05-15T16:00:00Z",
		"KHR_TEST_TELEMETRY_NOW=2026-05-15T15:00:00Z",
		"KHR_TEST_INTEGRITY_NOW=2026-07-01T12:00:00Z",
		"KHR_TEST_INTEGRITY_CHAIN_STUB=1",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go run: %v\n%s", err, out)
	}
	dline, err := os.ReadFile(dig)
	if err != nil {
		t.Fatal(err)
	}
	hex := strings.TrimSpace(string(dline))
	if len(hex) != 64 {
		t.Fatalf("digest hex length: %q", hex)
	}
	var mf map[string]any
	raw, err := os.ReadFile(man)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(raw, &mf); err != nil {
		t.Fatal(err)
	}
	if mf["bundleSha256"] != hex {
		t.Fatalf("manifest sha != digest: %v vs %s", mf["bundleSha256"], hex)
	}
	if mf["signingMode"] != "none" || mf["signaturePresent"] != false {
		t.Fatalf("unexpected signing fields: %+v", mf)
	}
	if mf["mutationsForbidden"] != true || mf["sourceMode"] != "collect-evidence" {
		t.Fatalf("contract fields: %+v", mf)
	}
}

func TestCollectEvidenceLocalDevWithoutKeyFails(t *testing.T) {
	t.Setenv("KHR_TEST_CGROUP_VERSION", "v2")
	root := moduleRoot(t)
	cfg := filepath.Join(root, "examples", "khr", "khr-linux-agent-config.yaml")
	cell := filepath.Join(root, "examples", "khr", "evidence", "collect-evidence-input-cell.json")
	cmdDir := filepath.Join(root, "cmd", "khr-linux-agent")
	r := t.TempDir()
	man := filepath.Join(r, "m.json")
	args := []string{
		"run", ".",
		"-mode", "collect-evidence",
		"-config", cfg,
		"-cell-input", cell,
		"-cgroup-root", r,
		"-signing-mode", "local-dev",
		"-evidence-manifest-output", man,
	}
	cmd := exec.Command("go", args...)
	cmd.Dir = cmdDir
	cmd.Env = append(os.Environ(), "KHR_TEST_CGROUP_VERSION=v2")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected failure, got: %s", out)
	}
	if !strings.Contains(string(out), "signing-key-file") {
		t.Fatalf("stderr/stdout missing key hint: %s", out)
	}
}

func TestCollectEvidenceLocalDevWithKeyWritesSignature(t *testing.T) {
	t.Setenv("KHR_TEST_CGROUP_VERSION", "v2")
	t.Setenv("KHR_TEST_COLLECTED_AT", "2026-05-15T16:00:00Z")
	t.Setenv("KHR_TEST_TELEMETRY_NOW", "2026-05-15T15:00:00Z")
	t.Setenv("KHR_TEST_INTEGRITY_NOW", "2026-07-01T12:00:00Z")
	t.Setenv("KHR_TEST_INTEGRITY_CHAIN_STUB", "1")
	root := moduleRoot(t)
	cfg := filepath.Join(root, "examples", "khr", "khr-linux-agent-config.yaml")
	cell := filepath.Join(root, "examples", "khr", "evidence", "collect-evidence-input-cell.json")
	key := filepath.Join(root, "examples", "khr", "evidence-integrity", "local-dev-key.example")
	cmdDir := filepath.Join(root, "cmd", "khr-linux-agent")
	r := t.TempDir()
	if err := os.MkdirAll(filepath.Join(r, "karl.slice", "karl-shell-dev-linux-systemd-001.scope"), 0o755); err != nil {
		t.Fatal(err)
	}
	scope := filepath.Join(r, "karl.slice", "karl-shell-dev-linux-systemd-001.scope")
	writeMetricFixtures(t, scope)
	man := filepath.Join(r, "manifest.json")
	args := []string{
		"run", ".",
		"-mode", "collect-evidence",
		"-config", cfg,
		"-cell-input", cell,
		"-cgroup-root", r,
		"-signing-mode", "local-dev",
		"-signing-key-file", key,
		"-evidence-manifest-output", man,
	}
	cmd := exec.Command("go", args...)
	cmd.Dir = cmdDir
	cmd.Env = append(os.Environ(),
		"KHR_TEST_CGROUP_VERSION=v2",
		"KHR_TEST_COLLECTED_AT=2026-05-15T16:00:00Z",
		"KHR_TEST_TELEMETRY_NOW=2026-05-15T15:00:00Z",
		"KHR_TEST_INTEGRITY_NOW=2026-07-01T12:00:00Z",
		"KHR_TEST_INTEGRITY_CHAIN_STUB=1",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go run: %v\n%s", err, out)
	}
	raw, err := os.ReadFile(man)
	if err != nil {
		t.Fatal(err)
	}
	var mf map[string]any
	if err := json.Unmarshal(raw, &mf); err != nil {
		t.Fatal(err)
	}
	if mf["signaturePresent"] != true {
		t.Fatalf("expected signature: %+v", mf)
	}
	sig, _ := mf["signature"].(string)
	if sig == "" {
		t.Fatal("missing signature")
	}
}
