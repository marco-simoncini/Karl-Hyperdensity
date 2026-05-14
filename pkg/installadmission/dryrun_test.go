package installadmission

import "testing"

func TestEvaluateDryRunResults(t *testing.T) {
	doc := map[string]interface{}{
		"dryRunResults": []interface{}{
			map[string]interface{}{
				"dryRunMode": "server", "serverSideDryRunPassed": true,
				"evidenceRefs": []interface{}{"server-dryrun"}, "claimBoundary": "server dry-run result",
			},
			map[string]interface{}{
				"dryRunMode": "client", "clientSideDryRunPassed": true,
				"evidenceRefs": []interface{}{"client-dryrun"}, "claimBoundary": "client dry-run result",
			},
		},
	}
	server, client, err := evaluateDryRunResults(doc)
	if err != nil {
		t.Fatal(err)
	}
	if !server || !client {
		t.Fatal("expected both server and client dry-run passed")
	}
}

func TestValidateDryRunRequestsRejectsProductionApplied(t *testing.T) {
	doc := map[string]interface{}{
		"dryRunRequests": []interface{}{
			map[string]interface{}{
				"dryRunRequestId": "req-1", "productionInstallApplied": true,
				"evidenceRefs": []interface{}{"req"}, "claimBoundary": "dry-run request",
			},
		},
	}
	if err := validateDryRunRequests(doc); err == nil {
		t.Fatal("expected productionInstallApplied rejection")
	}
}
