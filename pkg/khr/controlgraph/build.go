package controlgraph

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/actionapproval"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/certregistry"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/lanediscovery"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/policygates"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/provenance"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/resourcefuture"
)

// BuildInput aggregates read-only sources for graph export.
type BuildInput struct {
	Sprint         string
	ClusterContext string
	ObservedAt     time.Time
	Discovery      lanediscovery.Result
	Registry       *certregistry.Registry
	Simulation     *resourcefuture.SimulationResult
	Approvals      []actionapproval.ActionApproval
	Now            time.Time
}

// Build constructs the unified control graph from cluster observations.
func Build(in BuildInput) Graph {
	if in.Now.IsZero() {
		in.Now = time.Now().UTC()
	}
	if in.ObservedAt.IsZero() {
		in.ObservedAt = in.Now
	}
	rootCorr := rootCorrelationID(in.ClusterContext, in.ObservedAt)
	g := Graph{
		GraphID:                   GraphID,
		Sprint:                    in.Sprint,
		ClusterContext:            in.ClusterContext,
		ObservedAt:                in.ObservedAt.UTC().Format(time.RFC3339),
		CorrelationID:             rootCorr,
		ReadOnly:                  true,
		NoAutonomousOrchestration: true,
		NoMutation:                true,
		NoApply:                   true,
		Lineage: LineageSummary{
			RootCorrelationID: rootCorr,
			CorrelationIDs:    []string{rootCorr},
		},
	}
	nodeIndex := map[string]string{}

	addNode := func(n Node) {
		if _, ok := nodeIndex[n.ID]; ok {
			return
		}
		if n.CorrelationID == "" {
			n.CorrelationID = rootCorr
		}
		if n.Provenance.ProvenanceID == "" && in.ClusterContext != "" {
			lane := attrString(n.Attributes, "lane")
			if lane == "" {
				lane = "native-live"
			}
			n.Provenance = provenance.NewRecord("khr-control-graph-node", provenance.SourceContext{
				Cluster: in.ClusterContext, Namespace: "khr-runtime-sandbox", Lane: lane,
			}, rootCorr, []byte(n.Ref+"|"+n.Kind), in.ObservedAt)
		}
		if n.State == "" {
			n.State = StateObserved
		}
		g.Nodes = append(g.Nodes, n)
		nodeIndex[n.ID] = n.Kind
	}
	addEdge := func(from, to, rel string) {
		if from == "" || to == "" {
			return
		}
		g.Edges = append(g.Edges, Edge{
			ID:            edgeID(from, to, rel),
			From:          from,
			To:            to,
			Relationship:  rel,
			CorrelationID: rootCorr,
		})
	}

	obs := in.ObservedAt.UTC().Format(time.RFC3339)

	for _, h := range in.Discovery.DiscoveredHosts {
		id := nodeID(KindHost, h.HostID)
		addNode(Node{
			ID: id, Kind: KindHost, Ref: h.HostID,
			CorrelationID: rootCorr, ObservedAt: obs,
			Attributes: map[string]any{
				"nodeName": h.NodeName, "ready": h.Ready, "provider": h.Provider,
			},
		})
	}

	for _, sh := range in.Discovery.DiscoveredShells {
		id := nodeID(KindShell, sh.Ref)
		addNode(Node{
			ID: id, Kind: KindShell, Ref: sh.Ref,
			CorrelationID: lineageID(rootCorr, "shell", sh.Ref), ObservedAt: obs,
			Attributes: map[string]any{
				"namespace": sh.Namespace, "name": sh.Name,
				"osFamily": sh.OSFamily, "vmType": sh.VMType,
			},
		})
	}

	cellByRef := map[string]string{}
	for _, c := range in.Discovery.DiscoveredCells {
		id := nodeID(KindCell, c.Ref)
		cellByRef[c.Ref] = id
		corr := lineageID(rootCorr, "cell", c.Ref)
		addNode(Node{
			ID: id, Kind: KindCell, Ref: c.Ref,
			CorrelationID: corr, LineageParentID: nodeID(KindShell, c.ShellRef),
			ObservedAt: obs,
			Attributes: map[string]any{
				"namespace": c.Namespace, "name": c.Name, "running": c.Running,
				"nodeName": c.NodeName, "osFamily": c.OSFamily,
			},
		})
		if c.ShellRef != "" {
			addEdge(nodeID(KindShell, c.ShellRef), id, RelProjects)
		}
		if c.NodeName != "" {
			hostRef := "host/" + c.NodeName
			for _, h := range in.Discovery.DiscoveredHosts {
				if h.NodeName == c.NodeName {
					hostRef = h.HostID
					break
				}
			}
			addEdge(nodeID(KindHost, hostRef), id, RelHosts)
		}
	}

	portByRef := map[string]string{}
	for _, p := range in.Discovery.DiscoveredResourcePorts {
		id := nodeID(KindResourcePort, p.Ref)
		portByRef[p.Ref] = id
		addNode(Node{
			ID: id, Kind: KindResourcePort, Ref: p.Ref,
			CorrelationID: lineageID(rootCorr, "port", p.Ref),
			LineageParentID: nodeID(KindCell, p.CellRef),
			ObservedAt: obs,
			Attributes: map[string]any{
				"lane": p.Lane, "classification": p.Classification,
				"providerBinding": p.ProviderBinding,
				"liveScale": p.LiveScaleCapabilityObserved,
			},
		})
		if p.CellRef != "" {
			addEdge(cellByRef[p.CellRef], id, RelBinds)
		}
		leaseID := nodeID(KindResourceLease, p.Ref+"/lease")
		addNode(Node{
			ID: leaseID, Kind: KindResourceLease,
			Ref: p.Ref + "/ResourceLease/observed",
			CorrelationID: lineageID(rootCorr, "lease", p.Ref),
			LineageParentID: id, ObservedAt: obs,
			Attributes: map[string]any{"dryRunOnly": true, "simulationOnly": true},
		})
		addEdge(id, leaseID, RelBinds)
	}

	if in.Registry != nil {
		for _, e := range in.Registry.Entries {
			certID := nodeID(KindCertification, e.LaneID)
			stale := !certregistry.IsFresh(e, in.Now)
			state := StateObserved
			if stale {
				state = StateStale
			}
			certNode := Node{
				ID: certID, Kind: KindCertification, Ref: e.EvidenceRef,
				CorrelationID: lineageID(rootCorr, "cert", e.LaneID),
				State: state, Stale: stale, ObservedAt: e.LastCertifiedAt,
				Provenance: e.Provenance,
				Attributes: map[string]any{
					"laneId": e.LaneID, "certificationState": e.CertificationState,
					"continuityScore": e.ContinuityScore,
				},
			}
			addNode(certNode)
			for ref, pid := range portByRef {
				if strings.Contains(ref, "native-live") || e.LaneID == lanediscovery.LaneNativeLive {
					addEdge(certID, pid, RelCertifies)
					break
				}
			}
		}
	}

	gateID := nodeID(KindPolicyGate, "native-live-gates")
	gates := policygates.DefaultNativeLiveGates()
	addNode(Node{
		ID: gateID, Kind: KindPolicyGate, Ref: "policy/native-live",
		CorrelationID: lineageID(rootCorr, "gates", "native-live"),
		ObservedAt: obs,
		Attributes: map[string]any{
			"noRestart": gates.NoRestart, "evidenceFreshnessRequired": gates.EvidenceFreshnessRequired,
		},
	})
	if in.Registry != nil {
		for _, e := range in.Registry.Entries {
			addEdge(gateID, nodeID(KindCertification, e.LaneID), RelGates)
		}
	}

	if in.Simulation != nil {
		for _, live := range in.Simulation.LiveInPlaceEligibility {
			rfID := nodeID(KindResourceFuture, live.TargetRef)
			stale := live.StaleEvidence
			state := StateObserved
			if stale {
				state = StateStale
			}
			addNode(Node{
				ID: rfID, Kind: KindResourceFuture, Ref: live.TargetRef,
				CorrelationID: lineageID(rootCorr, "future", live.TargetRef),
				State: state, Stale: stale, ObservedAt: in.Simulation.ObservedAt,
				Attributes: map[string]any{
					"eligible": live.Eligible, "lane": live.Lane,
					"eligibilityState": live.EligibilityState,
				},
			})
			linked := false
			for _, p := range in.Discovery.DiscoveredResourcePorts {
				if p.CellRef == live.TargetRef {
					if pid, ok := portByRef[p.Ref]; ok {
						addEdge(pid, rfID, RelForecasts)
						linked = true
						break
					}
				}
			}
			if !linked {
				if pid, ok := portByRef[live.TargetRef]; ok {
					addEdge(pid, rfID, RelForecasts)
				}
			}
		}
	}

	for _, a := range in.Approvals {
		aid := nodeID(KindActionApproval, a.ActionID)
		stale := false
		state := StateObserved
		if a.ApprovalState == actionapproval.StatePending {
			if exp, err := time.Parse(time.RFC3339, a.ExpiresAt); err == nil && in.Now.After(exp) {
				stale = true
				state = StateStale
			}
		}
		addNode(Node{
			ID: aid, Kind: KindActionApproval, Ref: a.ActionID,
			CorrelationID: lineageID(rootCorr, "approval", a.ActionID),
			ObservedAt: a.ExpiresAt,
			State: state, Stale: stale,
			Provenance: a.Provenance,
			Attributes: map[string]any{
				"approvalState": a.ApprovalState, "laneId": a.LaneID,
			},
		})
		rfID := nodeID(KindResourceFuture, a.ResourceFutureRef)
		if _, ok := nodeIndex[rfID]; ok {
			addEdge(aid, rfID, RelApproves)
		} else {
			// try target from ref path
			for _, n := range g.Nodes {
				if n.Kind == KindResourceFuture && strings.Contains(a.ResourceFutureRef, n.Ref) {
					addEdge(aid, n.ID, RelApproves)
					break
				}
			}
		}
	}

	attachGraphProvenance(&g, in, rootCorr)
	addProvenanceEdges(&g, addEdge)
	AnnotateHealth(&g)
	sum, _ := VerifyLineageIntegrity(g, in.Now)
	g.ProvenanceValidation = sum
	g.LineageIntegrity = sum.LineageIntegrity
	return g
}

func attrString(attrs map[string]any, key string) string {
	if attrs == nil {
		return ""
	}
	if v, ok := attrs[key].(string); ok {
		return v
	}
	return ""
}

func rootCorrelationID(cluster string, at time.Time) string {
	h := sha256.Sum256([]byte(cluster + "|" + at.Format(time.RFC3339Nano)))
	return "khr-corr-" + hex.EncodeToString(h[:6])
}

func lineageID(root, kind, ref string) string {
	h := sha256.Sum256([]byte(root + "|" + kind + "|" + ref))
	return "khr-lineage-" + hex.EncodeToString(h[:6])
}

func nodeID(kind, ref string) string {
	h := sha256.Sum256([]byte(kind + "|" + ref))
	return "khr-node-" + kind + "-" + hex.EncodeToString(h[:8])
}

func edgeID(from, to, rel string) string {
	h := sha256.Sum256([]byte(from + "|" + to + "|" + rel))
	return "khr-edge-" + hex.EncodeToString(h[:8])
}

// AnnotateHealth runs validation and updates graph health + node flags.
func AnnotateHealth(g *Graph) {
	if g == nil {
		return
	}
	issues, orphans, stale := Validate(*g, time.Now().UTC())
	orphanSet := map[string]struct{}{}
	for _, id := range orphans {
		orphanSet[id] = struct{}{}
	}
	staleSet := map[string]struct{}{}
	for _, id := range stale {
		staleSet[id] = struct{}{}
	}
	for i := range g.Nodes {
		if _, ok := orphanSet[g.Nodes[i].ID]; ok {
			g.Nodes[i].Orphan = true
			g.Nodes[i].State = StateOrphan
		}
		if _, ok := staleSet[g.Nodes[i].ID]; ok {
			g.Nodes[i].Stale = true
			if g.Nodes[i].State != StateOrphan {
				g.Nodes[i].State = StateStale
			}
		}
	}
	g.Health = GraphHealth{
		NodeCount:     len(g.Nodes),
		EdgeCount:     len(g.Edges),
		OrphanCount:   len(orphans),
		StaleCount:    len(stale),
		OrphanNodeIDs: orphans,
		StaleNodeIDs:  stale,
		Consistent:    len(issues) == 0,
		Issues:        issues,
	}
	corrSet := map[string]struct{}{g.CorrelationID: {}}
	for _, n := range g.Nodes {
		corrSet[n.CorrelationID] = struct{}{}
	}
	ids := make([]string, 0, len(corrSet))
	for id := range corrSet {
		ids = append(ids, id)
	}
	g.Lineage.CorrelationIDs = ids
}
