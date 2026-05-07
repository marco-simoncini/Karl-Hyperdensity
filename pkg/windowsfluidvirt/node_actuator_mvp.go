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

type KARLNodeFluidActuatorAction string

const (
	NodeActuatorActionDryRun        KARLNodeFluidActuatorAction = "dry-run"
	NodeActuatorActionApply         KARLNodeFluidActuatorAction = "apply"
	NodeActuatorActionRollback      KARLNodeFluidActuatorAction = "rollback"
	NodeActuatorActionReturnToFloor KARLNodeFluidActuatorAction = "return-to-floor"
)

type KARLNodeFluidActuatorRequest struct {
	RequestID       string                      `json:"requestId"`
	RequestVersion  string                      `json:"requestVersion"`
	Action          KARLNodeFluidActuatorAction `json:"action"`
	Namespace       string                      `json:"namespace"`
	VMName          string                      `json:"vmName"`
	VMUID           string                      `json:"vmUid"`
	VMIUID          string                      `json:"vmiUid"`
	PodName         string                      `json:"podName"`
	PodUID          string                      `json:"podUid"`
	NodeName        string                      `json:"nodeName"`
	QemuPID         string                      `json:"qemuPid"`
	QemuStartTime   string                      `json:"qemuStartTime"`
	CgroupPath      string                      `json:"cgroupPath"`
	Controller      string                      `json:"controller"`
	PreviousCPUMax  string                      `json:"previousCpuMax"`
	RequestedCPUMax string                      `json:"requestedCpuMax"`
	RollbackCPUMax  string                      `json:"rollbackCpuMax"`
	MinCPUMax       string                      `json:"minCpuMax"`
	MaxCPUMax       string                      `json:"maxCpuMax"`
	TTLSeconds      int64                       `json:"ttlSeconds"`
	CreatedAt       string                      `json:"createdAt"`
	ExpiresAt       string                      `json:"expiresAt"`
	Reason          string                      `json:"reason"`
	Risk            string                      `json:"risk"`
	EvidenceRefs    []string                    `json:"evidenceRefs"`
	PolicyVersion   string                      `json:"policyVersion"`
	AttestationRef  string                      `json:"attestationRef,omitempty"`
}

type KARLNodeFluidActuatorAllowlist struct {
	AllowlistID            string                        `json:"allowlistId"`
	NodeName               string                        `json:"nodeName"`
	Namespace              string                        `json:"namespace"`
	VMName                 string                        `json:"vmName"`
	VMUID                  string                        `json:"vmUid"`
	PodUID                 string                        `json:"podUid"`
	QemuPID                string                        `json:"qemuPid"`
	QemuStartTime          string                        `json:"qemuStartTime"`
	AllowedCgroupPath      string                        `json:"allowedCgroupPath"`
	AllowedControllers     []string                      `json:"allowedControllers"`
	MinCPUMax              string                        `json:"minCpuMax"`
	MaxCPUMax              string                        `json:"maxCpuMax"`
	AllowParentCgroupWrite bool                          `json:"allowParentCgroupWrite"`
	AllowArbitraryWrite    bool                          `json:"allowArbitraryWrite"`
	AllowSymlinkTraversal  bool                          `json:"allowSymlinkTraversal"`
	AllowedActions         []KARLNodeFluidActuatorAction `json:"allowedActions"`
	ExpiresAt              string                        `json:"expiresAt"`
	CreatedAt              string                        `json:"createdAt"`
}

type KARLNodeFluidActuatorDecision string

const (
	NodeActuatorDecisionAccepted        KARLNodeFluidActuatorDecision = "accepted"
	NodeActuatorDecisionRejected        KARLNodeFluidActuatorDecision = "rejected"
	NodeActuatorDecisionApplied         KARLNodeFluidActuatorDecision = "applied"
	NodeActuatorDecisionRolledBack      KARLNodeFluidActuatorDecision = "rolled_back"
	NodeActuatorDecisionReturnedToFloor KARLNodeFluidActuatorDecision = "returned_to_floor"
	NodeActuatorDecisionBlocked         KARLNodeFluidActuatorDecision = "blocked"
)

type KARLNodeFluidActuatorResult struct {
	ResultID             string                        `json:"resultId"`
	RequestID            string                        `json:"requestId"`
	Action               KARLNodeFluidActuatorAction   `json:"action"`
	Decision             KARLNodeFluidActuatorDecision `json:"decision"`
	DryRun               bool                          `json:"dryRun"`
	MutationPerformed    bool                          `json:"mutationPerformed"`
	PreviousCPUMax       string                        `json:"previousCpuMax,omitempty"`
	RequestedCPUMax      string                        `json:"requestedCpuMax,omitempty"`
	ObservedBeforeCPUMax string                        `json:"observedBeforeCpuMax,omitempty"`
	ObservedAfterCPUMax  string                        `json:"observedAfterCpuMax,omitempty"`
	RollbackCPUMax       string                        `json:"rollbackCpuMax,omitempty"`
	ReturnToFloorCPUMax  string                        `json:"returnToFloorCpuMax,omitempty"`
	QemuPID              string                        `json:"qemuPid"`
	QemuStartTime        string                        `json:"qemuStartTime"`
	PodUID               string                        `json:"podUid"`
	NodeName             string                        `json:"nodeName"`
	EvidenceRefs         []string                      `json:"evidenceRefs,omitempty"`
	Blockers             []string                      `json:"blockers,omitempty"`
	AuditHash            string                        `json:"auditHash"`
	CreatedAt            string                        `json:"createdAt"`
}

type NodeFluidActuatorMVPInput struct {
	Action          KARLNodeFluidActuatorAction
	Request         KARLNodeFluidActuatorRequest
	Allowlist       KARLNodeFluidActuatorAllowlist
	KillSwitchPath  string
	EvaluationTime  time.Time
	EvidenceOutPath string
	DryRun          bool
}

type NodeFluidActuatorValidation struct {
	TargetFile string
	Blockers   []string
}

func EvaluateNodeFluidActuatorMVP(input NodeFluidActuatorMVPInput) (KARLNodeFluidActuatorResult, error) {
	evaluationTime := normalizeTime(input.EvaluationTime)
	result := KARLNodeFluidActuatorResult{
		ResultID:          "node-fluid-result-" + shortHash(input.Request.RequestID+"|"+string(input.Action)+"|"+evaluationTime.Format(time.RFC3339)),
		RequestID:         input.Request.RequestID,
		Action:            input.Action,
		Decision:          NodeActuatorDecisionBlocked,
		DryRun:            input.DryRun || input.Action == NodeActuatorActionDryRun,
		PreviousCPUMax:    input.Request.PreviousCPUMax,
		RequestedCPUMax:   input.Request.RequestedCPUMax,
		RollbackCPUMax:    input.Request.RollbackCPUMax,
		QemuPID:           input.Request.QemuPID,
		QemuStartTime:     input.Request.QemuStartTime,
		PodUID:            input.Request.PodUID,
		NodeName:          input.Request.NodeName,
		EvidenceRefs:      input.Request.EvidenceRefs,
		CreatedAt:         evaluationTime.Format(time.RFC3339),
		MutationPerformed: false,
	}

	validation := ValidateNodeFluidActuatorRequest(input)
	if len(validation.Blockers) > 0 {
		result.Decision = NodeActuatorDecisionRejected
		result.Blockers = validation.Blockers
		return finalizeNodeActuatorResult(result)
	}

	before, err := os.ReadFile(validation.TargetFile)
	if err != nil {
		result.Decision = NodeActuatorDecisionBlocked
		result.Blockers = []string{"cpu_max_read_failed"}
		return finalizeNodeActuatorResult(result)
	}
	result.ObservedBeforeCPUMax = strings.TrimSpace(string(before))
	if result.ObservedBeforeCPUMax != input.Request.PreviousCPUMax {
		result.Decision = NodeActuatorDecisionRejected
		result.Blockers = []string{"previous_cpu_max_mismatch"}
		return finalizeNodeActuatorResult(result)
	}

	desired := desiredCPUMaxForAction(input.Action, input.Request)
	if desired == "" {
		result.Decision = NodeActuatorDecisionRejected
		result.Blockers = []string{"invalid_desired_cpu_max"}
		return finalizeNodeActuatorResult(result)
	}

	if result.DryRun {
		result.Decision = NodeActuatorDecisionAccepted
		result.ObservedAfterCPUMax = result.ObservedBeforeCPUMax
		if input.Action == NodeActuatorActionReturnToFloor {
			result.ReturnToFloorCPUMax = desired
		}
		return finalizeNodeActuatorResult(result)
	}

	if err := os.WriteFile(validation.TargetFile, []byte(desired+"\n"), 0o600); err != nil {
		result.Decision = NodeActuatorDecisionBlocked
		result.Blockers = []string{"cpu_max_write_failed"}
		return finalizeNodeActuatorResult(result)
	}
	after, err := os.ReadFile(validation.TargetFile)
	if err != nil {
		result.Decision = NodeActuatorDecisionBlocked
		result.Blockers = []string{"cpu_max_read_after_failed"}
		return finalizeNodeActuatorResult(result)
	}
	result.MutationPerformed = true
	result.ObservedAfterCPUMax = strings.TrimSpace(string(after))
	if result.ObservedAfterCPUMax != desired {
		result.Decision = NodeActuatorDecisionBlocked
		result.Blockers = []string{"cpu_max_after_readback_mismatch"}
		return finalizeNodeActuatorResult(result)
	}

	switch input.Action {
	case NodeActuatorActionApply:
		result.Decision = NodeActuatorDecisionApplied
	case NodeActuatorActionRollback:
		result.Decision = NodeActuatorDecisionRolledBack
	case NodeActuatorActionReturnToFloor:
		result.Decision = NodeActuatorDecisionReturnedToFloor
		result.ReturnToFloorCPUMax = desired
	default:
		result.Decision = NodeActuatorDecisionBlocked
		result.Blockers = []string{"unsupported_action"}
	}
	return finalizeNodeActuatorResult(result)
}

func ValidateNodeFluidActuatorRequest(input NodeFluidActuatorMVPInput) NodeFluidActuatorValidation {
	req := input.Request
	al := input.Allowlist
	blockers := make([]string, 0, 24)

	if req.RequestID == "" || req.RequestVersion == "" || req.Namespace == "" || req.VMName == "" || req.VMUID == "" || req.VMIUID == "" || req.PodUID == "" || req.QemuPID == "" || req.QemuStartTime == "" || req.CgroupPath == "" {
		blockers = append(blockers, "request_identity_incomplete")
	}
	if req.Controller != "cpu.max" {
		blockers = append(blockers, "controller_not_allowed")
	}
	if req.Action == "" || !isAllowedAction(req.Action) {
		blockers = append(blockers, "invalid_action")
	}
	if input.Action != "" && input.Action != req.Action {
		blockers = append(blockers, "requested_action_mismatch")
	}
	if !containsAction(al.AllowedActions, req.Action) {
		blockers = append(blockers, "action_not_allowlisted")
	}
	if al.NodeName != req.NodeName || al.Namespace != req.Namespace || al.VMName != req.VMName || al.VMUID != req.VMUID || al.PodUID != req.PodUID || al.QemuPID != req.QemuPID || al.QemuStartTime != req.QemuStartTime {
		blockers = append(blockers, "allowlist_identity_mismatch")
	}
	if al.AllowedCgroupPath != req.CgroupPath {
		blockers = append(blockers, "allowlist_cgroup_mismatch")
	}
	if al.AllowParentCgroupWrite {
		blockers = append(blockers, "allow_parent_cgroup_write_must_be_false")
	}
	if al.AllowArbitraryWrite {
		blockers = append(blockers, "allow_arbitrary_write_must_be_false")
	}
	if al.AllowSymlinkTraversal {
		blockers = append(blockers, "allow_symlink_traversal_must_be_false")
	}
	if !contains(al.AllowedControllers, "cpu.max") {
		blockers = append(blockers, "cpu_max_not_allowlisted")
	}
	if input.KillSwitchPath != "" && isKillSwitchBlocked(input.KillSwitchPath) {
		blockers = append(blockers, "kill_switch_blocked")
	}
	if req.TTLSeconds <= 0 {
		blockers = append(blockers, "ttl_required")
	}
	if req.RollbackCPUMax == "" {
		blockers = append(blockers, "rollback_target_missing")
	}
	if req.MinCPUMax == "" || req.MaxCPUMax == "" {
		blockers = append(blockers, "request_bounds_missing")
	}
	if input.EvidenceOutPath != "" && !isSafeAuditPath(input.EvidenceOutPath) {
		blockers = append(blockers, "unsafe_audit_output_path")
	}

	createdAt, createdErr := time.Parse(time.RFC3339, req.CreatedAt)
	expiresAt, expiresErr := time.Parse(time.RFC3339, req.ExpiresAt)
	evaluationTime := normalizeTime(input.EvaluationTime)
	if createdErr != nil || expiresErr != nil {
		blockers = append(blockers, "request_time_invalid")
	} else {
		if expiresAt.Before(createdAt) {
			blockers = append(blockers, "request_expiry_before_creation")
		}
		ttlFromFields := createdAt.Add(time.Duration(req.TTLSeconds) * time.Second)
		if !expiresAt.Equal(ttlFromFields) {
			blockers = append(blockers, "request_expiry_ttl_mismatch")
		}
		if evaluationTime.After(expiresAt.UTC()) {
			blockers = append(blockers, "stale_request")
		}
	}

	if al.ExpiresAt != "" {
		allowlistExpiresAt, err := time.Parse(time.RFC3339, al.ExpiresAt)
		if err != nil || evaluationTime.After(allowlistExpiresAt.UTC()) {
			blockers = append(blockers, "allowlist_expired")
		}
	}

	if !cpuMaxWithinBounds(req.RequestedCPUMax, req.MinCPUMax, req.MaxCPUMax) {
		blockers = append(blockers, "requested_cpu_max_outside_request_bounds")
	}
	if !cpuMaxWithinBounds(req.RollbackCPUMax, req.MinCPUMax, req.MaxCPUMax) {
		blockers = append(blockers, "rollback_cpu_max_outside_request_bounds")
	}
	if !cpuMaxWithinBounds(req.RequestedCPUMax, al.MinCPUMax, al.MaxCPUMax) {
		blockers = append(blockers, "requested_cpu_max_outside_allowlist_bounds")
	}
	if !cpuMaxWithinBounds(req.RollbackCPUMax, al.MinCPUMax, al.MaxCPUMax) {
		blockers = append(blockers, "rollback_cpu_max_outside_allowlist_bounds")
	}

	if !isSafeCgroupPath(req.CgroupPath, req.CgroupPath, al.AllowSymlinkTraversal) {
		blockers = append(blockers, "symlink_or_path_escape_detected")
	}
	if !isSafeCgroupPath(req.CgroupPath, al.AllowedCgroupPath, al.AllowSymlinkTraversal) {
		blockers = append(blockers, "cgroup_path_not_exact_allowlist")
	}

	targetFile := filepath.Clean(filepath.Join(req.CgroupPath, req.Controller))
	expectedFile := filepath.Clean(filepath.Join(req.CgroupPath, "cpu.max"))
	if targetFile != expectedFile {
		blockers = append(blockers, "target_file_must_be_cpu_max")
	}
	if info, err := os.Lstat(targetFile); err == nil && (info.Mode()&os.ModeSymlink) != 0 {
		blockers = append(blockers, "symlink_or_path_escape_detected")
	}
	if strings.TrimSuffix(targetFile, "/cpu.max") != req.CgroupPath {
		blockers = append(blockers, "parent_cgroup_write_forbidden")
	}

	beforeRaw, readErr := os.ReadFile(targetFile)
	if readErr != nil {
		blockers = append(blockers, "cpu_max_read_failed")
	} else if strings.TrimSpace(string(beforeRaw)) != req.PreviousCPUMax {
		blockers = append(blockers, "previous_cpu_max_mismatch")
	}

	return NodeFluidActuatorValidation{
		TargetFile: targetFile,
		Blockers:   dedupe(blockers),
	}
}

func LoadNodeFluidActuatorMVPRequest(path string) (KARLNodeFluidActuatorRequest, error) {
	var request KARLNodeFluidActuatorRequest
	if err := loadNodeActuatorJSON(path, &request); err != nil {
		return KARLNodeFluidActuatorRequest{}, err
	}
	return request, nil
}

func LoadNodeFluidActuatorMVPAllowlist(path string) (KARLNodeFluidActuatorAllowlist, error) {
	var allowlist KARLNodeFluidActuatorAllowlist
	if err := loadNodeActuatorJSON(path, &allowlist); err != nil {
		return KARLNodeFluidActuatorAllowlist{}, err
	}
	return allowlist, nil
}

func loadNodeActuatorJSON(path string, target any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}

func desiredCPUMaxForAction(action KARLNodeFluidActuatorAction, request KARLNodeFluidActuatorRequest) string {
	switch action {
	case NodeActuatorActionDryRun, NodeActuatorActionApply:
		return request.RequestedCPUMax
	case NodeActuatorActionRollback:
		return request.RollbackCPUMax
	case NodeActuatorActionReturnToFloor:
		return request.PreviousCPUMax
	default:
		return ""
	}
}

func isAllowedAction(action KARLNodeFluidActuatorAction) bool {
	switch action {
	case NodeActuatorActionDryRun, NodeActuatorActionApply, NodeActuatorActionRollback, NodeActuatorActionReturnToFloor:
		return true
	default:
		return false
	}
}

func containsAction(actions []KARLNodeFluidActuatorAction, action KARLNodeFluidActuatorAction) bool {
	for _, current := range actions {
		if current == action {
			return true
		}
	}
	return false
}

func isSafeAuditPath(path string) bool {
	clean := filepath.Clean(path)
	if strings.Contains(clean, "..") {
		return false
	}
	return strings.HasSuffix(clean, ".json")
}

func isSafeCgroupPath(path, expected string, allowSymlinkTraversal bool) bool {
	if filepath.Clean(path) != filepath.Clean(expected) || strings.Contains(path, "..") {
		return false
	}
	if allowSymlinkTraversal {
		return true
	}
	if info, err := os.Lstat(path); err == nil && (info.Mode()&os.ModeSymlink) != 0 {
		return false
	}
	return true
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

func finalizeNodeActuatorResult(result KARLNodeFluidActuatorResult) (KARLNodeFluidActuatorResult, error) {
	result.AuditHash = ""
	hash, err := computeDeterministicHash(result)
	if err != nil {
		return KARLNodeFluidActuatorResult{}, fmt.Errorf("compute actuator audit hash: %w", err)
	}
	result.AuditHash = hash
	return result, nil
}

func normalizeTime(t time.Time) time.Time {
	if t.IsZero() {
		return time.Now().UTC()
	}
	return t.UTC()
}
