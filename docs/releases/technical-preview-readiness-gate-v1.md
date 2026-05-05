# Technical Preview Readiness Gate v1

Gate ID: `hyperdensity_technical_preview_readiness_gate_v1`  
Release track: `technical_preview`

Use this checklist as pass/fail gate for Technical Preview release posture.

## 1) Governance gates

- [ ] Policy Pack present.
- [ ] Policy Pack Consistency is `true`.
- [ ] Release Support Matrix present.
- [ ] Evidence Bundle present.
- [ ] Live Resource Authority present.

## 2) Safety gates

- [ ] Enforcement disabled.
- [ ] Autonomous apply `false`.
- [ ] Production mutation `false`.
- [ ] Evidence namespace only for proof paths.
- [ ] Dry-run required for relevant paths.
- [ ] Rollback required for supported apply paths.
- [ ] Cleanup required for temporary evidence objects.
- [ ] Warning-clean required before claiming success.
- [ ] No raw runtime controls exposed.
- [ ] Windows out-of-scope/frozen.

## 3) Creation path gates

- [ ] Profile Pack present.
- [ ] Supported profiles present.
- [ ] Shell Claim Generator aligned with Profile Pack.
- [ ] Dry-Run Create validated.
- [ ] Evidence Create/History path validated.
- [ ] Live create is not default.

## 4) Admission/remediation gates

- [ ] Admission Guard in audit mode.
- [ ] Classification present.
- [ ] Mutate Preview present.
- [ ] Enforce Simulation present.
- [ ] Apply Dry-Run present.
- [ ] No-mutation verified.
- [ ] Unsafe wording rejected.

## 5) Resource Exchange gates

- [ ] Donor liquidity present.
- [ ] Receiver demand present.
- [ ] Transfer plan present.
- [ ] Transfer dry-run present.
- [ ] Staged apply proof present.
- [ ] Chained apply proof present.
- [ ] Rollback proof present.
- [ ] Stage history present.

## 6) Live Resource Authority gates

- [ ] Runtime drivers present.
- [ ] No raw controls exposed.
- [ ] Container CPU/RAM driver present.
- [ ] VM CPU driver present.
- [ ] VM RAM virtio-mem/QMP driver present.
- [ ] Generic KubeVirt memory template mutation wording rejected.
- [ ] Runtime lease/overlay model present.
- [ ] Verify/rollback/reconcile semantics present.

## 7) Documentation gates

- [ ] Operator Runbook present.
- [ ] Release Notes present.
- [ ] Demo Guide present.
- [ ] Support Matrix present.
- [ ] Evidence Bundle present.

## 8) Test/build gates

- [ ] `go test ./pkg/server -count=1` (when Dashboard changed).
- [ ] `npm run typecheck` (when Dashboard changed).
- [ ] bridge guard (when Dashboard changed).
- [ ] `./scripts/validate.sh` in Karl-Hyperdensity.
- [ ] live `GET /api/hyperdensity/parent-fabric` is HTTP 200 when runtime/API changed.

## 9) Blockers

Any single item below yields `blocked`:

- enforcement enabled
- autonomous apply enabled
- production mutation enabled
- Windows marked supported
- unsafe wording approved
- generic VM RAM template mutation approved
- missing rollback requirement
- missing no-mutation proof
- missing support matrix
- policy consistency false
- raw runtime controls exposed
- failed tests
- failed validation

## 10) Final readiness states

- `ready_for_internal_evidence_release`
- `ready_for_founder_investor_demo`
- `ready_for_private_technical_preview`
- `not_ready`
- `blocked`

State interpretation:
- `ready_for_*` only if all mandatory categories pass.
- `not_ready` if incomplete but not hard-blocked.
- `blocked` if any blocker is present.
