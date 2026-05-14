package marketcontroller

import "testing"

func testSnapshot() Snapshot {
	return Snapshot{
		TopKDonors: 3, TopKReceivers: 3,
		FullDonorCount: 12, FullReceiverCount: 8,
		CurrentMovedIdleValue: 0.0042,
		CurrentEligibleIdleValue: 0.085,
		CurrentCompressionRate: 0.04941176470588235,
		TargetCompressionRate: 0.25,
		RateLimitRemaining: 5,
		Donors: []DonorCandidate{
			{DonorShellID: "shell-container-donor-a", Resource: "cpu", EligibleIdleAmount: "500m", EligibleIdleValue: 0.056, RollbackReady: true, SloGuardAvailable: true, NoRegressionAvailable: true, RiskScore: 0.1, GuaranteePotential: 0.8},
			{DonorShellID: "shell-container-donor-b", Resource: "cpu", EligibleIdleAmount: "300m", EligibleIdleValue: 0.018, RollbackReady: true, SloGuardAvailable: true, RiskScore: 0.2, GuaranteePotential: 0.5},
			{DonorShellID: "shell-windows-hyper", Resource: "cpu", EligibleIdleAmount: "0", EligibleIdleValue: 0.008, WindowsEvidenceGated: true},
			{DonorShellID: "shell-reference-sample", ReferenceOnly: true, EligibleIdleValue: 0.01},
			{DonorShellID: "shell-synthetic-fleet", SyntheticShadow: true, EligibleIdleValue: 0.015},
			{DonorShellID: "shell-protected-core-001", Protected: true, EligibleIdleValue: 0.025},
		},
		Receivers: []ReceiverCandidate{
			{ReceiverShellID: "shell-container-replica-b", Resource: "cpu", RequestedAmount: "500m", PotentialValueCapture: 0.0042, SloProfilePresent: true},
			{ReceiverShellID: "shell-container-replica-c", Resource: "cpu", RequestedAmount: "400m", PotentialValueCapture: 0.011},
			{ReceiverShellID: "shell-container-replica-d", Resource: "cpu", RequestedAmount: "200m", PotentialValueCapture: 0.005},
		},
	}
}

func TestRunTickNoFullNxN(t *testing.T) {
	res, err := RunTick(testSnapshot())
	if err != nil {
		t.Fatal(err)
	}
	if !res.NoFullNxNPairing {
		t.Fatal("must avoid full N×N")
	}
	if res.EvaluatedPairCount > 3*3 {
		t.Fatal("evaluated pairs exceed top-K bound")
	}
	if res.AvoidedPairCount != res.FullPairSpace-res.EvaluatedPairCount {
		t.Fatal("avoided pair count mismatch")
	}
}

func TestRunTickExcludesSyntheticReference(t *testing.T) {
	res, err := RunTick(testSnapshot())
	if err != nil {
		t.Fatal(err)
	}
	for _, d := range res.TopDonors {
		id, _ := d["donorShellId"].(string)
		if id == "shell-reference-sample" || id == "shell-synthetic-fleet" {
			t.Fatalf("synthetic/reference donor in top-K: %s", id)
		}
	}
}

func TestRunTickWindowsRemediationOnly(t *testing.T) {
	s := testSnapshot()
	res, err := RunTick(s)
	if err != nil {
		t.Fatal(err)
	}
	for _, a := range res.GeneratedActions {
		if a["donorShellId"] == "shell-windows-hyper" {
			if a["executionScopeRecommendation"] != "remediation_only" {
				t.Fatal("windows must be remediation_only")
			}
		}
	}
}

func TestRunTickKillSwitchBlocks(t *testing.T) {
	s := testSnapshot()
	s.KillSwitchActive = true
	res, err := RunTick(s)
	if err != nil {
		t.Fatal(err)
	}
	for _, a := range res.GeneratedActions {
		if a["executionScopeRecommendation"] != "blocked" {
			t.Fatal("kill switch must block scope")
		}
	}
}

func TestRunTickGeneratedActionsHaveRequiredFields(t *testing.T) {
	res, err := RunTick(testSnapshot())
	if err != nil {
		t.Fatal(err)
	}
	if len(res.GeneratedActions) < 1 {
		t.Fatal("expected generated actions")
	}
	for _, a := range res.GeneratedActions {
		for _, key := range []string{"actionId", "donorShellId", "receiverShellId", "resource", "amount"} {
			if a[key] == nil || a[key] == "" {
				t.Fatalf("missing %s", key)
			}
		}
		if boolOr(a["generalProductionAutoAllowed"]) || boolOr(a["productionAutoWithPolicy"]) {
			t.Fatal("general production auto forbidden")
		}
	}
}

func TestRunTickFuturesHaveExpiration(t *testing.T) {
	res, err := RunTick(testSnapshot())
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range res.GeneratedFutures {
		if f["expiration"] == nil || f["expiration"] == "" {
			t.Fatal("future missing expiration")
		}
		inv, ok := f["invalidationReasons"].([]interface{})
		if !ok || len(inv) == 0 {
			t.Fatal("future missing invalidation reasons")
		}
	}
}

func TestRunTickAddressesDonorA(t *testing.T) {
	res, err := RunTick(testSnapshot())
	if err != nil {
		t.Fatal(err)
	}
	if !res.HighSavingsOpportunityAddressed {
		t.Fatal("shell-container-donor-a should be addressed")
	}
	if res.ProjectedCompressionRate <= 0.04941176470588235 {
		t.Fatal("projected compression should improve over Sprint 11B")
	}
}
