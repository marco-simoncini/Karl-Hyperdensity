# Blocker Taxonomy

Source: `pkg/windowsfluidvirt/blockers.go`

All mandatory blocker IDs are implemented with:

- stable id
- severity
- category
- human-readable message
- remediable flag
- resulting phase (`BLOCKED`/`QUARANTINED`)
- evidence requirement

Mandatory blockers implemented:

- qmp_socket_unavailable
- guest_agent_unavailable
- karl_agent_fluid_module_missing
- pending_reboot_detected
- qemu_pid_changed
- last_boot_changed
- machine_guid_changed
- live_migration_required
- vmi_recreate_required
- virt_launcher_pod_changed
- node_changed
- memory_driver_unverified
- memory_return_not_safe
- cpu_topology_not_confirmed
- guest_memory_not_confirmed
- rollback_not_ready
- return_to_floor_not_ready
- qmp_ack_missing
- guest_ack_missing
- hotplug_error_detected
- critical_windows_event_detected
- dashboard_443_touch_risk
- candidate_8888_unavailable
- windows_agent_repo_not_present_in_target_repos
