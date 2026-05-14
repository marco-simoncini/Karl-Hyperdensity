package installgate

import "testing"

func TestLintManifestDirPassesForInstallGateBundle(t *testing.T) {
	res, err := LintManifestDir("../../deploy/kubernetes/controller/install-gate")
	if err != nil {
		t.Fatal(err)
	}
	if !res.Passed {
		t.Fatalf("expected pass got blockers=%v", res.Blockers)
	}
}

func TestLintManifestDirRejectsUnsafeBundle(t *testing.T) {
	res, err := LintManifestDir("../../deploy/kubernetes")
	if err == nil && res.Passed {
		t.Fatal("expected non-passing result on non-install-gate root")
	}
}
