# Technical Preview Documentation Pack v1

`documentationPackId=hyperdensity_technical_preview_documentation_pack_v1`  
`documentationPackVersion=v1`  
`releaseTrack=technical_preview`

## Status

- status: `operational` (documentation-only pack)
- intended use: coherent operator/founder/investor TP narrative with bounded claims

## Linked source surfaces

- `supportMatrixId=hyperdensity_release_support_matrix_v1`
- `evidenceBundleId=hyperdensity_evidence_bundle_demo_scenario_pack_v1`
- `liveResourceAuthorityId=hyperdensity_live_resource_authority_v1`
- `policyPackId=hyperdensity_policy_pack_v1`
- `profilePackId=hyperdensity_shell_claim_templates_profile_pack_v1`

## Core pack documents

- [Operator Runbook v1](../runbooks/operator-runbook-v1.md)
- [Technical Preview Release Notes v1](technical-preview-release-notes-v1.md)
- [Technical Preview Readiness Gate v1](technical-preview-readiness-gate-v1.md)
- [Technical Preview Demo Guide v1](../demos/technical-preview-demo-guide-v1.md)

## Related contract and boundary docs

- [Release Support Matrix v1](../contracts/release-support-matrix-v1.md)
- [Evidence Bundle Demo Scenario Pack v1](../contracts/evidence-bundle-demo-scenario-pack-v1.md)
- [Live Resource Authority v1](../contracts/live-resource-authority-v1.md)
- [Policy Pack v1](../contracts/policy-pack-v1.md)
- [Shell Claim Template/Profile Pack v1](../contracts/shell-claim-template-profile-pack-v1.md)

## Intended audience

- platform operators
- technical founders
- technical diligence reviewers
- internal release governance owners

## Release boundary

- Technical Preview only.
- Evidence-backed support claims only.
- No GA claim.
- No enforcement.
- No autonomous apply.
- No production workload mutation.
- Windows out-of-scope/frozen.

## What changed in this pack

- Consolidated practical operator runbook.
- Consolidated externally-readable TP release notes.
- Consolidated pass/fail readiness gate checklist.
- Consolidated 10-15 minute demo script with fallback guidance.

## How to validate docs

1. Run `./scripts/validate.sh`.
2. Confirm all required documentation files exist.
3. Confirm safety wording remains bounded:
   - no Windows support claim
   - no generic VM RAM template mutation claim
   - no enabled enforcement/autonomous/production mutation
   - TP wording, not GA wording
4. Confirm release boundary IDs match canonical contract IDs listed above.

## Next documentation actions

- Expand FAQ for operator troubleshooting cases observed in live TP operations.
- Add richer evidence refresh cadence guidance before beta/GA transitions.
- Add Control Room UI documentation when redesign milestone lands.
