# policy_attestation_model

Model: `WindowsFluidPolicyAttestation`

Attestation scope:

- subject type and subject ref
- policy version
- evidence refs
- blocker snapshot
- invariant snapshot
- decision snapshot
- attestor component metadata
- signature envelope

Current phase constraints:

- signature mode is `unsigned-dev` (or `future-signable` for compatibility)
- signature value remains empty
- no key management and no real signatures
- no secrets in attestation payload
