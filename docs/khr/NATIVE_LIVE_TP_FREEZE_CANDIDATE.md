# Native-Live Lane TP Freeze Candidate (KHR-AF)

| Field | Value |
|-------|-------|
| **Sprint** | KHR-AF |
| **Lane** | `native-live` |
| **Classification** | `native-live` |
| **Provider** | `khr.native` |
| **Certification ID** | `khr-native-live-certification-v1` |
| **Cluster** | `karl-metal-01@ovh` / `khr-runtime-sandbox` |
| **Verdict** | **APPROVED — TP freeze candidate** (semantics + sandbox evidence) |

**NOT production ready.** **NOT GA.** `certified-preview` / `status: certified` is **sandbox evidence only**. No autonomous orchestration.

---

## Native-live definition (frozen)

A workload is **native-live** when all apply:

| Criterion | Requirement |
|-----------|-------------|
| Provider | `khr.native` (Linux cgroup — not KubeVirt) |
| Namespace | `khr-runtime-sandbox` |
| Labels | `khr.karl.io/sandbox=true`; optional `khr.karl.io/native-live=true` |
| Naming | Prefix `khr-native-live-` recommended |
| Exclusion | Not `virt-launcher` / not VM compatibility path |

---

## Live-in-place semantics

| Term | Meaning |
|------|---------|
| **Live-in-place** | CPU/RAM envelope change via cgroup without pod restart/rollout/recreate |
| **liveInPlaceEligible** | ResourceFuture simulation says scale path is eligible |
| **nativeLiveEligible** | Dashboard aggregate: any lane classified `native-live` without interruption |
| **applyMode** | `live-in-place` when native-live path; else `compatibility-fallback` |

---

## Invariants (frozen — must hold for certification PASS)

| Invariant | Expected value | Evidence field |
|-----------|----------------|----------------|
| `noRestart` | `true` | `certification-summary.json` → `invariants` |
| `noRollout` | `true` | same |
| `noRecreate` | `true` | same |
| `interruptionWindowMs` | `0` | `invariants` + `metrics` |
| `interruptionDetected` | `false` | `invariants` |
| `restartCountDelta` | `0` | `metrics` (per run) |
| `rolloutCount` | `0` | `metrics` |
| `recreateDetected` | `false` | `metrics` |

---

## continuityScore expectations

| Score | Field | TP expectation |
|-------|-------|----------------|
| Resource continuity | `resourceContinuityScore` | `1.0` for PASS |
| Session continuity | `sessionContinuityScore` | `1.0` for PASS |
| **Aggregate** | `continuityScore` | **`1.0`** for freeze anchor |
| Panel flags | `continuityProof.*Preserved` | all `true` |

---

## liveScaleConfidence

| Value | Meaning |
|-------|---------|
| `high` | Certification PASS on native-live lane (frozen anchor) |
| `medium` | Partial observation / simulation only |
| `low` | Compatibility or unknown |

Frozen anchor evidence: `liveScaleConfidence: "high"` in `certification-summary.json`.

---

## Certification requirements

| Requirement | Artifact |
|-------------|----------|
| Multi-run certify | `scripts/khr_native_live_certify.sh` |
| Summary anchor | `docs/evidence/khr-native-live-lane/certification-summary.json` |
| `readOnly: true` | Summary |
| `noAutonomousOrchestration: true` | Summary |
| Registry row | `certified-preview` in Dashboard projection (not GA) |
| Baseline match | `baselineMatch: true` |

---

## Evidence / provenance requirements

| Evidence | Path |
|----------|------|
| Native-live certification | `docs/evidence/khr-native-live-lane/certification-summary.json` |
| Certification runs | `docs/evidence/khr-native-live-lane/certification/*/` |
| Provenance program | `docs/evidence/khr-provenance/summary.json` |
| Certification registry | `docs/evidence/khr-certification-registry/summary.json` |
| ResourceFuture (native-live run) | `certification/*/run-*/run-metrics.json` → `liveInPlaceEligible: true` |
| ResourcePort CR preview | `certification/*/run-*/cr-preview/resourceport-*-native-live-*` |

Policy gates: simulation predicates only — see `CERTIFICATION_REGISTRY_AND_POLICY_GATES.md`.

---

## Supported lanes (TP)

| Lane | Classification | Live scale |
|------|----------------|------------|
| `native-live` | `native-live` | cgroup CPU/RAM in sandbox |
| `live-in-place-capable` | simulation | eligible when certified |

---

## Unsupported lanes / semantics

| Path | Reason |
|------|--------|
| `kubevirt.compatibility` as native-live | Wrong provider |
| Production namespaces | Blocklist |
| GA / production-ready certification | **Forbidden** |
| Autonomous certify loop | **Forbidden** |
| ISO default host-runtime enable | **Forbidden** |
| Windows native-live cert parity | **Gap** (P2) |
| Multus as target fabric | **Unsupported** |

---

## compatibility-fallback semantics

When workload is VM/KubeVirt or lacks native-live labels:

| Signal | Projection |
|--------|------------|
| `classification` | `compatibility-fallback` |
| `liveInPlaceEligible` | `false` |
| `compatibilityFallbackForecast` | `true` |
| `restartRisk` | may be `medium`/`high` |
| `interruptionDetected` | may be `true` on compatibility path |

Native-live freeze does **not** claim compatibility paths are live-in-place.

---

## Freeze vs experimental

| Frozen (TP) | Experimental |
|-------------|--------------|
| Lane ID, classification, invariants | RAM apply on non-Linux |
| continuityScore / liveScaleConfidence vocabulary | Scheduled re-certify automation |
| Certification summary JSON shape | Production enable |
| Dashboard projection fields | Controller reconciliation |

---

## Beta blockers

| ID | Severity | Blocker |
|----|----------|---------|
| NL-B-01 | P0 | NOT production ready |
| NL-B-02 | P0 | No autonomous orchestration |
| NL-B-03 | P1 | Windows native-live certification |
| NL-B-04 | P1 | Inventory native-live observation ingest |
| NL-B-05 | P2 | Multi-cluster cert federation |

---

## Explicit non-goals

- Production enablement or GA
- Autonomous orchestration / approval→apply
- systemd enable on ISO
- Dashboard rewrite
- Removing KubeVirt compatibility

Validation: `./scripts/khr_native_live_freeze_check.sh`

---

## Related

- `NATIVE_LIVE_LANE_PROTOTYPE.md`, `NATIVE_LIVE_CERTIFICATION.md`, `SHELL_CONTINUITY_SEMANTICS.md`
- `RESOURCEPORT_TP_FREEZE_CANDIDATE.md`
- Dashboard `NATIVE_LIVE_PROJECTION_FREEZE_CANDIDATE.md`
