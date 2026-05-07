package windowsfluidvirt

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type NodeFluidActuatorMVPMode string

const (
	NodeActuatorModeDryRun        NodeFluidActuatorMVPMode = "dry-run"
	NodeActuatorModeApply         NodeFluidActuatorMVPMode = "apply"
	NodeActuatorModeRollback      NodeFluidActuatorMVPMode = "rollback"
	NodeActuatorModeReturnToFloor NodeFluidActuatorMVPMode = "return-to-floor"
)

type NodeFluidActuatorMVPRequest struct {
	RequestID       string   `json:"requestId"`
	Namespace       string   `json:"namespace"`
	VMName          string   `json:"vmName"`
	VMUID           string   `json:"vmUid"`
	VMIUID          string   `json:"vmiUid"`
	PodName         string   `json:"podName"`
	PodUID          string   `json:"podUid"`
	NodeName        string   `json:"nodeName"`
	QemuPID         string   `json:"qemuPid"`
	QemuStartTime   string   `json:"qemuStartTime"`
	CgroupPath      string   `json:"cgroupPath"`
	PreviousCPUMax  string   `json:"previousCpuMax"`
	RequestedCPUMax string   `json:"requestedCpuMax"`
	RollbackCPUMax  string   `json:"rollbackCpuMax"`
	TTLSeconds      int64    `json:"ttlSeconds"`
	Reason          string   `json:"reason"`
	EvidenceRefs    []string `json:"evidenceRefs"`
	PolicyVersion   string   `json:"policyVersion"`
	CreatedAt       string   `json:"createdAt"`
}

type NodeFluidActuatorMVPAllowlist struct {
	NodeName               string   `json:"nodeName"`
	Namespace              string   `json:"namespace"`
	VMName                 string   `json:"vmName"`
	PodUID                 string   `json:"podUid"`
	QemuPID                string   `json:"qemuPid"`
	QemuStartTime          string   `json:"qemuStartTime"`
	CgroupPath             string   `json:"cgroupPath"`
	AllowedControllers     []string `json:"allowedControllers"`
	MinCPUMax              string   `json:"minCpuMax"`
	MaxCPUMax              string   `json:"maxCpuMax"`
	AllowParentCgroupWrite bool     `json:"allowParentCgroupWrite"`
	AllowArbitraryWrite    bool     `json:"allowArbitraryWrite"`
}

type NodeFluidActuatorMVPEvidence struct {
	ActuatorID          string                   `json:"actuatorId"`
	ActuatorVersion     string                   `json:"actuatorVersion"`
	Mode                NodeFluidActuatorMVPMode `json:"mode"`
	DryRun              bool                     `json:"dryRun"`
	RequestID           string                   `json:"requestId"`
	Allowed             bool                     `json:"allowed"`
	Blockers            []string                 `json:"blockers"`
	PolicyDecision      string                   `json:"policyDecision"`
	TargetFile          string                   `json:"targetFile"`
	BeforeCPUMax        string                   `json:"beforeCpuMax,omitempty"`
	AppliedCPUMax       string                   `json:"appliedCpuMax,omitempty"`
	AfterCPUMax         string                   `json:"afterCpuMax,omitempty"`
	ReturnToFloorCPUMax string                   `json:"returnToFloorCpuMax,omitempty"`
	CreatedAt           string                   `json:"createdAt"`
	Notes               []string                 `json:"notes,omitempty"`
}

type NodeFluidActuatorMVPInput struct {
	Mode           NodeFluidActuatorMVPMode
	Request        NodeFluidActuatorMVPRequest
	Allowlist      NodeFluidActuatorMVPAllowlist
	KillSwitch     string
	EvaluationTime time.Time
	DryRun         bool
}

func EvaluateNodeFluidActuatorMVP(input NodeFluidActuatorMVPInput) (NodeFluidActuatorMVPEvidence, error) {
	evidence := NodeFluidActuatorMVPEvidence{
		ActuatorID:      "karl-node-fluid-actuator-mvp",
		ActuatorVersion: "v1",
		Mode:            input.Mode,
		DryRun:          input.DryRun || input.Mode == NodeActuatorModeDryRun,
		RequestID:       input.Request.RequestID,
		Allowed:         false,
		PolicyDecision:  "blocked",
		CreatedAt:       normalizeTime(input.EvaluationTime).Format(time.RFC3339),
	}
	targetFile, blockers := validateRequestAgainstAllowlist(input)
	evidence.TargetFile = targetFile
	if input.KillSwitch != "" && isKillSwitchBlocked(input.KillSwitch) {
		blockers = append(blockers, "kill_switch_blocked")
	}
	if len(blockers) > 0 {
		evidence.Blockers = dedupe(blockers)
		return evidence, nil
	}
	before, err := os.ReadFile(targetFile)
	if err != nil {
		evidence.Blockers = dedupe(append(evidence.Blockers, "cpu_max_read_failed"))
		evidence.Notes = append(evidence.Notes, err.Error())
		return evidence, nil
	}
	evidence.BeforeCPUMax = strings.TrimSpace(string(before))
	desired := desiredCPUMaxForMode(input.Mode, input.Request)
	if desired == "" {
		evidence.Blockers = dedupe(append(evidence.Blockers, "invalid_desired_cpu_max"))
		return evidence, nil
	}
	if evidence.DryRun {
		evidence.Allowed = true
		evidence.PolicyDecision = "dry-run-accepted"
		evidence.AppliedCPUMax = desired
		evidence.AfterCPUMax = evidence.BeforeCPUMax
		return evidence, nil
	}
	if err := os.WriteFile(targetFile, []byte(desired+"\n"), 0o600); err != nil {
		evidence.Blockers = dedupe(append(evidence.Blockers, "cpu_max_write_failed"))
		evidence.Notes = append(evidence.Notes, err.Error())
		return evidence, nil
	}
	after, err := os.ReadFile(targetFile)
	if err != nil {
		evidence.Blockers = dedupe(append(evidence.Blockers, "cpu_max_read_after_failed"))
		evidence.Notes = append(evidence.Notes, err.Error())
		return evidence, nil
	}
	evidence.Allowed = true
	evidence.PolicyDecision = "applied"
	evidence.AppliedCPUMax = desired
	evidence.AfterCPUMax = strings.TrimSpace(string(after))
	if input.Mode == NodeActuatorModeReturnToFloor {
		evidence.ReturnToFloorCPUMax = desired
	}
	return evidence, nil
}

func LoadNodeFluidActuatorMVPRequest(path string) (NodeFluidActuatorMVPRequest, error) {
	var request NodeFluidActuatorMVPRequest
	if err := loadJSON(path, &request); err != nil {
		return NodeFluidActuatorMVPRequest{}, err
	}
	return request, nil
}

func LoadNodeFluidActuatorMVPAllowlist(path string) (NodeFluidActuatorMVPAllowlist, error) {
	var allowlist NodeFluidActuatorMVPAllowlist
	if err := loadJSON(path, &allowlist); err != nil {
		return NodeFluidActuatorMVPAllowlist{}, err
	}
	return allowlist, nil
}

func loadJSON(path string, target any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}

func validateRequestAgainstAllowlist(input NodeFluidActuatorMVPInput) (string, []string) {
	req := input.Request
	al := input.Allowlist
	blockers := make([]string, 0, 12)
	if req.RequestID == "" || req.Namespace == "" || req.VMName == "" || req.PodUID == "" || req.QemuPID == "" || req.QemuStartTime == "" || req.CgroupPath == "" {
		blockers = append(blockers, "request_identity_incomplete")
	}
	if al.NodeName != req.NodeName || al.Namespace != req.Namespace || al.VMName != req.VMName || al.PodUID != req.PodUID || al.QemuPID != req.QemuPID || al.QemuStartTime != req.QemuStartTime {
		blockers = append(blockers, "allowlist_identity_mismatch")
	}
	if al.CgroupPath != req.CgroupPath {
		blockers = append(blockers, "allowlist_cgroup_mismatch")
	}
	if al.AllowParentCgroupWrite {
		blockers = append(blockers, "allow_parent_cgroup_write_must_be_false")
	}
	if al.AllowArbitraryWrite {
		blockers = append(blockers, "allow_arbitrary_write_must_be_false")
	}
	if !contains(al.AllowedControllers, "cpu.max") {
		blockers = append(blockers, "cpu_max_not_allowlisted")
	}
	evaluationTime := normalizeTime(input.EvaluationTime)
	createdAt, err := time.Parse(time.RFC3339, req.CreatedAt)
	if err != nil {
		blockers = append(blockers, "request_created_at_invalid")
	} else if req.TTLSeconds > 0 && evaluationTime.Sub(createdAt.UTC()) > time.Duration(req.TTLSeconds)*time.Second {
		blockers = append(blockers, "stale_request")
	}
	if req.PreviousCPUMax == "" || req.RequestedCPUMax == "" || req.RollbackCPUMax == "" {
		blockers = append(blockers, "cpu_max_fields_incomplete")
	}
	if !cpuMaxWithinBounds(req.RequestedCPUMax, al.MinCPUMax, al.MaxCPUMax) {
		blockers = append(blockers, "requested_cpu_max_out_of_bounds")
	}
	if !cpuMaxWithinBounds(req.RollbackCPUMax, al.MinCPUMax, al.MaxCPUMax) {
		blockers = append(blockers, "rollback_cpu_max_out_of_bounds")
	}
	targetFile := filepath.Clean(filepath.Join(req.CgroupPath, "cpu.max"))
	if strings.Contains(targetFile, "..") {
		blockers = append(blockers, "symlink_or_path_escape_detected")
	}
	if !strings.HasSuffix(targetFile, "cpu.max") {
		blockers = append(blockers, "invalid_target_file")
	}
	return targetFile, dedupe(blockers)
}

func desiredCPUMaxForMode(mode NodeFluidActuatorMVPMode, request NodeFluidActuatorMVPRequest) string {
	switch mode {
	case NodeActuatorModeDryRun:
		return request.RequestedCPUMax
	case NodeActuatorModeApply:
		return request.RequestedCPUMax
	case NodeActuatorModeRollback, NodeActuatorModeReturnToFloor:
		return request.RollbackCPUMax
	default:
		return ""
	}
}

func isKillSwitchBlocked(path string) bool {
	raw, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	value := strings.TrimSpace(strings.ToLower(string(raw)))
	return value == "true" || value == "1" || value == "block" || value == "blocked"
}

func cpuMaxWithinBounds(value, minValue, maxValue string) bool {
	quota, period, err := parseCPUMax(value)
	if err != nil {
		return false
	}
	minQuota, minPeriod, minErr := parseCPUMax(minValue)
	maxQuota, maxPeriod, maxErr := parseCPUMax(maxValue)
	if minErr != nil || maxErr != nil {
		return false
	}
	if period != minPeriod || period != maxPeriod {
		return false
	}
	return quota >= minQuota && quota <= maxQuota
}

func parseCPUMax(value string) (int64, int64, error) {
	parts := strings.Fields(strings.TrimSpace(value))
	if len(parts) != 2 {
		return 0, 0, errors.New("cpu.max must contain quota and period")
	}
	var quota int64
	var period int64
	if _, err := fmt.Sscanf(parts[0], "%d", &quota); err != nil {
		return 0, 0, err
	}
	if _, err := fmt.Sscanf(parts[1], "%d", &period); err != nil {
		return 0, 0, err
	}
	return quota, period, nil
}

func normalizeTime(t time.Time) time.Time {
	if t.IsZero() {
		return time.Now().UTC()
	}
	return t.UTC()
}
