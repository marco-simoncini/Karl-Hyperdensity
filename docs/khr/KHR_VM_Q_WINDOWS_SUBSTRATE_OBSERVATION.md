# KHR-VM-Q Windows Substrate Observation

Hyperdensity observes KHR-VM-Q as a planning-only Windows substrate readiness signal for Cell/Shell materialization. The observed status is `windows-substrate-ready` only when image, TPM, firmware, and disk attach checks all pass.

ResourcePort and ResourceLease remain planning-only. `resourceLeaseApply=false` and `resourcePortPersistentLoop=false`.
