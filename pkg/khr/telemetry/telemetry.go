// Package telemetry is a no-op stub for host metrics export (Sprint 5).
package telemetry

// Sink is a placeholder emitter.
type Sink struct{}

func (s *Sink) Describe() string {
	return "telemetry: stub (no export)"
}
