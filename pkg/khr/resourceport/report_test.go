package resourceport

import (
	"testing"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/host"
)

func TestReportCandidateMock(t *testing.T) {
	cfg := &host.Config{}
	cfg.Spec.HostID = "host-1"
	cfg.Spec.LinuxOnly = true
	cfg.Spec.SandboxMode = true
	c := ReportCandidate(cfg, "karl-sandbox/Shell/x", "karl-sandbox/Cell/x", "karl-sandbox", "port-1")
	if c.Kind != "ResourcePort" || c.Spec.Provider != "khr.native" {
		t.Fatalf("candidate=%+v", c)
	}
	if c.Status.Phase != "Observed" {
		t.Fatal("status phase")
	}
}
