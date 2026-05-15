package main

import (
	"encoding/json"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestCLIRecommendActionsLocalDeterministic(t *testing.T) {
	root := moduleRoot(t)
	ts := "2026-06-01T10:00:00Z"
	in1 := filepath.Join(root, "examples", "grandepadre", "recommendation", "recommendation-input-ready.yaml")
	cmdDir := filepath.Join(root, "cmd", "khr-linux-agent")
	cmd := exec.Command("go", "run", ".",
		"-mode", "recommend-actions-local",
		"-ingest-request-input", in1,
		"-recommendation-generated-at", ts,
	)
	cmd.Dir = cmdDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go run: %v\n%s", err, out)
	}
	var rep map[string]interface{}
	if err := json.Unmarshal(out, &rep); err != nil {
		t.Fatal(err)
	}
	if rep["mutationsForbidden"] != true || rep["applyAllowed"] != false {
		t.Fatalf("policy fields: %#v", rep)
	}
	slate, ok := rep["actionSlate"].(map[string]interface{})
	if !ok {
		t.Fatal("missing actionSlate")
	}
	if slate["generatedAt"] != ts {
		t.Fatalf("generatedAt=%v", slate["generatedAt"])
	}
}
