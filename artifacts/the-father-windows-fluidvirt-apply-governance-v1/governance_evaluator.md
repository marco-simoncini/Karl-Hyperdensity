# governance_evaluator

Evaluator: `EvaluateWindowsFluidApplyGovernance`

Input:

- admission decision
- runtime evidence bundle
- policy pack (optional, conservative default)
- requested governance future action
- deterministic evaluation time (optional)

Output:

- apply governance contract
- transition proof
- runtime invariant set
- pre-apply revalidation contract
- policy attestation
- final governance phase
- next safe step

Safety guarantees:

- mutation/apply always disabled
- no execution path to runtime apply
- no `ACTIVE` or `APPLYING` outputs
