# Certification Registry and Policy Gates (KHR-V)

| Field | Value |
|-------|-------|
| **Sprint** | KHR-V |
| **Mode** | Read-only registry + simulation gates |
| **Production** | **Not enabled** — no autonomous orchestration |

---

## Certification Registry

Read-only JSON registry (`khr-certification-registry-v1`) listing lanes that passed native-live certification.

| Field | Description |
|-------|-------------|
| `laneId` / `laneType` | Lane identifier (e.g. `native-live`) |
| `providerBinding` | Provider (e.g. `khr.native`) |
| `certificationState` | `certified-preview` \| `failed` |
| `continuityScore` | `0.0`–`1.0` |
| `liveScaleConfidence` | `high` \| `medium` \| `low` |
| `lastCertifiedAt` | RFC3339 timestamp |
| `evidenceRef` | Path/URI to certification evidence |
| `validForSeconds` | Freshness window |
| `attestation` | Gate inputs (noRestart, shell continuity, rollback observed, …) |

Generate from latest certification:

```bash
go run ./cmd/khr-cert-registry \
  -cert=docs/evidence/khr-native-live-lane/certification-summary.json \
  -evidence-ref=docs/evidence/khr-native-live-lane/certification/<id> \
  -out=docs/evidence/khr-certification-registry/registry.json
```

---

## Policy Gates

| Gate | Meaning |
|------|---------|
| `noRestart` | Certification attests no restart |
| `noRollout` | No rollout |
| `noRecreate` | No recreate |
| `noInterruption` | No interruption window |
| `shellContinuityRequired` | Shell continuity preserved |
| `rollbackRequired` | Rollback path observed in certification |
| `evidenceFreshnessRequired` | `lastCertifiedAt` within `validForSeconds` |

Package: `pkg/khr/policygates`.

---

## ResourceFuture integration

When `-cert-registry` is passed to `resourcefuture-simulate`, `liveInPlaceEligibility` includes:

| Field | Meaning |
|-------|---------|
| `eligibilityState` | `eligible` \| `blocked` |
| `blockedReason` | Human-readable block |
| `staleEvidence` | Certification expired |
| `uncertifiedLane` | Lane not in registry |

**Eligible:** certified `native-live` lane with fresh evidence and all gates passing.

**Blocked:** stale evidence, uncertified compatibility lanes, or failed gate attestation.

Simulation remains **read-only** (`noApply`, `noMutation`, `noAutonomousOrchestration`).

---

## Evidence script

```bash
./scripts/khr_cert_registry_policy_gates.sh
```

Produces registry, gated simulation, stale simulation, and uncertified-lane simulation under `docs/evidence/khr-certification-registry/`.
