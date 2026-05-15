# Cluster Probe Summary

status = cluster_probe_run

Read-only probe executed against current context.

Context:

- `karl-metal-01@ovh`

Evidence log:

- `build_or_test_logs/cluster-probe.log`

Observed relevant objects (sanitized summary):

- Existing Windows VM candidate observed: `master-win11` in namespace `karl`
- Pool replicas observed as context only: `win11-pool-0`, `win11-pool-1`
- No VMIM objects returned by the probe command output

Classification:

- `existing_windows_vm_candidate`
- `replica_pool_context_only`

No mutating cluster command was executed.
