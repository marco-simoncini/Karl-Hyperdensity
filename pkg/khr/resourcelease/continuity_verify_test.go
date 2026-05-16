package resourcelease

import (
	"testing"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/shellcontinuity"
)

func TestApplyContinuityProof_PreservesVerification(t *testing.T) {
	before := shellcontinuity.SnapshotFromWorkload("ns", "pod", "uid-1", "cid-1")
	after := shellcontinuity.SnapshotFromWorkload("ns", "pod", "uid-1", "cid-1")
	proof := shellcontinuity.Compare(before, after)
	v := VerificationOutcome{State: VerificationStatePass, NoRestart: true}
	ApplyContinuityProof(&v, proof)
	if !v.ShellContinuityPreserved || !v.AppContinuityPreserved || !v.SessionContinuityPreserved {
		t.Fatalf("verification=%+v", v)
	}
	if v.State != VerificationStatePass {
		t.Fatalf("state=%s", v.State)
	}
}

func TestApplyContinuityProof_InterruptionFailsVerification(t *testing.T) {
	before := shellcontinuity.SnapshotFromWorkload("ns", "pod", "uid-1", "cid-1")
	after := shellcontinuity.SnapshotFromWorkload("ns", "pod", "uid-2", "cid-2")
	proof := shellcontinuity.Compare(before, after)
	v := VerificationOutcome{State: VerificationStatePass}
	ApplyContinuityProof(&v, proof)
	if v.State != VerificationStateFail {
		t.Fatalf("state=%s", v.State)
	}
}
