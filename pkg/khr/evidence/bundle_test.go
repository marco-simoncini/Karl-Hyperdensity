package evidence

import (
	"testing"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/cgroup"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/crdv1alpha1"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/discovery"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/resourcelease"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/telemetry"
)

func cellDemo() *crdv1alpha1.Cell {
	return &crdv1alpha1.Cell{
		APIVersion: "runtime.karl.io/v1alpha1",
		Kind:       "Cell",
		Metadata: crdv1alpha1.ObjectMeta{
			Name:      "demo-cell",
			Namespace: "karl-sandbox",
		},
	}
}

func TestSummarizeReadyNoDryRun(t *testing.T) {
	disc := &discovery.CgroupDiscoveryOutput{
		AgentID:        "a1",
		SelectedPath:   "/cg/scope",
		BlockedReasons: []string{},
		Warnings:       []string{},
	}
	tel := TelemetrySnapshot{
		Skipped:       false,
		TelemetryMode: "read-only",
		CgroupPath:    "/cg/scope",
		Evidence: telemetry.Evidence{
			Confidence:     "high",
			Warnings:       []string{},
			BlockedReasons: []string{},
		},
	}
	dry := DryRunSkippedPayload("no lease or resource port inputs provided for optional dry-run")
	b := BuildCollectEvidenceBundle("0.0.1-sprint13", "a1", cellDemo(), disc, tel, dry, "")
	if !b.EvidenceSummary.ReadyForGrandePadre {
		t.Fatalf("expected ready")
	}
	if b.EvidenceSummary.Confidence != "high" {
		t.Fatalf("confidence: got %q", b.EvidenceSummary.Confidence)
	}
	if len(b.EvidenceSummary.BlockedReasons) != 0 {
		t.Fatalf("blocked: %v", b.EvidenceSummary.BlockedReasons)
	}
}

func TestSummarizeBlockedNoSelectedPath(t *testing.T) {
	disc := &discovery.CgroupDiscoveryOutput{
		AgentID:        "a1",
		SelectedPath:   "",
		BlockedReasons: []string{"no discoverable cgroup directory matched heuristics or providerHandle under scannedRoot (read-only discovery only)"},
		Warnings:       []string{"cannot resolve symlinks"},
	}
	cellRef := &telemetry.CellRef{Kind: "Cell", Name: "demo-cell", Namespace: "karl-sandbox"}
	tel := TelemetrySnapshotSkipped("", cellRef, "discovery did not resolve selectedPath; telemetry skipped")
	dry := DryRunSkippedPayload("no lease or resource port inputs provided for optional dry-run")
	b := BuildCollectEvidenceBundle("0.0.1-sprint13", "a1", cellDemo(), disc, tel, dry, "")
	if b.EvidenceSummary.ReadyForGrandePadre {
		t.Fatalf("expected not ready")
	}
	if b.EvidenceSummary.Confidence != "low" {
		t.Fatalf("confidence: got %q", b.EvidenceSummary.Confidence)
	}
	if len(b.EvidenceSummary.BlockedReasons) < 2 {
		t.Fatalf("expected aggregated blocked reasons, got %v", b.EvidenceSummary.BlockedReasons)
	}
}

func TestWarningsAggregationLeasePortPartial(t *testing.T) {
	disc := &discovery.CgroupDiscoveryOutput{
		AgentID:        "a1",
		SelectedPath:   "/x",
		BlockedReasons: nil,
		Warnings:       []string{"discovery warn"},
	}
	tel := TelemetrySnapshot{
		Skipped:       false,
		TelemetryMode: "read-only",
		CgroupPath:    "/x",
		Evidence: telemetry.Evidence{
			Confidence:     "high",
			Warnings:       []string{"telemetry warn"},
			BlockedReasons: nil,
		},
	}
	partial := "collect-evidence: dry-run skipped: both -lease-input and -resource-port-input are required when including ResourceLease simulation"
	dry := DryRunSkippedPayload(partial)
	b := BuildCollectEvidenceBundle("0.0.1-sprint13", "a1", cellDemo(), disc, tel, dry, partial)
	w := b.EvidenceSummary.Warnings
	if len(w) < 3 {
		t.Fatalf("warnings: %v", w)
	}
}

func TestSummarizeDryRunBlockedAggregates(t *testing.T) {
	disc := &discovery.CgroupDiscoveryOutput{
		AgentID:        "a1",
		SelectedPath:   "/cg/scope",
		BlockedReasons: nil,
	}
	tel := TelemetrySnapshot{
		Skipped:       false,
		TelemetryMode: "read-only",
		CgroupPath:    "/cg/scope",
		Evidence: telemetry.Evidence{
			Confidence:     "high",
			BlockedReasons: nil,
		},
	}
	lease := resourcelease.DryRunResult{Allowed: false, Blocked: true, Reason: "contract conflict (test)"}
	dry := DryRunPayloadFromResult(lease, cgroupPlanStub(), true, false, false, nil)
	b := BuildCollectEvidenceBundle("0.0.1-sprint13", "a1", cellDemo(), disc, tel, dry, "")
	if b.EvidenceSummary.ReadyForGrandePadre {
		t.Fatalf("expected not ready when dry-run blocked")
	}
	want := "resource lease dry-run blocked: contract conflict (test)"
	for _, br := range b.EvidenceSummary.BlockedReasons {
		if br == want {
			return
		}
	}
	t.Fatalf("blocked reasons: %v", b.EvidenceSummary.BlockedReasons)
}

func cgroupPlanStub() cgroup.EnvelopePlan {
	return cgroup.EnvelopePlan{CgroupVersion: cgroup.V2, WouldWrite: false}
}
