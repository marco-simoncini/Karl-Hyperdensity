// Command khr-action-approval manages read-only operator action approval evidence (KHR-W).
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/actionapproval"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/certregistry"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/policygates"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/resourcefuture"
)

func main() {
	cmd := flag.String("cmd", "", "generate|approve|reject|expire|check")
	simPath := flag.String("simulation", "", "resourcefuture simulation JSON")
	regPath := flag.String("registry", "", "certification registry JSON")
	approvalPath := flag.String("approval", "", "action approval JSON path")
	certRef := flag.String("cert-ref", "", "certification evidence ref")
	out := flag.String("out", "", "output path (default stdout)")
	by := flag.String("by", "", "operator id for approve/reject")
	reason := flag.String("reason", "", "reject reason")
	ttl := flag.Int64("ttl-seconds", actionapproval.DefaultTTLSeconds, "pending approval TTL")
	sprint := flag.String("sprint", "KHR-W", "sprint label")
	flag.Parse()

	if *cmd == "" {
		fmt.Fprintln(os.Stderr, "-cmd is required: generate|approve|reject|expire|check")
		os.Exit(2)
	}
	now := time.Now().UTC()

	switch *cmd {
	case "generate":
		if *simPath == "" || *regPath == "" {
			fmt.Fprintln(os.Stderr, "-simulation and -registry required")
			os.Exit(2)
		}
		sim, err := loadSimulation(*simPath)
		if err != nil {
			fatal(err)
		}
		reg, err := certregistry.LoadJSON(*regPath)
		if err != nil {
			fatal(err)
		}
		ref := *certRef
		if ref == "" {
			ref = *regPath
		}
		pending, err := actionapproval.GeneratePending(actionapproval.GenerateInput{
			Simulation: sim, Registry: reg, Gates: policygates.DefaultNativeLiveGates(),
			CertificationRef: ref, TTLSeconds: *ttl, Now: now,
		})
		if err != nil {
			fatal(err)
		}
		bundle := actionapproval.Bundle{
			WorkflowID: actionapproval.WorkflowID, Sprint: *sprint,
			GeneratedAt: now.Format(time.RFC3339), ReadOnly: true,
			NoAutonomousOrchestration: true, Approvals: pending,
		}
		emit(outPath(*out), bundle)

	case "approve", "reject", "expire", "check":
		if *approvalPath == "" {
			fmt.Fprintln(os.Stderr, "-approval required")
			os.Exit(2)
		}
		a, err := actionapproval.LoadJSON(*approvalPath)
		if err != nil {
			fatal(err)
		}
		var reg *certregistry.Registry
		if *regPath != "" {
			r, err := certregistry.LoadJSON(*regPath)
			if err != nil {
				fatal(err)
			}
			reg = &r
		}
		gates := policygates.DefaultNativeLiveGates()

		switch *cmd {
		case "approve":
			if *by == "" {
				fmt.Fprintln(os.Stderr, "-by required for approve")
				os.Exit(2)
			}
			a, err = actionapproval.Approve(a, reg, gates, *by, now)
		case "reject":
			if *by == "" {
				*by = "operator"
			}
			a, err = actionapproval.Reject(a, *by, *reason, now)
		case "expire":
			a = actionapproval.SimulateExpire(a, now)
		case "check":
			err = actionapproval.CanApprove(a, reg, gates, now)
			if err != nil {
				fmt.Fprintf(os.Stderr, "check failed: %v\n", err)
				os.Exit(1)
			}
			emit(outPath(*out), map[string]any{"ok": true, "actionId": a.ActionID})
			return
		}
		if err != nil {
			fatal(err)
		}
		emit(outPath(*out), a)
	default:
		fmt.Fprintf(os.Stderr, "unknown cmd %q\n", *cmd)
		os.Exit(2)
	}
}

func loadSimulation(path string) (resourcefuture.SimulationResult, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return resourcefuture.SimulationResult{}, err
	}
	var sim resourcefuture.SimulationResult
	return sim, json.Unmarshal(data, &sim)
}

func outPath(p string) string {
	if p == "" {
		return "-"
	}
	return p
}

func emit(path string, v any) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fatal(err)
	}
	if path == "-" {
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
