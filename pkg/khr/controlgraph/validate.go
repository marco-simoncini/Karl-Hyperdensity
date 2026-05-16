package controlgraph

import (
	"fmt"
	"strings"
	"time"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/actionapproval"
)

// Validate checks graph consistency and returns issues, orphan IDs, stale IDs.
func Validate(g Graph, now time.Time) (issues []string, orphanIDs, staleIDs []string) {
	nodeByID := map[string]Node{}
	incoming := map[string]int{}
	outgoing := map[string]int{}
	for _, n := range g.Nodes {
		nodeByID[n.ID] = n
	}
	for _, e := range g.Edges {
		if _, ok := nodeByID[e.From]; !ok {
			issues = append(issues, fmt.Sprintf("edge %s references missing from node %s", e.ID, e.From))
		}
		if _, ok := nodeByID[e.To]; !ok {
			issues = append(issues, fmt.Sprintf("edge %s references missing to node %s", e.ID, e.To))
		}
		outgoing[e.From]++
		incoming[e.To]++
	}

	// Correlation integrity: all nodes must share root prefix lineage or root id family.
	root := g.CorrelationID
	for _, n := range g.Nodes {
		if n.CorrelationID == "" {
			issues = append(issues, "node "+n.ID+" missing correlationId")
		}
		if root != "" && n.CorrelationID != root && !strings.HasPrefix(n.CorrelationID, "khr-lineage-") && !strings.HasPrefix(n.CorrelationID, "khr-corr-") {
			issues = append(issues, "node "+n.ID+" correlationId outside lineage family")
		}
	}

	// Orphan detection
	for _, n := range g.Nodes {
		switch n.Kind {
		case KindResourcePort:
			if incoming[n.ID] == 0 {
				orphanIDs = append(orphanIDs, n.ID)
				issues = append(issues, "orphan ResourcePort "+n.Ref)
			}
		case KindCell:
			if incoming[n.ID] == 0 && n.LineageParentID == "" {
				orphanIDs = append(orphanIDs, n.ID)
				issues = append(issues, "orphan Cell "+n.Ref)
			}
		case KindActionApproval:
			if outgoing[n.ID] == 0 {
				orphanIDs = append(orphanIDs, n.ID)
				issues = append(issues, "orphan ActionApproval "+n.Ref)
			}
		case KindResourceFuture:
			if incoming[n.ID] == 0 {
				orphanIDs = append(orphanIDs, n.ID)
				issues = append(issues, "orphan ResourceFuture "+n.Ref)
			}
		}
	}

	// Stale detection
	for _, n := range g.Nodes {
		if n.Stale {
			staleIDs = append(staleIDs, n.ID)
			continue
		}
		if n.Kind == KindActionApproval {
			if state, _ := n.Attributes["approvalState"].(string); state == actionapproval.StatePending {
				if exp, err := time.Parse(time.RFC3339, n.ObservedAt); err == nil && now.After(exp) {
					staleIDs = append(staleIDs, n.ID)
				}
			}
		}
		if n.Kind == KindCertification && n.State == StateStale {
			staleIDs = append(staleIDs, n.ID)
		}
	}

	// Required kinds for native-live sandbox graph
	hasHost, hasCell := false, false
	for _, n := range g.Nodes {
		if n.Kind == KindHost {
			hasHost = true
		}
		if n.Kind == KindCell {
			hasCell = true
		}
	}
	if len(g.Nodes) > 0 && (!hasHost || !hasCell) {
		issues = append(issues, "graph missing Host or Cell anchor nodes")
	}

	return issues, orphanIDs, staleIDs
}

// CheckCorrelationIntegrity verifies lineage IDs are present on connected chains.
func CheckCorrelationIntegrity(g Graph) error {
	issues, _, _ := Validate(g, time.Now().UTC())
	for _, iss := range issues {
		if strings.Contains(iss, "correlationId") {
			return fmt.Errorf("%s", iss)
		}
	}
	return nil
}
