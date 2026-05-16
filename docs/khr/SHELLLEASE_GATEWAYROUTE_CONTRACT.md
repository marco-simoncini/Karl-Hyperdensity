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

## Non-goals

- No token minting, no rdp-GW apply, no session broker controller
