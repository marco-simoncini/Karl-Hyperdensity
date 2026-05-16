package windowslane

import "testing"

func TestCapabilitiesForProviderHostRuntime(t *testing.T) {
	c, err := CapabilitiesForProvider(ProviderWindowsHostRuntime)
	if err != nil {
		t.Fatal(err)
	}
	if !c.CPULiveScaleSupported || !c.RAMLiveScaleSupported || c.RequiresRestart {
		t.Fatalf("caps=%+v", c)
	}
	if err := ValidateCapabilities(c, ProviderWindowsHostRuntime); err != nil {
		t.Fatal(err)
	}
}

func TestCapabilitiesForProviderKubevirtBlocked(t *testing.T) {
	c, err := CapabilitiesForProvider(ProviderKubevirtCompatibility)
	if err != nil {
		t.Fatal(err)
	}
	if c.CPULiveScaleSupported || !c.RequiresRestart {
		t.Fatalf("caps=%+v", c)
	}
}

func TestBlockMemoryOnCompatibility(t *testing.T) {
	b := BlockMemoryOnCompatibility()
	if !b.Blocked || b.BlockedState != BlockedRequiresRestart || !b.NoApply {
		t.Fatalf("b=%+v", b)
	}
}
