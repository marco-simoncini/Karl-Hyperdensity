// Package telemetry provides read-only cgroup v2 evidence for KHR Linux (Sprint 8).
package telemetry

// Sink is a legacy placeholder emitter (unused by CLI paths).
type Sink struct{}

func (s *Sink) Describe() string {
	return "telemetry: cgroup v2 reader (read-only)"
}
