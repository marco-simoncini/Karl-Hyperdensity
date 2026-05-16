# KHR Access Graph contract (read-only — KHR-AO)

| Field | Value |
|-------|-------|
| **Sprint** | KHR-AO |
| **Contract set** | `khr-tp-contract-v1` |
| **Mode** | Read-only Technical Preview |
| **Production** | **NOT production ready** |

---

## Purpose

Formalize the first **read-only Access Graph** for Wave 2:

`UserIdentity` → `ShellLease` → `GatewayRoute` → `GatewaySession` → `AppSession` / `WindowsApp`

No revoke, disconnect automation, auth enforcement, or mutating session flows.

---

## Nodes

| Kind | Description |
|------|-------------|
| `UserIdentityRef` | Directory/OIDC subject |
| `ShellLeaseRef` | Time-bounded Shell access reservation |
| `GatewayRouteRef` | Session path from lease to protocol endpoint |
| `GatewaySessionRef` | Active gateway/RDP session projection |
| `AppSessionRef` | Application/RemoteApp session correlation |
| `WindowsAppRef` | Published Windows application descriptor |

Each node has stable `id` (`<Kind>:<key>`) and `ref` payload (read-only).

---

## Edges

| Kind | From → To | Semantics |
|------|-----------|-----------|
| `leases` | UserIdentity → ShellLease | Subject holds lease reservation |
| `routes` | ShellLease → GatewayRoute | Lease selects gateway route |
| `opens` | GatewayRoute → GatewaySession | Route opens gateway session |
| `runs` | GatewaySession → AppSession | Session runs app context |
| `runs` | AppSession → WindowsApp | App session runs published app |
| `runs` | GatewaySession → WindowsApp | Direct desktop/RemoteApp (no AppSession node) |
| `owns` | UserIdentity → GatewaySession | Compatibility ownership projection |
| `governedBy` | GatewaySession → ShellLease | Session bounded by lease policy |

---

## Compatibility mapping

Legacy keys may appear in `compatibility[]` (read-only audit):

| Legacy | Target node kind |
|--------|------------------|
| `poolId` | `ShellLeaseRef` (compatibility-only) |
| `replicaId` | `GatewaySessionRef` |
| `remoteApplicationProgram` | `WindowsAppRef` |

**poolId** is legacy-only; **GatewayRoute** remains the target routing model.

---

## Invariants (KHR-AO)

| ID | Invariant |
|----|-----------|
| AG-01 | Graph export is **read-only** (`readOnly: true`, `mutating: false`) |
| AG-02 | **noRevoke** and **noDisconnect** are true on export; no mutating revoke/disconnect fields |
| AG-03 | **No auth enforcement** — graph is observability/correlation only |
| AG-04 | **Compatibility mapping allowed** — legacy ids may appear in `compatibility[]` |
| AG-05 | **No CRD or controller changes** in KHR-AO |
| AG-06 | **No production enable** claims |

---

## Consumers

| Repo | Artifact |
|------|----------|
| **rdp-GW** | `GET /karl-gw/v1/accessgraph/session?id=...` |
| **KARL-APP** | `AccessGraphDescriptor`, `AppLaunchIntentDescriptor` |
| **karl-directoryservice** | `IdentityGraphNode`, `IdentityEntitlementEdge` |
| **Karl-Hyperdensity** | This contract (source of truth) |

---

## Related

- `SHELLLEASE_GATEWAYROUTE_CONTRACT.md`
- rdp-GW `SHELLLEASE_GATEWAYROUTE_COMPATIBILITY.md`
