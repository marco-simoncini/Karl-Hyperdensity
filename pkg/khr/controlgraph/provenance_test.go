package controlgraph

import (
	"testing"
	"time"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/certregistry"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/lanediscovery"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/nativelive"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/provenance"
)

func TestVerifyLineageIntegrityPass(t *testing.T) {
	evidence := []byte(`{"cert":true}`)
	reg := certregistry.GenerateFromSummaryWithEvidence("KHR-Y", nativelive.CertificationSummary{
		Lane: lanediscovery.LaneNativeLive, Status: nativelive.CertificationCertified,
		Invariants: nativelive.Invariants{NoRestart: true, NoRollout: true, NoRecreate: true},
	}, "evidence/cert", 3600, time.Now().UTC(), evidence, "karl-metal-01@ovh")
	g := Build(BuildInput{
		Sprint: "KHR-Y", ClusterContext: "karl-metal-01@ovh",
		Discovery: sampleDiscovery(), Registry: &reg, Now: time.Now().UTC(),
	})
	sum, err := VerifyLineageIntegrity(g, time.Now().UTC())
	if err != nil || !sum.LineageIntegrity {
		t.Fatalf("sum=%+v err=%v", sum, err)
	}
}

func TestLineageMismatchWarning(t *testing.T) {
	g := Build(BuildInput{
		Sprint: "KHR-Y", ClusterContext: "karl-metal-01@ovh",
		Discovery: sampleDiscovery(), Now: time.Now().UTC(),
	})
	for i := range g.Nodes {
		if g.Nodes[i].Kind == KindCell {
			g.Nodes[i].Provenance.SourceCluster = "other-cluster"
		}
	}
	sum, _ := VerifyLineageIntegrity(g, time.Now().UTC())
	if sum.LineageIntegrity {
		t.Fatal("expected lineage integrity failure")
	}
	if sum.ProvenanceState != provenance.StateMismatch {
		t.Fatalf("state=%q", sum.ProvenanceState)
	}
}
