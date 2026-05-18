# KHR-DASH-F — Karl Hyperdensity diagnostic observation (Full VM parity)

- `hyperdensityFullVmParityDiagnosticObservationGenerated=true`
- `hyperdensityPlacement=diagnosticOnly`
- `canonicalVmSurface=Karl Instances` (Dashboard)
- `cellShellResourcePortSessionTargetGovernance=metadata-only`
- `resourceLeaseApply=false`
- `resourcePortPersistentLoop=false`

KHR-DASH-F adds no Hyperdensity primary VM surface and no Hyperdensity write path.
Hyperdensity continues to publish governance metadata (Cell/Shell/ResourcePort/SessionTarget)
that the Dashboard reads via the KHR Karl Instance overview projection. No KHR-DASH-F
runtime mutation occurs in Hyperdensity-managed namespaces.

Live evidence: `KHR/evidence/latest/dashboard-full-vm-parity/hyperdensity-full-vm-parity-diagnostic-observation.json`.
