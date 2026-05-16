# ShellLease + GatewayRoute access foundation (KHR-E)

| Field | Value |
|-------|-------|
| **ShellLease API** | `runtime.karl.io/v1alpha1` |
| **GatewayRoute API** | `gateway.karl.io/v1alpha1` (rdp-GW / gateway alignment) |
| **Controllers** | None |

---

## ShellLease

Time-bounded access reservation binding a **Shell** to a user/tenant.

| Spec field | Required | Description |
|------------|----------|-------------|
| `shellRef` | yes | Shell name (+ optional namespace) |
| `userRef` | no | Subject identity |
| `tenant` | no | Tenancy namespace |
| `leaseMode` | yes | `ephemeral`, `persistent`, `scheduled` |
| `accessProfile` | no | Policy profile id |
| `expiresAt` | no | RFC3339 expiry |

| Status field | Description |
|--------------|-------------|
| `phase` | `Pending`, `Active`, `Expired`, `Revoked`, `Observed` |
| `conditions` | Standard conditions |

---

## GatewayRoute

Session path from **ShellLease** to gateway protocol endpoint.

| Spec field | Required | Description |
|------------|----------|-------------|
| `shellLeaseRef` | no* | Bound lease (*required for KHR-E examples) |
| `protocol` | yes | `rdp`, `rdp-remoteapp`, `web`, `ssh` |
| `targetRef` | no | Backend target (VM, app, URL) |
| `gateway` | no | Gateway instance/class ref |
| `policyRefs` | no | Policy bundle refs |
| `tokenRef` | no | Opaque token issuance ref |

Legacy fields retained: `gatewayClass`, `shellRef`, `rdp`, `tokenPolicy`.

| Status field | Description |
|--------------|-------------|
| `phase` | `Pending`, `Ready`, `Failed`, `Revoked`, `Observed` |
| `conditions` | Standard conditions |

---

## Artifacts

- `api/crds/runtime.karl.io/shelllease.yaml`
- `api/crds/gateway.karl.io/gatewayroute.yaml`
- `docs/contracts/khr/shelllease.schema.json`
- `docs/contracts/khr/gatewayroute.schema.json`
- `docs/contracts/khr/examples/shelllease-demo.json`
- `docs/contracts/khr/examples/gatewayroute-rdp-demo.json`
- `docs/contracts/khr/examples/gatewayroute-remoteapp-demo.json`

---

## rdp-GW consumer expectations (KHR-AM)

| Expectation | Detail |
|-------------|--------|
| **Mode** | Read-only compatibility resolvers only |
| **Endpoints** | `GET /karl-gw/v1/shell/resolve`, `GET /karl-gw/v1/gatewayroute/resolve` |
| **Legacy** | `GET /api/replica/free?id=<poolId>`, `ReplicaProvider.SelectReplica`, `ReplicaSession` unchanged |
| **Identity** | **GatewayRoute** = target session-path model; **poolId** = legacy via `ShellPoolRef.compatibilityOnly` |
| **Mapping** | `ReplicaSession` → `GatewaySession`; RDP host/port → `CellEndpoint`; RemoteApp → `WindowsAppRef` |
| **Provider** | KubeVirt / vmpool remains compatibility provider for replica selection |

### GatewayRoute read-only resolver

rdp-GW `GatewayRouteResolver` (stub) returns `ResolveGatewayRouteResponse` with `readOnly: true`, `mutating: false`, `productionReady: false`. No cluster informer or CR apply in KHR-AM.

### ShellLease read-only resolver

rdp-GW `ShellLeaseResolver` (stub) returns `ResolveShellResponse` with compatibility `GatewaySession`. Optional `poolId` query alias maps legacy pool only.

### Explicitly not in KHR-AM

- No session create via KHR endpoints
- No lease **revoke** or session **disconnect** mutating actions
- No rdp-GW production enable or autonomous orchestration

See rdp-GW: `docs/khr/RDPGW_KHR_ALIGNMENT_PLAN.md`, `docs/khr/SHELLLEASE_GATEWAYROUTE_COMPATIBILITY.md`.

---

## Identity / session / app semantics (KHR-AN)

Wave 2 adds **read-only correlation** across ShellLease, user session, and application session — no mutating auth or disconnect/revoke automation.

| Concept | Semantics |
|---------|-----------|
| **UserSessionRef** | Gateway/user session bound to `subjectId`, optional legacy `poolId`, ShellLease correlation string |
| **AppSessionRef** | RemoteApp / Windows app session via `WindowsAppRef` |
| **GatewaySessionIdentity** | Bundle: user + optional app + ShellLease + `GatewaySession` |
| **UserIdentityRef** | Directory/OIDC subject (directoryservice stub) |

### Correlation expectations

| Consumer | Endpoint / artifact | Behavior |
|----------|---------------------|----------|
| **rdp-GW** | `GET /karl-gw/v1/session/resolve?id=...` | `readOnly`, `noDisconnect`, `noRevoke` |
| **KARL-APP** | `AppSessionDescriptor`, `WindowsAppDescriptor` | AppShell compatibility; no session mutation |
| **karl-directoryservice** | `SessionIdentityCorrelation` | ShellLease ↔ user; OIDC-compatible fields |
| **Hyperdensity** | ShellLease `userRef` in CR spec | Contract source; no controller change in KHR-AN |

### Compatibility identity flow

```
OIDC login (unchanged) → gateway session (legacy/rdp-GW)
       → poolId (legacy) ──compat──► ShellPoolRef / ShellLeaseRef
       → RemoteApp fields ──compat──► WindowsAppRef / AppSessionRef
       → session/resolve (read-only) ──► UserSessionRef + AppSessionRef correlation
```

**poolId** remains legacy-only. **GatewayRoute** + **ShellLease** are target models for session-path and access reservation.

### Explicitly not in KHR-AN

- No new auth enforcement in directoryservice
- No token mutation flows
- No automated revoke/disconnect
- No production enable claims

---

## Non-goals

- No token minting, no rdp-GW apply, no session broker controller
- No mutating revoke/disconnect in KHR-AM (read-only wave only)
