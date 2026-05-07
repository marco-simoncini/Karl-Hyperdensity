package windowsfluidvirt

import (
	"path/filepath"
	"testing"
)

func TestWindowsHyperdensityComplianceFixtureMatrix(t *testing.T) {
	fixtures := []string{
		"compliance-standalone-ready.json",
		"compliance-pool-child-ready.json",
		"compliance-missing-fluidshell.blocked.json",
		"compliance-missing-qmp.blocked.json",
		"compliance-missing-ram-balloon.blocked.json",
		"compliance-missing-cpu-actuator.blocked.json",
		"compliance-pool-scaling-mechanism.blocked.json",
		"compliance-cgroup-path-mismatch.blocked.json",
		"compliance-qemu-pid-mismatch.blocked.json",
	}
	for _, name := range fixtures {
		t.Run(name, func(t *testing.T) {
			fixture, err := LoadWindowsHyperdensityComplianceReplayFixture(windowsComplianceFixtureAbsPath(t, name))
			if err != nil {
				t.Fatalf("load fixture: %v", err)
			}
			result := EvaluateWindowsHyperdensityReadyCompliance(fixture.Input)
			if result.CompliancePhase != fixture.ExpectedPhase {
				t.Fatalf("phase mismatch expected=%s got=%s blockers=%v", fixture.ExpectedPhase, result.CompliancePhase, result.Blockers)
			}
			for _, blocker := range fixture.ExpectedBlockers {
				assertHas(t, result.Blockers, blocker)
			}
			for _, remediation := range fixture.ExpectedRemediations {
				assertHas(t, result.RemediationActions, remediation)
			}
		})
	}
}

func TestNodeActuatorSafetyFixtures(t *testing.T) {
	fixtures := []string{
		"actuator-stale-request.blocked.json",
	}
	for _, name := range fixtures {
		t.Run(name, func(t *testing.T) {
			fixture, err := LoadNodeActuatorSafetyReplayFixture(windowsComplianceFixtureAbsPath(t, name))
			if err != nil {
				t.Fatalf("load safety fixture: %v", err)
			}
			model := DefaultNodeFluidActuatorSafetyModel()
			if fixture.Model != nil {
				model = *fixture.Model
			}
			result := EvaluateNodeFluidActuatorSafety(fixture.Input, model)
			if result.Allowed != fixture.ExpectedAllowed {
				t.Fatalf("allowed mismatch expected=%v got=%v blockers=%v", fixture.ExpectedAllowed, result.Allowed, result.Blockers)
			}
			for _, blocker := range fixture.ExpectedBlockers {
				assertHas(t, result.Blockers, blocker)
			}
		})
	}
}

func TestWindowsCpuLeaseFixtures(t *testing.T) {
	fixtures := []string{
		"lease-updown-valid.accepted.json",
		"lease-vcpu-hotplug.rejected.json",
		"lease-vm-spec-patch.rejected.json",
	}
	for _, name := range fixtures {
		t.Run(name, func(t *testing.T) {
			fixture, err := LoadWindowsCpuLeaseReplayFixture(windowsComplianceFixtureAbsPath(t, name))
			if err != nil {
				t.Fatalf("load lease fixture: %v", err)
			}
			result := EvaluateWindowsCpuEntitlementLease(fixture.Lease)
			if result.Status != fixture.ExpectedStatus {
				t.Fatalf("status mismatch expected=%s got=%s blockers=%v", fixture.ExpectedStatus, result.Status, result.Blockers)
			}
			for _, blocker := range fixture.ExpectedBlockers {
				assertHas(t, result.Blockers, blocker)
			}
		})
	}
}

func windowsComplianceFixtureAbsPath(t *testing.T, name string) string {
	t.Helper()
	return filepath.Join(admissionRepoRoot(t), "examples", "windows-fluid-compliance-fixtures", name)
}
