package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestCLIIndexEvidenceLocalJSON(t *testing.T) {
	root := moduleRoot(t)
	inPath := filepath.Join(root, "examples", "grandepadre", "evidence-store", "ingest-request-ready.yaml")
	outPath := filepath.Join(t.TempDir(), "idx.json")
	cmdDir := filepath.Join(root, "cmd", "khr-linux-agent")
	cmd := exec.Command("go", "run", ".",
		"-mode", "index-evidence-local",
		"-ingest-request-input", inPath,
		"-index-output", outPath,
		"-query", "ready",
	)
	cmd.Dir = cmdDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go run: %v\n%s", err, out)
	}
	var rep map[string]interface{}
	if err := json.Unmarshal(out, &rep); err != nil {
		t.Fatalf("stdout json: %v", err)
	}
	if rep["mutationsForbidden"] != true {
		t.Fatalf("mutationsForbidden: %v", rep["mutationsForbidden"])
	}
	body, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(body, &rep); err != nil {
		t.Fatal(err)
	}
	if rep["indexedCount"].(float64) != 1 {
		t.Fatalf("indexedCount: %v", rep["indexedCount"])
	}
}
