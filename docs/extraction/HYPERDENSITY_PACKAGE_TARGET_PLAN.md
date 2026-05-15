# Hyperdensity package target plan (Sprint 16)

Maps **Dashboard runtime** (authoritative today) to **future Go packages** under `Karl-Hyperdensity`. Names are targets — **no new packages are created in Sprint 16**.

---

## Proposed packages

### `pkg/hyperdensity/contracts`

- **Enters:** Snapshot DTOs shared by API/UI (parent pool, overlap, execution/decision/arbitration summaries); request/response structs for parent-fabric GET/POST; executive cockpit envelope; versioned JSON tags aligned with `schemas/`.
- **Does not enter:** HTTP handlers, auth, Kubernetes clients, ConfigMap persistence.
- **Dashboard deps to remove later:** Duplicated TypeScript interfaces re-defined ad-hoc in Go — replace with generated or hand-maintained shared structs consumed by Dashboard via submodule or `go generate` export (TBD).

### `pkg/hyperdensity/blockers`

- **Enters:** Canonical blocker IDs and severities (Linux + Windows exclusion gates such as `no_windows_lane`, lease violation codes); lookup tables; “gate contribution” structs distilled from collectors.
- **Does not enter:** Live cluster probes.
- **Order:** **Early** — stabilizes strings used across Dashboard tests and Hyperdensity CRD docs.

### `pkg/hyperdensity/evidence`

- **Enters:** Read-only evidence channel contracts (readonly JSON channel, kubernetes API readonly), quality ranks (`live_parent_fabric_observed`), merge rules for bootstrap vs live.
- **Does not enter:** Side effects, admission servers.

### `pkg/hyperdensity/equilibrium`

- **Enters:** Donor/receiver math, overlap summaries, resource movement **pure functions** (no IO).
- **Does not enter:** Market tick execution, auto-apply.

### `pkg/hyperdensity/parentfabric`

- **Enters:** Merge pipeline **pure** layer: inventory + policy seed + freeze/deny reasons → snapshot fragment.
- **Does not enter:** HTTP, OIDC, stores.

### `pkg/hyperdensity/execution`

- **Enters:** Validation of execution POST bodies, idempotency keys, dry-run expansion, “would mutate” classification **without** performing mutations.
- **Does not enter:** Actual PATCH/apply (stays Dashboard or future controller).

### `pkg/hyperdensity/cohorts`

- **Enters:** VM Linux CPU/memory narrow cohort builders, guest-assisted **read models** (inputs → structured cohort output).
- **Does not enter:** KubeVirt API clients (initially).

### `pkg/hyperdensity/history`

- **Enters:** Append-only event types, serialization for usage history / ledger (interfaces first).
- **Does not enter:** etcd/SQL backing — Dashboard keeps storage until cutover.

### `pkg/hyperdensity/claimpolicy`

- **Enters:** Auto-scope policy document schema, effective policy resolution, deny reason strings.
- **Does not enter:** ConfigMap write path.

### `pkg/hyperdensity/windowslane`

- **Enters:** Windows exclusion constants, cross-links to `pkg/windowsfluidvirt` contracts, “planning-only until promoted” markers.
- **Does not enter:** Production mutation paths; **reuse** `pkg/windowsfluidvirt` for deep FluidVirt logic to avoid a second implementation.

---

## Extraction order (recommended)

1. **`contracts` + `blockers`** — smallest blast radius; unlock golden JSON tests.
2. **`evidence` + `parentfabric` (pure merge)`** — decouple UI polling from merge math.
3. **`cohorts` (VM Linux)`** — large but self-contained files in Dashboard.
4. **`equilibrium` + `claimpolicy`** — depends on stable blocker IDs.
5. **`execution` (validate-only)** — ties POST semantics.
6. **`history`** — needs storage decision.
7. **`windowslane` facade** — wraps existing `windowsfluidvirt`; align naming with Dashboard `hyperdensity_parent_fabric_windows_*`.

---

## Explicit non-goals

- Do **not** move Next.js components into Hyperdensity.
- Do **not** collapse Grande Padre into parent-fabric runtime — different trust boundary (local evidence vs cluster live).
