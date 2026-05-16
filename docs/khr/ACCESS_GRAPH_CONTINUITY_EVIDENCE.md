# Access Graph continuity evidence (KHR-AS)

| Field | Value |
|-------|-------|
| **Sprint** | KHR-AS |
| **Contract set** | `khr-tp-contract-v1` |
| **Cluster** | `karl-metal-01@ovh` |

---

## Expected evidence bundle

Produced by **rdp-GW** `scripts/khr_accessgraph_continuity_evidence.sh`:

| Artifact | Description |
|----------|-------------|
| `accessgraph-session.json` | Full access graph + continuity export |
| `validation.json` | Automated field/kind checks |
| `summary.json` | Run metadata + pass/fail |
| `run.log` | Execution log |

Path: `rdp-GW/docs/evidence/khr-accessgraph-continuity/<runId>/`

---

## Required summary fields

| Field | Value |
|-------|-------|
| `status` | `PASS` |
| `source` | `live` or `fixture-readonly` |
| `readOnly` | `true` |
| `mutating` | `false` |
| `noDisconnect` / `noRevoke` | `true` |
| `continuityObserved` | `true` |
| `noSessionMutation` | `true` |
| `productionReady` | `false` |

---

## Validation criteria

End-to-end path evidenced:

`UserIdentityRef` → `ShellLeaseRef` → `GatewayRouteRef` → `GatewaySessionRef` → `AppSessionRef` / `WindowsAppRef`

Plus continuity: `continuityLineageId`, `sessionCorrelationId`, `continuityObserved`.

---

## TP limitations

| Limitation | Detail |
|------------|--------|
| **Read-only** | Evidence does not prove live cluster routing |
| **Fixture fallback** | `source=fixture-readonly` when rdp-GW unreachable |
| **No revoke/disconnect** | Evidence validates absence of mutating flags |
| **No production** | `productionReady: false` always |
| **No CRD apply** | No runtime mutation in KHR-AS |

---

## Bundle check

```bash
./scripts/khr_access_graph_continuity_bundle_check.sh
```

Accepts sibling `../rdp-GW` evidence or committed fixture under Hyperdensity `docs/evidence/khr-accessgraph-continuity-relay/`.

---

## Related

- `CONTINUITY_LINEAGE_CONTRACT.md`
- `ACCESS_GRAPH_CONTRACT.md`
- rdp-GW `RDPGW_ACCESS_GRAPH_CONTINUITY_EVIDENCE.md`
