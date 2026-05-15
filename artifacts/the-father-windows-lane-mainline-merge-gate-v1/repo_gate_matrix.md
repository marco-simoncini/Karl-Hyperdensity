# repo_gate_matrix

- Karl-Hyperdensity
  - branch: `The-Father-FluidVirt-Windows-Integration`
  - head: `a60516a`
  - gate status: ready with conditions
  - merge candidate: true
  - notes: core/model/boundary/replay contracts present, Go tests green

- Karl-Inventory
  - branch: `The-Father-FluidVirt-Windows-Integration`
  - head: `f4076ed`
  - gate status: ready with conditions
  - merge candidate: true
  - notes: fluidShell guest witness evidence-only operational, .NET validation passed

- Karl-Dashboard
  - branch source: stale Windows lane (excluded)
  - gate status: excluded
  - merge candidate: false
  - notes: visibility work must restart from official `The-Father` branch

- Karl-OS-ISO
  - branch source: Windows lane deferred
  - gate status: deferred
  - merge candidate: false
  - notes: no packaging/systemd/daemonset in this gate
