# Shell continuity semantics (KHR-U)

| Field | Value |
|-------|-------|
| **Scope** | Native-live sandbox lane — extends KHR-T resource continuity |
| **Production** | **Not enabled** — preview observation only |
| **Automation** | **None** |

---

## Model

| Concept | Description |
|---------|-------------|
| `shellSessionId` | Stable ID derived from namespace + pod + pod UID (+ container) |
| `appSessionId` | Shell session + `/app` suffix (pause/container app proxy) |
| `userSessionId` | Shell session + `/user` (user session proxy) |
| `continuityState` | `preserved` \| `interrupted` \| `unknown` |
| `continuityEvidence` | Read-only bundle on verification/certification JSON |

Package: `pkg/khr/shellcontinuity`

---

## ResourceLease verification

`VerificationOutcome` (guarded apply) observes:

| Check | Field |
|-------|-------|
| Resource | `noRestart`, `noRollout`, `noRecreate` |
| Session | `sessionContinuityPreserved` |
| Shell | `shellContinuityPreserved` |
| App | `appContinuityPreserved` |
| Evidence | `continuityEvidence` |

When `ContinuityBefore` / `ContinuityAfter` snapshots are supplied to guarded apply, interruption fails verification.

---

## Native-live certification

`run-metrics.json` includes shell/session continuity proof.  
`certification-summary.json` adds:

| Field | Meaning |
|-------|---------|
| `continuityProof.resourceContinuityPreserved` | No restart/rollout/recreate |
| `continuityProof.shellContinuityPreserved` | Shell session IDs stable |
| `continuityProof.appContinuityPreserved` | App session IDs stable |
| `scores.resourceContinuityScore` | Resource-only score `0.0`–`1.0` |
| `scores.sessionContinuityScore` | Shell/app/session score |
| `scores.continuityScore` | `min(resource, session)` |

Regression guard fails on shell/app/session interruption.

---

## Evidence pipeline

```bash
export KHR_RUNTIME_CLUSTER_CONTEXT=karl-metal-01@ovh
./scripts/khr_native_live_certify.sh
```

Artifacts per run:

- `continuity-before.json` / `continuity-after.json`
- `continuity-proof.json` (shell + session continuity proof)

---

## Related

- `docs/khr/NATIVE_LIVE_CERTIFICATION.md` (KHR-T)
- `docs/khr/NATIVE_LIVE_LANE_PROTOTYPE.md` (KHR-S)
