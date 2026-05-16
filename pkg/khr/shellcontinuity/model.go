// Package shellcontinuity defines read-only Shell/User/App continuity semantics (KHR-U).
package shellcontinuity

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

const (
	StatePreserved    = "preserved"
	StateInterrupted  = "interrupted"
	StateUnknown      = "unknown"

	SourceSandboxObservation = "sandbox-pod-observation"
)

// Snapshot is a point-in-time workload identity for continuity compare.
type Snapshot struct {
	Namespace    string `json:"namespace"`
	PodName      string `json:"podName"`
	PodUID       string `json:"podUID"`
	ContainerID  string `json:"containerID,omitempty"`
	ShellRef     string `json:"shellRef,omitempty"`
	CellRef      string `json:"cellRef,omitempty"`
	ShellSessionID string `json:"shellSessionId"`
	AppSessionID   string `json:"appSessionId"`
	UserSessionID  string `json:"userSessionId,omitempty"`
}

// Evidence is continuity proof attached to verification/certification (read-only).
type Evidence struct {
	ShellSessionID        string `json:"shellSessionId,omitempty"`
	AppSessionID          string `json:"appSessionId,omitempty"`
	UserSessionID         string `json:"userSessionId,omitempty"`
	ShellContinuityState  string `json:"shellContinuityState"`
	AppContinuityState    string `json:"appContinuityState"`
	UserSessionContinuity string `json:"userSessionContinuityState,omitempty"`
	ContinuityState       string `json:"continuityState"`
	Source                string `json:"source"`
	ReadOnly              bool   `json:"readOnly"`
}

// Proof is before/after comparison output.
type Proof struct {
	ShellContinuityPreserved       bool     `json:"shellContinuityPreserved"`
	AppContinuityPreserved         bool     `json:"appContinuityPreserved"`
	UserSessionContinuityPreserved bool     `json:"userSessionContinuityPreserved"`
	SessionContinuityPreserved     bool     `json:"sessionContinuityPreserved"`
	InterruptionDetected           bool     `json:"interruptionDetected"`
	Evidence                       Evidence `json:"continuityEvidence"`
}

// DeriveSessionIDs computes stable session identifiers from workload identity.
func DeriveSessionIDs(ns, podName, podUID, containerID string) (shell, app, user string) {
	base := fmt.Sprintf("%s/%s/%s", ns, podName, podUID)
	if strings.TrimSpace(containerID) != "" {
		base = base + "/" + normalizeContainerID(containerID)
	}
	sum := sha256.Sum256([]byte(base))
	shell = "shell-session-" + hex.EncodeToString(sum[:8])
	app = shell + "/app"
	user = shell + "/user"
	return shell, app, user
}

func normalizeContainerID(id string) string {
	id = strings.TrimSpace(id)
	if idx := strings.LastIndex(id, "://"); idx >= 0 {
		return id[idx+3:]
	}
	return id
}

// SnapshotFromWorkload builds a snapshot with derived session IDs.
func SnapshotFromWorkload(ns, podName, podUID, containerID string) Snapshot {
	shell, app, user := DeriveSessionIDs(ns, podName, podUID, containerID)
	cellRef := fmt.Sprintf("%s/Cell/%s", ns, podName)
	return Snapshot{
		Namespace: ns, PodName: podName, PodUID: podUID, ContainerID: containerID,
		ShellRef: fmt.Sprintf("%s/Shell/%s", ns, podName), CellRef: cellRef,
		ShellSessionID: shell, AppSessionID: app, UserSessionID: user,
	}
}

// Compare evaluates continuity between before and after snapshots.
func Compare(before, after Snapshot) Proof {
	p := Proof{
		ShellContinuityPreserved:       sessionEqual(before.ShellSessionID, after.ShellSessionID) && uidEqual(before, after),
		AppContinuityPreserved:         sessionEqual(before.AppSessionID, after.AppSessionID) && containerEqual(before, after),
		UserSessionContinuityPreserved: sessionEqual(before.UserSessionID, after.UserSessionID) && uidEqual(before, after),
	}
	p.SessionContinuityPreserved = p.UserSessionContinuityPreserved
	p.InterruptionDetected = !(p.ShellContinuityPreserved && p.AppContinuityPreserved && p.SessionContinuityPreserved)
	p.Evidence = Evidence{
		ShellSessionID:        after.ShellSessionID,
		AppSessionID:          after.AppSessionID,
		UserSessionID:         after.UserSessionID,
		ShellContinuityState:  stateFor(p.ShellContinuityPreserved),
		AppContinuityState:    stateFor(p.AppContinuityPreserved),
		UserSessionContinuity: stateFor(p.UserSessionContinuityPreserved),
		ContinuityState:       aggregateState(p),
		Source:                SourceSandboxObservation,
		ReadOnly:              true,
	}
	return p
}

func sessionEqual(a, b string) bool {
	return strings.TrimSpace(a) != "" && a == b
}

func uidEqual(a, b Snapshot) bool {
	return strings.TrimSpace(a.PodUID) != "" && a.PodUID == b.PodUID
}

func containerEqual(a, b Snapshot) bool {
	if strings.TrimSpace(a.ContainerID) == "" || strings.TrimSpace(b.ContainerID) == "" {
		return uidEqual(a, b)
	}
	return normalizeContainerID(a.ContainerID) == normalizeContainerID(b.ContainerID)
}

func stateFor(preserved bool) string {
	if preserved {
		return StatePreserved
	}
	return StateInterrupted
}

func aggregateState(p Proof) string {
	if p.ShellContinuityPreserved && p.AppContinuityPreserved && p.SessionContinuityPreserved {
		return StatePreserved
	}
	if p.InterruptionDetected {
		return StateInterrupted
	}
	return StateUnknown
}
