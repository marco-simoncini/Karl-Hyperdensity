package host

import (
	"encoding/json"
	"testing"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/crdv1alpha1"
)

func testHostConfig() *Config {
	cfg := &Config{}
	cfg.Spec.HostID = "karl-host-runtime-karl-metal-01"
	cfg.Spec.LinuxOnly = true
	cfg.Spec.SandboxMode = true
	cfg.Spec.AllowedNamespaces = []string{"khr-runtime-sandbox"}
	cfg.Spec.AllowedLabels = map[string]string{"khr.karl.io/sandbox": "true"}
	return cfg
}

func TestBuildHostStatusSandbox(t *testing.T) {
	h := BuildHostStatus(testHostConfig(), "karl-metal-01", []crdv1alpha1.ObjectRef{
		{Name: "khr-runtime-sandbox-port", Namespace: "khr-runtime-sandbox"},
	})
	if h.Kind != "Host" || h.Spec.RuntimeMode != "sandbox" {
		t.Fatalf("host=%+v", h)
	}
	if h.Status.SafetyMode != "sandbox" || h.Status.Phase != "Observed" {
		t.Fatalf("status=%+v", h.Status)
	}
	if h.Status.RuntimeVersion == "" {
		t.Fatal("runtimeVersion required")
	}
}

func TestBuildHostStatusNoProductionMutation(t *testing.T) {
	h := BuildHostStatus(testHostConfig(), "karl-metal-01", nil)
	for _, c := range h.Status.Conditions {
		if c.Type == "SandboxOnly" && c.Status != "True" {
			t.Fatalf("conditions=%+v", h.Status.Conditions)
		}
	}
	var caps CapabilitiesReport
	if err := json.Unmarshal(h.Status.Capabilities, &caps); err != nil {
		t.Fatal(err)
	}
	if len(caps.BlockedSurfaces) == 0 {
		t.Fatal("expected blocked surfaces")
	}
}
