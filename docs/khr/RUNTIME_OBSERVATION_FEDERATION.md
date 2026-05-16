# Runtime Observation Federation (KHR-AU)

| Field | Value |
|-------|-------|
| **Sprint** | KHR-AU |
| **Contract set** | `khr-tp-contract-v1` |
| **Cluster** | `karl-metal-01@ovh` |
| **Mode** | Read-only aggregation — **no orchestration** |

---

## Purpose

Federate continuity and access-graph **observations** from multiple repos into one read-only bundle. Supports Technical Preview evidence review without mutating runtime, sessions, or policy.

**Explicit non-goals:** revoke, disconnect, auth enforcement, production enable, CRD apply, autonomous orchestration.

---

## Core types

### ObservationSource

Identifies which component produced an observation (metadata only).

| `observationSource` | Repo | Typical `trustLevel` |
|---------------------|------|----------------------|
| `rdpgw-live` | rdp-GW | `live-readonly` |
| `rdpgw-fixture` | rdp-GW | `fixture-readonly` |
| `inventory-observed` | Karl-Inventory | `inventory-observed` |
| `dashboard-projected` | Karl-Dashboard | `projected-readonly` |
| `karl-app-projected` | KARL-APP | `projected-readonly` |
| `directoryservice-projected` | karl-directoryservice | `projected-readonly` |

### ObservationTrustLevel

| Value | Meaning |
|-------|---------|
| `fixture-readonly` | Golden / offline copy; no live attestation |
| `live-readonly` | Live HTTP read against sandbox/test gateway |
| `inventory-observed` | Host/session posture from Inventory export (observation only) |
| `projected-readonly` | Dashboard / APP / DS descriptor or fixture projection |

Trust levels are **never upgraded** during merge (fixture cannot become live).

### FederationCorrelation

Cross-source correlation keys (observational):

| Field | Role |
|-------|------|
| `federationCorrelationId` | Primary federation key (defaults to `sessionCorrelationId` when present) |
| `continuityLineageId` | Lineage anchor across gateway graph |
| `sessionCorrelationId` | Session binding across identity → gateway → app |

### FederationBundle

Aggregated read-only artifact:

| Field | Description |
|-------|-------------|
| `observationSources[]` | Per-repo observations with trust + paths |
| `federationCorrelation` | Merged correlation block |
| `trustLevels` | Map repo → trust level |
| `primaryTrustLevel` | Highest precedence trust present (see merge rules) |
| `consistency` | Minimal cross-source checks |

Output path: `docs/evidence/khr-runtime-observation-federation/<runId>/federation-summary.json`

---

## Merge rules

### 1. Lineage correlation

When multiple sources expose `continuityLineageId`, federation **requires** exact string match for `consistency.lineageCorrelationMatch=true`. Mismatch → `status=FAIL` (observation inconsistency, not enforcement).

### 2. Session correlation

Same for `sessionCorrelationId` / `federationCorrelationId` → `consistency.sessionCorrelationMatch`.

### 3. Continuity precedence

For `primaryTrustLevel` and human-readable precedence only (no runtime effect):

```
live-readonly > inventory-observed > projected-readonly > fixture-readonly
```

Within rdp-GW: `rdpgw-live` beats `rdpgw-fixture`.

Continuity flags (`continuityObserved`, `noRevoke`, `noDisconnect`, `mutating:false`) must remain consistent across sources; any `mutating:true` or `productionReady:true` fails federation.

---

## Source inputs

| Repo | Input artifact |
|------|----------------|
| rdp-GW | `docs/evidence/khr-accessgraph-continuity/<runId>/summary.json` + optional `accessgraph-session.json` |
| Karl-Inventory | `examples/khr/federation-observation-stub.json` |
| Karl-Dashboard | `examples/khr-dashboard/access-graph-continuity-summary.json` |
| KARL-APP | continuity descriptor / federation doc mapping |
| karl-directoryservice | `examples/khr/identity-federation-observation-stub.json` |

---

## Federation check script

```bash
./scripts/khr_runtime_observation_federation_check.sh
```

Environment:

| Variable | Default |
|----------|---------|
| `KHR_FEDERATION_RUN_ID` | UTC timestamp |
| `KHR_RDP_GW_PATH` | sibling `../rdp-GW` |
| `KHR_INVENTORY_PATH` | sibling `../Karl-Inventory` |
| `KHR_DASHBOARD_PATH` | sibling `../Karl-Dashboard` |
| `KHR_KARL_APP_PATH` | sibling `../KARL-APP` |
| `KHR_DIRECTORYSERVICE_PATH` | sibling `../karl-directoryservice` |

---

## TP invariants

| ID | Invariant |
|----|-----------|
| FED-01 | Federation is **read-only** |
| FED-02 | No revoke / disconnect / session mutation |
| FED-03 | No auth enforcement via federation |
| FED-04 | No production enable (`productionReady: false`) |
| FED-05 | No CRD or controller changes in KHR-AU |
| FED-06 | Trust levels exposed, never upgraded |

---

## Related

- `ACCESS_GRAPH_CONTINUITY_EVIDENCE.md`
- `CONTINUITY_LINEAGE_CONTRACT.md`
- rdp-GW `RDPGW_FEDERATION_OBSERVATION.md`
- Per-repo federation docs (Inventory, Dashboard, KARL-APP, directoryservice)
