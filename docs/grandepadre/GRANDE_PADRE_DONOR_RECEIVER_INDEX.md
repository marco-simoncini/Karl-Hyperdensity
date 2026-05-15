# Grande Padre — donor / receiver index (Sprint 13)

## Donor candidates

Computed from evidence rows that are:

- `readyForGrandePadre == true`,
- `confidence == high` (case-insensitive),
- `trustTier` **not** `IntegrityFailed`,
- `trustTier` **not** `DevOnly` (local-dev integrity is **non-production** for donor eligibility),
- have a non-nil `cellRef`.

Duplicates by namespace/name are collapsed.

## Receiver / remediation candidates

Derived from **blocked-like** rows (`!ready`, `blockedReasons` non-empty, or `IntegrityFailed`). Cells without a ref are skipped for the candidate list but can still appear in aggregate remediable indexes.

## Hyperdensity mapping (intent only)

In a future live resource market, donors might supply headroom and receivers absorb remediation or deficit signals. **Sprint 13** only labels candidates from static evidence JSON/YAML — it does **not** match live supply, pricing, or policy engines.

## Limits

- No cross-tenant pairing, no feasibility checks against live telemetry beyond what is already embedded in the bundle.
- `DevOnly` bundles never become donors even if confidence is high.
