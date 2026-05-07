# Guest Evidence Integration

Guest evidence integration remains aligned with `Karl-Inventory/modules.fluidShell`.

Validated in tests:

- pending reboot -> `pending_reboot_detected`
- missing machine guid hash -> `machine_guid_changed`
- missing last boot -> `last_boot_changed`
- memory adapter unverified -> `memory_driver_unverified`
- return-to-floor false -> `return_to_floor_not_ready`
- guest ack false -> `guest_ack_missing`

No Inventory source changes were required for this milestone.
