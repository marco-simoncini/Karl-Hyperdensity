# Unified KHR Control Graph (KHR-X)

| Field | Value |
|-------|-------|
| **Sprint** | KHR-X |
| **Purpose** | Shell/Cell/Port-first control-plane model (not VM/VMI-first) |
| **Mutation** | **None** — read-only graph export |

---

## Graph entities

| Kind | Role |
|------|------|
| `Host` | Cluster node anchor |
| `Shell` | User/session shell projection |
| `Cell` | Workload cell (container/session unit) |
| `ResourcePort` | Capability port on a cell |
| `ResourceLease` | Observed lease (dry-run / simulation) |
| `ResourceFuture` | Simulation forecast node |
| `Certification` | Lane certification registry entry |
| `PolicyGate` | Policy gate configuration |
| `ActionApproval` | Operator approval evidence |

---

## Relationships

| Edge | Meaning |
|------|---------|
| `hosts` | Host → Cell |
| `projects` | Shell → Cell |
| `binds` | Cell → ResourcePort → ResourceLease |
| `forecasts` | ResourcePort → ResourceFuture |
| `certifies` | Certification → ResourcePort |
| `gates` | PolicyGate → Certification |
| `approves` | ActionApproval → ResourceFuture |

---

## Lineage / correlation

- `correlationId` — root export correlation (`khr-corr-…`)
- `lineageParentId` — parent node in lineage chain
- Per-node `correlationId` — `khr-lineage-…` derived from root + kind + ref

---

## Health signals

| Signal | Detection |
|--------|-----------|
| `orphan` | ResourcePort without Cell, ActionApproval without ResourceFuture, etc. |
| `stale` | Expired certification or pending approval past TTL |
| `consistent` | No validation issues |

---

## CLI

```bash
go run ./cmd/khr-control-graph \
  -config=examples/khr/runtime-sandbox/karl-host-runtime-config-resourcefuture-simulate.yaml \
  -cluster-context=karl-metal-01@ovh \
  -registry=docs/evidence/khr-certification-registry/registry.json \
  -simulation=simulation-gated.json \
  -approvals=pending-bundle.json \
  -out=graph.json
```

Evidence: `./scripts/khr_control_graph_evidence.sh`
