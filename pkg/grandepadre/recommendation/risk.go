// Package recommendation is a dry-run action slate skeleton (Sprint 13).
package recommendation

// Risk is a coarse risk label for recommendations (not admission).
type Risk string

const (
	RiskLow     Risk = "low"
	RiskMedium  Risk = "medium"
	RiskHigh    Risk = "high"
	RiskBlocked Risk = "blocked"
)
