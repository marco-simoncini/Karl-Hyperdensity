# Contractkit consumer — decision record (template)

**Copy this file into the consuming repository** (e.g. `docs/adr/CONTRACTKIT_CONSUMER_DECISION_RECORD.md`) and fill every section. Do not submodule this template.

---

## Repository

- **Git remote / name:** <!-- e.g. marco-simoncini/Karl-Inventory -->

## Branch

- **Target integration branch:** <!-- e.g. KHR -->

## Reason for consuming contractkit

- **Problem statement:** <!-- shared IDs, manifest parity, claimpolicy vocabulary, … -->
- **Alternatives considered:** <!-- docs-only link, vendored copy, generated code, … -->

## Subpackages used

- [ ] `github.com/.../pkg/hyperdensity/contractkit/blockers`
- [ ] `github.com/.../pkg/hyperdensity/contractkit/contracts`
- [ ] `github.com/.../pkg/hyperdensity/contractkit/claimpolicy`

## Runtime or test-only

- **Classification:** <!-- test-only | mixed (declare below) -->
- **Production import paths allowed:** <!-- e.g. only …/blockers in package X/Y; none -->

## CI environment

- **GOPRIVATE:** <!-- e.g. github.com/marco-simoncini/* -->
- **GONOSUMDB:** <!-- e.g. github.com/marco-simoncini/* -->
- **GONOPROXY (optional):** <!-- … -->

## Exact module pin

- **Pinned version:** <!-- e.g. v0.1.9-khr-m1-m19 -->
- **Git tag on Hyperdensity:** <!-- e.g. pkg/hyperdensity/contractkit/v0.1.9-khr-m1-m19 -->

## Audit scripts installed

- [ ] `audit_contractkit_module_pin.sh` (from Hyperdensity template)
- [ ] `go mod verify` in CI
- [ ] Runtime import audit / allowlist (if `blockers` in production)

## Owner

- **Primary owner / team:** <!-- … -->

## Rollback plan

- **How to revert:** <!-- remove require, revert go.sum, drop tests, … -->
- **Blast radius:** <!-- which binaries/jobs affected -->

## Review date

- **Date:** <!-- YYYY-MM-DD -->

## Approval

- **Reviewers:** <!-- names / handles -->
- **Sign-off:** <!-- link to PR or meeting notes -->
