package sidecar

import (
	"errors"
	"testing"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/windowsfluidvirt"
)

func TestReadOnlyExecutorRejectsMutatingCommand(t *testing.T) {
	executor := NewReadOnlyExecutor(&StaticTransport{
		Connected: true,
		Responses: map[string]map[string]any{},
		Errors:    map[string]error{},
	})
	if _, err := executor.Execute("device_add"); err == nil {
		t.Fatal("expected mutating command to be rejected")
	}
}

func TestReadOnlyExecutorAllowsReadOnlyCommand(t *testing.T) {
	executor := NewReadOnlyExecutor(&StaticTransport{
		Connected: true,
		Responses: map[string]map[string]any{
			"query-status": {"status": "running"},
		},
		Errors: map[string]error{},
	})
	result, err := executor.Execute("query-status")
	if err != nil {
		t.Fatalf("expected read-only command to succeed: %v", err)
	}
	if result["status"] != "running" {
		t.Fatalf("unexpected query result: %+v", result)
	}
}

func TestDiscoverEvidenceHandshakeAndQueries(t *testing.T) {
	transport := &StaticTransport{
		Connected: true,
		Responses: map[string]map[string]any{
			"query-cpus-fast":           {"maxCpus": int64(8)},
			"query-hotpluggable-cpus":   {"maxCpus": int64(8)},
			"query-memory-devices":      {"devices": 2},
			"query-memory-size-summary": {"base": 1024},
		},
		Errors: nil,
	}
	executor := NewReadOnlyExecutor(transport)
	evidence := executor.DiscoverEvidence("/tmp/qmp.sock", "9011", "v-test")
	if !evidence.QMPConnected || !evidence.QMPCapabilitiesNegotiated || !evidence.CPUTopologyObserved || !evidence.MemoryDevicesObserved {
		t.Fatalf("unexpected evidence: %+v", evidence)
	}
	if !evidence.QMPReadOnly {
		t.Fatal("expected qmpReadOnly=true")
	}
}

func TestDiscoverEvidenceMissingSocketMapsToBlocker(t *testing.T) {
	executor := NewReadOnlyExecutor(&StaticTransport{Connected: false})
	evidence := executor.DiscoverEvidence("/missing.sock", "", "v-test")
	blockers := windowsfluidvirt.ValidateQmpReadiness(evidence)
	hasSocketBlocker := false
	for _, blocker := range blockers {
		if blocker == windowsfluidvirt.BlockerQMPSocketUnavailable {
			hasSocketBlocker = true
			break
		}
	}
	if !hasSocketBlocker {
		t.Fatalf("expected qmp_socket_unavailable blocker, got %v", blockers)
	}
}

func TestDiscoverEvidenceQmpErrorMapsToAckMissing(t *testing.T) {
	transport := &StaticTransport{
		Connected: true,
		Responses: map[string]map[string]any{
			"qmp_capabilities": {"ok": true},
		},
		Errors: map[string]error{
			"query-status": errors.New("broken qmp command"),
		},
	}
	executor := NewReadOnlyExecutor(transport)
	evidence := executor.DiscoverEvidence("/tmp/qmp.sock", "9011", "v-test")
	blockers := windowsfluidvirt.ValidateQmpReadiness(evidence)
	hasAckBlocker := false
	for _, blocker := range blockers {
		if blocker == windowsfluidvirt.BlockerQMPAckMissing {
			hasAckBlocker = true
			break
		}
	}
	if !hasAckBlocker {
		t.Fatalf("expected qmp_ack_missing blocker, got %v", blockers)
	}
}
