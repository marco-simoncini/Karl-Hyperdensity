package controlgraph

import (
	"testing"
	"time"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/certregistry"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/lanediscovery"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/nativelive"
)

func sampleDiscovery() lanediscovery.Result {
	return lanediscovery.Result{
		ClusterContext: "karl-metal-01@ovh",
		ObservedAt:     time.Now().UTC().Format(time.RFC3339),
		DiscoveredHosts: []lanediscovery.DiscoveredHost{{
			HostID: "host/node-1", NodeName: "node-1", Ready: true,
		}},
		DiscoveredShells: []lanediscovery.DiscoveredShell{{
			Ref: "khr-runtime-sandbox/Shell/target", Namespace: "khr-runtime-sandbox", Name: "target",
		}},
		DiscoveredCells: []lanediscovery.DiscoveredCell{{
			Ref: "khr-runtime-sandbox/Cell/target", ShellRef: "khr-runtime-sandbox/Shell/target",
			Namespace: "khr-runtime-sandbox", Name: "target", NodeName: "node-1", Running: true,
		}},
		DiscoveredResourcePorts: []lanediscovery.DiscoveredResourcePort{{
			Ref: "khr-runtime-sandbox/ResourcePort/target-port",
			ShellRef: "khr-runtime-sandbox/Shell/target",
			CellRef:  "khr-runtime-sandbox/Cell/target",
			Lane: lanediscovery.LaneNativeLive, Classification: lanediscovery.ClassificationNativeLive,
		}},
	}
}

func TestBuildGraphConsistency(t *testing.T) {
	g := Build(BuildInput{
		Sprint: "KHR-X", ClusterContext: "karl-metal-01@ovh",
		Discovery: sampleDiscovery(), Now: time.Now().UTC(),
	})
	if g.Health.NodeCount < 4 || g.Health.EdgeCount < 2 {
		t.Fatalf("graph too small nodes=%d edges=%d", g.Health.NodeCount, g.Health.EdgeCount)
	}
	if err := CheckCorrelationIntegrity(g); err != nil {
		t.Fatal(err)
	}
}

func TestOrphanDetection(t *testing.T) {
	g := Build(BuildInput{Sprint: "KHR-X", Discovery: sampleDiscovery(), Now: time.Now().UTC()})
	// inject orphan port
	g.Nodes = append(g.Nodes, Node{
		ID: nodeID(KindResourcePort, "orphan/port"), Kind: KindResourcePort, Ref: "orphan/port",
		CorrelationID: g.CorrelationID, State: StateObserved,
	})
	AnnotateHealth(&g)
	if g.Health.OrphanCount == 0 {
		t.Fatal("expected orphan detection")
	}
}

func TestStaleNodeDetection(t *testing.T) {
	reg := certregistry.GenerateFromSummary("KHR-X", nativelive.CertificationSummary{
		CertificationID: "khr-native-live-certification-v1",
		Lane:            lanediscovery.LaneNativeLive,
		Status:          "certified",
		Invariants:      nativelive.Invariants{NoRestart: true, NoRollout: true, NoRecreate: true},
		Scores:          nativelive.CertificationScores{ContinuityScore: 1},
		ContinuityProof: nativelive.ContinuityCertificationProof{ShellContinuityPreserved: true},
	}, "evidence/cert", 3600, time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC))
	g := Build(BuildInput{
		Sprint: "KHR-X", Discovery: sampleDiscovery(), Registry: &reg, Now: time.Now().UTC(),
	})
	if g.Health.StaleCount == 0 {
		t.Fatal("expected stale certification node")
	}
}

func TestNoMutationFlags(t *testing.T) {
	g := Build(BuildInput{Sprint: "KHR-X", Discovery: sampleDiscovery(), Now: time.Now().UTC()})
	if !g.ReadOnly || !g.NoApply || !g.NoMutation || !g.NoAutonomousOrchestration {
		t.Fatalf("safety=%+v", g)
	}
}
