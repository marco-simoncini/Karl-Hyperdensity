# pre_apply_revalidation_contract

Model: `WindowsFluidPreApplyRevalidationContract`

Requires fresh evidence for:

- kubevirt identity
- qmp evidence
- guest evidence
- rollback readiness
- return-to-floor readiness
- kill-switch readiness

Requires unchanged comparisons for:

- node
- virt-launcher pod
- qemu pid
- last boot
- machine identity
- no VMIM
- no live migration
- no recreate

Outputs:

- `REVALIDATION_READY`
- `REVALIDATION_BLOCKED`
- `REVALIDATION_QUARANTINED`
- `REVALIDATION_STALE`
