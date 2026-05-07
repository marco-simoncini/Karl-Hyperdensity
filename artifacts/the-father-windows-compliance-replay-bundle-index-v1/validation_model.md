# Validation Model

Validator: `ValidateWindowsComplianceReplayBundleIndex`

Checks:

- non-empty `bundleVersion`
- `runCount == len(runs)`
- required run hashes present
- deterministic `runHash` recomputation
- `previousRunHash` link consistency
- `firstRunHash` / `latestRunHash` consistency
- attestation mode allowed (`unsigned-dev` or `future-signable`)
- attestation signature value empty
