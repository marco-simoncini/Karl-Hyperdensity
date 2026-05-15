package windowsfluidvirt

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
)

type WindowsFluidVirtAuditHashChain struct {
	ChainID                string `json:"chainId"`
	ChainVersion           string `json:"chainVersion"`
	HashAlgorithm          string `json:"hashAlgorithm"`
	Deterministic          bool   `json:"deterministic"`
	EventCount             int    `json:"eventCount"`
	RootHash               string `json:"rootHash"`
	TerminalHash           string `json:"terminalHash"`
	Verified               bool   `json:"verified"`
	ContainsSecretMaterial bool   `json:"containsSecretMaterial"`
	Mutable                bool   `json:"mutable"`
	ReplayOnly             bool   `json:"replayOnly"`
	ClaimBoundary          string `json:"claimBoundary"`
}

func BuildWindowsFluidVirtAuditHashChain(events []WindowsFluidVirtComplianceReplayedEvent) (WindowsFluidVirtAuditHashChain, []WindowsFluidVirtComplianceReplayedEvent, error) {
	if len(events) == 0 {
		return WindowsFluidVirtAuditHashChain{}, nil, fmt.Errorf("at least one replayed event is required")
	}
	replayed := append([]WindowsFluidVirtComplianceReplayedEvent(nil), events...)
	sort.Slice(replayed, func(i, j int) bool {
		return replayed[i].DeterministicOrder < replayed[j].DeterministicOrder
	})
	previousHash := ""
	for i := range replayed {
		replayed[i].PreviousHash = previousHash
		replayed[i].EventHash = hashReplayEvent(replayed[i], previousHash)
		previousHash = replayed[i].EventHash
	}
	chain := WindowsFluidVirtAuditHashChain{
		ChainID:                "windows_fluidvirt_compliance_audit_hash_chain_v1",
		ChainVersion:           "v1",
		HashAlgorithm:          "sha256",
		Deterministic:          true,
		EventCount:             len(replayed),
		RootHash:               replayed[0].EventHash,
		TerminalHash:           replayed[len(replayed)-1].EventHash,
		Verified:               true,
		ContainsSecretMaterial: false,
		Mutable:                false,
		ReplayOnly:             true,
		ClaimBoundary:          "audit_hash_chain_replay_only",
	}
	return chain, replayed, nil
}

func VerifyWindowsFluidVirtAuditHashChain(chain WindowsFluidVirtAuditHashChain, events []WindowsFluidVirtComplianceReplayedEvent) bool {
	if chain.HashAlgorithm != "sha256" || !chain.Deterministic || chain.Mutable || !chain.ReplayOnly {
		return false
	}
	rebuilt, rebuiltEvents, err := BuildWindowsFluidVirtAuditHashChain(events)
	if err != nil {
		return false
	}
	if chain.EventCount != rebuilt.EventCount || chain.RootHash != rebuilt.RootHash || chain.TerminalHash != rebuilt.TerminalHash {
		return false
	}
	for idx, event := range rebuiltEvents {
		if event.EventHash != events[idx].EventHash || event.PreviousHash != events[idx].PreviousHash {
			return false
		}
		if event.ContainsSecretMaterial || event.TouchesRuntime {
			return false
		}
	}
	return true
}

func hashReplayEvent(event WindowsFluidVirtComplianceReplayedEvent, previousHash string) string {
	raw := fmt.Sprintf(
		"%s|%s|%s|%d|%s|%t|%t|%s",
		event.EventID,
		event.EventType,
		event.EventState,
		event.DeterministicOrder,
		previousHash,
		event.ContainsSecretMaterial,
		event.TouchesRuntime,
		event.ClaimBoundary,
	)
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}
