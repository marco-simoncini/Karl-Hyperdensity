# Executors Package

Documentation-first home for executor contracts.

Current live truths:
- VM Linux CPU: guest-assisted via libvirt/QGA
- VM Linux RAM: runtime overlay via virtio-mem requested-size
- Container Linux CPU/RAM: pod resize + cgroup evidence

Phase 0 does not move executor runtime code out of `Karl-Dashboard`.
