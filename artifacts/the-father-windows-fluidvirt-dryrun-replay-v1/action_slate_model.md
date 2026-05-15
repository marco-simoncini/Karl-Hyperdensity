# action_slate_model

Model: `WindowsFluidActionSlate`

Mandatory safety fields:

- `runtimeMode=in-place-qmp`
- `mutationAllowed=false`
- `applyAllowed=false`
- proof flags (`qmpReady`, `guestAckReady`, `sameQemuProof`, `sameNodeProof`, `samePodProof`, `noRebootProof`)
- rollback and return-to-floor readiness

Action types:

- `certify-shell`
- `prepare-cpu-lease`
- `prepare-memory-lease`
- `evidence-refresh`
- `quarantine`
- `blocked`

No runtime command payloads are present in the action slate.
