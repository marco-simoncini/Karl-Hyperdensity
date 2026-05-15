package windowsfluidvirt

import (
	"encoding/json"
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

func TestRealEvidenceMasterWin11Ready(t *testing.T) {
	fixture := mustLoadComplianceFixture(t, "master-win11-real-evidence.ready.json")
	result := EvaluateWindowsHyperdensityReadyCompliance(fixture.Input)
	logComplianceReplayResult(t, "master-win11-real-evidence", result)
	if result.CompliancePhase != ComplianceHyperdensityReadyWindowsShell {
		t.Fatalf("expected ready phase got=%s blockers=%v", result.CompliancePhase, result.Blockers)
	}
}

func TestRealEvidenceMasterWin11PoolChildReady(t *testing.T) {
	fixture := mustLoadComplianceFixture(t, "master-win11-pool-child-real-evidence.ready.json")
	result := EvaluateWindowsHyperdensityReadyCompliance(fixture.Input)
	logComplianceReplayResult(t, "master-win11-pool-child-real-evidence", result)
	if result.CompliancePhase != ComplianceHyperdensityReadyWindowsShell {
		t.Fatalf("expected ready phase got=%s blockers=%v", result.CompliancePhase, result.Blockers)
	}
}

func TestRealEvidenceMasterWin11PoolScalingMechanismBlocked(t *testing.T) {
	fixture := mustLoadComplianceFixture(t, "master-win11-pool-scaling-mechanism.blocked.json")
	result := EvaluateWindowsHyperdensityReadyCompliance(fixture.Input)
	logComplianceReplayResult(t, "master-win11-pool-scaling-mechanism", result)
	if result.CompliancePhase != ComplianceBlockedWithRemediation {
		t.Fatalf("expected blocked phase got=%s blockers=%v", result.CompliancePhase, result.Blockers)
	}
	assertHas(t, result.Blockers, BlockerPoolScalingAsMechanism)
}

func TestRealEvidenceMissingCPUActuatorBlocked(t *testing.T) {
	fixture := mustLoadComplianceFixture(t, "compliance-missing-cpu-actuator.blocked.json")
	result := EvaluateWindowsHyperdensityReadyCompliance(fixture.Input)
	logComplianceReplayResult(t, "missing-cpu-actuator", result)
	if result.CompliancePhase != ComplianceBlockedWithRemediation {
		t.Fatalf("expected blocked phase got=%s blockers=%v", result.CompliancePhase, result.Blockers)
	}
	assertHas(t, result.Blockers, BlockerNodeFluidActuatorUnavailable)
}

func TestRealEvidenceMissingRAMBalloonBlocked(t *testing.T) {
	fixture := mustLoadComplianceFixture(t, "compliance-missing-ram-balloon.blocked.json")
	result := EvaluateWindowsHyperdensityReadyCompliance(fixture.Input)
	logComplianceReplayResult(t, "missing-ram-balloon", result)
	if result.CompliancePhase != ComplianceBlockedWithRemediation {
		t.Fatalf("expected blocked phase got=%s blockers=%v", result.CompliancePhase, result.Blockers)
	}
	assertHas(t, result.Blockers, BlockerRAMBalloonUnavailable)
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

func mustLoadComplianceFixture(t *testing.T, name string) WindowsHyperdensityComplianceReplayFixture {
	t.Helper()
	fixture, err := LoadWindowsHyperdensityComplianceReplayFixture(windowsComplianceFixtureAbsPath(t, name))
	if err != nil {
		t.Fatalf("load fixture %s: %v", name, err)
	}
	return fixture
}

func logComplianceReplayResult(t *testing.T, scenario string, result EvaluateWindowsHyperdensityReadyComplianceOutput) {
	t.Helper()
	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("marshal replay result: %v", err)
	}
	t.Logf("scenario=%s replay=%s", scenario, string(data))
}
