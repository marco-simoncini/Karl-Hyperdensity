package observationwindow

func ReferenceSurface() (map[string]interface{}, error) {
	return referenceSurfaceDoc(), nil
}

func referenceSurfaceDoc() map[string]interface{} {
	ticks := []interface{}{
		tickRecord("tick-1", 1, "2026-05-14T18:00:00Z", "2026-05-14T18:00:45Z", 12, 4, 3, 1, 1, 3, 2, 1100, 0.16),
		tickRecord("tick-2", 2, "2026-05-14T18:05:00Z", "2026-05-14T18:05:42Z", 12, 4, 3, 1, 0, 3, 2, 1150, 0.17),
		tickRecord("tick-3", 3, "2026-05-14T18:10:00Z", "2026-05-14T18:10:40Z", 13, 4, 3, 0, 0, 3, 2, 1200, 0.18),
		tickRecord("tick-4", 4, "2026-05-14T18:15:00Z", "2026-05-14T18:15:38Z", 13, 5, 3, 0, 0, 3, 2, 1250, 0.19),
	}
	return map[string]interface{}{
		"milestone": Milestone, "surfaceVersion": "v1", "observationWindowId": "obs-win-s20-v1",
		"generatedAt": "2026-05-14T18:16:00Z",
		"sourceRealInputMarketTickRef": "configmap://karl-system/hyperdensity-runtime-market-real-inputs-no-apply-v1",
		"controllerMode": "production_canary_only",
		"productionObservationWindowEnabled": true, "productionObservationMode": true, "noApplyMode": true,
		"multiTickWindowObserved": true, "observationWindowPassed": true,
		"tickCount": float64(4), "successfulTickCount": float64(4), "failedTickCount": float64(0),
		"firstTickAt": "2026-05-14T18:00:00Z", "lastTickAt": "2026-05-14T18:15:38Z",
		"windowDurationSeconds": float64(938),
		"realInputTicksObserved": true, "inputStabilityAnalyzed": true, "donorStabilityAnalyzed": true,
		"receiverPressurePersistenceAnalyzed": true, "idleOpportunityPersistenceAnalyzed": true,
		"blockerDecayAnalyzed": true, "staleInputDecayAnalyzed": true,
		"actionSlateRefreshed": true, "resourceFuturesRefreshed": true,
		"projectedCompressionTrendComputed": true, "projectedValueTrendComputed": true,
		"realizedValueSeparated": true, "noApplySafetyWindowVerified": true, "dashboardProjectionUpdated": true,
		"productionMovementExecuted": false, "broadProductionMutationExecuted": false,
		"realizedMovementCount": float64(0), "realizedMovedIdleValue": float64(0), "realizedCompressionDelta": 0.0,
		"projectedMovedIdleValue": float64(1250), "projectedCompressionDelta": 0.03,
		"projectedValueOpportunity": float64(4800), "persistentIdleValue": float64(3200),
		"persistentPressureValue": float64(1600), "blockerDecayRate": 0.75, "staleInputDecayRate": 1.0,
		"candidateStabilityRate": 0.82, "futureRefreshRate": 0.5,
		"generalProductionAutoAllowed": false, "productionAutoWithPolicy": false,
		"universalGuaranteedSavingsAllowed": false, "universalGuaranteedSavingsClaimed": false,
		"estimatedIdleCountedAsMoved": false, "projectedCompressionCountedAsRealized": false,
		"projectedValueCountedAsRealized": false, "syntheticFleetCountedAsProduction": false,
		"referenceFleetCountedAsProduction": false, "dashboardExecutor": false,
		"fluidvirtPolicyAuthority": false, "fluidvirtObservationWindowAuthority": false,
		"inventoryRuntimeExecutor": false,
		"tickSequence": map[string]interface{}{
			"tickSequenceId": "seq-1", "tickCount": float64(4), "ticks": ticks,
			"monotonicSequence": true, "noMissingTicks": true, "allTicksNoApply": true,
			"evidenceRefs": []interface{}{"artifact://tick-sequence"}, "claimBoundary": "tick sequence only",
		},
		"inputStabilityAnalysis": map[string]interface{}{
			"inputStabilityAnalysisId": "isa-1", "stableShellCount": float64(11), "changedShellCount": float64(1),
			"disappearedShellCount": float64(0), "newShellCount": float64(1), "stabilityRate": 0.92,
			"inputStabilityAnalyzed": true,
			"evidenceRefs": []interface{}{"artifact://input-stability"}, "claimBoundary": "input stability only",
		},
		"donorStabilityAnalysis": map[string]interface{}{
			"donorStabilityAnalysisId": "dsa-1", "persistentDonorCount": float64(3), "newDonorCount": float64(1),
			"removedDonorCount": float64(0), "donorStabilityRate": 0.75,
			"topPersistentDonors": []interface{}{
				map[string]interface{}{"shellId": "d1", "ticksPresent": float64(4)},
				map[string]interface{}{"shellId": "d2", "ticksPresent": float64(4)},
			},
			"evidenceRefs": []interface{}{"artifact://donor-stability"}, "claimBoundary": "donor stability only",
		},
		"receiverPressurePersistence": map[string]interface{}{
			"receiverPressurePersistenceId": "rpp-1", "persistentPressureReceiverCount": float64(2),
			"transientPressureReceiverCount": float64(1), "pressurePersistenceRate": 0.67,
			"topPersistentReceivers": []interface{}{
				map[string]interface{}{"shellId": "r1", "ticksPresent": float64(4)},
			},
			"evidenceRefs": []interface{}{"artifact://pressure-persistence"}, "claimBoundary": "receiver pressure persistence only",
		},
		"idleOpportunityPersistence": map[string]interface{}{
			"idleOpportunityPersistenceId": "iop-1", "persistentIdleValue": float64(3200),
			"transientIdleValue": float64(400), "persistentIdleShellCount": float64(3),
			"idleOpportunityPersistenceRate": 0.8,
			"topPersistentIdleShells": []interface{}{
				map[string]interface{}{"shellId": "sh-1", "persistentIdleValue": float64(1600)},
			},
			"evidenceRefs": []interface{}{"artifact://idle-persistence"}, "claimBoundary": "idle opportunity persistence only",
		},
		"blockerDecayAnalysis": map[string]interface{}{
			"blockerDecayAnalysisId": "bda-1", "initialBlockerCount": float64(4), "finalBlockerCount": float64(1),
			"blockersResolved": float64(3), "blockersAdded": float64(0), "blockerDecayRate": 0.75,
			"persistentBlockers": []interface{}{map[string]interface{}{"blockerCode": "stale_input"}},
			"resolvedBlockers": []interface{}{
				map[string]interface{}{"blockerCode": "missing_slo"},
				map[string]interface{}{"blockerCode": "missing_rollback"},
			},
			"evidenceRefs": []interface{}{"artifact://blocker-decay"}, "claimBoundary": "blocker decay only",
		},
		"staleInputDecayAnalysis": map[string]interface{}{
			"staleInputDecayAnalysisId": "sida-1", "initialStaleInputCount": float64(2), "finalStaleInputCount": float64(0),
			"staleInputsInvalidated": float64(2), "staleInputDecayRate": 1.0,
			"evidenceRefs": []interface{}{"artifact://stale-decay"}, "claimBoundary": "stale input decay only",
		},
		"actionSlateRefreshEvidence": map[string]interface{}{
			"actionSlateRefreshId": "asr-1", "initialActionCount": float64(3), "finalActionCount": float64(3),
			"refreshedActionCount": float64(2), "stableActionCount": float64(1), "expiredActionCount": float64(0),
			"selectedForExecutionCount": float64(0), "allActionsNoApply": true,
			"evidenceRefs": []interface{}{"artifact://action-refresh"}, "claimBoundary": "action slate refresh only",
		},
		"resourceFutureRefreshEvidence": map[string]interface{}{
			"futureRefreshId": "frr-1", "initialFutureCount": float64(2), "finalFutureCount": float64(2),
			"refreshedFutureCount": float64(1), "expiredFutureCount": float64(0), "futureRefreshRate": 0.5,
			"evidenceRefs": []interface{}{"artifact://future-refresh"}, "claimBoundary": "future refresh only",
		},
		"projectedCompressionTrend": map[string]interface{}{
			"projectedCompressionTrendId": "pct-1", "firstProjectedCompressionRate": 0.16,
			"lastProjectedCompressionRate": 0.19, "projectedCompressionDelta": 0.03,
			"trendDirection": "increasing", "projectedCompressionCountedAsRealized": false,
			"evidenceRefs": []interface{}{"artifact://compression-trend"}, "claimBoundary": "projected compression trend only",
		},
		"projectedValueTrend": map[string]interface{}{
			"projectedValueTrendId": "pvt-1", "firstProjectedValue": float64(1100),
			"lastProjectedValue": float64(1250), "projectedValueDelta": float64(150),
			"trendDirection": "increasing", "projectedValueCountedAsRealized": false,
			"evidenceRefs": []interface{}{"artifact://value-trend"}, "claimBoundary": "projected value trend only",
		},
		"realizedValueSeparation": map[string]interface{}{
			"realizedValueSeparationId": "rvs-1", "realizedMovementCount": float64(0),
			"realizedMovedIdleValue": float64(0), "realizedCompressionDelta": 0.0,
			"movementEvidencePresent": false, "projectedValueExcludedFromRealized": true,
			"estimatedValueExcludedFromRealized": true, "syntheticValueExcludedFromRealized": true,
			"referenceValueExcludedFromRealized": true, "realizedValueSeparated": true,
			"evidenceRefs": []interface{}{"artifact://realized-separation"}, "claimBoundary": "realized value separation only",
		},
		"noApplySafetyWindow": map[string]interface{}{
			"safetyWindowId": "safety-1", "productionMovementExecuted": false, "broadProductionMutationExecuted": false,
			"selectedForExecutionCount": float64(0), "generalProductionAutoAllowed": false, "productionAutoWithPolicy": false,
			"dashboardExecutor": false, "fluidvirtObservationWindowAuthority": false,
			"rawRuntimeControlsExposed": false, "inventoryRuntimeExecutor": false,
			"evidenceRefs": []interface{}{"artifact://safety-window"}, "claimBoundary": "no-apply safety window only",
		},
		"observationMetrics": map[string]interface{}{
			"metricsId": "metrics-1", "emittedAt": "2026-05-14T18:16:00Z",
			"windowTickCount": float64(4), "successfulTickCount": float64(4), "failedTickCount": float64(0),
			"inputStabilityRateGauge": 0.92, "donorStabilityRateGauge": 0.75,
			"pressurePersistenceRateGauge": 0.67, "blockerDecayRateGauge": 0.75,
			"staleInputDecayRateGauge": 1.0, "candidateStabilityRateGauge": 0.82,
			"futureRefreshRateGauge": 0.5, "projectedCompressionDeltaGauge": 0.03,
			"projectedMovedIdleValueGauge": float64(1250), "realizedMovedIdleValueGauge": float64(0),
			"generalProductionAutoEnabledGauge": float64(0), "productionAutoWithPolicyEnabledGauge": float64(0),
			"evidenceRefs": []interface{}{"artifact://window-metrics"}, "claimBoundary": "observation window metrics only",
		},
		"observationEvents": []interface{}{
			map[string]interface{}{
				"eventId": "evt-1", "eventType": "observation_window_started", "occurredAt": "2026-05-14T18:00:00Z",
				"sourceComponent": "Karl-Hyperdensity", "involvedObjectKind": "ConfigMap",
				"involvedObjectName": "hyperdensity-production-observation-window-v1",
				"summary": "observation window started", "evidenceRefs": []interface{}{"artifact://evt-1"}, "claimBoundary": "observation event only",
			},
			map[string]interface{}{
				"eventId": "evt-2", "eventType": "observation_window_passed", "occurredAt": "2026-05-14T18:16:00Z",
				"sourceComponent": "Karl-Hyperdensity", "involvedObjectKind": "ConfigMap",
				"involvedObjectName": "hyperdensity-production-observation-window-v1",
				"summary": "no-apply observation window passed", "evidenceRefs": []interface{}{"artifact://evt-2"}, "claimBoundary": "observation event only",
			},
		},
		"sourceOfTruth": map[string]interface{}{
			"runtimeMutation": "FluidVirt", "marketDecision": "Karl-Hyperdensity",
			"identitySignals": "Karl-Inventory", "operatorUI": "Karl-Dashboard (projection only)",
		},
		"safetyInvariants": []interface{}{
			"no production movement in observation window", "projected value never counted as realized",
			"dashboard projection read-only",
		},
		"claimBoundaries": []interface{}{
			"observation window generates projected trends only", "no broad production mutation claim",
			"realized value requires movement evidence",
		},
		"blockers": []interface{}{},
		"warnings": []interface{}{},
		"referenceScenarios": []interface{}{
			map[string]interface{}{"scenarioId": "success", "observationWindowPassed": true, "tickCount": float64(4)},
			map[string]interface{}{"scenarioId": "too_few_ticks", "tickCount": float64(2), "observationWindowPassed": false},
			map[string]interface{}{"scenarioId": "tick_production_movement", "productionMovementExecuted": true, "observationWindowPassed": false},
			map[string]interface{}{"scenarioId": "action_selected", "selectedForExecutionCount": float64(1), "observationWindowPassed": false},
			map[string]interface{}{"scenarioId": "projected_realized", "projectedCompressionCountedAsRealized": true, "observationWindowPassed": false},
			map[string]interface{}{"scenarioId": "estimated_moved", "estimatedIdleCountedAsMoved": true, "observationWindowPassed": false},
			map[string]interface{}{"scenarioId": "synthetic_reference_prod", "syntheticFleetCountedAsProduction": true, "observationWindowPassed": false},
			map[string]interface{}{"scenarioId": "general_auto", "generalProductionAutoAllowed": true, "observationWindowPassed": false},
			map[string]interface{}{"scenarioId": "prod_auto_policy", "productionAutoWithPolicy": true, "observationWindowPassed": false},
			map[string]interface{}{"scenarioId": "dashboard_executor", "dashboardExecutor": true, "observationWindowPassed": false},
		},
	}
}

func tickRecord(id string, seq int, started, completed string, shells, donors, receivers, stale, invalidated, actions, futures int, projectedValue float64, projectedDelta float64) map[string]interface{} {
	return map[string]interface{}{
		"tickId": id, "sequenceNumber": float64(seq), "startedAt": started, "completedAt": completed,
		"status": "passed", "observedShellCount": float64(shells), "eligibleDonorCount": float64(donors),
		"pressureReceiverCount": float64(receivers), "staleInputCount": float64(stale),
		"invalidatedInputCount": float64(invalidated), "generatedActionCount": float64(actions),
		"generatedFutureCount": float64(futures), "projectedMovedIdleValue": projectedValue,
		"projectedCompressionDelta": projectedDelta, "productionMovementExecuted": false,
		"selectedForExecutionCount": float64(0),
		"evidenceRefs": []interface{}{"artifact://" + id}, "claimBoundary": "tick record only",
	}
}
