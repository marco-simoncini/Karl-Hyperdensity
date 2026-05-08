# Integration Branch Strategy

## Strategy Decision
Do not merge `The-Father-Windows` directly.

## Branch Model
1. Freeze `The-Father-Windows` as **proof/evidence source branch**.
2. Create integration branch from mainline baseline:
   - recommended: `The-Father-FluidVirt-Windows-Integration` from `main` (`Karl-Hyperdensity`).
3. Port only approved files via cherry-pick/manual port per PR sequence.
4. Keep each PR scoped to one contract/capability block with explicit safety tests.

## Why This Strategy
- avoids large drift blast radius from 699-file Windows delta
- avoids stale dashboard import chain
- keeps Linux mainline critical path clean
- allows progressive gate reviews and rollback-friendly adoption

## Merge Rules
- no direct merge commit from Windows branch
- no artifact bulk import
- no runtime executable enablement without dedicated gate milestone
- no Dashboard UI sourcing from Windows branch
