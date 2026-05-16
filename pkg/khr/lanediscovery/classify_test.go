package lanediscovery

import "testing"

func TestClassifySandboxLinuxContainer(t *testing.T) {
	lane, _, class, live, block := ClassifyWorkload(WorkloadHint{
		Name: "khr-runtime-linux-target", Namespace: "khr-runtime-sandbox",
		OSFamily: "linux", VMType: "container", Running: true, SandboxPod: true,
	})
	if lane != LaneLinuxContainerCgroup || class != ClassificationLiveInPlaceCapable || !live || block != nil {
		t.Fatalf("lane=%s class=%s live=%v block=%+v", lane, class, live, block)
	}
}

func TestClassifyWindowsVMCompatibilityFallback(t *testing.T) {
	_, _, class, _, block := ClassifyWorkload(WorkloadHint{
		Name: "win11-pool-0", Namespace: "karl", OSFamily: "windows", VMType: "vm", Running: true,
	})
	if class != ClassificationObservationOnly || block == nil || block.State == "" {
		t.Fatalf("class=%s block=%+v", class, block)
	}
}

func TestInferOSFamilyFromName(t *testing.T) {
	if InferOSFamily("master-win11", nil) != "windows" {
		t.Fatal("expected windows")
	}
	if InferOSFamily("linux-vm-hd-donor-v1", nil) != "linux" {
		t.Fatal("expected linux")
	}
}
