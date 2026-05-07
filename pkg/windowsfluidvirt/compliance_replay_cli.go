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
