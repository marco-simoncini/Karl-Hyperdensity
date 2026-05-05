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
	}
	for _, id := range ids {
		if id == "" {
			t.Fatal("contract id must not be empty")
		}
	}
}
