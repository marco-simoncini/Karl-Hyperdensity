package host

import "testing"

func testConfig() *Config {
	cfg := &Config{}
	cfg.Spec.HostID = "h1"
	cfg.Spec.LinuxOnly = true
	cfg.Spec.SandboxMode = true
	cfg.Spec.AllowedNamespaces = []string{"karl-sandbox"}
	cfg.Spec.AllowedLabels = map[string]string{"karl.io/khr-sandbox": "true"}
	return cfg
}

func TestSandboxApplyBlockedByDefault(t *testing.T) {
	gate := SandboxApplyAllowed(testConfig(), "karl-sandbox", map[string]string{"karl.io/khr-sandbox": "true"})
	if gate.Allowed {
		t.Fatal("apply must be blocked by default")
	}
}

func TestSandboxApplyAllowedWhenEnabled(t *testing.T) {
	cfg := testConfig()
	cfg.Spec.SandboxApplyEnabled = true
	gate := SandboxApplyAllowed(cfg, "karl-sandbox", map[string]string{"karl.io/khr-sandbox": "true"})
	if !gate.Allowed {
		t.Fatalf("gate=%+v", gate)
	}
}

func TestSandboxApplyNamespaceAllowlist(t *testing.T) {
	cfg := testConfig()
	cfg.Spec.SandboxApplyEnabled = true
	gate := SandboxApplyAllowed(cfg, "other-ns", map[string]string{"karl.io/khr-sandbox": "true"})
	if gate.Allowed {
		t.Fatal("namespace must be blocked")
	}
}
