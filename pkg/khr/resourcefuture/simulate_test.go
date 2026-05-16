package resourcefuture

import (
	"testing"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/host"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/lanediscovery"
)

func testSimCfg() *host.Config {
	cfg := &host.Config{}
	cfg.Spec.HostID = "test-host"
	cfg.Spec.SandboxMode = true
	cfg.Spec.LaneDiscoveryEnabled = true
	cfg.Spec.ResourceFutureSimulationEnabled = true
	return cfg
}

func TestBuildForecastsLinuxSandbox(t *testing.T) {
	cells := []lanediscovery.DiscoveredCell{{
		Ref: "khr-runtime-sandbox/Cell/target", Namespace: "khr-runtime-sandbox",
		Name: "target", OSFamily: "linux", VMType: "container", Running: true,
	}}
	ports := []lanediscovery.DiscoveredResourcePort{{
		Ref: "khr-runtime-sandbox/ResourcePort/target-port",
		CellRef: cells[0].Ref, Lane: lanediscovery.LaneLinuxContainerCgroup,
		Classification: lanediscovery.ClassificationLiveInPlaceCapable,
		LiveScaleCapabilityObserved: true,
	}}
	plans, _, blocked, _, live, restart, bundle := buildForecasts(cells, ports, PolicyContext{})
	if len(plans) == 0 || len(bundle.CPUScale) == 0 {
		t.Fatalf("plans=%d cpu=%d", len(plans), len(bundle.CPUScale))
	}
	if !live[0].Eligible {
		t.Fatalf("live=%+v", live)
	}
	if restart[0].Required {
		t.Fatal("sandbox should not require restart")
	}
	if len(blocked) > 0 {
		t.Fatalf("unexpected blocked=%+v", blocked)
	}
}

func TestBuildForecastsWindowsBlockedRAM(t *testing.T) {
	cells := []lanediscovery.DiscoveredCell{{
		Ref: "karl/Cell/win11", Namespace: "karl", Name: "win11",
		OSFamily: "windows", VMType: "vm", Running: true,
	}}
	ports := []lanediscovery.DiscoveredResourcePort{{
		CellRef: cells[0].Ref, Lane: lanediscovery.LaneWindowsVMSession,
		Classification: lanediscovery.ClassificationCompatibilityFallback,
	}}
	plans, _, blocked, compat, live, restart, bundle := buildForecasts(cells, ports, PolicyContext{})
	if !restart[0].Required {
		t.Fatal("windows restart expected")
	}
	if live[0].Eligible {
		t.Fatal("windows should not be live-in-place eligible")
	}
	ramBlocked := false
	for _, p := range plans {
		if p.Resource == "memory" && p.Blocked {
			ramBlocked = true
		}
	}
	if !ramBlocked {
		t.Fatalf("plans=%+v", plans)
	}
	if len(compat) == 0 || len(bundle.WindowsCompatibility) == 0 {
		t.Fatalf("compat=%d win=%d", len(compat), len(bundle.WindowsCompatibility))
	}
	_ = blocked
}

func TestDefaultSafetyPolicyNoMutation(t *testing.T) {
	s := DefaultSafetyPolicy()
	if !s.ReadOnly || !s.NoApply || !s.NoMutation || !s.SimulationOnly || !s.NoAutonomousOrchestration {
		t.Fatalf("s=%+v", s)
	}
}

func TestValidateGateRequiresFlag(t *testing.T) {
	cfg := testSimCfg()
	cfg.Spec.ResourceFutureSimulationEnabled = false
	g := validateGate(Options{Config: cfg, ClusterContext: "karl-metal-01@ovh"})
	if g.Allowed {
		t.Fatal("expected blocked gate")
	}
}
