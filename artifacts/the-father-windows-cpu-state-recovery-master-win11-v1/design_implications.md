# Design Implications

Why first `+CPU` was not guest-confirmed:

- QEMU/libvirt successfully created hotplugged CPU object (`vcpu6`), but Windows guest logical processor count stayed at 6.
- Guest view indicates CPU objects beyond active baseline remain non-functional (`PnP Status=Unknown` on higher ACPI processor instances).

Why `-CPU` timed out:

- unplug attempted against a CPU that guest never transitioned into a usable online logical processor.
- guest agent cannot assist CPU state transitions (`guest-set-vcpus` disabled, `can-offline=false` for visible CPUs).

Hyperdensity model violations observed:

- runtime actual diverged from observed guest actual (`7` vs `6`);
- return-to-floor not deterministic in-place;
- rollback path failed live (both `setvcpus --live` and `device_del` ineffective).

Technical paths still valid for future work:

- explicit device-id protocol around QMP add/del with strict readiness checks;
- prearmed topology tuned for Windows CPU grouping/online behavior;
- Windows image/driver/hypervisor compatibility review for ACPI CPU hotplug and unplug;
- Windows-side CPU online evidence adapter before declaring runtime actual;
- policy option to block CPU downscale on Windows profiles lacking deterministic unplug support.

Non-proposed paths:

- no replicas/pools;
- no migration;
- no reboot as Hyperdensity runtime mechanism.
