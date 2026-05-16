package resourcelease

import (
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/shellcontinuity"
)

// ApplyContinuityProof merges shell/session/app continuity into guarded-apply verification.
func ApplyContinuityProof(v *VerificationOutcome, proof shellcontinuity.Proof) {
	if v == nil {
		return
	}
	v.SessionContinuityPreserved = proof.SessionContinuityPreserved
	v.ShellContinuityPreserved = proof.ShellContinuityPreserved
	v.AppContinuityPreserved = proof.AppContinuityPreserved
	v.ContinuityEvidence = proof.Evidence
	if proof.InterruptionDetected {
		v.State = VerificationStateFail
	}
}

// VerifyContinuityRequired returns false when continuity proof shows interruption.
func VerifyContinuityRequired(proof shellcontinuity.Proof) bool {
	return proof.ShellContinuityPreserved && proof.AppContinuityPreserved && proof.SessionContinuityPreserved
}
