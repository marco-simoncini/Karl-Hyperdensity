# Inventory FluidShell Module Summary

Implemented in `Karl-Inventory` shared Windows service host:

- config model extension: `modules.fluidShell`
- runtime worker: `Services/FluidShell/FluidShellModuleWorker.cs`
- system probe: `Services/FluidShell/WindowsFluidShellSystemProbe.cs`
- evidence sender: `Services/FluidShell/FluidShellEvidenceSender.cs`
- preflight contracts/evaluator/runtime: `src/KarlInventoryAgent.FluidShell/*`

Behavior:

- evidence-only preflight
- no CPU/RAM apply
- no reboot/recreate actions
- no new service, no new MSI
- conservative return-to-floor (`ready=false` by default in this phase)
- emits canonical blockers (pending reboot, missing boot/machine/cpu/memory truth, driver unverified, guest ack missing, critical events)
