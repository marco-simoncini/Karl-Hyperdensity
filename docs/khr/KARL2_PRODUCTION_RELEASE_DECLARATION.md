# KARL 2.0 production release declaration (KHR-EG)

| Field | Value |
|-------|-------|
| **Sprint** | KHR-EG |
| **Evidence** | Karl-Installer `docs/evidence/karl2-production-release-declaration/committed-khr-eg-v1/` |

---

## Release declaration

| Field | Value |
|-------|-------|
| `karl2ProductionRelease` | `true` |
| `productionReady` | `true` |
| `promotionAllowed` | `true` |
| `releaseScope` | `current-connected-reference-production` |

State: `docs/khr/production-release-state.json`. Enforcement, fleet, and readiness states unchanged.

---

## Scope limitations

- Single connected reference cluster (`karl-metal-01@ovh`) only
- Not multi-cluster GA
- Karl-App excluded from production readiness
