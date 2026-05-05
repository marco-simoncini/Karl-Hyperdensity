# Technical Preview Demo Guide v1

Demo ID: `hyperdensity_technical_preview_demo_guide_v1`  
Target runtime: `technical_preview`

## 1) Demo goal

Show Grande Padre as a governed live resource market, not a generic dashboard and not a simplistic autoscaler.

## 2) Suggested duration

10-15 minutes.

## 3) Audience

- technical founder
- CTO
- platform engineer
- investor technical diligence
- internal operator

## 4) Demo prerequisites

- Parent Fabric deployed and reachable.
- Live authenticated GET access available.
- Evidence namespace available.
- Artifact bundle available for fallback.
- No production mutation.
- No autonomous apply.
- No enforcement.

## 5) Demo flow

### A) Opening narrative (1 minute)

- Positioning: governed live resource exchange/control plane.
- Clarify: not a Kubernetes dashboard, not a simple autoscaler.
- Show top-level release boundary and TP mode.

### B) Creation path (2 minutes)

- Open Shell Claim Template/Profile Pack.
- Show supported profiles.
- State: no raw resource creation.
- Show Shell Claim generator and dry-run-before-create.

### C) Governance/remediation path (3 minutes)

- Show Admission Guard classification (`wouldReject` style posture).
- Show Mutate Preview output.
- Show Enforce Simulation output.
- Show Mutate Preview Apply Dry-Run result and no-mutation posture.

### D) Resource Exchange path (3 minutes)

- Show donor liquidity and receiver pressure.
- Show transfer plan recommendation.
- Show transfer dry-run.
- Show staged/chained apply proof + rollback/history evidence references.

### E) Live Resource Authority (2 minutes)

- Show single KARL authority surface for live CPU/RAM.
- Explain runtime drivers are internal adapters.
- Show Linux container CPU/RAM lane.
- Show Linux VM CPU/RAM evidence-backed lane.
- Show VM RAM runtime overlay wording and no raw runtime control exposure.

### F) Release boundary (2 minutes)

- Show Release Support Matrix.
- Show Evidence Bundle approved/rejected claims.
- Highlight Windows out-of-scope/frozen.
- Reiterate no production mutation, no autonomous apply, no enforcement.

## 6) Presenter notes

- "Hyperdensity does not wait to scale; it maintains a prevalidated resource market."
- "The product is not Pod Resize or QMP; the product is KARL Live Resource Authority."
- "No raw resource creation. Only Hyperdensity-ready shell creation."
- "Simulation and dry-run are not enforcement."
- "Technical Preview is evidence-backed, not universal GA."

## 7) Fallbacks

- If live API fails:
  - use saved artifacts
  - show evidence bundle extracts
  - state artifact freshness explicitly
  - do not claim live validation
- If Resource Exchange proof is stale:
  - present it as historical proof only
  - do not claim fresh live proof
- If VM support is questioned:
  - distinguish evidence-backed object-specific VM path from generic VM support claims
