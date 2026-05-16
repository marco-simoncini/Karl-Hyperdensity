package flightrecorder

import (
	"sync"
	"time"
)

// Event is one flight-recorder entry (in-memory only).
type Event struct {
	At                    string `json:"at"`
	Phase                 string `json:"phase"`
	Message               string `json:"message"`
	Data                  string `json:"data,omitempty"`
	RuntimeSessionID      string `json:"runtimeSessionId,omitempty"`
	HostRuntimeInstanceID string `json:"hostRuntimeInstanceId,omitempty"`
	CorrelationID         string `json:"correlationId,omitempty"`
}

var (
	mu                    sync.Mutex
	events                []Event
	runtimeSessionID      string
	hostRuntimeInstanceID string
	correlationID         string
)

// SessionContext carries correlation fields for flight recorder entries (KHR-N).
type SessionContext struct {
	RuntimeSessionID      string
	HostRuntimeInstanceID string
	CorrelationID         string
}

// InitContext sets session correlation fields for subsequent Record calls.
func InitContext(ctx SessionContext) {
	mu.Lock()
	defer mu.Unlock()
	runtimeSessionID = ctx.RuntimeSessionID
	hostRuntimeInstanceID = ctx.HostRuntimeInstanceID
	correlationID = ctx.CorrelationID
}

// SetCorrelation updates the correlation id for the current operation chain.
func SetCorrelation(id string) {
	mu.Lock()
	defer mu.Unlock()
	correlationID = id
}

// Record appends an event to the in-memory flight recorder.
func Record(phase, message, data string) {
	mu.Lock()
	defer mu.Unlock()
	events = append(events, Event{
		At:                    time.Now().UTC().Format(time.RFC3339Nano),
		Phase:                 phase,
		Message:               message,
		Data:                  data,
		RuntimeSessionID:      runtimeSessionID,
		HostRuntimeInstanceID: hostRuntimeInstanceID,
		CorrelationID:         correlationID,
	})
}

// Snapshot returns a copy of recorded events.
func Snapshot() []Event {
	mu.Lock()
	defer mu.Unlock()
	out := make([]Event, len(events))
	copy(out, events)
	return out
}

// Reset clears the recorder (tests).
func Reset() {
	mu.Lock()
	defer mu.Unlock()
	events = nil
	runtimeSessionID = ""
	hostRuntimeInstanceID = ""
	correlationID = ""
}
