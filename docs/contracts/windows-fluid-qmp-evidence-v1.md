# windows-fluid-qmp-evidence-v1

Read-only QMP evidence contract emitted by `karl-fluid-sidecar`.

## Required fields

- `sidecarVersion`
- `qmpConnected`
- `qmpGreetingObserved`
- `qmpCapabilitiesNegotiated`
- `qmpSocketPath`
- `qemuPid`
- `qemuProcessStartTime` (optional when unknown)
- `cpuTopologyObserved`
- `maxCpusObserved`
- `hotpluggableCpusObserved`
- `memoryDevicesObserved`
- `memoryBackendsObserved`
- `qmpCommandsExecuted`
- `qmpReadOnly`
- `qmpErrors`
- `timestamps`

## Contract rules

- `qmpReadOnly` must stay `true`.
- `qmpCommandsExecuted` can only include read-only allowlisted commands.
- Any mutating command attempt must be rejected and not executed.
- `qmpConnected=false` maps to `qmp_socket_unavailable` / `qmp_ack_missing` blockers depending on context.
