# windows-fluid-kubevirt-identity-v1

Read-only identity evidence model for KubeVirt-backed Windows Fluid Shell runtime continuity.

## Fields

- `vmName`, `vmNamespace`, `vmUid`
- `vmiName`, `vmiUid`, `vmiPhase`
- `virtLauncherPodName`, `virtLauncherPodUid`
- `nodeName`
- `podRestartCount`
- `containerIds` (when available)
- `qemuPid` (from sidecar evidence when available)
- `qmpSocketPath` (when available)
- `liveMigrationObjectsObserved`
- `vmimObjectsObserved`
- `migrationRequired`
- `recreateRequired`
- `rolloutObserved`
- `timestamps`

## Continuity proofs

- `sameNode` requires same `nodeName`.
- `sameVirtLauncherPod` requires same `virtLauncherPodUid`.
- `sameQemuProcess` requires same `qemuPid`.
- `noMigration` requires no VMIM/live migration evidence.
- `noRecreate` requires stable `vmiUid` and no recreate/rollout markers.
