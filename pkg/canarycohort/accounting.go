package canarycohort

import "fmt"

func validateAccountingAggregation(doc map[string]interface{}) error {
	rca := doc["realizedCompressionAggregation"].(map[string]interface{})
	if num(doc["aggregateRealizedMovedIdleValue"]) > 0 && !boolv(rca["movementEvidencePresent"]) {
		return fmt.Errorf("aggregateRealizedMovedIdleValue > 0 without movement evidence")
	}
	if num(doc["aggregateRealizedCompressionDelta"]) > 0 && !boolv(rca["movementEvidencePresent"]) {
		return fmt.Errorf("aggregateRealizedCompressionDelta > 0 without movement evidence")
	}
	if num(doc["aggregateRealizedValue"]) > 0 && !boolv(doc["realizedValueAggregated"]) {
		return fmt.Errorf("aggregateRealizedValue > 0 without aggregation evidence")
	}
	sep := doc["projectedRealizedSeparation"].(map[string]interface{})
	if boolv(sep["projectedCompressionCountedAsRealized"]) || boolv(sep["projectedValueCountedAsRealized"]) {
		return fmt.Errorf("projected must not count as realized")
	}
	if boolv(doc["productionCanaryCohortExecuted"]) && !boolv(sep["separationPassed"]) {
		return fmt.Errorf("separationPassed must be true")
	}
	return nil
}
