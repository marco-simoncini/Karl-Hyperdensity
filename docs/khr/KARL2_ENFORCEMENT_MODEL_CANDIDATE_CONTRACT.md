# Enforcement model candidate contract (KHR-EB)

| Field | Value |
|-------|-------|
| **Sprint** | KHR-EB |
| **Evidence** | Karl-Installer `docs/evidence/karl2-enforcement-model-authorization-gate/committed-khr-eb-v1/` |

---

## Candidate enforcement model

| Field | Value |
|-------|-------|
| `enforcementMode` | `candidate-dry-run` |
| `allowedScope` | `single-reference-cluster` |
| `enforcementEnabled` | false |
| `policyMutation` | false |

---

## Scope boundary

Hyperdensity enforcement apply paths remain disabled. Dashboard KHR backend projection and rdp-GW accessgraph/session remain read-only. Policy matrix dry-run validates deny/allow decisions without cluster mutation.

---

## Authorization

`KARL2_ENFORCEMENT_MODEL_I_UNDERSTAND` guard required for any future enforcement execution sprint. `guardSatisfied=false` and `enforcementModelAuthorized=false` in KHR-EB.
