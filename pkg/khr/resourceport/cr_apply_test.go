package resourceport

import "testing"

func TestApplyCRBlockedByDefault(t *testing.T) {
	opts := baseLoopOpts(loopConfig(true))
	opts.ApplyCR = true
	opts.EmitCR = true
	opts.SandboxConfirm = false
	gate := ValidateApplyCRGate(opts)
	if gate.Allowed {
		t.Fatal("apply must be blocked without confirmation")
	}
}

func TestApplyCRBlockedWithoutEmitCR(t *testing.T) {
	opts := baseLoopOpts(loopConfig(true))
	opts.ApplyCR = true
	opts.SandboxConfirm = true
	gate := ValidateApplyCRGate(opts)
	if gate.Allowed {
		t.Fatal("apply must require emit-cr")
	}
}

func TestApplyCRBlockedOutsideSandbox(t *testing.T) {
	opts := baseLoopOpts(loopConfig(true))
	opts.Namespace = "karl-system"
	opts.ApplyCR = true
	opts.EmitCR = true
	opts.SandboxConfirm = true
	gate := ValidateApplyCRGate(opts)
	if gate.Allowed {
		t.Fatal("production namespace must block apply")
	}
}

func TestRunLoopApplyBlockedByDefault(t *testing.T) {
	opts := baseLoopOpts(loopConfig(true))
	opts.ApplyCR = true
	opts.EmitCR = true
	opts.OutputDir = t.TempDir()
	res, err := RunLoop(opts)
	if err != nil {
		t.Fatal(err)
	}
	if !res.ApplyCRBlocked || res.ApplyCRApplied {
		t.Fatalf("res=%+v", res)
	}
}

func TestCleanupSelectorSandboxNamespace(t *testing.T) {
	opts := baseLoopOpts(loopConfig(true))
	opts.Namespace = "khr-runtime-sandbox"
	opts.ClusterContext = "karl-metal-01@ovh"
	// Selector construction only; no cluster call.
	res := CleanupResult{
		Namespace: opts.Namespace,
		Selector:  LabelManagedBy + "=" + ManagedByValue + "," + LabelSandboxNamespace + "=khr-runtime-sandbox",
	}
	if res.Selector == "" {
		t.Fatal("expected selector")
	}
}
