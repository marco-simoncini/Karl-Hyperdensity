# KHR Access Graph contract (read-only — KHR-AO)

| Field | Value |
|-------|-------|
| **Sprint** | KHR-AO / KHR-AP |
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
| **Karl-Dashboard** | `accessGraphSummary` on KHR projection / TP readiness |
| **Karl-Hyperdensity** | This contract (source of truth) |

---

## Dashboard projection expectations (KHR-AP)

Cockpit exposes **read-only** metadata only — no new mutating actions, no graph editor UI.

| Field | Value |
|-------|-------|
| `accessGraphAvailable` | `true` when rdp-GW contract is advertised |
| `accessGraphEndpoint` | `/karl-gw/v1/accessgraph/session?id={sessionId}` |
| `accessGraphNodeTypes` | Six node kinds (see Nodes table) |
| `accessGraphReadOnly` | always `true` |
| `noRevoke` / `noDisconnect` | always `true` on projection block |

Parent Fabric `khrProjection.accessGraphSummary` and `tpReadinessSummary.accessGraphSummary` document the path:

**User → ShellLease → GatewayRoute → GatewaySession → App/WindowsApp**

See Karl-Dashboard: `docs/khr/DASHBOARD_ACCESS_GRAPH_PROJECTION.md`.

---

## rdp-GW accessgraph endpoint contract (KHR-AP)

| Item | Detail |
|------|--------|
| **Method** | `GET` |
| **Path** | `/karl-gw/v1/accessgraph/session` |
| **Query** | `id` (required) |
| **Body** | `nodes[]`, `edges[]`, `compatibility[]` |
| **Stability** | Golden: `cmd/rdpgw/khr/testdata/khr_accessgraph_session_golden.json` |
| **Docs** | rdp-GW `docs/khr/RDPGW_ACCESS_GRAPH_API.md` |

---

## KARL-APP launch intent relation (KHR-AP)

| Descriptor | Relation |
|------------|----------|
| `AppLaunchIntentDescriptor` | Entry point for app launch correlation; links to `AccessGraphDescriptor` via shared `sessionKey` |
| `WindowsAppDescriptor` | Maps to `WindowsAppRef` node in graph export |
| `AccessGraphDescriptor` | Client-side mirror of rdp-GW graph JSON |

No session mutation; launch intent is documentation/stub only in KHR-AP.

---

## Directory identity relation (KHR-AP)

| Model | Relation |
|-------|----------|
| `IdentityGraphNode` (`UserIdentityRef`) | Root of directory-side graph |
| `IdentityEntitlementEdge` (`leases`) | User → ShellLease entitlement |
| `IdentityEntitlementEdge` (`routes`) | ShellLease → ShellClass / route policy (compatibility) |
| `build_identity_access_graph_stub` | Aligns with Hyperdensity ShellLease `userRef` semantics |

No auth enforcement; OIDC flows unchanged.

---

## Related

- `SHELLLEASE_GATEWAYROUTE_CONTRACT.md`
- rdp-GW `SHELLLEASE_GATEWAYROUTE_COMPATIBILITY.md`
