# Runtime Overlay Model

## Model

Runtime overlay means the active resource authority is applied at runtime without redefining declared baseline/spec as the live control path.

- Declared/template resources represent baseline envelope.
- Runtime desired values are controlled by Hyperdensity executors.
- Observed usage remains independent telemetry.

## VM Linux RAM Overlay

Live VM Linux RAM path is based on virtio-mem runtime control:
- QMP/QOM `requested-size` is the live mutation lane.
- VM template/spec mutation is not the live path.
- Template/spec stays as declared envelope and baseline.

## VM Linux CPU Executor

CPU live lane is guest-assisted:
- libvirt + QGA assisted vCPU up/down
- guest online/offline convergence checks
- same-runtime guard and rollback checks

## Container Linux Live Resize

Container lane uses:
- pod resize for runtime intent
- cgroup/runtime evidence for applied truth
- no rollout/no restart guards
