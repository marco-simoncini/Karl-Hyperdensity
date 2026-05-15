# WindowsHyperdensityTarget Model

Model implemented in `pkg/windowsfluidvirt/product_model.go`.

Highlights:

- standalone and pool-child target kinds
- runtime mode pinned to `prearmed-fluid-envelope-v2`
- CPU mechanism pinned to `cgroup-v2-cpu-max`
- RAM mechanism pinned to `qmp-balloon`
- continuity guarantees embedded
- pool scaling, logical CPU scaling claim, and vCPU hotplug claim explicitly blocked
