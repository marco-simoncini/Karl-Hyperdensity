# Runtime Gate Design

Implemented in `pkg/windowsfluidvirt/runtime_gate.go`.

## Core interfaces

- `EvaluateFluidRuntimeGate`
- `EvaluateKubeVirtIdentityContinuity`
- `EvaluateNoMigrationProof`
- `EvaluateNoRecreateProof`
- `ValidateQmpReadiness`
- `EvaluateGuestReadiness`
- `EvaluateFluidShellCertificationReadiness`

## Annotation gates enforced

- `hyperdensity.karl.io/fluid-runtime=true`
- `hyperdensity.karl.io/no-live-migration=true`
- `hyperdensity.karl.io/no-reboot=required`
- `hyperdensity.karl.io/no-recreate=required`
- `hyperdensity.karl.io/runtime-mode=in-place-qmp`
- `hyperdensity.karl.io/single-node-compatible=true`

## Conditions emitted

- `fluidRuntimeReady`
- `qmpReady`
- `guestAckReady`
- `noMigrationRequired`
- `noRebootProof`
- `sameQemuProcess`
- `sameNode`
- `sameVirtLauncherPod`
- `returnToFloorReady`
- `rollbackReady`

READY is emitted only when all mandatory conditions are true.
