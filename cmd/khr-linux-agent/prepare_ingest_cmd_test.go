package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestCLIPrepareIngestRequestWritesYAML(t *testing.T) {
	t.Setenv("KHR_TEST_COLLECTED_AT", "2026-05-15T16:00:00Z")
	t.Setenv("KHR_TEST_TELEMETRY_NOW", "2026-05-15T15:00:00Z")
	root := moduleRoot(t)
	raw, err := os.ReadFile(filepath.Join(root, "examples", "khr", "evidence", "collect-evidence-output-ready.json"))
	if err != nil {
		t.Fatal(err)
	}
	replaced := strings.ReplaceAll(string(raw), "__CGROUP_ROOT__", "/tmp/khr-evidence-example-root")
	dir := t.TempDir()
	bundlePath := filepath.Join(dir, "bundle.json")
	if err := os.WriteFile(bundlePath, []byte(replaced), 0o600); err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(root, "examples", "khr", "evidence-integrity", "manifest-none.json")
	digestPath := filepath.Join(root, "examples", "khr", "evidence-integrity", "digest.txt")
	outPath := filepath.Join(dir, "ingest.yaml")
	cmdDir := filepath.Join(root, "cmd", "khr-linux-agent")
	cmd := exec.Command("go", "run", ".",
		"-mode", "prepare-ingest-request",
		"-bundle-input", bundlePath,
		"-manifest-input", manifestPath,
		"-digest-input", digestPath,
		"-ingest-request-output", outPath,
		"-ingest-request-format", "yaml",
		"-dry-run-only=true",
	)
	cmd.Dir = cmdDir
	cmd.Env = append(os.Environ(),
		"KHR_TEST_COLLECTED_AT=2026-05-15T16:00:00Z",
		"KHR_TEST_TELEMETRY_NOW=2026-05-15T15:00:00Z",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go run: %v\n%s", err, out)
	}
	body, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(body), "kind: EvidenceIngestRequest") && !strings.Contains(string(body), "EvidenceIngestRequest") {
		t.Fatalf("unexpected output: %s", body)
	}
	if !strings.Contains(string(body), "dryRunOnly: true") {
		t.Fatalf("expected dryRunOnly true: %s", body)
	}
}

func TestCLIPrepareIngestRequestAliasEvidenceFlags(t *testing.T) {
	t.Setenv("KHR_TEST_COLLECTED_AT", "2026-05-15T16:00:00Z")
	t.Setenv("KHR_TEST_TELEMETRY_NOW", "2026-05-15T15:00:00Z")
	root := moduleRoot(t)
	raw, err := os.ReadFile(filepath.Join(root, "examples", "khr", "evidence", "collect-evidence-output-ready.json"))
	if err != nil {
		t.Fatal(err)
	}
	replaced := strings.ReplaceAll(string(raw), "__CGROUP_ROOT__", "/tmp/khr-evidence-example-root")
	dir := t.TempDir()
	bundlePath := filepath.Join(dir, "bundle.json")
	if err := os.WriteFile(bundlePath, []byte(replaced), 0o600); err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(root, "examples", "khr", "evidence-integrity", "manifest-none.json")
	digestPath := filepath.Join(root, "examples", "khr", "evidence-integrity", "digest.txt")
	outPath := filepath.Join(dir, "ingest-alias.yaml")
	cmdDir := filepath.Join(root, "cmd", "khr-linux-agent")
	cmd := exec.Command("go", "run", ".",
		"-mode", "prepare-ingest-request",
		"-evidence-output", bundlePath,
		"-evidence-manifest-output", manifestPath,
		"-evidence-digest-output", digestPath,
		"-ingest-request-output", outPath,
		"-ingest-request-format", "yaml",
		"-dry-run-only=true",
	)
	cmd.Dir = cmdDir
	cmd.Env = append(os.Environ(),
		"KHR_TEST_COLLECTED_AT=2026-05-15T16:00:00Z",
		"KHR_TEST_TELEMETRY_NOW=2026-05-15T15:00:00Z",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go run: %v\n%s", err, out)
	}
	body, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(body), "kind: EvidenceIngestRequest") && !strings.Contains(string(body), "EvidenceIngestRequest") {
		t.Fatalf("unexpected output: %s", body)
	}
}
