# Hyperdensity KHR-native VM lifecycle gap map (KHR-VM-A)

| Field | Value |
|-------|-------|
| Sprint | KHR-VM-A |
| Evidence | Karl-Installer `docs/evidence/khr-vm-kubevirt-dependency-audit/committed-khr-vm-a-v1/` |
| Mode | audit-only |

## Gap focus

- Existing Shell/Cell/ResourcePort/ResourceLease contracts vs full VM lifecycle needs.
- Missing provider-neutral VM intent and reconciliation contracts.
- Missing VM lifecycle mapping between Dashboard, Inventory, and rdp-GW.
- Compatibility-provider model requirements while KubeVirt remains retained.

## Outcome

KHR-native VM lifecycle is **not** marked ready in this sprint. Gap inventory is complete and bounded.
