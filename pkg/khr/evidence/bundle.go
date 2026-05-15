// Package evidence builds local JSON evidence bundles for KHR Linux (Sprint 9).
package evidence

import (
	"os"
	"strings"
	"time"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/audit"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/cgroup"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/crdv1alpha1"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/discovery"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/resourcelease"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/telemetry"
)

// CollectEvidenceBundle is the unified JSON envelope for `collect-evidence` mode.
type CollectEvidenceBundle struct {
	Tool               string             `json:"tool"`
	Version            string             `json:"version"`
	Mode               string             `json:"mode"`
	AgentID            string             `json:"agentId"`
	CollectedAt        string             `json:"collectedAt"`
	CellRef            *telemetry.CellRef `json:"cellRef,omitempty"`
	Discovery          DiscoverySnapshot  `json:"discovery"`
	Telemetry          TelemetrySnapshot  `json:"telemetry"`
	DryRun             DryRunPayload      `json:"dryRun"`
	EvidenceSummary    EvidenceSummary    `json:"evidenceSummary"`
	MutationsForbidden bool               `json:"mutationsForbidden"`
}

// DiscoverySnapshot strips CLI wrapper fields from discovery output.
type DiscoverySnapshot struct {
	AgentID            string   `json:"agentId"`
	CgroupVersion      string   `json:"cgroupVersion"`
	DiscoveryMode      string   `json:"discoveryMode"`
	ScannedRoot        string   `json:"scannedRoot"`
	AllowedPathPrefix  string   `json:"allowedPathPrefix"`
	CandidatePaths     []string `json:"candidatePaths"`
	SelectedPath       string   `json:"selectedPath,omitempty"`
	BlockedReasons     []string `json:"blockedReasons"`
	Warnings           []string `json:"warnings"`
	MutationsForbidden bool     `json:"mutationsForbidden"`
}

// TelemetrySnapshot is telemetry evidence suitable for embedding in a bundle.
type TelemetrySnapshot struct {
	Skipped            bool                    `json:"skipped,omitempty"`
	SkipReason         string                  `json:"skipReason,omitempty"`
	TelemetryMode      string                  `json:"telemetryMode"`
	CgroupPath         string                  `json:"cgroupPath"`
	AllowedPathPrefix  string                  `json:"allowedPathPrefix"`
	CellRef            *telemetry.CellRef      `json:"cellRef,omitempty"`
	Metrics            telemetry.MetricsBundle `json:"metrics"`
	Evidence           telemetry.Evidence      `json:"evidence"`
	MutationsForbidden bool                    `json:"mutationsForbidden"`
}

// DryRunPayload is optional ResourceLease simulation embedded in the bundle.
type DryRunPayload struct {
	Skipped                 bool                        `json:"skipped,omitempty"`
	SkipReason              string                      `json:"skipReason,omitempty"`
	ResourceLeaseDryRun     *resourcelease.DryRunResult `json:"resourceLeaseDryRun,omitempty"`
	CgroupEnvelopePlan      *cgroup.EnvelopePlan        `json:"cgroupEnvelopePlan,omitempty"`
	MutationsForbidden      bool                        `json:"mutationsForbidden,omitempty"`
	UnsafeApplyFlagPresent  bool                        `json:"unsafeApplyFlagPresent,omitempty"`
	FutureApplyGateRequired bool                        `json:"futureApplyGateRequired,omitempty"`
	Audit                   []audit.Record              `json:"audit,omitempty"`
}

// EvidenceSummary aggregates cross-step readiness for downstream consumers.
type EvidenceSummary struct {
	Confidence            string   `json:"confidence"`
	ReadyForGrandePadre   bool     `json:"readyForGrandePadre"`
	BlockedReasons        []string `json:"blockedReasons"`
	Warnings              []string `json:"warnings"`
	RecommendedNextAction string   `json:"recommendedNextAction"`
}

// BuildCollectEvidenceBundle assembles the bundle and computes evidenceSummary.
func BuildCollectEvidenceBundle(version, agentID string, cell *crdv1alpha1.Cell, disc *discovery.CgroupDiscoveryOutput, tel TelemetrySnapshot, dry DryRunPayload, leasePortPartialWarning string) *CollectEvidenceBundle {
	b := &CollectEvidenceBundle{
		Tool:               "khr-linux-agent",
		Version:            version,
		Mode:               "collect-evidence",
		AgentID:            agentID,
		CollectedAt:        collectedAtNow().Format(time.RFC3339),
		Discovery:          DiscoverySnapshotFrom(disc),
		Telemetry:          tel,
		DryRun:             dry,
		MutationsForbidden: true,
	}
	if cell != nil {
		b.CellRef = &telemetry.CellRef{
			APIVersion: cell.APIVersion,
			Kind:       cell.Kind,
			Namespace:  cell.Metadata.Namespace,
			Name:       cell.Metadata.Name,
		}
	}
	b.EvidenceSummary = Summarize(b, leasePortPartialWarning)
	return b
}

// DiscoverySnapshotFrom copies discovery output without nested CLI fields.
func DiscoverySnapshotFrom(d *discovery.CgroupDiscoveryOutput) DiscoverySnapshot {
	if d == nil {
		return DiscoverySnapshot{BlockedReasons: []string{"discovery output is nil"}}
	}
	return DiscoverySnapshot{
		AgentID:            d.AgentID,
		CgroupVersion:      d.CgroupVersion,
		DiscoveryMode:      d.DiscoveryMode,
		ScannedRoot:        d.ScannedRoot,
		AllowedPathPrefix:  d.AllowedPathPrefix,
		CandidatePaths:     append([]string{}, d.CandidatePaths...),
		SelectedPath:       d.SelectedPath,
		BlockedReasons:     append([]string{}, d.BlockedReasons...),
		Warnings:           append([]string{}, d.Warnings...),
		MutationsForbidden: d.MutationsForbidden,
	}
}

// TelemetrySnapshotFrom strips CLI wrapper fields from telemetry output.
func TelemetrySnapshotFrom(t *telemetry.ReadTelemetryOutput) TelemetrySnapshot {
	if t == nil {
		return TelemetrySnapshot{
			Skipped:            true,
			SkipReason:         "telemetry output is nil",
			TelemetryMode:      "read-only",
			Evidence:           telemetry.BuildEvidence(nil, []string{"telemetry missing"}, telemetry.MetricsBundle{}),
			MutationsForbidden: true,
		}
	}
	return TelemetrySnapshot{
		Skipped:            false,
		TelemetryMode:      t.TelemetryMode,
		CgroupPath:         t.CgroupPath,
		AllowedPathPrefix:  t.AllowedPathPrefix,
		CellRef:            t.CellRef,
		Metrics:            t.Metrics,
		Evidence:           t.Evidence,
		MutationsForbidden: t.MutationsForbidden,
	}
}

// TelemetrySnapshotSkipped builds a read-only telemetry section when discovery did not resolve a path.
func TelemetrySnapshotSkipped(allowPrefix string, cellRef *telemetry.CellRef, reason string) TelemetrySnapshot {
	m := telemetry.MetricsBundle{}
	br := []string{reason}
	ev := telemetry.BuildEvidence(nil, br, m)
	return TelemetrySnapshot{
		Skipped:            true,
		SkipReason:         reason,
		TelemetryMode:      "read-only",
		CgroupPath:         "",
		AllowedPathPrefix:  strings.TrimSpace(allowPrefix),
		CellRef:            cellRef,
		Metrics:            m,
		Evidence:           ev,
		MutationsForbidden: true,
	}
}

// DryRunPayloadFromResult copies fields from the agent dry-run result shape (passed as components from agent).
func DryRunPayloadFromResult(lease resourcelease.DryRunResult, plan cgroup.EnvelopePlan, mut, unsafe, gate bool, aud []audit.Record) DryRunPayload {
	l := lease
	p := plan
	return DryRunPayload{
		Skipped:                 false,
		ResourceLeaseDryRun:     &l,
		CgroupEnvelopePlan:      &p,
		MutationsForbidden:      mut,
		UnsafeApplyFlagPresent:  unsafe,
		FutureApplyGateRequired: gate,
		Audit:                   append([]audit.Record{}, aud...),
	}
}

// DryRunSkippedPayload marks dry-run as not executed.
func DryRunSkippedPayload(reason string) DryRunPayload {
	return DryRunPayload{
		Skipped:    true,
		SkipReason: reason,
	}
}

// Summarize aggregates warnings and blocked reasons and sets readiness fields.
func Summarize(b *CollectEvidenceBundle, leasePortPartialWarning string) EvidenceSummary {
	var warns []string
	var blocked []string

	warns = append(warns, b.Discovery.Warnings...)
	warns = append(warns, b.Telemetry.Evidence.Warnings...)
	if leasePortPartialWarning != "" {
		warns = append(warns, leasePortPartialWarning)
	}

	blocked = append(blocked, b.Discovery.BlockedReasons...)
	blocked = append(blocked, b.Telemetry.Evidence.BlockedReasons...)

	if !b.DryRun.Skipped && b.DryRun.ResourceLeaseDryRun != nil {
		if b.DryRun.ResourceLeaseDryRun.Blocked {
			blocked = append(blocked, "resource lease dry-run blocked: "+strings.TrimSpace(b.DryRun.ResourceLeaseDryRun.Reason))
		}
	}

	conf := aggregateConfidence(b)
	dryOK := b.DryRun.Skipped || (b.DryRun.ResourceLeaseDryRun != nil && b.DryRun.ResourceLeaseDryRun.Allowed && !b.DryRun.ResourceLeaseDryRun.Blocked)
	ready := len(dedupe(blocked)) == 0 &&
		strings.TrimSpace(b.Discovery.SelectedPath) != "" &&
		!b.Telemetry.Skipped &&
		len(b.Telemetry.Evidence.BlockedReasons) == 0 &&
		dryOK

	action := recommendAction(b, ready, leasePortPartialWarning)

	return EvidenceSummary{
		Confidence:            conf,
		ReadyForGrandePadre:   ready,
		BlockedReasons:        dedupe(blocked),
		Warnings:              dedupe(warns),
		RecommendedNextAction: action,
	}
}

func aggregateConfidence(b *CollectEvidenceBundle) string {
	conf := "high"
	if b.Discovery.SelectedPath == "" || len(b.Discovery.BlockedReasons) > 0 {
		conf = minConf(conf, "low")
	}
	if b.Telemetry.Skipped {
		conf = minConf(conf, "low")
	} else {
		conf = minConf(conf, b.Telemetry.Evidence.Confidence)
	}
	if !b.DryRun.Skipped && b.DryRun.ResourceLeaseDryRun != nil {
		if b.DryRun.ResourceLeaseDryRun.Blocked {
			conf = minConf(conf, "low")
		} else if !b.DryRun.ResourceLeaseDryRun.Allowed {
			conf = minConf(conf, "medium")
		}
	}
	return conf
}

func minConf(a, b string) string {
	rank := map[string]int{"low": 0, "medium": 1, "high": 2}
	ra, okA := rank[strings.ToLower(strings.TrimSpace(a))]
	rb, okB := rank[strings.ToLower(strings.TrimSpace(b))]
	if !okA {
		ra = 2
	}
	if !okB {
		rb = 2
	}
	if ra < rb {
		return a
	}
	return b
}

func recommendAction(b *CollectEvidenceBundle, ready bool, leasePortPartial string) string {
	if leasePortPartial != "" {
		return "Provide both -lease-input and -resource-port-input to include ResourceLease dry-run simulation in the bundle."
	}
	if ready {
		return "Bundle is consistent for downstream Grande Padre ingest; attach cluster admission and blast-radius context before any apply planning."
	}
	if strings.TrimSpace(b.Discovery.SelectedPath) == "" {
		return "Fix Cell.providerHandle / slice layout or expand discovery scan root; re-run collect-evidence until discovery.selectedPath is non-empty."
	}
	if b.Telemetry.Skipped || len(b.Telemetry.Evidence.BlockedReasons) > 0 {
		return "Verify cgroup v2 metric files are readable under the selected path and that allow-path-prefix policy matches the host layout."
	}
	if !b.DryRun.Skipped && b.DryRun.ResourceLeaseDryRun != nil && b.DryRun.ResourceLeaseDryRun.Blocked {
		return "Resolve ResourceLease / ResourcePort contract conflicts shown in dryRun.resourceLeaseDryRun before scheduling apply."
	}
	return "Review evidenceSummary.blockedReasons and warnings; remediate host or contract issues then re-run collect-evidence."
}

func dedupe(in []string) []string {
	seen := make(map[string]struct{})
	var out []string
	for _, s := range in {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	if len(out) == 0 {
		return []string{}
	}
	return out
}

func collectedAtNow() time.Time {
	s := os.Getenv("KHR_TEST_COLLECTED_AT")
	if s == "" {
		return time.Now().UTC()
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		t2, err2 := time.Parse(time.RFC3339Nano, s)
		if err2 != nil {
			return time.Now().UTC()
		}
		return t2.UTC()
	}
	return t.UTC()
}
