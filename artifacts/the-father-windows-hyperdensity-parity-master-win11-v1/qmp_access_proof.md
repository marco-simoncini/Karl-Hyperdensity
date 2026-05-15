# QMP Access Proof

## Result

`WINDOWS_HYPERDENSITY_PARITY_BLOCKED_BY_QMP`

## Evidence

- No `VirtualMachineInstance` exists for `master-win11` in `karl`.
- No `virt-launcher` pod exists for `master-win11`.
- Without a running `virt-launcher`, no QMP socket can be attested.

## Allowed commands executed

- `kubectl get vmi master-win11 -n karl -o yaml` -> NotFound
- `kubectl get pod -n karl -l kubevirt.io=virt-launcher,vm.kubevirt.io/name=master-win11 -o wide` -> no resources

## Mutating commands executed

None.
