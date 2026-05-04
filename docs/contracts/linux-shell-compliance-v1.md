# linux-shell-compliance-v1

Contract ID: `hyperdensity_linux_shell_compliance_v1`

Defines operational compliance for Linux VM and Linux container references.

Required guarantees:
- CPU up live
- CPU down live
- RAM up live
- RAM down live
- no reboot
- no VMI recreate
- no rollout
- no destructive migration
- same-runtime proof
- rollback proof

Any restart-bound path is non-compliant.
