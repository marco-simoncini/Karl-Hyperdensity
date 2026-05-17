# Inventory live ingest — Hyperdensity consumer expectations (KHR-BW)

Karl-Inventory read-only ingest of KHR evidence bundles. **No Hyperdensity runtime or CRD changes** in KHR-BW.

| Field | Value |
|-------|-------|
| **Sprint** | KHR-BW |
| **Beta blocker** | `inventory-live-ingest` — **dashboard-visible** when projection passes (KHR-BX) |
| **Inventory evidence** | `docs/evidence/khr-inventory-live-ingest/committed-khr-bw-v1/` |

---

## Expected observation fields

| Field | Source | Hyperdensity / Dashboard use |
|-------|--------|------------------------------|
| `observationId` | `khr-inventory-observation/<runId>` | Correlation id |
| `observationSource` | `live-readonly` \| `fixture-readonly` | Same as `inventoryObservationSource` |
| `observedAt` | RFC3339 timestamp | Evidence freshness |
| `snapshotRef` | `{ runId, summaryPath }` | Anchor to Reference Snapshot v1 |
| `inventoryObservationSource` | `live-readonly` \| `fixture-readonly` | Federation trust `inventory-observed` |
| `postureObserved` | Post-install + stub posture | Runtime posture federation |
| `scopeReadinessObserved` | Snapshot `scopeReadiness` | TP readiness correlation |
| `continuityObserved` | Federation + governance | Continuity bundles |
| `enforcement` | always `false` | No enforcement proof |
| `mutationObserved` | always `false` | No mutation proof |

### Dashboard projection (KHR-BX)

Karl-Dashboard `GET /api/hyperdensity/khr-backend/projection` exposes read-only `inventoryObservationSummary`:

| Field | Required |
|-------|----------|
| `available` | bool |
| `source` | observation source |
| `postureObserved` | bool |
| `scopeReadinessObserved` | bool |
| `continuityObserved` | bool |
| `enforcement` | `false` |
| `mutationObserved` | `false` |

Fixture: `Karl-Dashboard/examples/khr-dashboard/inventory-live-ingest-summary.json`. **No Cockpit UI** in KHR-BX.

---

## Ingest inputs (read-only)

| Artifact | Hyperdensity path |
|----------|-------------------|
| Reference snapshot v1 | `docs/evidence/khr-tp-reference-snapshot-v1/committed-khr-bt-v1/snapshot-summary.json` |
| Federation | `docs/evidence/khr-runtime-observation-federation/*/federation-summary.json` |
| Post-install bundle | `docs/evidence/khr-tp-post-install-bundle/summary.json` |

---

## Explicit invariants (no enforcement)

| Invariant | Required value |
|-----------|----------------|
| `readOnly` | `true` |
| `mutating` | `false` |
| `enforcement` | `false` |
| `applyObserved` | `false` |
| `productionReady` | `false` |
| `autonomousOrchestration` | `false` |

---

## Beta readiness impact

| Before KHR-BW | After KHR-BW |
|---------------|--------------|
| Export `stub` only | Read-only ingest from committed evidence |
| Live ingest blocker | **Partially reduced** — file ingest, not cluster agent |
| Federation | May cite `inventory-observed` with ingest evidence |

---

## Related

- Karl-Inventory: `INVENTORY_LIVE_INGEST_PLAN.md`
- `RUNTIME_OBSERVATION_FEDERATION.md`
- `KHR_BETA_READINESS_PLAN.md`
