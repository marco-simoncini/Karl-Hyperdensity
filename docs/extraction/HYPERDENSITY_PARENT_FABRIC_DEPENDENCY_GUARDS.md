# Hyperdensity Parent Fabric — dependency guard requirements (Sprint 44)

Future **`pkg/hyperdensity/parentfabric/**`** pure packages must remain **extractable** and **testable** without the OpenShift console runtime. This document lists **intended** static guard rules (to be enforced in a later sprint via `go vet` build tags, `grep` CI, or `go list -deps` audits).

**Sprint 44** does **not** add CI enforcement here — requirements only.

---

## Parentfabric pure packages — default deny list

Pure trees (`parentfabric`, `summary`, `governance`, `evidence`, `recommendation` — see **`HYPERDENSITY_PARENT_FABRIC_EXTRACTION_BOUNDARY.md`**) **must not** import:

| Pattern | Rationale |
|---------|-----------|
| **`k8s.io/*`** | Kubernetes client / API machinery belongs in Dashboard adapters or future operator code, not pure-core. |
| **`kubevirt.io/*`** | KubeVirt provider coupling; **no removal** of KubeVirt from product — only avoid entangling **pure** packages. |
| **`github.com/gorilla/*`** | HTTP / websocket stack tied to console server. |
| **`net/http`** | Unless isolated in a clearly named **`adapter/http`** package behind an interface (explicit sprint). |
| **Dashboard repo paths** | e.g. `github.com/openshift/console/pkg/server` — pure packages must not call back into the console. |

---

## Dashboard consumption rules

- Dashboard may import Hyperdensity `parentfabric` pure packages **only** with an **allowlist + sprint** (same discipline as contractkit runtime expansion).
- **contractkit** `claimpolicy` / `contracts` **runtime freeze** stays **independent** — see Dashboard **`HYPERDENSITY_CONTRACTKIT_RUNTIME_IMPORT_FREEZE_M17.md`**; do not conflate with Parent Fabric extraction gates.

---

## Optional enforcement (future sprint)

| Guard | Mechanism sketch |
|-------|------------------|
| **Import deny** | `go list -deps` / custom script in Hyperdensity CI on `pkg/hyperdensity/parentfabric/...` |
| **HTTP leak** | `grep -R "net/http" pkg/hyperdensity/parentfabric` fail in CI |
| **K8s leak** | `grep -R "k8s.io/" pkg/hyperdensity/parentfabric` fail in CI |

---

## Related

- `HYPERDENSITY_PARENT_FABRIC_EXTRACTION_BOUNDARY.md`
- `HYPERDENSITY_PARENT_FABRIC_EXTRACTION_PHASES.md`
- Dashboard `docs/hyperdensity/HYPERDENSITY_PARENT_FABRIC_EXTRACTION_BOUNDARY_M28.md`
