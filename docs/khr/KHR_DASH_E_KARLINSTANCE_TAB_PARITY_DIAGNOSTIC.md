# KHR-DASH-E — Karl Instance tab parity diagnostic (Karl-Hyperdensity observation)

Karl-Hyperdensity remains **diagnostic only**. KHR-DASH-E moves Karl Instance tab parity onto KHR-native adapters and synthetic projections; Hyperdensity does not host the canonical VM surface and does not become a canonical tab surface.

## Observations
- `canonicalSurface = Karl Instances` (Hyperdensity is diagnostic governance metadata only).
- Cell / Shell / ResourcePort / SessionTarget remain governance projections, not enforced bindings.
- `resourceLeaseApply = false`.
- `resourcePortPersistentLoop = false`.
- Tab parity coverage for KHR-native Karl Instances is sourced from the Dashboard `khr.io/tab-parity-projection` annotation on the synthetic VirtualMachine — Hyperdensity may import it for diagnostics but never as an enforcement matrix.

## Forbidden
KubeVirt runtime use, CDI, virtctl, ResourceLease apply, ResourcePort persistent loop, secret commit, Windows/ISO commit, arbitrary shell execution — all remain `false`.
