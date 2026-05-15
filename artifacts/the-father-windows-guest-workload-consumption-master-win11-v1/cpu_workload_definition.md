# CPU Workload Definition

Guest workload executed via QGA `guest-exec` + inline PowerShell.

- bounded duration per run: ~20s
- workerCount: `Environment.ProcessorCount` (6)
- each worker executes CPU-bound loop (`Math.Sqrt`) until deadline
- output JSON includes:
  - start/end timestamps
  - durationMs
  - logicalProcessors
  - workerCount
  - iterations
  - iterationsPerSecond
  - errors
- no binary installs, no persistent services, no reboot
- cleanup: jobs removed at end of each run
