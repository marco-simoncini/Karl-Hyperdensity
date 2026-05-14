package installgate

import "testing"

func validSurface() map[string]interface{} {
	return map[string]interface{}{
		"milestone":                           Milestone,
		"productionInstallGateEnabled":        true,
		"hardenedManifestsDefined":            true,
		"installFailClosed":                   true,
		"manifestLintPassed":                  true,
		"rbacHardeningPassed":                 true,
		"serviceAccountBoundaryPassed":        true,
		"deploymentSpecBoundaryPassed":        true,
		"probesConfigured":                    true,
		"metricsEndpointConfigured":           true,
		"leaderElectionWired":                 true,
		"durableStateWired":                   true,
		"upgradeSafetyDefined":                true,
		"rollbackSafetyDefined":               true,
		"installAuditEnabled":                 true,
		"generalProductionAutoAllowed":        false,
		"productionAutoWithPolicy":            false,
		"universalGuaranteedSavingsAllowed":   false,
		"universalGuaranteedSavingsClaimed":   false,
		"estimatedIdleCountedAsMoved":         false,
		"projectedCompressionCountedAsRealized": false,
		"syntheticFleetCountedAsProduction":   false,
		"referenceFleetCountedAsProduction":   false,
		"dashboardExecutor":                   false,
		"fluidvirtPolicyAuthority":            false,
		"inventoryRuntimeExecutor":            false,
		"productionInstallAllowed":            true,
		"installBlockers":                     []interface{}{},
		"claimBoundaries":                     []interface{}{"boundary"},
		"installCandidate": map[string]interface{}{
			"referenceOnly": false, "notProductionInstall": false, "evidenceRefs": []interface{}{"deploy"},
		},
		"rbacHardening": map[string]interface{}{
			"clusterAdmin": false, "namespaceScoped": true, "wildcardResourceAllowed": false, "wildcardVerbAllowed": false,
			"podsExecAllowed": false, "nodesWriteAllowed": false, "rawRuntimeControlsAllowed": false, "directLibvirtAllowed": false, "directCgroupAllowed": false,
		},
		"deploymentSpecBoundary": map[string]interface{}{
			"livenessProbe": true, "readinessProbe": true, "startupProbe": true, "resourceRequests": true, "resourceLimits": true, "runAsNonRoot": true,
			"privileged": false, "allowPrivilegeEscalation": false, "hostPID": false, "hostNetwork": false, "hostIPC": false,
			"generalProductionAutoAllowedEnv": false, "productionAutoWithPolicyEnv": false,
		},
		"metricsEndpoint": map[string]interface{}{
			"generalProductionAutoEnabledGauge": float64(0), "productionAutoWithPolicyEnabledGauge": float64(0),
		},
		"leaderElectionWiring": map[string]interface{}{
			"enabled": true, "rbacConfigured": true, "haProductionProven": false,
		},
		"productionInstallDecision": map[string]interface{}{
			"decision": "install_allowed", "nextAllowedState": "production_canary_only",
		},
	}
}

func TestValidateReferenceFile(t *testing.T) {
	if err := ValidateReferenceFile("../../examples/controller-deployment-hardening-production-install-gate-reference.json"); err != nil {
		t.Fatal(err)
	}
}

func TestRejectProductionAutoEnabled(t *testing.T) {
	doc := validSurface()
	doc["generalProductionAutoAllowed"] = true
	if err := ValidateSurface(doc); err == nil {
		t.Fatal("expected rejection for generalProductionAutoAllowed")
	}
}

func TestRejectClusterAdmin(t *testing.T) {
	doc := validSurface()
	doc["rbacHardening"].(map[string]interface{})["clusterAdmin"] = true
	if err := ValidateSurface(doc); err == nil {
		t.Fatal("expected clusterAdmin rejection")
	}
}

func TestRejectInstallAllowedWithBlockers(t *testing.T) {
	doc := validSurface()
	doc["installBlockers"] = []interface{}{map[string]interface{}{"blockerCode": "unsafe_rbac"}}
	if err := ValidateSurface(doc); err == nil {
		t.Fatal("expected rejection for blockers with install allowed")
	}
}

func TestRejectReferenceCandidateAsInstallAllowed(t *testing.T) {
	doc := validSurface()
	doc["installCandidate"].(map[string]interface{})["referenceOnly"] = true
	if err := ValidateSurface(doc); err == nil {
		t.Fatal("expected rejection for reference-only install candidate")
	}
}

func TestRejectForbiddenMetricsGauge(t *testing.T) {
	doc := validSurface()
	doc["metricsEndpoint"].(map[string]interface{})["generalProductionAutoEnabledGauge"] = float64(1)
	if err := ValidateSurface(doc); err == nil {
		t.Fatal("expected rejection for forbidden metrics gauge")
	}
}

func TestRejectMissingProbe(t *testing.T) {
	doc := validSurface()
	doc["deploymentSpecBoundary"].(map[string]interface{})["startupProbe"] = false
	if err := ValidateSurface(doc); err == nil {
		t.Fatal("expected rejection for missing startupProbe")
	}
}
