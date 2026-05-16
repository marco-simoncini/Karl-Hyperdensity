package host

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"sync"
	"time"
)

const runtimeSessionPrefix = "khr-session-"

// RuntimeSession identifies a karl-host-runtime process and correlates operations.
type RuntimeSession struct {
	RuntimeSessionID      string `json:"runtimeSessionId"`
	HostRuntimeInstanceID string `json:"hostRuntimeInstanceId"`
	CorrelationID         string `json:"correlationId,omitempty"`
	StartedAt             string `json:"startedAt"`
	HostID                string `json:"hostId,omitempty"`
}

var (
	sessionMu          sync.Mutex
	processSession     *RuntimeSession
	processStart       time.Time
	correlationSeq     int
	pinnedCorrelation  string
)

// InitRuntimeSession returns stable session IDs for this process (idempotent).
func InitRuntimeSession(cfg *Config) RuntimeSession {
	sessionMu.Lock()
	defer sessionMu.Unlock()
	return initRuntimeSessionLocked(cfg)
}

func initRuntimeSessionLocked(cfg *Config) RuntimeSession {
	if processSession != nil {
		out := *processSession
		out.CorrelationID = nextCorrelationIDLocked()
		return out
	}
	processStart = time.Now().UTC()
	hostID := "karl-host-unknown"
	if cfg != nil && cfg.Spec.HostID != "" {
		hostID = cfg.Spec.HostID
	}
	hostname, _ := os.Hostname()
	inst := sha256.Sum256([]byte(hostID + "|" + hostname + "|" + processStart.Format(time.RFC3339Nano)))
	instanceID := "khr-inst-" + hex.EncodeToString(inst[:8])
	sessionID := runtimeSessionPrefix + hex.EncodeToString(inst[8:16])
	processSession = &RuntimeSession{
		RuntimeSessionID:      sessionID,
		HostRuntimeInstanceID: instanceID,
		StartedAt:             processStart.Format(time.RFC3339),
		HostID:                hostID,
	}
	out := *processSession
	out.CorrelationID = nextCorrelationIDLocked()
	return out
}

// CurrentRuntimeSession returns the active session with a fresh correlation ID.
func CurrentRuntimeSession() RuntimeSession {
	sessionMu.Lock()
	defer sessionMu.Unlock()
	if processSession == nil {
		return initRuntimeSessionLocked(nil)
	}
	out := *processSession
	out.CorrelationID = nextCorrelationIDLocked()
	return out
}

// SetCorrelationID overrides the next correlation id prefix for an operation chain.
func SetCorrelationID(id string) {
	sessionMu.Lock()
	defer sessionMu.Unlock()
	if id != "" {
		pinnedCorrelation = id
		if processSession != nil {
			processSession.CorrelationID = id
		}
	}
}

func nextCorrelationIDLocked() string {
	if pinnedCorrelation != "" {
		id := pinnedCorrelation
		pinnedCorrelation = ""
		return id
	}
	correlationSeq++
	if processSession == nil {
		return fmt.Sprintf("khr-corr-%d", correlationSeq)
	}
	return fmt.Sprintf("%s-op-%d", processSession.RuntimeSessionID, correlationSeq)
}

// ResetRuntimeSession clears session state (tests only).
func ResetRuntimeSession() {
	sessionMu.Lock()
	defer sessionMu.Unlock()
	processSession = nil
	processStart = time.Time{}
	correlationSeq = 0
	pinnedCorrelation = ""
}
