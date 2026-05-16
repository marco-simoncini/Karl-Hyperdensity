package controlgraph

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/provenance"
)

// VerifyLineageIntegrity checks graph lineage and provenance consistency (KHR-Y).
func VerifyLineageIntegrity(g Graph, now time.Time) (provenance.ValidationSummary, error) {
	sum := provenance.ValidationSummary{
		ModelID: provenance.ModelID, Sprint: g.Sprint,
		ValidatedAt: now.UTC().Format(time.RFC3339),
		ReadOnly: true, NoAutonomousOrchestration: true, NoMutation: true, NoApply: true,
		GraphLineageVerified: true, LineageIntegrity: true,
		ProvenanceState:      provenance.StateTrusted,
	}
	if g.Provenance.EvidenceFingerprint != "" && provenance.IsStaleProvenance(g.Provenance, now, 0) {
		sum.TrustWarnings = append(sum.TrustWarnings, "stale graph provenance")
		sum.ProvenanceState = provenance.StateStale
		sum.LineageIntegrity = false
	}
	lineageSeen := map[string]string{}
	for _, n := range g.Nodes {
		if n.Provenance.LineageHash == "" {
			continue
		}
		if prev, ok := lineageSeen[n.Provenance.LineageHash]; ok && prev != n.Kind {
			// same lineage hash across kinds is ok
		}
		lineageSeen[n.Provenance.LineageHash] = n.Kind
		if n.Provenance.SourceCluster != "" && g.ClusterContext != "" &&
			n.Provenance.SourceCluster != g.ClusterContext {
			sum.TrustWarnings = append(sum.TrustWarnings, "lineage cluster mismatch on "+n.Ref)
			sum.LineageIntegrity = false
			sum.GraphLineageVerified = false
		}
	}
	// Provenance edges (RelProvenance) assert shared lineage groups; structural edges may span derived hashes.
	for _, e := range g.Edges {
		if e.Relationship != RelProvenance {
			continue
		}
		from, fok := nodeByID(g, e.From)
		to, tok := nodeByID(g, e.To)
		if !fok || !tok {
			continue
		}
		if from.Provenance.LineageHash != "" && to.Provenance.LineageHash != "" &&
			from.Provenance.LineageHash != to.Provenance.LineageHash {
			sum.TrustWarnings = append(sum.TrustWarnings,
				fmt.Sprintf("provenance lineage mismatch: %s -> %s", e.From, e.To))
			sum.LineageIntegrity = false
		}
	}
	if !sum.LineageIntegrity {
		sum.ProvenanceState = provenance.StateMismatch
	}
	if g.Provenance.EvidenceFingerprint != "" {
		sum.EvidenceFingerprint = g.Provenance.EvidenceFingerprint
	}
	return sum, nil
}

func nodeByID(g Graph, id string) (Node, bool) {
	for _, n := range g.Nodes {
		if n.ID == id {
			return n, true
		}
	}
	return Node{}, false
}

func attachGraphProvenance(g *Graph, in BuildInput, rootCorr string) {
	if g == nil {
		return
	}
	payload, _ := json.Marshal(struct {
		GraphID string `json:"graphId"`
		Nodes   int    `json:"nodes"`
		Edges   int    `json:"edges"`
		Corr    string `json:"correlationId"`
	}{
		GraphID: g.GraphID, Nodes: len(g.Nodes), Edges: len(g.Edges), Corr: rootCorr,
	})
	g.Provenance = provenance.NewRecord("khr-control-graph", provenance.SourceContext{
		Cluster: in.ClusterContext, Namespace: "khr-runtime-sandbox", Lane: "native-live",
	}, rootCorr, payload, in.ObservedAt)
}

func addProvenanceEdges(g *Graph, addEdge func(from, to, rel string)) {
	if g == nil {
		return
	}
	byLineage := map[string][]string{}
	for _, n := range g.Nodes {
		if n.Provenance.LineageHash != "" {
			byLineage[n.Provenance.LineageHash] = append(byLineage[n.Provenance.LineageHash], n.ID)
		}
	}
	for _, ids := range byLineage {
		for i := 1; i < len(ids); i++ {
			addEdge(ids[0], ids[i], RelProvenance)
		}
	}
}
