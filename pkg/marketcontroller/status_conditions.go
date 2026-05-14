package marketcontroller

import "time"

// Condition types for Sprint 14.
const (
	ConditionReady                        = "Ready"
	ConditionReconciling                  = "Reconciling"
	ConditionDegraded                     = "Degraded"
	ConditionBackpressure               = "Backpressure"
	ConditionStateStoreReady            = "StateStoreReady"
	ConditionLeaderElectionReady        = "LeaderElectionReady"
	ConditionGeneralProductionAutoDisabled = "GeneralProductionAutoDisabled"
)

func newStatusCondition(condType, status, reason, message string) map[string]interface{} {
	return map[string]interface{}{
		"conditionId":          "cond-" + condType,
		"type":                 condType,
		"status":               status,
		"reason":               reason,
		"message":              message,
		"observedGeneration":   1,
		"lastTransitionTime":   time.Now().UTC().Format(time.RFC3339),
		"severity":             conditionSeverity(condType, status),
		"evidenceRefs":         []interface{}{"hyperdensity-durable-controller-reconciler-v1"},
		"claimBoundary":        "controller status condition",
	}
}

func conditionSeverity(condType, status string) string {
	if status != "True" {
		return "info"
	}
	switch condType {
	case ConditionDegraded:
		return "warning"
	default:
		return "info"
	}
}

func baseStatusConditions() []map[string]interface{} {
	return []map[string]interface{}{
		newStatusCondition(ConditionStateStoreReady, "True", "ConfigMapStoreReady", "kubernetes_configmap store implemented"),
		newStatusCondition(ConditionLeaderElectionReady, "True", "LeaderElectionConfigured", "leader election reference configured; not HA production proven"),
		newStatusCondition(ConditionGeneralProductionAutoDisabled, "True", "PolicyEnforced", "general production auto disabled"),
		newStatusCondition(ConditionBackpressure, "False", "NoBackpressure", "no backpressure active"),
	}
}

func setReadyConditions(conditions []map[string]interface{}) []map[string]interface{} {
	out := append([]map[string]interface{}{}, conditions...)
	out = appendOrReplace(out, newStatusCondition(ConditionReady, "True", "ReconcileSuccess", "reconcile completed successfully"), "type")
	out = appendOrReplace(out, newStatusCondition(ConditionReconciling, "False", "Idle", "not reconciling"), "type")
	out = appendOrReplace(out, newStatusCondition(ConditionDegraded, "False", "Healthy", "not degraded"), "type")
	return out
}

func setDegradedConditions(conditions []map[string]interface{}, reason, message string) []map[string]interface{} {
	out := append([]map[string]interface{}{}, conditions...)
	out = appendOrReplace(out, newStatusCondition(ConditionReady, "False", reason, message), "type")
	out = appendOrReplace(out, newStatusCondition(ConditionDegraded, "True", reason, message), "type")
	out = appendOrReplace(out, newStatusCondition(ConditionReconciling, "False", "Degraded", "reconcile degraded"), "type")
	return out
}
