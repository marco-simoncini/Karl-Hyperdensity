# Signed Attestation Trust Model

Model: `WindowsFluidSignedAttestationTrustModel`

- Subjects include dry-run, admission, governance, invariant set, revalidation, guard, kill switch.
- Trust boundaries include controller, sidecar, guest agent, policy plane, audit store.
- Signer roles are future-only and not implemented in this phase.
- Replay protection and freshness checks are required before any future unlock.
