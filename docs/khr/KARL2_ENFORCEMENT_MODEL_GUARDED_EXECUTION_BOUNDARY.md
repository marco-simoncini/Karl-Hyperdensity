# Enforcement model guarded execution boundary (KHR-EC)

| Field | Value |
|-------|-------|
| **Sprint** | KHR-EC |
| **State** | `docs/khr/enforcement-model-state.json` |
| **Evidence** | Karl-Installer `docs/evidence/karl2-enforcement-model-guarded-execution/committed-khr-ec-v1/` |

---

## Execution posture

| Field | Value |
|-------|-------|
| `enforcementModelAuthorized` | true |
| `enforcementEnabled` | true |
| `enforcementMode` | deny-only |
| `allowedScope` | single-reference-cluster |

Read-only projection and accessgraph/session surfaces unchanged. Forbidden actions remain denied per KHR-EB matrix.

---

## Rollback

Restore `enforcementEnabled=false` by removing enforcement posture markers from dashboard deployment.
