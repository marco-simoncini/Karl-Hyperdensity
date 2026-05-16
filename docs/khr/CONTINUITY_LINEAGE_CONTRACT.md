# Continuity Lineage contract (read-only — KHR-AR)

| Field | Value |
|-------|-------|
| **Sprint** | KHR-AR |
| **Contract set** | `khr-tp-contract-v1` |
| **Mode** | Observational only |
| **Production** | **NOT production ready** |

---

## Purpose

End-to-end **continuity lineage/correlation** across identity → lease → route → session → app, for Technical Preview evidence. No enforcement, disconnect, revoke, or forced reconnect.

---

## Identifiers

| Field | Description |
|-------|-------------|
| `continuityLineageId` | Stable lineage id for a session export (`khr-lineage-{sessionKey}`) |
| `sessionCorrelationId` | Correlation id binding gateway session to identity (`khr-correlation-{sessionKey}`) |
| `identityContinuityRef` | Identity-side continuity anchor |
| `shellContinuityRef` | Shell/ShellLease continuity anchor |
| `appContinuityRef` | AppSession / WindowsApp continuity anchor |

---

## Lineage chain

```
UserIdentityRef
  → IdentityBinding
  → ShellEntitlement
  → ShellLeaseRef
  → GatewayRouteRef
  → GatewaySessionRef
  → AppSessionRef / WindowsAppRef
```

Continuity edges (observational):

| Edge | Semantics |
|------|-----------|
| `correlates` | Lineage root correlates session + identity |
| `continues` | Adjacent nodes share uninterrupted TP observation path |

---

## Response fields (rdp-GW)

| Field | Value (KHR-AR stub) |
|-------|---------------------|
| `continuityState` | `observed` |
| `continuityObserved` | `true` |
| `continuityLineageId` | per session |
| `sessionCorrelationId` | per session |
| `lineageRefs` | list of continuity ref strings |

---

## TP invariants

| ID | Invariant |
|----|-----------|
| CL-01 | Continuity is **observational only** — no runtime mutation |
| CL-02 | **No interruption** claims in KHR-AR stub (`interruptionDetected: false`) |
| CL-03 | **No revoke** / **no forced reconnect** |
| CL-04 | **No enforcement** — lineage does not trigger policy apply |
| CL-05 | **Compatibility mapping** allowed (legacy poolId, replicaId) |
| CL-06 | **No CRD/runtime** changes in KHR-AR |
| CL-07 | **No production enable** |

---

## Consumers

| Repo | Artifact |
|------|----------|
| **rdp-GW** | `GET /karl-gw/v1/accessgraph/session` continuity fields + edges |
| **Karl-Dashboard** | `continuitySummary` projection |
| **karl-directoryservice** | `IdentityContinuityRef`, `SessionCorrelationRef` |
| **KARL-APP** | `AppContinuityDescriptor`, `SessionContinuityDescriptor` |
| **Karl-Hyperdensity** | This contract |

---

## Related

- `IDENTITY_ACCESS_LEASE_CONTRACT.md`
- `ACCESS_GRAPH_CONTRACT.md`
