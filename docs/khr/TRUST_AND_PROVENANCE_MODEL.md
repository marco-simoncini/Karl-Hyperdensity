# Trust and provenance model (KHR-Y)

| Field | Value |
|-------|-------|
| **Sprint** | KHR-Y |
| **Mode** | Read-only trust semantics |
| **Production** | **Not** zero-trust production enforcement |

---

## Provenance record

Every evidence artifact, certification registry row, approval, and control-graph export may carry:

| Field | Description |
|-------|-------------|
| `provenanceId` | Stable provenance record id |
| `generatedBy` | Tool/component (e.g. `khr-cert-registry`) |
| `generatedAt` | RFC3339 timestamp |
| `sourceCluster` | Cluster context |
| `sourceNamespace` | Sandbox namespace |
| `sourceLane` | Lane id (e.g. `native-live`) |
| `lineageHash` | Lineage anchor hash |
| `evidenceFingerprint` | `sha256:` content fingerprint |

---

## Verification

| Check | Package |
|-------|---------|
| Certification registry fingerprint | `certregistry.VerifyIntegrity` |
| Approval vs certification provenance | `provenance.VerifyApprovalProvenance` |
| Control graph lineage | `controlgraph.VerifyLineageIntegrity` |
| Stale provenance | `provenance.IsStaleProvenance` |

---

## CLI

```bash
go run ./cmd/khr-provenance-validate \
  -cert=docs/evidence/khr-native-live-lane/certification-summary.json \
  -registry=docs/evidence/khr-certification-registry/registry.json \
  -approval=pending.json \
  -graph=control-graph.json \
  -out=provenance-validation-summary.json
```

Evidence: `./scripts/khr_provenance_evidence.sh`

---

## States

| State | Meaning |
|-------|---------|
| `trusted` | Fingerprints and lineage align |
| `mismatch` | Fingerprint or lineage mismatch |
| `stale` | Provenance older than freshness window |
