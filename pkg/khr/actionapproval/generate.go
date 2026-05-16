package actionapproval

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/certregistry"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/lanediscovery"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/policygates"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/resourcefuture"
)

// GenerateInput configures pending approval generation.
type GenerateInput struct {
	Simulation       resourcefuture.SimulationResult
	Registry         certregistry.Registry
	Gates            policygates.Gates
	CertificationRef string
	TTLSeconds       int64
	Now              time.Time
}

// GeneratePending builds pending approvals for eligible ResourceFuture targets only.
func GeneratePending(in GenerateInput) ([]ActionApproval, error) {
	if in.Now.IsZero() {
		in.Now = time.Now().UTC()
	}
	if in.TTLSeconds <= 0 {
		in.TTLSeconds = DefaultTTLSeconds
	}
	expires := in.Now.Add(time.Duration(in.TTLSeconds) * time.Second).UTC().Format(time.RFC3339)
	reg := &in.Registry
	var out []ActionApproval
	seen := map[string]struct{}{}
	for _, live := range in.Simulation.LiveInPlaceEligibility {
		if !live.Eligible || live.Lane != lanediscovery.LaneNativeLive {
			continue
		}
		gate := policygates.Evaluate(live.Lane, reg, in.Gates, in.Now)
		if !gate.Eligible {
			continue
		}
		id := actionID(live.TargetRef, live.Lane)
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		ns, name := splitTargetRef(live.TargetRef)
		leaseRef := fmt.Sprintf("%s/ResourceLease/%s-scale-approval", ns, name)
		rfRef := "simulation/ResourceFuture/" + strings.ReplaceAll(live.TargetRef, "/", "-")
		prov := in.Registry.Provenance
		if e := reg.FindByLane(live.Lane); e != nil && e.Provenance.ProvenanceID != "" {
			prov = e.Provenance
		}
		out = append(out, ActionApproval{
			ActionID:          id,
			ResourceFutureRef: rfRef,
			ResourceLeaseRef:  leaseRef,
			LaneID:            live.Lane,
			CertificationRef:  in.CertificationRef,
			PolicyGateResult:  gate,
			ApprovalState:     StatePending,
			ExpiresAt:         expires,
			ReadOnly:          true,
			NoApply:           true,
			NoMutation:        true,
			NoAutonomousOrchestration: true,
			Provenance:        prov,
		})
	}
	return out, nil
}

func actionID(targetRef, lane string) string {
	h := sha256.Sum256([]byte(targetRef + "|" + lane))
	return "khr-approval-" + hex.EncodeToString(h[:8])
}

func splitTargetRef(ref string) (namespace, name string) {
	parts := strings.Split(ref, "/")
	if len(parts) >= 3 {
		return parts[0], parts[2]
	}
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return "khr-runtime-sandbox", "target"
}
