package installadmission

import "fmt"

func evaluateDryRunResults(doc map[string]interface{}) (serverPassed bool, clientPassed bool, err error) {
	items, ok := doc["dryRunResults"].([]interface{})
	if !ok || len(items) == 0 {
		return false, false, fmt.Errorf("dryRunResults required")
	}
	for _, item := range items {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		if strv(m["claimBoundary"]) == "" || !hasNonEmptyStringList(m["evidenceRefs"]) {
			return false, false, fmt.Errorf("dryRunResult missing evidenceRefs or claimBoundary")
		}
		mode := strv(m["dryRunMode"])
		switch mode {
		case "server":
			serverPassed = boolv(m["serverSideDryRunPassed"])
		case "client":
			clientPassed = boolv(m["clientSideDryRunPassed"])
		case "client_and_server":
			serverPassed = boolv(m["serverSideDryRunPassed"])
			clientPassed = boolv(m["clientSideDryRunPassed"])
		}
	}
	return serverPassed, clientPassed, nil
}

func validateDryRunRequests(doc map[string]interface{}) error {
	items, ok := doc["dryRunRequests"].([]interface{})
	if !ok || len(items) == 0 {
		return fmt.Errorf("dryRunRequests required")
	}
	for _, item := range items {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		if strv(m["claimBoundary"]) == "" || !hasNonEmptyStringList(m["evidenceRefs"]) {
			return fmt.Errorf("dryRunRequest missing evidenceRefs or claimBoundary")
		}
		if boolv(m["productionInstallApplied"]) {
			return fmt.Errorf("productionInstallApplied true in dry-run request")
		}
	}
	return nil
}
