// Command khr-provenance-validate runs KHR-Y provenance validation on artifacts.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/actionapproval"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/certregistry"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/controlgraph"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/policygates"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/provenance"
)

func main() {
	certPath := flag.String("cert", "", "certification-summary.json")
	registryPath := flag.String("registry", "", "certification registry.json")
	approvalPath := flag.String("approval", "", "action approval json")
	graphPath := flag.String("graph", "", "control graph json")
	out := flag.String("out", "", "validation summary output")
	sprint := flag.String("sprint", "KHR-Y", "sprint label")
	flag.Parse()

	now := time.Now().UTC()
	sum := provenance.ValidationSummary{
		ModelID: provenance.ModelID, Sprint: *sprint,
		ValidatedAt: now.Format(time.RFC3339),
		ReadOnly: true, NoAutonomousOrchestration: true,
		NoMutation: true, NoApply: true,
		ProvenanceState: provenance.StateTrusted,
		LineageIntegrity: true, RegistryIntegrity: true,
		ApprovalProvenanceValid: true, GraphLineageVerified: true,
	}

	var certBytes []byte
	if *certPath != "" {
		var err error
		certBytes, err = os.ReadFile(*certPath)
		if err != nil {
			fatal(err)
		}
	}

	if *registryPath != "" && len(certBytes) > 0 {
		reg, err := certregistry.LoadJSON(*registryPath)
		if err != nil {
			fatal(err)
		}
		if err := certregistry.VerifyIntegrity(reg, certBytes); err != nil {
			sum.RegistryIntegrity = false
			sum.TrustWarnings = append(sum.TrustWarnings, err.Error())
		}
		if reg.Provenance.EvidenceFingerprint != "" {
			sum.EvidenceFingerprint = reg.Provenance.EvidenceFingerprint
		}
	}

	if *approvalPath != "" && *registryPath != "" {
		a, err := actionapproval.LoadJSON(*approvalPath)
		if err != nil {
			fatal(err)
		}
		reg, err := certregistry.LoadJSON(*registryPath)
		if err != nil {
			fatal(err)
		}
		entry := reg.FindByLane(a.LaneID)
		if entry != nil && a.Provenance.ProvenanceID != "" && entry.Provenance.ProvenanceID != "" {
			if err := provenance.VerifyApprovalProvenance(a.Provenance, entry.Provenance); err != nil {
				sum.ApprovalProvenanceValid = false
				sum.TrustWarnings = append(sum.TrustWarnings, err.Error())
			}
		}
		if err := actionapproval.CanApprove(a, &reg, policygates.DefaultNativeLiveGates(), now); err != nil {
			sum.ApprovalProvenanceValid = false
			sum.TrustWarnings = append(sum.TrustWarnings, err.Error())
		}
	}

	if *graphPath != "" {
		data, err := os.ReadFile(*graphPath)
		if err != nil {
			fatal(err)
		}
		var g controlgraph.Graph
		if err := json.Unmarshal(data, &g); err != nil {
			fatal(err)
		}
		gsum, _ := controlgraph.VerifyLineageIntegrity(g, now)
		sum.GraphLineageVerified = gsum.GraphLineageVerified
		sum.LineageIntegrity = gsum.LineageIntegrity
		sum.TrustWarnings = append(sum.TrustWarnings, gsum.TrustWarnings...)
	}

	if len(sum.TrustWarnings) > 0 {
		sum.ProvenanceState = provenance.StateMismatch
	}
	emit(*out, sum)
}

func emit(path string, v any) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fatal(err)
	}
	if path == "" {
		fmt.Println(string(data))
		return
	}
	if err := os.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		fatal(err)
	}
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
