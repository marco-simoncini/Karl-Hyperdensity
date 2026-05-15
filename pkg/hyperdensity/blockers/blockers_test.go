package blockers

import "testing"

func TestNoWindowsLaneKnown(t *testing.T) {
	if !Known(IDNoWindowsLane) {
		t.Fatal("no_windows_lane must be known")
	}
	sev, ok := Severity(IDNoWindowsLane)
	if !ok || sev != SeverityCritical {
		t.Fatalf("severity: ok=%v sev=%q", ok, sev)
	}
}

func TestCatalogStableSeverity(t *testing.T) {
	required := []struct {
		id  string
		sev string
	}{
		{IDNoWindowsLane, SeverityCritical},
		{IDNoProductionMutation, SeverityCritical},
		{IDKeepWindowsLaneDisabled, SeverityHigh},
		{IDWindowsDisabled, SeverityInfo},
		{IDDryRunOnly, SeverityHigh},
		{IDRuntimeApplyDisabled, SeverityCritical},
		{IDUnsupportedBroadVMExecution, SeverityHigh},
		{IDUnsupportedBroadMemoryExec, SeverityHigh},
		{IDUnsupportedMultiContainerWiden, SeverityHigh},
		{IDUnsupportedBroadAutomation, SeverityHigh},
	}
	for _, tc := range required {
		if !Known(tc.id) {
			t.Fatalf("missing catalog entry %q", tc.id)
		}
		sev, ok := Severity(tc.id)
		if !ok || sev != tc.sev {
			t.Fatalf("%q: want severity %q, got %q ok=%v", tc.id, tc.sev, sev, ok)
		}
	}
	cat := Catalog()
	if len(cat) != len(required) {
		t.Fatalf("catalog len=%d want %d", len(cat), len(required))
	}
	for i := 1; i < len(cat); i++ {
		if cat[i].ID <= cat[i-1].ID {
			t.Fatalf("catalog not sorted: %q after %q", cat[i].ID, cat[i-1].ID)
		}
	}
}
