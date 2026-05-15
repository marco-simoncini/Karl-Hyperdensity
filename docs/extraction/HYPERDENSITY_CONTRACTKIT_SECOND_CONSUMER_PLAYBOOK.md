# Hyperdensity contractkit — second consumer adoption playbook (Sprint 43)

## Purpose

Operational guide for introducing a **second** (or subsequent) **Go** consumer of the nested module  
`github.com/marco-simoncini/Karl-Hyperdensity/pkg/hyperdensity/contractkit` **without** modifying any repo other than the consumer itself and the canonical docs in **Karl-Hyperdensity** / **Karl-Dashboard** when policy is updated.

Nothing in this playbook changes HTTP APIs, Parent Fabric runtime behavior, JSON ordering, or execution/apply paths.

---

## 1. When to become a consumer

- **Only if** the repository **needs** shared Hyperdensity **IDs** (e.g. blocker catalog constants), **contracts** (manifest / release metadata), or **claimpolicy** vocabulary for parity, validation, or alignment with Dashboard/Hyperdensity extraction work.
- **Avoid** adding a `go.mod` dependency **solely for documentation** — prefer linking to Hyperdensity docs or duplicating minimal constants with an explicit local ADR if the cost of a module pin is unjustified.
- **Default for new consumers:** **test-only** imports (`*_test.go`, CI-only tools) until an ADR and sprint explicitly approve **runtime** use of a subpath.

---

## 2. Allowed subpackages

### `contractkit/blockers`

- **Runtime:** Possible **only** with a **dedicated sprint**, an explicit **ADR**, a **local allowlist** (or equivalent static audit) listing every production file that may import `blockers`, and review sign-off. **Karl-Dashboard** is the reference pattern (`contractkit_runtime_import_allowlist.txt` + audit scripts).
- **Test-only:** Always permitted for parity, golden tests, and compile-time guards.

### `contractkit/contracts`

- **Default:** **Test-only** (parity, manifest version guards, release metadata assertions).
- **Runtime:** **Only** after a **dedicated ADR + sprint** that updates Hyperdensity policy, consumer audits, and any parity matrices.

### `contractkit/claimpolicy`

- **Default:** **Test-only** (surface mapping, traceability, token guards).
- **Runtime:** **Forbidden** unless a **dedicated sprint** explicitly reopens the import surface (today: **no** such sprint for Dashboard; see `HYPERDENSITY_CONTRACTKIT_RUNTIME_IMPORT_FREEZE_M17.md` on Dashboard).

---

## 3. Mandatory consumer files

| Artifact | Action |
|----------|--------|
| **`audit_contractkit_module_pin.sh`** | Copy from **`docs/extraction/templates/audit_contractkit_module_pin.sh`** (Hyperdensity); set `EXPECTED_CONTRACTKIT_VERSION` / extend forbidden pins per consumer policy. |
| **Consumer ADR or local doc** | Use **`docs/extraction/templates/CONTRACTKIT_CONSUMER_DECISION_RECORD.md`** (filled in the consumer repo). |
| **CI env** | **`GOPRIVATE=github.com/marco-simoncini/*`**, **`GONOSUMDB=github.com/marco-simoncini/*`**; optional **`GONOPROXY=…`** — see **`HYPERDENSITY_CONTRACTKIT_CONSUMER_CI_HARDENING.md`**. |
| **`go mod verify`** | Run in CI **where compatible** (after `go.sum` is complete). |
| **Pin + import audits** | Adapt from Karl-Dashboard `kubernetes-console/scripts/hyperdensity/` (pin audit, optional runtime-import audit for `blockers`-only). |

---

## 4. Required review checklist

Use this list in PR description or release gate:

- [ ] **`go.mod`** uses **exact** current stable pin (**`v0.1.9-khr-m1-m19`** or whatever **`ContractKitModuleVersion`** is after the next Hyperdensity sprint), **no** pseudo-version for `contractkit`.
- [ ] **No superseded** pins (`v0.1.5-khr-m1-m16`, `v0.1.7-khr-m1-m18`) and no undeclared downgrade to older non-current pins.
- [ ] **No undeclared runtime import** of `contracts` / `claimpolicy` (or of `blockers` without allowlist + sprint).
- [ ] **Test-only** default respected unless ADR says otherwise.
- [ ] **Parity / audit scripts** green in CI (pin audit, `go mod verify` if used, tests, runtime import audit if applicable).
- [ ] **`HYPERDENSITY_CONTRACTKIT_CONSUMER_AUDIT.md`** updated on Hyperdensity when a new consumer ships (table row + date).

---

## 5. Example adoption paths (conceptual only)

These repos are **not** modified by this sprint; paths illustrate **how** a future consumer might justify adoption.

| Repo | Suggested path |
|------|----------------|
| **Karl-Inventory** | Possible **test-only** consumer later for **host/resource evidence** semantics aligned with Hyperdensity claim vocabulary; start with `claimpolicy` or `contracts` in tests only. |
| **Karl-Warden** | Possible **policy/guard** alignment; **test-only first** (`contracts` / `claimpolicy` in CI), runtime only after explicit sprint + ADR. |
| **FluidVirt** | Possible **provider / R&D** consumer; **test-only first** to avoid coupling KubeVirt runtime to Dashboard parity rules. |
| **Karl-Installer** / **Karl-OS-ISO** | **Avoid** becoming a Go module consumer **unless** native Go code truly needs `contractkit`; installer/ISO flows often do not justify a nested module pin — prefer docs links or generated artifacts. |

---

## Related

- `HYPERDENSITY_CONTRACTKIT_CONSUMER_POLICY.md`
- `HYPERDENSITY_CONTRACTKIT_CONSUMER_AUDIT.md`
- `HYPERDENSITY_CONTRACTKIT_CONSUMER_CI_HARDENING.md`
- `templates/CONTRACTKIT_CONSUMER_DECISION_RECORD.md`
- `templates/audit_contractkit_module_pin.sh`
- `HYPERDENSITY_PARENT_FABRIC_EXTRACTION_BOUNDARY.md`
- `HYPERDENSITY_PARENT_FABRIC_EXTRACTION_PHASES.md`
- `HYPERDENSITY_PARENT_FABRIC_DEPENDENCY_GUARDS.md`
- Dashboard `docs/hyperdensity/HYPERDENSITY_CONTRACTKIT_REFERENCE_CONSUMER_CHECKLIST_M26.md`
