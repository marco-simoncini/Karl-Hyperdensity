package agent

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadAndValidateConfig(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "cfg.yaml")
	content := `
apiVersion: khr.karl.io/v1alpha1
kind: KHRLinuxAgentConfig
metadata:
  name: test
spec:
  agentId: test-agent
  linuxOnly: true
  logLevel: debug
  cellWatchMode: disabled
`
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	cfg, err := LoadConfig(p)
	if err != nil {
		t.Fatal(err)
	}
	if errs := ValidateConfig(cfg); len(errs) != 0 {
		t.Fatalf("expected valid, got %v", errs)
	}
}

func TestValidateConfigRejectsNonLinux(t *testing.T) {
	cfg := &Config{}
	cfg.Spec.AgentID = "x"
	cfg.Spec.LinuxOnly = false
	errs := ValidateConfig(cfg)
	if len(errs) == 0 {
		t.Fatal("expected error")
	}
}
