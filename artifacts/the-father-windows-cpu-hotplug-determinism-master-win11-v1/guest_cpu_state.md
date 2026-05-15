# Guest CPU State

Guest evidence collected via QGA:

- guest ACK: true (`guest-ping` success)
- `guest-get-vcpus`: 6 logical CPUs online (`logical-id` 0..5)
- Windows telemetry (PowerShell via guest-exec):
  - `LogicalProcessors`: `6`
  - `LastBootRaw`: `/Date(1778177812500)/`
  - `MachineGuidHash`: `d1fdf2ce69932d7ac2f9e3497be7f845d4a3ef9140c45bd0d7b0ea47c195508e`
  - `PendingReboot`: `false`
  - `CriticalEvents24h`: `0`

Guest-side conclusion:

- Guest never confirmed `+CPU` to `7`, despite QMP/libvirt showing `7`.
