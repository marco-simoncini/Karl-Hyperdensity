// Command khr-cert-registry generates read-only certification registry JSON (KHR-V).
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/certregistry"
)

func main() {
	certPath := flag.String("cert", "", "path to certification-summary.json")
	evidenceRef := flag.String("evidence-ref", "", "evidence reference path or URI")
	validFor := flag.Int64("valid-for-seconds", certregistry.DefaultValidForSeconds, "evidence freshness window")
	sprint := flag.String("sprint", "KHR-V", "sprint label")
	certifiedAt := flag.String("certified-at", "", "RFC3339 lastCertifiedAt (default: now)")
	out := flag.String("out", "", "output registry JSON path (default: stdout)")
	flag.Parse()
	if *certPath == "" {
		fmt.Fprintln(os.Stderr, "-cert is required")
		os.Exit(2)
	}
	summary, err := certregistry.LoadSummaryJSON(*certPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	at := time.Now().UTC()
	if *certifiedAt != "" {
		at, err = time.Parse(time.RFC3339, *certifiedAt)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
	ref := *evidenceRef
	if ref == "" {
		ref = *certPath
	}
	reg := certregistry.GenerateFromSummary(*sprint, summary, ref, *validFor, at)
	data, err := json.MarshalIndent(reg, "", "  ")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if *out == "" {
		fmt.Println(string(data))
		return
	}
	if err := os.WriteFile(*out, append(data, '\n'), 0o644); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
