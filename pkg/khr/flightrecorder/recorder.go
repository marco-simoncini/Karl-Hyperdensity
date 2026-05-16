package flightrecorder

import (
	"sync"
	"time"
)

// Event is one flight-recorder entry (in-memory only).
type Event struct {
	At      string `json:"at"`
	Phase   string `json:"phase"`
	Message string `json:"message"`
	Data    string `json:"data,omitempty"`
}

var (
	mu     sync.Mutex
	events []Event
)

// Record appends an event to the in-memory flight recorder.
func Record(phase, message, data string) {
	mu.Lock()
	defer mu.Unlock()
	events = append(events, Event{
		At:      time.Now().UTC().Format(time.RFC3339Nano),
		Phase:   phase,
		Message: message,
		Data:    data,
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
}
