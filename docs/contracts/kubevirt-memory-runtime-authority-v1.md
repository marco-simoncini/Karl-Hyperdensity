# kubevirt-memory-runtime-authority-v1

Contract ID: `kubevirt_memory_runtime_authority_v1`

Defines runtime authority for VM Linux memory overlay.

- Declared spec/template remains envelope and baseline.
- Runtime overlay is authority for live requested memory.
- Mechanism: virtio-mem QMP/QOM `requested-size`.
- Live path is explicitly not template/spec mutation.

Platform ownership remains in `Karl-OS-ISO` for patch/productization lanes.
