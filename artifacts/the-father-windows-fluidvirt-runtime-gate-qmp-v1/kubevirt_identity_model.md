# KubeVirt Identity Model

Implemented in `pkg/windowsfluidvirt/runtime_identity.go`.

Fields included:

- `vmName`, `vmNamespace`, `vmUid`
- `vmiName`, `vmiUid`, `vmiPhase`
- `virtLauncherPodName`, `virtLauncherPodUid`
- `nodeName`
- `podRestartCount`
- `containerIds`
- `qemuPid`
- `qmpSocketPath`
- `liveMigrationObjectsObserved`
- `vmimObjectsObserved`
- `migrationRequired`
- `recreateRequired`
- `rolloutObserved`
- `timestamps`

Proof evaluators:

- `sameNode`
- `sameVirtLauncherPod`
- `sameQemuProcess`
- `noMigration`
- `noRecreate`
- `noRollout`
