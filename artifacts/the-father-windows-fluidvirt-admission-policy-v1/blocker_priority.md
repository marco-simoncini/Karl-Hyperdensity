# blocker_priority

Priority tiers implemented:

- `P0_QUARANTINE`
  - `qemu_pid_changed`
  - `last_boot_changed`
  - `machine_guid_changed`
  - `node_changed`
  - `virt_launcher_pod_changed`
  - `hotplug_error_detected`
  - `critical_windows_event_detected`
- `P1_HARD_BLOCK`
  - `qmp_socket_unavailable`
  - `guest_agent_unavailable`
  - `karl_agent_fluid_module_missing`
  - `pending_reboot_detected`
  - `live_migration_required`
  - `vmi_recreate_required`
  - `qmp_ack_missing`
  - `guest_ack_missing`
  - `rollback_not_ready`
  - `return_to_floor_not_ready`
  - `memory_return_not_safe`
- `P2_CAPABILITY_BLOCK`
  - `memory_driver_unverified`
  - `cpu_topology_not_confirmed`
  - `guest_memory_not_confirmed`
- `P3_ENVIRONMENT_BLOCK`
  - `dashboard_443_touch_risk`
  - `candidate_8888_unavailable`
