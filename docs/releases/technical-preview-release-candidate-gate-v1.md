# Technical Preview Release Candidate Gate v1

`releaseCandidateGateId=hyperdensity_technical_preview_release_candidate_gate_v1`  
`releaseTrack=technical_preview`  
`gateMode=go_no_go_validation`  
`targetReadinessState=ready_for_private_technical_preview`

## Canonical dependency IDs

- `supportMatrixId=hyperdensity_release_support_matrix_v1`
- `evidenceBundleId=hyperdensity_evidence_bundle_demo_scenario_pack_v1`
- `liveResourceAuthorityId=hyperdensity_live_resource_authority_v1`
- `documentationPackId=hyperdensity_technical_preview_documentation_pack_v1`
- `policyPackId=hyperdensity_policy_pack_v1`
- `profilePackId=hyperdensity_shell_claim_templates_profile_pack_v1`

## Global safety posture

- `enforcementMode=disabled`
- `autonomousApplyAllowed=false`
- `productionMutationAllowed=false`
- `windowsSupportState=out_of_scope_frozen`

## finalGateChecks[]

The release candidate gate requires all categories below to pass:

1. governance_gate
2. safety_gate
3. creation_path_gate
4. admission_remediation_gate
5. resource_exchange_gate
6. live_resource_authority_gate
7. technical_preview_claim_gate
8. documentation_gate
9. test_build_gate

## blockerChecks[]

The release candidate gate is `blocked` if any blocker is true:

- live parent fabric GET fails
- policy consistency false
- enforcement enabled
- autonomous apply enabled
- production mutation enabled
- windows marked supported
- raw runtime control exposed
- generic kubevirt vm memory template mutation approved
- unsafe wording approved
- technical preview described as ga
- missing support matrix
- missing evidence bundle
- missing live resource authority
- missing documentation pack
- missing rollback proof or no-mutation proof references
- failing tests
- failed validation
- missing required docs

## releaseDecision

Allowed values:

- `ready_for_private_technical_preview`
- `ready_for_founder_investor_demo_only`
- `ready_for_internal_evidence_release_only`
- `blocked`

Decision criteria:

- choose `ready_for_private_technical_preview` only when all required gate categories pass and no blockers are present.
- choose `blocked` if any blocker is present.
- intermediary readiness states are allowed when governance/safety holds but one or more release-wide categories are intentionally deferred.

Required release wording when ready:

"KARL Hyperdensity is ready for a private Technical Preview with evidence-backed Linux container and Linux VM support boundaries, operator-controlled safety posture, and no production/autonomous/enforcement claims."

## releaseBlocker

Single top-priority blocker identifier when decision is `blocked`, otherwise empty.

## requiredFollowUps[]

At minimum track:

- refresh historical proof references before beta/ga claims
- continue object-specific VM evidence expansion without broadening claims
- keep control room ui redesign separate from release boundary contracts
- preserve windows lane frozen state until explicit program activation

## summary

This gate determines go/no-go for private Technical Preview readiness.  
It does not add product capability and does not widen support claims.  
It enforces bounded evidence-backed release truth across policy, support, authority, documentation, tests, and live payload validation.
