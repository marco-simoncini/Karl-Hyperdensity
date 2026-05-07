# Actuator Request / Allowlist

Request and allowlist used:

- `raw_logs_sanitized/actuator_request_set_floor.json`
- `raw_logs_sanitized/actuator_request_cpu_up.json`
- `raw_logs_sanitized/actuator_allowlist_master-win11.json`
- `raw_logs_sanitized/kill_switch_state.txt`

Validated identity bindings:

- namespace: `karl`
- vm: `master-win11`
- pod UID: `7b6a904a-1c9a-4a44-9b37-1dc737304773`
- node: `karl-lab-metal-01`
- host qemu pid: `2480927`
- qemu start: `Thu May 7 18:58:03 2026`
- cgroup target: `/host-sys-fs-cgroup/kubepods.slice/kubepods-burstable.slice/kubepods-burstable-pod7b6a904a_1c9a_4a44_9b37_1dc737304773.slice/cri-containerd-6b3736cf90a6e2f41147f27714926ac68e57ff1a2d2b22da2375a4a74fcf0c87.scope`
