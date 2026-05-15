package windowsfluidvirt

import "testing"

func TestDerivePlanningSafetyFlagsForbiddenWindowsClaims(t *testing.T) {
	target := EvaluateWindowsHyperdensityTarget(baseReadyTarget())
	lease := PrepareWindowsFluidResourceLease(target, baseCombinedLease())
	slate := BuildWindowsFluidActionSlate(target, lease)

	t.Run("vcpu hotplug on target", func(t *testing.T) {
		tgt := baseReadyTarget()
		tgt.VCPUHotplugRequested = true
		tgt = EvaluateWindowsHyperdensityTarget(tgt)
		s := DerivePlanningSafetyWithTargetLease(slate, tgt, lease)
		if !s.VCPUHotplugClaim || s.PlanningOnly {
			t.Fatalf("expected hotplug claim and non-planning-only, got %+v", s)
		}
	})

	t.Run("logical cpu scaling on target", func(t *testing.T) {
		tgt := baseReadyTarget()
		tgt.LogicalCPUScalingClaimed = true
		tgt = EvaluateWindowsHyperdensityTarget(tgt)
		s := DerivePlanningSafetyWithTargetLease(slate, tgt, lease)
		if !s.LogicalCPUScalingClaim || s.PlanningOnly {
			t.Fatalf("expected logical scaling claim and non-planning-only, got %+v", s)
		}
	})

	t.Run("pool scaling on target", func(t *testing.T) {
		tgt := baseReadyTarget()
		tgt.PoolScalingRequested = true
		tgt = EvaluateWindowsHyperdensityTarget(tgt)
		s := DerivePlanningSafetyWithTargetLease(slate, tgt, lease)
		if !s.PoolScalingClaim || s.PlanningOnly {
			t.Fatalf("expected pool scaling claim and non-planning-only, got %+v", s)
		}
	})

	t.Run("pool-scaling mechanism on lease", func(t *testing.T) {
		l := baseCombinedLease()
		l.RequestedMechanism = "pool-scaling"
		l = PrepareWindowsFluidResourceLease(target, l)
		s := DerivePlanningSafetyWithTargetLease(slate, target, l)
		if !s.PoolScalingClaim || s.PlanningOnly {
			t.Fatalf("expected pool scaling via lease mechanism, got %+v", s)
		}
	})
}

func TestWindowsLaneDoesNotAssertProductionHyperdensityReadyWhenBlocked(t *testing.T) {
	tgt := baseReadyTarget()
	tgt.PoolScalingRequested = true
	tgt = EvaluateWindowsHyperdensityTarget(tgt)
	if tgt.HyperdensityReady {
		t.Fatal("HyperdensityReady must stay false when lane blockers present (no production-ready claim)")
	}
}
