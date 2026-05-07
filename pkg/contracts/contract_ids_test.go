package contracts

import "testing"

func TestContractIDsAreNotEmpty(t *testing.T) {
	ids := []string{
		LinuxShellComplianceV1ID,
		ResourceEquilibriumV1ID,
		FleetEquilibriumOnboardingV1ID,
		ShellFactoryV1ID,
		ShellClaimV1ID,
		ShellClaimTemplateProfilePackV1ID,
		ReleaseSupportMatrixV1ID,
		EvidenceBundleDemoScenarioPackV1ID,
		LiveResourceAuthorityV1ID,
		ActionSlateV1ID,
		GuardedAutoSandboxV1ID,
		AutoRollbackControllerV1ID,
		BlastRadiusPolicyV1ID,
		WindowsFluidShellV1ID,
		FluidResourceLeaseV1ID,
		WindowsFluidEvidenceV1ID,
		WindowsFluidBlockerV1ID,
		WindowsFluidQMPEvidenceV1ID,
		WindowsFluidRuntimeBundleV1ID,
		WindowsFluidDryRunEvaluationV1ID,
		WindowsFluidAdmissionPolicyV1ID,
		WindowsFluidApplyGovernanceV1ID,
	}
	for _, id := range ids {
		if id == "" {
			t.Fatal("contract id must not be empty")
		}
	}
}
