# windows-fluidvirt-readonly-cluster-discovery-v1

Read-only runbook for discovering Windows FluidVirt candidates on `karl-metal-01@ovh`.

## Safety constraints

- Read-only commands only.
- Never run `kubectl apply`, `kubectl patch`, `kubectl rollout`, `helm upgrade`.
- Never print tokens, kubeconfig content, secret payloads, bearer credentials.

## Read-only command set

```bash
kubectl get vm,vmi,pods,vmim -A
kubectl describe vm <name> -n <namespace>
kubectl describe vmi <name> -n <namespace>
kubectl get pods -n <namespace> -o wide
kubectl get events -n <namespace> --sort-by=.lastTimestamp
```

## Classification outcomes

- `existing_windows_vm_candidate`
- `replica_pool_context_only`
- `unsupported_generic_windows_vm`
- `ready_for_fluid_shell_certification`
- `blocked_missing_annotations`
- `blocked_missing_guest_agent`
- `blocked_missing_qmp_sidecar`
- `blocked_live_migration_required`

## Evidence collection checklist

- VM, VMI, virt-launcher pod identity
- Node identity
- VMIM/live migration object absence
- Guest evidence availability (`modules.fluidShell`)
- QMP sidecar evidence availability (read-only)

## Notes on replica pool

- Replica pool can be present as context only.
- Replica pool must not be used as the FluidVirt success mechanism.
- Success is continuity and readiness proof on the same VM runtime.
