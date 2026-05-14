package marketcontroller

import "time"

// MetricsSnapshot holds reconciler counters.
type MetricsSnapshot struct {
	MetricsID                          string  `json:"metricsId"`
	EmittedAt                          string  `json:"emittedAt"`
	ReconcileTotal                     int     `json:"reconcileTotal"`
	ReconcileSuccessTotal              int     `json:"reconcileSuccessTotal"`
	ReconcileFailureTotal              int     `json:"reconcileFailureTotal"`
	ReconcileDurationP95Ms             int     `json:"reconcileDurationP95Ms"`
	StateLoadTotal                     int     `json:"stateLoadTotal"`
	StateSaveTotal                     int     `json:"stateSaveTotal"`
	StaleWriteRejectTotal              int     `json:"staleWriteRejectTotal"`
	IdempotentReplayTotal              int     `json:"idempotentReplayTotal"`
	ActionsCreatedTotal                int     `json:"actionsCreatedTotal"`
	ActionsInvalidatedTotal            int     `json:"actionsInvalidatedTotal"`
	LeasesExpiredTotal                 int     `json:"leasesExpiredTotal"`
	FuturesExpiredTotal                int     `json:"futuresExpiredTotal"`
	RealizedCompressionRate            float64 `json:"realizedCompressionRate"`
	ProjectedCompressionRate           float64 `json:"projectedCompressionRate"`
	GeneralProductionAutoEnabledGauge  float64 `json:"generalProductionAutoEnabledGauge"`
	ProductionAutoWithPolicyEnabledGauge float64 `json:"productionAutoWithPolicyEnabledGauge"`
}

// MetricsCollector accumulates reconciler metrics.
type MetricsCollector struct {
	snapshot MetricsSnapshot
}

func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{snapshot: MetricsSnapshot{
		MetricsID: "metrics-sprint14",
		GeneralProductionAutoEnabledGauge: 0,
		ProductionAutoWithPolicyEnabledGauge: 0,
	}}
}

func (m *MetricsCollector) RecordReconcile(success bool, durationMs int, result ReconcileResult) {
	m.snapshot.ReconcileTotal++
	if success {
		m.snapshot.ReconcileSuccessTotal++
	} else {
		m.snapshot.ReconcileFailureTotal++
	}
	if durationMs > m.snapshot.ReconcileDurationP95Ms {
		m.snapshot.ReconcileDurationP95Ms = durationMs
	}
	m.snapshot.StateLoadTotal++
	if result.NewResourceVersion != "" {
		m.snapshot.StateSaveTotal++
	}
	if result.StaleWriteRejected {
		m.snapshot.StaleWriteRejectTotal++
	}
	if result.IdempotentReplay {
		m.snapshot.IdempotentReplayTotal++
	}
	m.snapshot.ActionsCreatedTotal += result.ActionsCreated
	m.snapshot.ActionsInvalidatedTotal += result.ActionsInvalidated
	m.snapshot.LeasesExpiredTotal += result.LeasesExpired
	m.snapshot.FuturesExpiredTotal += result.FuturesExpired
	m.snapshot.RealizedCompressionRate = result.RealizedCompressionRate
	m.snapshot.ProjectedCompressionRate = result.ProjectedCompressionRate
	m.snapshot.EmittedAt = time.Now().UTC().Format(time.RFC3339)
}

func (m *MetricsCollector) Snapshot() map[string]interface{} {
	s := m.snapshot
	return map[string]interface{}{
		"metricsId": s.MetricsID, "emittedAt": s.EmittedAt,
		"reconcileTotal": s.ReconcileTotal, "reconcileSuccessTotal": s.ReconcileSuccessTotal,
		"reconcileFailureTotal": s.ReconcileFailureTotal, "reconcileDurationP95Ms": s.ReconcileDurationP95Ms,
		"stateLoadTotal": s.StateLoadTotal, "stateSaveTotal": s.StateSaveTotal,
		"staleWriteRejectTotal": s.StaleWriteRejectTotal, "idempotentReplayTotal": s.IdempotentReplayTotal,
		"actionsCreatedTotal": s.ActionsCreatedTotal, "actionsInvalidatedTotal": s.ActionsInvalidatedTotal,
		"leasesExpiredTotal": s.LeasesExpiredTotal, "futuresExpiredTotal": s.FuturesExpiredTotal,
		"realizedCompressionRate": s.RealizedCompressionRate, "projectedCompressionRate": s.ProjectedCompressionRate,
		"generalProductionAutoEnabledGauge": s.GeneralProductionAutoEnabledGauge,
		"productionAutoWithPolicyEnabledGauge": s.ProductionAutoWithPolicyEnabledGauge,
		"evidenceRefs": []interface{}{"hyperdensity-durable-controller-reconciler-v1"},
		"claimBoundary": "controller metrics snapshot",
	}
}
