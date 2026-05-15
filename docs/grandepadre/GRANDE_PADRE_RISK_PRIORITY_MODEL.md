# Grande Padre — risk & priority model (Sprint 13)

## Risk (`low` \| `medium` \| `high` \| `blocked`)

| Risk | Typical trigger |
|------|-----------------|
| `blocked` | `IntegrityFailed` digest/manifest alignment. |
| `high` | Blocked / not-ready evidence needing remediation attention. |
| `medium` | DevOnly trust tier, or low-confidence collect-more-evidence path. |
| `low` | Ready + high confidence observe / prepare-resourcelease suggestions. |

Risk is **not** a blast-radius computation; it is a coarse triage label for the local slate.

## Priority (`low` \| `medium` \| `high`)

Priorities order recommendations within a single slate output. They do **not** schedule work, open change windows, or interact with SLO burners.

## Interaction with Sprint 12 trust tiers

Trust tiers (`Unsigned`, `DevOnly`, `IntegrityVerified`, `IntegrityFailed`, `Unknown`) feed risk selection. `Unknown` signed modes are reserved for future PKI integration and should be treated conservatively by operators even if risk is not `blocked`.

## Policy reminder

Neither risk nor priority grants apply rights. **KHR apply gate** (future) plus Hyperdensity controllers remain the only legitimate path to mutations.
