# Grande Padre — action slate model (Sprint 13)

## ActionSlate

| Field | Role |
|-------|------|
| `generatedAt` | RFC3339 timestamp (injectable for tests via `-recommendation-generated-at`). |
| `source` | Always `local-evidence-store` for this skeleton. |
| `recommendations` | Ordered `ActionRecommendation` list. |
| `blocked` | Evidence rows classified as blocked for Grande Padre (subset of indexed rows). |
| `remediable` | `BlockedRemediableIndex` groups per cell. |
| `donorCandidates` | Cell refs eligible as donors (high confidence, ready, non-DevOnly, non-IntegrityFailed). |
| `receiverCandidates` | Cells needing remediation attention. |
| `summary` | Counts by recommendation, donor/receiver cardinality, and coarse `byRisk` / `byActionType` maps. |

## ActionRecommendation

Each row includes `dryRunOnly` (from `-recommend-dry-run-only`, default **true**) and **`applyAllowed: false`** — a hardwired reminder that Sprint 13 never authorizes apply.

### actionType semantics (skeleton)

| Type | Meaning |
|------|---------|
| `observe` | Continue read-only monitoring / evidence collection posture. |
| `remediate` | Address `blockedReasons` or integrity failures locally (operator workflow). |
| `prepare-resourcelease` | Suggests using existing KHR **dry-run** lease preparation; still not apply. |
| `collect-more-evidence` | Emitted when confidence is **low** to avoid over-confident market signals. |

## Safety

The slate must never be mistaken for admission: digest/signature trust tiers flow from Sprint 12 indexing and remain **integrity**, not **authorization**.
