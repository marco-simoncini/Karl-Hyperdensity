package sidecar

import (
	"errors"
	"fmt"
	"time"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/windowsfluidvirt"
)

type QMPTransport interface {
	Connect() error
	Close() error
	Execute(command string) (map[string]any, error)
}

type ReadOnlyExecutor struct {
	transport QMPTransport
}

func NewReadOnlyExecutor(transport QMPTransport) *ReadOnlyExecutor {
	return &ReadOnlyExecutor{transport: transport}
}

func (e *ReadOnlyExecutor) Execute(command string) (map[string]any, error) {
	if !windowsfluidvirt.IsQMPCommandReadOnly(command) {
		return nil, fmt.Errorf("command %s is forbidden in read-only sidecar", command)
	}
	return e.transport.Execute(command)
}

func (e *ReadOnlyExecutor) DiscoverEvidence(socketPath, qemuPID, sidecarVersion string) windowsfluidvirt.QMPEvidence {
	evidence := windowsfluidvirt.NewReadOnlyQMPEvidence(sidecarVersion, socketPath)
	evidence.QemuPID = qemuPID
	if err := e.transport.Connect(); err != nil {
		evidence.QMPConnected = false
		evidence.QMPErrors = append(evidence.QMPErrors, err.Error())
		return evidence
	}
	defer func() {
		_ = e.transport.Close()
	}()

	evidence.QMPConnected = true
	evidence.QMPGreetingObserved = true
	commands := []string{
		"qmp_capabilities",
		"query-status",
		"query-cpus-fast",
		"query-hotpluggable-cpus",
		"query-memory-devices",
		"query-memory-size-summary",
	}
	for _, command := range commands {
		result, err := e.Execute(command)
		evidence.QMPCommandsExecuted = append(evidence.QMPCommandsExecuted, command)
		if err != nil {
			evidence.QMPErrors = append(evidence.QMPErrors, err.Error())
			continue
		}
		if command == "qmp_capabilities" {
			evidence.QMPCapabilitiesNegotiated = true
		}
		switch command {
		case "query-cpus-fast", "query-hotpluggable-cpus":
			evidence.CPUTopologyObserved = true
			evidence.HotpluggableCPUsObserved = true
			if maxCPUs, ok := result["maxCpus"].(int64); ok {
				evidence.MaxCPUsObserved = maxCPUs
			}
		case "query-memory-devices":
			evidence.MemoryDevicesObserved = true
		case "query-memory-size-summary":
			evidence.MemoryBackendsObserved = true
		}
	}
	evidence.Timestamps = append(evidence.Timestamps, time.Now().UTC().Format(time.RFC3339))
	return evidence
}

type StaticTransport struct {
	Connected bool
	Responses map[string]map[string]any
	Errors    map[string]error
}

func (s *StaticTransport) Connect() error {
	if !s.Connected {
		return errors.New("qmp socket unavailable")
	}
	return nil
}

func (s *StaticTransport) Close() error { return nil }

func (s *StaticTransport) Execute(command string) (map[string]any, error) {
	if err, ok := s.Errors[command]; ok {
		return nil, err
	}
	if response, ok := s.Responses[command]; ok {
		return response, nil
	}
	return map[string]any{"ok": true}, nil
}
