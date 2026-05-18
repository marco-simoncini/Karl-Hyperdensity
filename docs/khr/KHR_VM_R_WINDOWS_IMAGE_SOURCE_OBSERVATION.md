# KHR-VM-R Windows Image Source Observation

Hyperdensity observes KHR-VM-R as a planning-only Windows image source readiness signal for Cell/Shell materialization. Image readiness requires a real local qcow2/raw path plus successful `qemu-img info`.

ResourcePort and ResourceLease remain planning-only. `resourceLeaseApply=false` and `resourcePortPersistentLoop=false`.
