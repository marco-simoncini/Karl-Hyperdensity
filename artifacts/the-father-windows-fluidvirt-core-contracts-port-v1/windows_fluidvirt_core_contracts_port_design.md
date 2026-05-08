# Windows FluidVirt Core Contracts Port Design v1

## Milestone
`hyperdensity_windows_fluidvirt_core_contracts_port_v1`

## Goal
Port only core Windows FluidVirt contracts/models to integration branch with safety-first defaults and no runtime execution path.

## Included Scope
- `pkg/windowsfluidvirt/product_model.go`
- `pkg/windowsfluidvirt/action_slate.go`
- `pkg/windowsfluidvirt/blockers.go`
- `docs/contracts/windows-fluidvirt-product-model-v1.md`
- minimal fixtures under `examples/windows-fluid-product-fixtures/`
- safety-focused tests under `pkg/windowsfluidvirt/*_test.go`

## Explicit Non-Goals
- no direct merge from `The-Father-Windows`
- no actuator/compliance/executor/controlled-apply implementation
- no runtime mutation enablement
- no autonomous apply enablement
- no dashboard/inventory/os-iso work

## Design Notes
- model naming aligned with Hyperdensity/Parent Fabric vocabulary (`WindowsFluidVirtProductModel`, `WindowsFluidVirtActionSlate`, `WindowsFluidVirtBlocker`)
- action slate is strictly planning/model only
- claim boundary and support boundary are encoded as explicit booleans defaulting to safe values
- blocker catalog encodes forbidden claims and missing gate evidence
