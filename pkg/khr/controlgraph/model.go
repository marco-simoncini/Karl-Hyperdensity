// Package controlgraph provides the unified KHR control-plane graph (KHR-X).
package controlgraph

const (
	GraphID = "khr-control-graph-v1"

	KindHost           = "Host"
	KindShell          = "Shell"
	KindCell           = "Cell"
	KindResourcePort   = "ResourcePort"
	KindResourceLease  = "ResourceLease"
	KindResourceFuture = "ResourceFuture"
	KindCertification  = "Certification"
	KindPolicyGate     = "PolicyGate"
	KindActionApproval = "ActionApproval"

	RelHosts      = "hosts"
	RelProjects   = "projects"
	RelBinds      = "binds"
	RelForecasts  = "forecasts"
	RelCertifies  = "certifies"
	RelGates      = "gates"
	RelApproves   = "approves"
	RelCorrelates = "correlates"

	StateObserved = "observed"
	StateStale    = "stale"
	StateOrphan   = "orphan"

	DefaultStaleAfterSeconds int64 = 86400
)

// Node is one entity in the unified control graph.
type Node struct {
	ID              string         `json:"id"`
	Kind            string         `json:"kind"`
	Ref             string         `json:"ref"`
	CorrelationID   string         `json:"correlationId"`
	LineageParentID string         `json:"lineageParentId,omitempty"`
	State           string         `json:"state"`
	Stale           bool           `json:"stale"`
	Orphan          bool           `json:"orphan"`
	ObservedAt      string         `json:"observedAt,omitempty"`
	Attributes      map[string]any `json:"attributes,omitempty"`
}

// Edge is a directed relationship between graph nodes.
type Edge struct {
	ID            string `json:"id"`
	From          string `json:"from"`
	To            string `json:"to"`
	Relationship  string `json:"relationship"`
	CorrelationID string `json:"correlationId"`
}

// LineageSummary aggregates correlation lineage for export.
type LineageSummary struct {
	RootCorrelationID string   `json:"rootCorrelationId"`
	CorrelationIDs    []string `json:"correlationIds"`
}

// GraphHealth summarizes consistency signals.
type GraphHealth struct {
	NodeCount      int      `json:"nodeCount"`
	EdgeCount      int      `json:"edgeCount"`
	OrphanCount    int      `json:"orphanCount"`
	StaleCount     int      `json:"staleCount"`
	OrphanNodeIDs  []string `json:"orphanNodeIds,omitempty"`
	StaleNodeIDs   []string `json:"staleNodeIds,omitempty"`
	Consistent     bool     `json:"consistent"`
	Issues         []string `json:"issues,omitempty"`
}

// Graph is a JSON snapshot of the unified KHR control plane.
type Graph struct {
	GraphID                   string         `json:"graphId"`
	Sprint                    string         `json:"sprint"`
	ClusterContext            string         `json:"clusterContext"`
	ObservedAt                string         `json:"observedAt"`
	CorrelationID             string         `json:"correlationId"`
	ReadOnly                  bool           `json:"readOnly"`
	NoAutonomousOrchestration bool           `json:"noAutonomousOrchestration"`
	NoMutation                bool           `json:"noMutation"`
	NoApply                   bool           `json:"noApply"`
	Nodes                     []Node         `json:"nodes"`
	Edges                     []Edge         `json:"edges"`
	Lineage                   LineageSummary `json:"lineage"`
	Health                    GraphHealth    `json:"health"`
}
