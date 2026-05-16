# Access Graph continuity evidence (KHR-AS / KHR-AT)

| Field | Value |
|-------|-------|
| **Sprint** | KHR-AS (fixture), KHR-AT (live-readonly) |
| **Contract set** | `khr-tp-contract-v1` |
| **Cluster** | `karl-metal-01@ovh` |

---

## Expected evidence bundle

Produced by **rdp-GW** `scripts/khr_accessgraph_continuity_evidence.sh`:

| Artifact | Description |
|----------|-------------|
| `accessgraph-session.json` | Full access graph + continuity export |
| `endpoint-probe.json` | Base URL, healthz/accessgraph HTTP status (live runs) |
| `validation.json` | Automated field/kind checks |
| `summary.json` | Run metadata + pass/fail + trust level |
| `run.log` | Execution log |

Path: `rdp-GW/docs/evidence/khr-accessgraph-continuity/<runId>/`

---

## Evidence trust levels (KHR-AT)

| `source` | `trustLevel` | Meaning |
|----------|--------------|---------|
| `live-readonly` | `live-readonly` | HTTP GET against sandbox/test rdp-GW; endpoint flags verified |
| `fixture-readonly` | `fixture` | Golden copy when gateway unreachable or `KHR_ACCESSGRAPH_EVIDENCE_FIXTURE_ONLY=true` |

Both levels may **PASS** bundle check. **Live-readonly is preferred** for TP continuity proof; fixture is acceptable for CI/offline.

Bundle check selects the **highest-trust** PASS summary under `docs/evidence/khr-accessgraph-continuity/` (live-readonly over fixture-readonly).

---

## Required summary fields

| Field | Value |
|-------|-------|
| `status` | `PASS` |
| `source` | `live-readonly` or `fixture-readonly` |
| `trustLevel` | `live-readonly` or `fixture` (must align with `source`) |
| `readOnly` | `true` |
| `mutating` | `false` |
| `noDisconnect` / `noRevoke` | `true` |
| `continuityObserved` | `true` |
| `noSessionMutation` | `true` |
| `productionReady` | `false` |

Live runs also include `baseUrl`, `healthzHttpStatus`, `accessGraphHttpStatus`.

---

## Live evidence (KHR-AT)

Sandbox rdp-GW (non-production):

- Runbook: `rdp-GW/docs/khr/RDPGW_SANDBOX_LIVE_EVIDENCE.md`
- Manifests: `rdp-GW/examples/khr/rdpgw-sandbox/` (`khr.karl.io/sandbox=true`)

```bash
export RDP_GW_BASE_URL=http://127.0.0.1:9443   # or port-forward URL
export KHR_RUNTIME_CLUSTER_CONTEXT=karl-metal-01@ovh
export KHR_ACCESSGRAPH_EVIDENCE_RUN_ID=live-sandbox-khr-at
./scripts/khr_accessgraph_continuity_evidence.sh
```

Require live (fail if only fixture): `KHR_REQUIRE_LIVE_EVIDENCE=true`.

---

## Validation criteria

End-to-end path evidenced:

`UserIdentityRef` → `ShellLeaseRef` → `GatewayRouteRef` → `GatewaySessionRef` → `AppSessionRef` / `WindowsAppRef`

Plus continuity: `continuityLineageId`, `sessionCorrelationId`, `continuityObserved`.

---

## TP limitations

| Limitation | Detail |
|------------|--------|
| **Read-only** | No revoke, disconnect, auth enforcement, or session mutation |
| **Fixture fallback** | `source=fixture-readonly` when rdp-GW unreachable |
| **Live optional in CI** | Fixture test script always runs; live preferred in sandbox |
| **No production** | `productionReady: false` always |
| **No CRD apply** | No Hyperdensity runtime changes in KHR-AT |

---

## Bundle check

```bash
./scripts/khr_access_graph_continuity_bundle_check.sh
```

Accepts sibling `../rdp-GW` evidence or committed relay under Hyperdensity `docs/evidence/khr-accessgraph-continuity-relay/`.

---

## Related

- `CONTINUITY_LINEAGE_CONTRACT.md`
- `ACCESS_GRAPH_CONTRACT.md`
- rdp-GW `RDPGW_SANDBOX_LIVE_EVIDENCE.md`
- Dashboard `DASHBOARD_CONTINUITY_PROJECTION.md`
