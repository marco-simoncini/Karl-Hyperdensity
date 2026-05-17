# Scope-4 Operational Governance (KHR-BG)

| Field | Value |
|-------|-------|
| **Sprint** | KHR-BG |
| **Capability** | TP Scope-4 guarded apply (certified-evidence-backed) |
| **Cluster** | `karl-metal-01@ovh` |
| **Live mutation in BG** | **None** |

---

## Purpose

Consolidate Scope-4 as a **Technical Preview certified capability** with operational governance: evidence lifecycle, operator workflow, and certification lifecycle. **No runtime expansion** — governance reads committed evidence only.

---

## Governance states

| State | Meaning | Operator action |
|-------|---------|-----------------|
| `certified` | Certification + dependencies PASS; within validity window | Read-only consumption; no apply |
| `stale` | Evidence or upstream summaries aged / soft-stale signals | **Revalidation required** before trusting projection |
| `expired` | Certification past `certificationExpiry` | **Revalidation required**; do not cite as current |
| `revoked` | Hard block (provenance/registry integrity fail) | Do not apply; fix upstream evidence |
| `regression-detected` | Continuity or rollback evidence fails | **Revalidation required**; investigate regression |

Scope-4 remains **manual**, **not active**, **not enabled**, **not autonomous** in all states except explicit future sprint sign-off.

---

## Evidence retention

| Tier | Path | Retention policy |
|------|------|------------------|
| **Anchor** | `docs/evidence/khr-tp-live-scope4-guarded-apply/committed-scope4-guarded-apply-khr-be/` | Retain for TP contract freeze; do not delete without ADR |
| **Certification** | `docs/evidence/khr-scope4-guarded-apply-certification/committed-scope4-certification-khr-bf/` | Retain with apply anchor |
| **Governance bundle** | `docs/evidence/khr-scope4-governance/<runId>/` | Timestamped runs optional; committed anchor `committed-scope4-governance-khr-bg` |
| **Upstream** | provenance, federation, approval, control-graph summaries | Latest committed summaries referenced by bundle |

No automatic pruning in KHR-BG.

---

## Certification lifecycle

| Phase | Artifact | Trigger |
|-------|----------|---------|
| Preflight | `scope4-preflight-summary.json` (KHR-BD) | Before apply sprint |
| Apply + rollback | KHR-BE evidence | Manual operator sprint |
| Certification | `certification-summary.json` (KHR-BF) | `khr_scope4_certification_check.sh` |
| Governance | `governance-summary.json` (KHR-BG) | `khr_scope4_governance_bundle.sh` |

**Expiry:** `certificationExpiry` = certification `at` + **180 days** (TP reference window; not production SLA).

**Stale semantics:** certification older than **90 days** OR provenance/federation not `ok`/`PASS` → `staleCertification=true`.

---

## Operator revalidation flow

1. Confirm `KHR_TP_LIVE_SCOPE4_I_UNDERSTAND_GUARDED_APPLY` remains required for any **new** apply (not used in governance sprint).
2. Re-run `./scripts/khr_scope4_certification_check.sh` (read-only).
3. Re-run `./scripts/khr_scope4_governance_bundle.sh`.
4. If `operatorRevalidationRequired=true`, do **not** update Dashboard/Inventory projections until PASS.
5. **Never** enable production namespaces or systemd host-runtime.

---

## Rollback evidence requirements

Governance **requires** KHR-BE:

- `rollback-summary.json` with `rollbackVerified=true`
- `verify-summary.json` with `rollbackVerified=true`
- Apply evidence retained for audit

---

## Continuity regression handling

If `continuityPreserved=false` or `regressionDetected=true`:

- Set `scope4GovernanceState=regression-detected`
- Block operator trust of Scope-4 certification
- Re-execute manual apply sprint only with explicit sign-off (future sprint)

---

## Governance bundle (read-only)

```bash
./scripts/khr_scope4_governance_bundle.sh
```

Output: `docs/evidence/khr-scope4-governance/committed-scope4-governance-khr-bg/governance-summary.json`

---

## Related

- `SCOPE4_GUARDED_APPLY_CERTIFICATION.md`
- `SCOPE4_FAILURE_SEMANTICS.md`
- `KHR_TP_LIVE_ENABLEMENT_PLAN.md`
