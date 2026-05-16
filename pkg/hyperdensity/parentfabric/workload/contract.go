package workload

import "encoding/json"

// KindSample is one kind → output row in the pure-candidates contract.
type KindSample struct {
	Kind     string `json:"kind"`
	Resource string `json:"resource,omitempty"`
	Term     string `json:"term,omitempty"`
	OK       bool   `json:"ok,omitempty"`
	Supported bool  `json:"supported,omitempty"`
}

// ContractDocument is the deterministic Sprint 52 workload pure-candidates snapshot.
type ContractDocument struct {
	ContractVersion               string       `json:"contractVersion"`
	WorkloadHelpersOverallVerdict string       `json:"workloadHelpersOverallVerdict"`
	RuntimeImportAllowed          bool         `json:"runtimeImportAllowed"`
	AppsWorkloadResourceSamples   []KindSample `json:"appsWorkloadResourceSamples"`
	PilotWorkloadTermSamples      []KindSample `json:"pilotWorkloadTermSamples"`
	ExecutionSupportsLiveApplyKindSamples []KindSample `json:"executionSupportsLiveApplyKindSamples"`
}

// DefaultContractDocument returns canonical samples aligned with Dashboard source.
func DefaultContractDocument() ContractDocument {
	kinds := []string{
		"Deployment",
		"StatefulSet",
		"DaemonSet",
		"Pod",
		"VirtualMachine",
		"VirtualMachineInstance",
		"ReplicaSet",
		"UnknownKind",
		"",
		" Deployment ",
	}
	doc := ContractDocument{
		ContractVersion:               WorkloadPackageVersion,
		WorkloadHelpersOverallVerdict: "copy-deferred",
		RuntimeImportAllowed:          false,
	}
	for _, kind := range kinds {
		res, ok := AppsWorkloadResource(kind)
		doc.AppsWorkloadResourceSamples = append(doc.AppsWorkloadResourceSamples, KindSample{
			Kind: kind, Resource: res, OK: ok,
		})
		term, tok := PilotWorkloadTerm(kind)
		doc.PilotWorkloadTermSamples = append(doc.PilotWorkloadTermSamples, KindSample{
			Kind: kind, Term: term, OK: tok,
		})
		doc.ExecutionSupportsLiveApplyKindSamples = append(doc.ExecutionSupportsLiveApplyKindSamples, KindSample{
			Kind: kind, Supported: ExecutionSupportsLiveApplyKind(kind),
		})
	}
	return doc
}

// CanonicalContractJSON returns marshaled contract with computed outputs.
func CanonicalContractJSON() ([]byte, error) {
	return json.Marshal(DefaultContractDocument())
}
