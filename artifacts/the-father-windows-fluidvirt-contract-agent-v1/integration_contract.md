# Integration Contract

Canonical integration document added in:

- `docs/contracts/windows-fluidvirt-integration-contract-v1.md`

Highlights:

- guest evidence mapping from `modules.fluidShell` to `WindowsFluidEvidence.guestEvidence`
- blocker projection from guest evidence to canonical `WindowsFluidBlocker`
- mandatory READY fields
- mandatory future ACTIVE fields
- optional fields
- redaction policy for machine identity (`machineGuidHash`)
- version compatibility through `agentModuleVersion`
