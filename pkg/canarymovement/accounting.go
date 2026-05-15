package canarymovement

import "fmt"

func validateAccountingSeparation(doc map[string]interface{}) error {
	if boolv(doc["broadProductionMutationExecuted"]) {
		return fmt.Errorf("broadProductionMutationExecuted must be false")
	}
	if boolv(doc["projectedCompressionCountedAsRealized"]) || boolv(doc["projectedValueCountedAsRealized"]) {
		return fmt.Errorf("projected must not count as realized")
	}
	if !boolv(doc["projectedRealizedSeparated"]) {
		return fmt.Errorf("projectedRealizedSeparated must be true when movement executed")
	}
	sep := doc["projectedRealizedSeparation"].(map[string]interface{})
	if boolv(doc["productionCanaryMovementExecuted"]) && !boolv(sep["separationPassed"]) {
		return fmt.Errorf("separationPassed must be true")
	}
	return nil
}

func evaluateAccounting(doc map[string]interface{}) error {
	if err := validateAccountingSeparation(doc); err != nil {
		return err
	}
	if num(doc["realizedCompressionDelta"]) > 0 && !boolv(doc["postVerifyPassed"]) {
		return fmt.Errorf("realizedCompressionDelta > 0 without post-verify")
	}
	if num(doc["realizedValue"]) > 0 && !boolv(doc["mutationObserved"]) {
		return fmt.Errorf("realizedValue > 0 without mutation observed")
	}
	return nil
}
