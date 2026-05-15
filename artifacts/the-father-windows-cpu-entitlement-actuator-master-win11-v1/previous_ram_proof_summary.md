# Previous RAM Proof Summary

Source artifact:

- `artifacts/the-father-windows-prearmed-fluid-envelope-master-win11-v1/`
- commit reference: `661527e`

Extracted baseline:

- QEMU PID/start: `96` / `Thu May 7 18:58:03 2026`
- pod/node: `virt-launcher-master-win11-kmwgg` / `karl-lab-metal-01`
- libvirt mode: `qemu:///session`
- cgroup path from QEMU process: `/`
- CPU tuning errors:
  - session mode `schedinfo` unsupported for tuning
  - `cpu.max` read-only in compute container mount
  - `taskset` not permitted
- RAM entitlement evidence:
  - up via QMP balloon to `13958643712`
  - down to floor `12884901888`
  - return-to-floor verified
- continuity evidence:
  - same QEMU and same Windows boot during prior RAM run
  - no VMIM / no migration

Reuse decision:

- RAM mechanism and evidence model remain valid and reusable for envelope decision, but CPU and RAM were not rerun together in one single sprint.
