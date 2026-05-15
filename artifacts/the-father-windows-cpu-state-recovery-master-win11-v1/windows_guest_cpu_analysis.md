# Windows Guest CPU Analysis

Guest agent capabilities:

- `guest-get-vcpus`: enabled
- `guest-set-vcpus`: disabled
- consequence: no guest-agent assisted online/offline vCPU control.

Guest CPU views:

- `guest-get-vcpus`: 6 entries (`logical-id 0..5`), all `online=true`, all `can-offline=false`
- `Win32_ComputerSystem.NumberOfLogicalProcessors`: `6`
- `Win32_Processor`:
  - one socket (`CPU0`)
  - `NumberOfCores=6`
  - `NumberOfLogicalProcessors=6`

PnP processor device view:

- processor instances `_0.._6`: `Status=OK`
- processor instances `_7.._B`: `Status=Unknown`
- interpretation: ACPI processor objects above active baseline exist but are not fully functional/online from guest perspective.

Events and system state:

- pending reboot: false
- critical events 24h: `0`
- CPU-related event sample in System log: kernel processor power capabilities for processors `0..5` only.

Root-cause hypothesis (guest side):

- Windows sees baseline 6 CPUs as active and does not transition hotplugged CPU into guest-usable logical processor count.
- Guest-side unplug control is unavailable (`can-offline=false` and `guest-set-vcpus` disabled), reinforcing unplug timeout behavior.
