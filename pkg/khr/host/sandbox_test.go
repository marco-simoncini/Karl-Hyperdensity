package host

import "testing"

func testConfig() *Config {
	cfg := &Config{}
	cfg.Spec.HostID = "h1"
	cfg.Spec.LinuxOnly = true
	cfg.Spec.SandboxMode = true
	cfg.Spec.AllowedNamespaces = []string{"khr-runtime-sandbox"}
	cfg.Spec.AllowedLabels = map[string]string{"khr.karl.io/sandbox": "true"}
	return cfg
}

func TestSandboxApplyBlockedByDefault(t *testing.T) {
	gate := SandboxApplyAllowed(testConfig(), "khr-runtime-sandbox", map[string]string{"khr.karl.io/sandbox": "true"})
	if gate.Allowed {
		t.Fatal("apply must be blocked by default")
	}
}

func TestSandboxApplyAllowedWhenEnabled(t *testing.T) {
	cfg := testConfig()
	cfg.Spec.SandboxApplyEnabled = true
	gate := SandboxApplyAllowed(cfg, "khr-runtime-sandbox", map[string]string{"khr.karl.io/sandbox": "true"})
	if !gate.Allowed {
		t.Fatalf("gate=%+v", gate)
	}
}

func TestSandboxApplyNamespaceAllowlist(t *testing.T) {
	cfg := testConfig()
	cfg.Spec.SandboxApplyEnabled = true
	gate := SandboxApplyAllowed(cfg, "other-ns", map[string]string{"khr.karl.io/sandbox": "true"})
	if gate.Allowed {
		t.Fatal("namespace must be blocked")
	}
}
