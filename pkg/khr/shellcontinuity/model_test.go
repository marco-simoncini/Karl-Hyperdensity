package shellcontinuity

import "testing"

func TestCompare_ContinuityPreserved(t *testing.T) {
	before := SnapshotFromWorkload("khr-runtime-sandbox", "khr-native-live-target-abc", "uid-1", "containerd://cid-1")
	after := SnapshotFromWorkload("khr-runtime-sandbox", "khr-native-live-target-abc", "uid-1", "containerd://cid-1")
	p := Compare(before, after)
	if !p.ShellContinuityPreserved || !p.AppContinuityPreserved || !p.SessionContinuityPreserved {
		t.Fatalf("proof=%+v", p)
	}
	if p.InterruptionDetected || p.Evidence.ContinuityState != StatePreserved {
		t.Fatalf("evidence=%+v", p.Evidence)
	}
}

func TestCompare_InterruptionDetected(t *testing.T) {
	before := SnapshotFromWorkload("khr-runtime-sandbox", "pod-a", "uid-1", "cid-1")
	after := SnapshotFromWorkload("khr-runtime-sandbox", "pod-a", "uid-2", "cid-2")
	p := Compare(before, after)
	if !p.InterruptionDetected {
		t.Fatal("expected interruption")
	}
	if p.Evidence.ShellContinuityState != StateInterrupted {
		t.Fatalf("shell state=%s", p.Evidence.ShellContinuityState)
	}
}

func TestDeriveSessionIDs_Deterministic(t *testing.T) {
	a, _, _ := DeriveSessionIDs("ns", "pod", "uid", "cid")
	b, _, _ := DeriveSessionIDs("ns", "pod", "uid", "cid")
	if a != b {
		t.Fatalf("non-deterministic: %s vs %s", a, b)
	}
}

func TestCompare_ContinuityRegression(t *testing.T) {
	before := SnapshotFromWorkload("ns", "p", "u1", "c1")
	after := SnapshotFromWorkload("ns", "p", "u2", "c1")
	if Compare(before, after).ShellContinuityPreserved {
		t.Fatal("pod uid change must break shell continuity")
	}
}
