# KHR Technical Preview — Operator Runbook (KHR-AC)

| Field | Value |
|-------|-------|
| **Sprint** | KHR-AC |
| **Audience** | Operators / platform engineers |
| **Production** | **NOT production ready** |

Single operator runbook for consuming the KHR Technical Preview package. All mutation paths are **sandbox/manual only** with explicit operator invocation. **No autonomous orchestration.** **No production enable.** **No systemd enable** on ISO by default.

---

## 1. Prerequisites

| Item | Requirement |
|------|-------------|
| Cluster | `karl-metal-01@ovh` (reference) |
| Context | `export KHR_RUNTIME_CLUSTER_CONTEXT=karl-metal-01@ovh` |
| Namespace | `khr-runtime-sandbox` |
| Labels | `khr.karl.io/sandbox=true` on namespace and workloads |
| Hyperdensity | KHR branch built; scripts executable |
| Dashboard (optional) | `HYPERDENSITY_KHR_TP_READINESS_ENABLED=true` for `GET /api/hyperdensity/tp-readiness` |
| ISO | CRDs installed; host-runtime **disabled** |
| Bootstrap contract | `KHR_BOOTSTRAP_CONSUMER_EXPECTATIONS.md` (KHR-AG) |

---

## 1b. Bootstrap verify (read-only, KHR-AG)

After ISO provision, run **from Karl-OS-ISO** (filesystem checks only — no apply, no enable):

```bash
cd Karl-OS-ISO
./scripts/khr_iso_tp_verify.sh
./scripts/guard_khr_iso_boundaries.sh
```

Wave-1 convergence map: ISO CRD foundation → optional Installer profile (`karl2-khr-technical-preview` / `hybrid-transition`) → Hyperdensity sandbox evidence. See `KHR_BOOTSTRAP_CONSUMER_EXPECTATIONS.md`, `Karl-OS-ISO/docs/khr/KHR_BOOTSTRAP_FLOW.md`, `Karl-Installer/docs/khr/KHR_INSTALLER_PROFILE_MATRIX.md`.

---

## 2. Preflight (read-only)

```bash
cd Karl-Hyperdensity
./scripts/guard_khr_docs_scope.sh
./scripts/khr_tp_package_check.sh
./scripts/khr_tp_operator_bundle.sh
export KHR_RUNTIME_NAMESPACE=khr-runtime-sandbox
./scripts/khr_runtime_sandbox_preflight.sh   # validates context/namespace/labels only
```

Do **not** set `sandboxApplyEnabled: true` unless running an approved sandbox evidence pipeline.

---

## 3. Evidence generation (sandbox)

Run from Hyperdensity repo root. Each script is **operator-initiated**; none auto-apply in production.

| Step | Script | Output |
|------|--------|--------|
| Lane discovery | `./scripts/khr_lane_discovery_evidence.sh` | `docs/evidence/khr-lane-discovery/` |
| ResourceFuture | `./scripts/khr_resourcefuture_evidence.sh` | `docs/evidence/khr-resourcefuture/` |
| Native-live lane | `./scripts/khr_native_live_lane_evidence.sh` | `docs/evidence/khr-native-live-lane/` |
| Native-live certification | `./scripts/khr_native_live_certify.sh` | `docs/evidence/khr-native-live-lane/certification/` |
| Certification registry | `./scripts/khr_cert_registry_policy_gates.sh` | `docs/evidence/khr-certification-registry/` |
| Action approval | `./scripts/khr_action_approval_evidence.sh` | `docs/evidence/khr-action-approval/` |
| Control graph | `./scripts/khr_control_graph_evidence.sh` | `docs/evidence/khr-control-graph/` |
| Provenance | `./scripts/khr_provenance_evidence.sh` | `docs/evidence/khr-provenance/` |

Full sandbox pipeline (opt-in live):

```bash
KHR_RUNTIME_SANDBOX_LIVE=1 ./scripts/khr_runtime_sandbox_execute.sh
```

---

## 4. Native-live certification

Reference: `docs/khr/NATIVE_LIVE_CERTIFICATION.md`

```bash
export KHR_RUNTIME_CLUSTER_CONTEXT=karl-metal-01@ovh
./scripts/khr_native_live_certify.sh
```

Verify anchor: `docs/evidence/khr-native-live-lane/certification-summary.json`  
State is **`certified-preview`** — **not GA**.

---

## 5. Provenance validation

Reference: `docs/khr/TRUST_AND_PROVENANCE_MODEL.md`

```bash
./scripts/khr_provenance_evidence.sh
```

Verify: `docs/evidence/khr-provenance/summary.json` with `readOnly: true`, `noAutonomousOrchestration: true`.

---

## 6. Policy gates

Reference: `docs/khr/CERTIFICATION_REGISTRY_AND_POLICY_GATES.md`

```bash
./scripts/khr_cert_registry_policy_gates.sh
```

Gates are **simulation predicates** only — not production enforcement.

---

## 7. Action approval evidence

Reference: `docs/khr/OPERATOR_ACTION_APPROVAL_WORKFLOW.md`

```bash
./scripts/khr_action_approval_evidence.sh
```

Approval is **local evidence only** — **no apply** on approve.

---

## 8. Control graph export

Reference: `docs/khr/KHR_CONTROL_GRAPH.md`

```bash
./scripts/khr_control_graph_evidence.sh
```

Export: `docs/evidence/khr-control-graph/control-graph.json` (read-only).

---

## 9. Cleanup

| Action | Command |
|--------|---------|
| Sandbox rollback | `./scripts/khr_runtime_sandbox_rollback.sh` |
| ResourcePort preview cleanup | `karl-host-runtime -mode=resourceport-cleanup` (sandbox context only) |
| Evidence review | Retain `docs/evidence/khr-tp-operator-bundle/<runId>/` from bundle script |

Never delete production namespace resources via KHR scripts.

---

## 10. Failure handling

| Symptom | Action |
|---------|--------|
| Preflight namespace mismatch | Fix `KHR_RUNTIME_NAMESPACE` to `khr-runtime-sandbox` |
| Certification regression | Stop pipeline; inspect `regressionDetected` in summary JSON |
| Provenance mismatch | Re-run registry + provenance scripts; do not approve actions |
| Policy gate blocked | Expected for stale/uncertified lanes — read `blockedReason` |
| Dashboard TP API disabled | Set `HYPERDENSITY_KHR_TP_READINESS_ENABLED=true` or use Hyperdensity bundle script |
| ISO host-runtime enabled accidentally | `systemctl disable karl-host-runtime`; restore example config defaults |

---

## 11. Operator bundle script

```bash
./scripts/khr_tp_operator_bundle.sh
```

Produces read-only index under `docs/evidence/khr-tp-operator-bundle/<runId>/`:

- `bundle-index.json` — evidence file inventory
- `run-summary.json` — guard/check results
- `blocker-summary.json` — P0/P1 blockers
- `next-actions.json` — suggested manual steps

Does **not** invoke guarded apply or autonomous workflows.

---

## 12. Cross-repo consumption

| Repo | Consumption |
|------|-------------|
| Hyperdensity | This runbook + `TECHNICAL_PREVIEW_PACKAGE.md` |
| Dashboard | `GET /api/hyperdensity/tp-readiness` (flag on) |
| Inventory | `scripts/khr_tp_observation_export.sh` (stub export) |
| ISO | `docs/khr/TECHNICAL_PREVIEW_ISO_GUIDE.md` post-install |

---

## Explicit statement

**NOT production ready.** Sandbox/manual evidence only. No autonomous orchestration. ISO runtime disabled by default.
