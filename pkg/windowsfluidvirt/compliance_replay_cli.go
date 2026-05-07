package windowsfluidvirt

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

type AttestationSignatureMode string

const (
	AttestationModeUnsignedDev    AttestationSignatureMode = "unsigned-dev"
	AttestationModeFutureSignable AttestationSignatureMode = "future-signable"
)

type ReplayPoolContext struct {
	IsPoolChild            bool `json:"isPoolChild"`
	RequestedAsMechanism   bool `json:"requestedAsMechanism"`
	TreatedAsProvisionOnly bool `json:"treatedAsProvisionOnly"`
}

type ReplayMutationFlags struct {
	RuntimeMutation bool `json:"runtimeMutation"`
	CPUApply        bool `json:"cpuApply"`
	RAMApply        bool `json:"ramApply"`
	ActuatorApply   bool `json:"actuatorApply"`
	ClusterCalls    bool `json:"clusterCalls"`
	QMPCalls        bool `json:"qmpCalls"`
}

type WindowsComplianceReplayOutput struct {
	ReplayID                    string              `json:"replayId"`
	InputRef                    string              `json:"inputRef"`
	EvaluationTime              string              `json:"evaluationTime"`
	CompliancePhase             string              `json:"compliancePhase"`
	VMRef                       string              `json:"vmRef"`
	Namespace                   string              `json:"namespace"`
	ShellRef                    string              `json:"shellRef"`
	EvidenceSummary             map[string]any      `json:"evidenceSummary"`
	Blockers                    []string            `json:"blockers"`
	RemediationActions          []string            `json:"remediationActions"`
	AutomatableActions          []string            `json:"automatableActions"`
	ManualActions               []string            `json:"manualActions"`
	Risk                        string              `json:"risk"`
	PoolContext                 ReplayPoolContext   `json:"poolContext"`
	PoolScalingMechanismBlocked bool                `json:"poolScalingMechanismBlocked"`
	HyperdensityReady           bool                `json:"hyperdensityReady"`
	EvidenceHash                string              `json:"evidenceHash"`
	ReplayHash                  string              `json:"replayHash"`
	AuditRefs                   []string            `json:"auditRefs"`
	MutationFlags               ReplayMutationFlags `json:"mutationFlags"`
}

type ReplayAttestor struct {
	ComponentName    string `json:"componentName"`
	ComponentVersion string `json:"componentVersion"`
}

type ReplayAttestationSignature struct {
	Mode  AttestationSignatureMode `json:"mode"`
	Value string                   `json:"value"`
}

type WindowsComplianceReplayAttestation struct {
	AttestationID       string                     `json:"attestationId"`
	ReplayID            string                     `json:"replayId"`
	SubjectType         string                     `json:"subjectType"`
	SubjectRef          string                     `json:"subjectRef"`
	PolicyVersion       string                     `json:"policyVersion"`
	EvaluatorVersion    string                     `json:"evaluatorVersion"`
	EvidenceHash        string                     `json:"evidenceHash"`
	ReplayHash          string                     `json:"replayHash"`
	DecisionSnapshot    map[string]any             `json:"decisionSnapshot"`
	BlockerSnapshot     []string                   `json:"blockerSnapshot"`
	RemediationSnapshot []string                   `json:"remediationSnapshot"`
	CreatedAt           string                     `json:"createdAt"`
	Attestor            ReplayAttestor             `json:"attestor"`
	Signature           ReplayAttestationSignature `json:"signature"`
}

type WindowsComplianceReplayCLIOutput struct {
	WindowsComplianceReplayOutput
	Attestation *WindowsComplianceReplayAttestation `json:"attestation,omitempty"`
	BundleIndex *WindowsComplianceReplayBundleIndex `json:"bundleIndex,omitempty"`
}

type WindowsComplianceReplayBundleRun struct {
	RunID                     string                   `json:"runId"`
	InputRef                  string                   `json:"inputRef"`
	OutputRef                 string                   `json:"outputRef"`
	AttestationRef            string                   `json:"attestationRef,omitempty"`
	EvidenceHash              string                   `json:"evidenceHash"`
	ReplayHash                string                   `json:"replayHash"`
	AttestationHash           string                   `json:"attestationHash,omitempty"`
	PreviousRunHash           string                   `json:"previousRunHash,omitempty"`
	RunHash                   string                   `json:"runHash"`
	EvaluationTime            string                   `json:"evaluationTime"`
	CompliancePhase           string                   `json:"compliancePhase"`
	HyperdensityReady         bool                     `json:"hyperdensityReady"`
	Blockers                  []string                 `json:"blockers"`
	RemediationActions        []string                 `json:"remediationActions"`
	AttestationMode           AttestationSignatureMode `json:"attestationMode,omitempty"`
	AttestationSignatureValue string                   `json:"attestationSignatureValue,omitempty"`
}

type WindowsComplianceReplayHashChain struct {
	ChainMode     string `json:"chainMode"`
	FirstRunHash  string `json:"firstRunHash,omitempty"`
	LatestRunHash string `json:"latestRunHash,omitempty"`
	ChainValid    bool   `json:"chainValid"`
	BrokenAtRunID string `json:"brokenAtRunId,omitempty"`
	Notes         string `json:"notes,omitempty"`
}

type WindowsComplianceReplayBundleIndex struct {
	BundleID                string                             `json:"bundleId"`
	BundleVersion           string                             `json:"bundleVersion"`
	CreatedAt               string                             `json:"createdAt"`
	SubjectRef              string                             `json:"subjectRef"`
	SubjectType             string                             `json:"subjectType"`
	RunCount                int                                `json:"runCount"`
	Runs                    []WindowsComplianceReplayBundleRun `json:"runs"`
	Chain                   WindowsComplianceReplayHashChain   `json:"chain"`
	AggregateStatus         string                             `json:"aggregateStatus"`
	LatestCompliancePhase   string                             `json:"latestCompliancePhase"`
	LatestHyperdensityReady bool                               `json:"latestHyperdensityReady"`
	AuditRefs               []string                           `json:"auditRefs"`
}

func EvaluateWindowsComplianceReplay(
	input EvaluateWindowsHyperdensityReadyComplianceInput,
	inputRef string,
	evaluationTime time.Time,
) (WindowsComplianceReplayOutput, error) {
	evaluationTime = evaluationTime.UTC()
	if evaluationTime.IsZero() {
		evaluationTime = time.Now().UTC()
	}
	compliance := EvaluateWindowsHyperdensityReadyCompliance(input)
	evidenceHash, err := computeDeterministicHash(input)
	if err != nil {
		return WindowsComplianceReplayOutput{}, err
	}
	replayID := buildReplayID(inputRef, input.Identity.Namespace, input.Identity.VMRef, evaluationTime)
	replay := WindowsComplianceReplayOutput{
		ReplayID:                    replayID,
		InputRef:                    inputRef,
		EvaluationTime:              evaluationTime.Format(time.RFC3339),
		CompliancePhase:             string(compliance.CompliancePhase),
		VMRef:                       input.Identity.VMRef,
		Namespace:                   input.Identity.Namespace,
		ShellRef:                    deriveShellRef(input.Identity.Namespace, input.Identity.VMRef),
		EvidenceSummary:             compliance.EvidenceSummary,
		Blockers:                    compliance.Blockers,
		RemediationActions:          compliance.RemediationActions,
		AutomatableActions:          compliance.AutomatableActions,
		ManualActions:               compliance.ManualActions,
		Risk:                        string(compliance.Risk),
		PoolContext:                 ReplayPoolContext(input.PoolContext),
		PoolScalingMechanismBlocked: contains(compliance.Blockers, BlockerPoolScalingAsMechanism),
		HyperdensityReady:           compliance.CompliancePhase == ComplianceHyperdensityReadyWindowsShell,
		EvidenceHash:                evidenceHash,
		AuditRefs: []string{
			"input://" + inputRef,
			"evidence-hash://" + evidenceHash,
			"replay-id://" + replayID,
		},
		MutationFlags: ReplayMutationFlags{
			RuntimeMutation: false,
			CPUApply:        false,
			RAMApply:        false,
			ActuatorApply:   false,
			ClusterCalls:    false,
			QMPCalls:        false,
		},
	}
	replayHash, err := computeReplayHash(replay)
	if err != nil {
		return WindowsComplianceReplayOutput{}, err
	}
	replay.ReplayHash = replayHash
	replay.AuditRefs = append(replay.AuditRefs, "replay-hash://"+replayHash)
	return replay, nil
}

func BuildWindowsComplianceReplayAttestation(
	replay WindowsComplianceReplayOutput,
	mode AttestationSignatureMode,
	evaluationTime time.Time,
) (WindowsComplianceReplayAttestation, error) {
	if mode != AttestationModeUnsignedDev && mode != AttestationModeFutureSignable {
		return WindowsComplianceReplayAttestation{}, fmt.Errorf("invalid attestation mode: %s", mode)
	}
	evaluationTime = evaluationTime.UTC()
	if evaluationTime.IsZero() {
		evaluationTime = time.Now().UTC()
	}
	if replay.ReplayHash == "" || replay.EvidenceHash == "" {
		return WindowsComplianceReplayAttestation{}, errors.New("replay hash and evidence hash are required")
	}
	attestationID := "windows-compliance-attestation-" + shortHash(replay.ReplayID+"|"+replay.ReplayHash+"|"+evaluationTime.Format(time.RFC3339))
	return WindowsComplianceReplayAttestation{
		AttestationID:    attestationID,
		ReplayID:         replay.ReplayID,
		SubjectType:      "windows-hyperdensity-ready-compliance-replay",
		SubjectRef:       replay.ShellRef,
		PolicyVersion:    "windows-hyperdensity-ready-compliance-v1",
		EvaluatorVersion: "windows-compliance-replay-cli-v1",
		EvidenceHash:     replay.EvidenceHash,
		ReplayHash:       replay.ReplayHash,
		DecisionSnapshot: map[string]any{
			"compliancePhase":             replay.CompliancePhase,
			"hyperdensityReady":           replay.HyperdensityReady,
			"poolScalingMechanismBlocked": replay.PoolScalingMechanismBlocked,
			"risk":                        replay.Risk,
		},
		BlockerSnapshot:     replay.Blockers,
		RemediationSnapshot: replay.RemediationActions,
		CreatedAt:           evaluationTime.Format(time.RFC3339),
		Attestor: ReplayAttestor{
			ComponentName:    "karl-fluid-compliance-replay",
			ComponentVersion: "v1",
		},
		Signature: ReplayAttestationSignature{
			Mode:  mode,
			Value: "",
		},
	}, nil
}

func BuildWindowsComplianceReplayBundleRun(
	replay WindowsComplianceReplayOutput,
	attestation *WindowsComplianceReplayAttestation,
	previousRunHash string,
	evaluationTime time.Time,
) (WindowsComplianceReplayBundleRun, error) {
	evaluationTime = evaluationTime.UTC()
	if evaluationTime.IsZero() {
		evaluationTime = time.Now().UTC()
	}
	if replay.EvidenceHash == "" || replay.ReplayHash == "" {
		return WindowsComplianceReplayBundleRun{}, errors.New("replay evidence/replay hash are required")
	}
	run := WindowsComplianceReplayBundleRun{
		RunID:              "windows-compliance-run-" + shortHash(replay.ReplayID+"|"+replay.ReplayHash+"|"+previousRunHash),
		InputRef:           replay.InputRef,
		OutputRef:          "stdout://replay",
		EvidenceHash:       replay.EvidenceHash,
		ReplayHash:         replay.ReplayHash,
		PreviousRunHash:    previousRunHash,
		EvaluationTime:     replay.EvaluationTime,
		CompliancePhase:    replay.CompliancePhase,
		HyperdensityReady:  replay.HyperdensityReady,
		Blockers:           replay.Blockers,
		RemediationActions: replay.RemediationActions,
	}
	if attestation != nil {
		hash, err := computeDeterministicHash(attestation)
		if err != nil {
			return WindowsComplianceReplayBundleRun{}, err
		}
		run.AttestationRef = "stdout://attestation"
		run.AttestationHash = hash
		run.AttestationMode = attestation.Signature.Mode
		run.AttestationSignatureValue = attestation.Signature.Value
	}
	runHash, err := computeBundleRunHash(run)
	if err != nil {
		return WindowsComplianceReplayBundleRun{}, err
	}
	run.RunHash = runHash
	return run, nil
}

func BuildWindowsComplianceReplayBundleIndex(
	subjectRef string,
	bundleVersion string,
	runs []WindowsComplianceReplayBundleRun,
	evaluationTime time.Time,
) (WindowsComplianceReplayBundleIndex, error) {
	evaluationTime = evaluationTime.UTC()
	if evaluationTime.IsZero() {
		evaluationTime = time.Now().UTC()
	}
	if bundleVersion == "" {
		return WindowsComplianceReplayBundleIndex{}, errors.New("bundleVersion is required")
	}
	if subjectRef == "" {
		return WindowsComplianceReplayBundleIndex{}, errors.New("subjectRef is required")
	}
	bundle := WindowsComplianceReplayBundleIndex{
		BundleID:      "windows-compliance-bundle-" + shortHash(subjectRef+"|"+bundleVersion+"|"+evaluationTime.Format(time.RFC3339)),
		BundleVersion: bundleVersion,
		CreatedAt:     evaluationTime.Format(time.RFC3339),
		SubjectRef:    subjectRef,
		SubjectType:   "windows-hyperdensity-ready-compliance-replay-bundle",
		RunCount:      len(runs),
		Runs:          runs,
		Chain: WindowsComplianceReplayHashChain{
			ChainMode:  "local-deterministic-hash-chain",
			ChainValid: true,
		},
		AuditRefs: []string{
			"subject://" + subjectRef,
			"bundle-version://" + bundleVersion,
		},
	}
	if len(runs) > 0 {
		bundle.Chain.FirstRunHash = runs[0].RunHash
		bundle.Chain.LatestRunHash = runs[len(runs)-1].RunHash
		last := runs[len(runs)-1]
		bundle.LatestCompliancePhase = last.CompliancePhase
		bundle.LatestHyperdensityReady = last.HyperdensityReady
		if last.HyperdensityReady {
			bundle.AggregateStatus = "latest-ready"
		} else {
			bundle.AggregateStatus = "latest-blocked-or-not-ready"
		}
	} else {
		bundle.AggregateStatus = "empty"
		bundle.Chain.Notes = "no runs in bundle"
	}
	validation := ValidateWindowsComplianceReplayBundleIndex(bundle)
	if !validation.Chain.ChainValid {
		return validation, errors.New("constructed bundle index failed validation")
	}
	return validation, nil
}

func ValidateWindowsComplianceReplayBundleIndex(index WindowsComplianceReplayBundleIndex) WindowsComplianceReplayBundleIndex {
	index.Chain.ChainMode = "local-deterministic-hash-chain"
	index.Chain.ChainValid = true
	index.Chain.BrokenAtRunID = ""
	index.Chain.Notes = ""

	if index.BundleVersion == "" {
		index.Chain.ChainValid = false
		index.Chain.Notes = "bundleVersion is required"
		return index
	}
	if index.RunCount != len(index.Runs) {
		index.Chain.ChainValid = false
		index.Chain.Notes = "runCount mismatch"
		return index
	}
	if len(index.Runs) == 0 {
		if index.Chain.FirstRunHash != "" || index.Chain.LatestRunHash != "" {
			index.Chain.ChainValid = false
			index.Chain.Notes = "empty runs must not define first/latest hashes"
		}
		return index
	}

	for i := range index.Runs {
		run := index.Runs[i]
		if run.EvidenceHash == "" || run.ReplayHash == "" {
			index.Chain.ChainValid = false
			index.Chain.BrokenAtRunID = run.RunID
			index.Chain.Notes = "missing evidenceHash or replayHash"
			return index
		}
		if i > 0 && run.PreviousRunHash != index.Runs[i-1].RunHash {
			index.Chain.ChainValid = false
			index.Chain.BrokenAtRunID = run.RunID
			index.Chain.Notes = "previousRunHash mismatch"
			return index
		}
		expectedRunHash, err := computeBundleRunHash(run)
		if err != nil {
			index.Chain.ChainValid = false
			index.Chain.BrokenAtRunID = run.RunID
			index.Chain.Notes = "run hash computation failed"
			return index
		}
		if run.RunHash != expectedRunHash {
			index.Chain.ChainValid = false
			index.Chain.BrokenAtRunID = run.RunID
			index.Chain.Notes = "runHash mismatch"
			return index
		}
		if run.AttestationSignatureValue != "" {
			index.Chain.ChainValid = false
			index.Chain.BrokenAtRunID = run.RunID
			index.Chain.Notes = "attestation signature value must be empty"
			return index
		}
		if run.AttestationMode != "" &&
			run.AttestationMode != AttestationModeUnsignedDev &&
			run.AttestationMode != AttestationModeFutureSignable {
			index.Chain.ChainValid = false
			index.Chain.BrokenAtRunID = run.RunID
			index.Chain.Notes = "attestation mode invalid"
			return index
		}
	}

	index.Chain.FirstRunHash = index.Runs[0].RunHash
	index.Chain.LatestRunHash = index.Runs[len(index.Runs)-1].RunHash
	if index.Chain.LatestRunHash != index.Runs[len(index.Runs)-1].RunHash {
		index.Chain.ChainValid = false
		index.Chain.BrokenAtRunID = index.Runs[len(index.Runs)-1].RunID
		index.Chain.Notes = "latestRunHash mismatch"
	}
	last := index.Runs[len(index.Runs)-1]
	index.LatestCompliancePhase = last.CompliancePhase
	index.LatestHyperdensityReady = last.HyperdensityReady
	if last.HyperdensityReady {
		index.AggregateStatus = "latest-ready"
	} else {
		index.AggregateStatus = "latest-blocked-or-not-ready"
	}
	return index
}

func ParseAttestationMode(raw string) (AttestationSignatureMode, error) {
	mode := AttestationSignatureMode(raw)
	if mode != AttestationModeUnsignedDev && mode != AttestationModeFutureSignable {
		return "", fmt.Errorf("unsupported attestation mode: %s", raw)
	}
	return mode, nil
}

func LoadComplianceReplayInput(path string) (EvaluateWindowsHyperdensityReadyComplianceInput, error) {
	fixture, err := LoadWindowsHyperdensityComplianceReplayFixture(path)
	if err == nil && fixture.Input.Identity.VMRef != "" {
		return fixture.Input, nil
	}
	data, readErr := os.ReadFile(path)
	if readErr != nil {
		return EvaluateWindowsHyperdensityReadyComplianceInput{}, readErr
	}
	var raw EvaluateWindowsHyperdensityReadyComplianceInput
	if unmarshalErr := json.Unmarshal(data, &raw); unmarshalErr != nil {
		return EvaluateWindowsHyperdensityReadyComplianceInput{}, err
	}
	if raw.Identity.VMRef == "" {
		return EvaluateWindowsHyperdensityReadyComplianceInput{}, errors.New("compliance replay input is missing identity.vmRef")
	}
	return raw, nil
}

func computeDeterministicHash(value any) (string, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), nil
}

func computeReplayHash(replay WindowsComplianceReplayOutput) (string, error) {
	replay.ReplayHash = ""
	replay.AuditRefs = nil
	return computeDeterministicHash(replay)
}

func buildReplayID(inputRef, namespace, vmRef string, evaluationTime time.Time) string {
	body := fmt.Sprintf("%s|%s|%s|%s", inputRef, namespace, vmRef, evaluationTime.Format(time.RFC3339))
	return "windows-compliance-replay-" + shortHash(body)
}

func shortHash(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:8])
}

func deriveShellRef(namespace, vmRef string) string {
	if namespace == "" {
		return "windows-shell/" + vmRef
	}
	return "windows-shell/" + namespace + "/" + vmRef
}

func computeBundleRunHash(run WindowsComplianceReplayBundleRun) (string, error) {
	run.RunHash = ""
	return computeDeterministicHash(run)
}
