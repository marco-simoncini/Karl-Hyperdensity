// Command khr-linux-agent is a dry-run-first Linux MVP skeleton (Sprint 5+; discovery Sprint 7).
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	gpdevidence "github.com/marco-simoncini/Karl-Hyperdensity/pkg/grandepadre/evidence"
	gprec "github.com/marco-simoncini/Karl-Hyperdensity/pkg/grandepadre/recommendation"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/agent"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/crdv1alpha1"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/evidenceingest"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/resourcelease"
)

func main() {
	mode := flag.String("mode", "", "one of: validate-config, dry-run, print-capabilities, discover-cgroups, read-telemetry, collect-evidence, prepare-ingest-request, index-evidence-local, recommend-actions-local")
	configPath := flag.String("config", "", "path to agent YAML/JSON config")
	cellInputPath := flag.String("cell-input", "", "optional path to Cell JSON (discover-cgroups, read-telemetry); required for collect-evidence")
	cgroupRoot := flag.String("cgroup-root", "", "optional cgroup scan root (default /sys/fs/cgroup) for discover-cgroups")
	allowPathPrefix := flag.String("allow-path-prefix", "", "optional path prefix policy (discover-cgroups, read-telemetry, collect-evidence)")
	telemetryCgroupPath := flag.String("cgroup-path", "", "resolved cgroup directory to sample (read-telemetry)")
	telemetryOutputPath := flag.String("telemetry-output", "", "optional path to write the same JSON as stdout (read-telemetry)")
	evidenceOutputPath := flag.String("evidence-output", "", "collect-evidence: optional write path for bundle JSON; prepare-ingest-request: read path when -bundle-input is omitted")
	evidenceManifestOut := flag.String("evidence-manifest-output", "", "collect-evidence: write manifest JSON; prepare-ingest-request: read path when -manifest-input is omitted")
	evidenceDigestOut := flag.String("evidence-digest-output", "", "collect-evidence: write digest line; prepare-ingest-request: read path when -digest-input is omitted")
	signingMode := flag.String("signing-mode", "none", "none|local-dev: local integrity signing for collect-evidence (local-dev is not production security)")
	signingKeyFile := flag.String("signing-key-file", "", "Ed25519 private key PEM when -signing-mode=local-dev (collect-evidence)")
	artifactID := flag.String("artifact-id", "", "optional artifact id recorded in evidence manifest (collect-evidence)")
	leasePath := flag.String("lease-input", "", "path to ResourceLease JSON (dry-run)")
	portPath := flag.String("resource-port-input", "", "path to ResourcePort JSON (dry-run)")
	cellCtxPath := flag.String("cell-context", "", "optional path to CellContext JSON (dry-run)")
	allowUnsafe := flag.Bool("allow-unsafe-apply", false, "non-operational in Sprint 6: emits audit only; never enables writes")
	cpuDelta := flag.String("cpu-delta", "", "optional cpu.max delta string for envelope dry-run plan (simulation only)")
	memDelta := flag.String("memory-delta", "", "optional memory.max delta string for envelope dry-run plan (simulation only)")
	bundleInputPath := flag.String("bundle-input", "", "prepare-ingest-request: path to collect-evidence bundle JSON")
	manifestInputPath := flag.String("manifest-input", "", "prepare-ingest-request: path to artifact manifest JSON")
	digestInputPath := flag.String("digest-input", "", "prepare-ingest-request: path to digest text file")
	ingestRequestOut := flag.String("ingest-request-output", "", "prepare-ingest-request: path to write EvidenceIngestRequest YAML/JSON")
	ingestRequestFmt := flag.String("ingest-request-format", "yaml", "prepare-ingest-request: yaml or json")
	ingestDryRunOnly := flag.Bool("dry-run-only", false, "prepare-ingest-request: set spec.dryRunOnly (simulation-only ingest)")
	ingestReqNamespace := flag.String("ingest-request-namespace", "", "prepare-ingest-request: metadata.namespace (default from bundle cellRef or karl-sandbox)")
	ingestReqName := flag.String("ingest-request-name", "", "prepare-ingest-request: metadata.name (default derived)")
	sourceNodeName := flag.String("source-node-name", "", "prepare-ingest-request: spec.source.nodeName")
	sourceHostID := flag.String("source-host-id", "", "prepare-ingest-request: spec.source.hostId")
	sourceTenant := flag.String("source-tenant", "", "prepare-ingest-request: spec.source.tenant")
	requireDigestMatch := flag.Bool("require-digest-match", true, "prepare-ingest-request: spec.policy.requireDigestMatch")
	allowUnsigned := flag.Bool("allow-unsigned", true, "prepare-ingest-request: spec.policy.allowUnsigned")
	allowLocalDevSignature := flag.Bool("allow-local-dev-signature", false, "prepare-ingest-request: spec.policy.allowLocalDevSignature (auto-enabled when manifest signingMode=local-dev)")
	indexOutputPath := flag.String("index-output", "", "index-evidence-local: optional path to write index report JSON")
	indexQuery := flag.String("query", "", "index-evidence-local: optional ready|blocked|by-confidence|by-cell")
	indexCellNamespace := flag.String("cell-namespace", "", "index-evidence-local: cell namespace for -query=by-cell")
	indexCellName := flag.String("cell-name", "", "index-evidence-local: cell name for -query=by-cell")
	indexConfidence := flag.String("confidence", "", "index-evidence-local: low|medium|high for -query=by-confidence")
	unsignedDigestTrust := flag.String("unsigned-digest-trust", "verified", "index-evidence-local / recommend-actions-local: verified|unsigned (IntegrityVerified vs Unsigned for digest-only bundles; see docs)")
	var ingestRequestInputs stringList
	flag.Var(&ingestRequestInputs, "ingest-request-input", "index-evidence-local: one file. recommend-actions-local: repeat flag for multiple files")
	ingestRequestDir := flag.String("ingest-request-dir", "", "recommend-actions-local: directory of EvidenceIngestRequest yaml/json (*.yaml, *.yml, *.json)")
	recommendationOutputPath := flag.String("recommendation-output", "", "recommend-actions-local: optional path to write recommendation JSON")
	recommendTenant := flag.String("tenant", "", "recommend-actions-local: optional filter on cell namespace")
	recommendDryRunOnly := flag.Bool("recommend-dry-run-only", true, "recommend-actions-local: dryRunOnly on each recommendation (default true; distinct from prepare -dry-run-only)")
	recommendationGeneratedAt := flag.String("recommendation-generated-at", "", "recommend-actions-local: optional RFC3339 clock for deterministic JSON (tests)")
	flag.Parse()

	out := map[string]interface{}{
		"tool":    "khr-linux-agent",
		"version": agent.AgentVersion,
		"mode":    *mode,
	}

	if *mode == "" {
		out["error"] = "missing -mode"
		emit(out, 2)
	}

	switch *mode {
	case "validate-config":
		if *configPath == "" {
			out["error"] = "missing -config"
			emit(out, 2)
		}
		cfg, err := agent.LoadConfig(*configPath)
		if err != nil {
			out["error"] = err.Error()
			emit(out, 2)
		}
		errs := agent.ValidateConfig(cfg)
		out["valid"] = len(errs) == 0
		out["validationErrors"] = errs
		emit(out, boolExit(len(errs) > 0))

	case "print-capabilities":
		if *configPath == "" {
			out["error"] = "missing -config"
			emit(out, 2)
		}
		cfg, err := agent.LoadConfig(*configPath)
		if err != nil {
			out["error"] = err.Error()
			emit(out, 2)
		}
		if errs := agent.ValidateConfig(cfg); len(errs) > 0 {
			out["validationErrors"] = errs
			out["error"] = "invalid config"
			emit(out, 2)
		}
		b, err := agent.PrintCapabilitiesJSON(cfg, *allowUnsafe)
		if err != nil {
			out["error"] = err.Error()
			emit(out, 2)
		}
		os.Stdout.Write(b)
		os.Stdout.Write([]byte("\n"))
		os.Exit(0)

	case "dry-run":
		if *configPath == "" || *leasePath == "" || *portPath == "" {
			out["error"] = "dry-run requires -config, -lease-input, and -resource-port-input"
			emit(out, 2)
		}
		cfg, err := agent.LoadConfig(*configPath)
		if err != nil {
			out["error"] = err.Error()
			emit(out, 2)
		}
		if errs := agent.ValidateConfig(cfg); len(errs) > 0 {
			out["validationErrors"] = errs
			out["error"] = "invalid config"
			emit(out, 2)
		}
		leaseRaw, err := os.ReadFile(*leasePath)
		if err != nil {
			out["error"] = err.Error()
			emit(out, 2)
		}
		portRaw, err := os.ReadFile(*portPath)
		if err != nil {
			out["error"] = err.Error()
			emit(out, 2)
		}
		var ctx *resourcelease.CellContext
		if *cellCtxPath != "" {
			raw, err := os.ReadFile(*cellCtxPath)
			if err != nil {
				out["error"] = err.Error()
				emit(out, 2)
			}
			ctx = &resourcelease.CellContext{}
			if err := json.Unmarshal(raw, ctx); err != nil {
				out["error"] = err.Error()
				emit(out, 2)
			}
		}
		cliOut, err := agent.RunDryRunCLI(leaseRaw, portRaw, ctx, *allowUnsafe, *cpuDelta, *memDelta)
		if err != nil {
			out["error"] = err.Error()
			emit(out, 2)
		}
		b, err := json.MarshalIndent(cliOut, "", "  ")
		if err != nil {
			out["error"] = err.Error()
			emit(out, 2)
		}
		os.Stdout.Write(b)
		os.Stdout.Write([]byte("\n"))
		os.Exit(0)

	case "discover-cgroups":
		if *configPath == "" {
			out["error"] = "discover-cgroups requires -config"
			emit(out, 2)
		}
		cfg, err := agent.LoadConfig(*configPath)
		if err != nil {
			out["error"] = err.Error()
			emit(out, 2)
		}
		if errs := agent.ValidateConfig(cfg); len(errs) > 0 {
			out["validationErrors"] = errs
			out["error"] = "invalid config"
			emit(out, 2)
		}
		var cell *crdv1alpha1.Cell
		if *cellInputPath != "" {
			raw, err := os.ReadFile(*cellInputPath)
			if err != nil {
				out["error"] = err.Error()
				emit(out, 2)
			}
			cell = &crdv1alpha1.Cell{}
			if err := json.Unmarshal(raw, cell); err != nil {
				out["error"] = err.Error()
				emit(out, 2)
			}
		}
		cliOut := agent.RunDiscoverCgroupsCLI(cfg, cell, *cgroupRoot, *allowPathPrefix)
		b, err := json.MarshalIndent(cliOut, "", "  ")
		if err != nil {
			out["error"] = err.Error()
			emit(out, 2)
		}
		os.Stdout.Write(b)
		os.Stdout.Write([]byte("\n"))
		os.Exit(0)

	case "read-telemetry":
		if *configPath == "" || *telemetryCgroupPath == "" {
			out["error"] = "read-telemetry requires -config and -cgroup-path"
			emit(out, 2)
		}
		cfg, err := agent.LoadConfig(*configPath)
		if err != nil {
			out["error"] = err.Error()
			emit(out, 2)
		}
		if errs := agent.ValidateConfig(cfg); len(errs) > 0 {
			out["validationErrors"] = errs
			out["error"] = "invalid config"
			emit(out, 2)
		}
		var cell *crdv1alpha1.Cell
		if *cellInputPath != "" {
			raw, err := os.ReadFile(*cellInputPath)
			if err != nil {
				out["error"] = err.Error()
				emit(out, 2)
			}
			cell = &crdv1alpha1.Cell{}
			if err := json.Unmarshal(raw, cell); err != nil {
				out["error"] = err.Error()
				emit(out, 2)
			}
		}
		cliOut, err := agent.RunReadTelemetryCLI(cfg, *telemetryCgroupPath, *allowPathPrefix, cell, *telemetryOutputPath)
		if err != nil {
			out["error"] = err.Error()
			emit(out, 2)
		}
		b, err := json.MarshalIndent(cliOut, "", "  ")
		if err != nil {
			out["error"] = err.Error()
			emit(out, 2)
		}
		os.Stdout.Write(b)
		os.Stdout.Write([]byte("\n"))
		os.Exit(0)

	case "collect-evidence":
		if *configPath == "" || *cellInputPath == "" {
			out["error"] = "collect-evidence requires -config and -cell-input"
			emit(out, 2)
		}
		cfg, err := agent.LoadConfig(*configPath)
		if err != nil {
			out["error"] = err.Error()
			emit(out, 2)
		}
		if errs := agent.ValidateConfig(cfg); len(errs) > 0 {
			out["validationErrors"] = errs
			out["error"] = "invalid config"
			emit(out, 2)
		}
		cellRaw, err := os.ReadFile(*cellInputPath)
		if err != nil {
			out["error"] = err.Error()
			emit(out, 2)
		}
		cell := &crdv1alpha1.Cell{}
		if err := json.Unmarshal(cellRaw, cell); err != nil {
			out["error"] = err.Error()
			emit(out, 2)
		}
		var leaseRaw, portRaw []byte
		if *leasePath != "" {
			leaseRaw, err = os.ReadFile(*leasePath)
			if err != nil {
				out["error"] = err.Error()
				emit(out, 2)
			}
		}
		if *portPath != "" {
			portRaw, err = os.ReadFile(*portPath)
			if err != nil {
				out["error"] = err.Error()
				emit(out, 2)
			}
		}
		var ctx *resourcelease.CellContext
		if *cellCtxPath != "" {
			raw, err := os.ReadFile(*cellCtxPath)
			if err != nil {
				out["error"] = err.Error()
				emit(out, 2)
			}
			ctx = &resourcelease.CellContext{}
			if err := json.Unmarshal(raw, ctx); err != nil {
				out["error"] = err.Error()
				emit(out, 2)
			}
		}
		var integ *agent.CollectEvidenceIntegrityOpts
		man := strings.TrimSpace(*evidenceManifestOut)
		dig := strings.TrimSpace(*evidenceDigestOut)
		sign := strings.TrimSpace(*signingMode)
		if man != "" || dig != "" || (sign != "" && sign != "none") {
			integ = &agent.CollectEvidenceIntegrityOpts{
				ManifestOutputPath: *evidenceManifestOut,
				DigestOutputPath:   *evidenceDigestOut,
				SigningMode:        *signingMode,
				SigningKeyFile:     *signingKeyFile,
				ArtifactID:         *artifactID,
			}
		}
		cliOut, err := agent.RunCollectEvidenceCLI(cfg, cell, *cgroupRoot, *allowPathPrefix, leaseRaw, portRaw, ctx, *evidenceOutputPath, *allowUnsafe, *cpuDelta, *memDelta, integ)
		if err != nil {
			out["error"] = err.Error()
			emit(out, 2)
		}
		b, err := json.MarshalIndent(cliOut, "", "  ")
		if err != nil {
			out["error"] = err.Error()
			emit(out, 2)
		}
		os.Stdout.Write(b)
		os.Stdout.Write([]byte("\n"))
		os.Exit(0)

	case "prepare-ingest-request":
		if strings.TrimSpace(*ingestRequestOut) == "" {
			out["error"] = "prepare-ingest-request requires -ingest-request-output"
			emit(out, 2)
		}
		bundlePath := strings.TrimSpace(*bundleInputPath)
		if bundlePath == "" {
			bundlePath = strings.TrimSpace(*evidenceOutputPath)
		}
		manifestPath := strings.TrimSpace(*manifestInputPath)
		if manifestPath == "" {
			manifestPath = strings.TrimSpace(*evidenceManifestOut)
		}
		digestPath := strings.TrimSpace(*digestInputPath)
		if digestPath == "" {
			digestPath = strings.TrimSpace(*evidenceDigestOut)
		}
		opts := evidenceingest.DefaultPrepareOptions()
		opts.Format = *ingestRequestFmt
		opts.DryRunOnly = *ingestDryRunOnly
		opts.Namespace = *ingestReqNamespace
		opts.Name = *ingestReqName
		opts.NodeName = *sourceNodeName
		opts.HostID = *sourceHostID
		opts.Tenant = *sourceTenant
		opts.RequireDigestMatch = *requireDigestMatch
		opts.AllowUnsigned = *allowUnsigned
		opts.AllowLocalDevSignature = *allowLocalDevSignature
		b, err := evidenceingest.PrepareIngestRequest(bundlePath, manifestPath, digestPath, opts)
		if err != nil {
			out["error"] = err.Error()
			emit(out, 2)
		}
		if err := evidenceingest.WriteFile(*ingestRequestOut, b); err != nil {
			out["error"] = err.Error()
			emit(out, 2)
		}
		os.Exit(0)

	case "index-evidence-local":
		paths := []string(ingestRequestInputs)
		if len(paths) == 0 {
			out["error"] = "index-evidence-local requires -ingest-request-input"
			emit(out, 2)
		}
		path := strings.TrimSpace(paths[0])
		raw, err := os.ReadFile(path)
		if err != nil {
			out["error"] = err.Error()
			emit(out, 2)
		}
		unsignedPol := gpdevidence.UnsignedDigestAsIntegrityVerified
		if strings.EqualFold(strings.TrimSpace(*unsignedDigestTrust), "unsigned") {
			unsignedPol = gpdevidence.UnsignedDigestAsUnsigned
		}
		rep, err := gpdevidence.RunLocalIndex(gpdevidence.NewStore(), raw, gpdevidence.LocalIndexParams{
			UnsignedLabel: unsignedPol,
			Query:         gpdevidence.ParseQueryKind(*indexQuery),
			CellNamespace: *indexCellNamespace,
			CellName:      *indexCellName,
			Confidence:    *indexConfidence,
		})
		if err != nil {
			out["error"] = err.Error()
			emit(out, 2)
		}
		b, err := json.MarshalIndent(rep, "", "  ")
		if err != nil {
			out["error"] = err.Error()
			emit(out, 2)
		}
		b = append(b, '\n')
		os.Stdout.Write(b)
		if p := strings.TrimSpace(*indexOutputPath); p != "" {
			if err := os.WriteFile(p, b, 0o600); err != nil {
				out["error"] = err.Error()
				emit(out, 2)
			}
		}
		os.Exit(0)

	case "recommend-actions-local":
		paths, err := gprec.CollectIngestPaths([]string(ingestRequestInputs), strings.TrimSpace(*ingestRequestDir))
		if err != nil {
			out["error"] = err.Error()
			emit(out, 2)
		}
		unsignedPol := gpdevidence.UnsignedDigestAsIntegrityVerified
		if strings.EqualFold(strings.TrimSpace(*unsignedDigestTrust), "unsigned") {
			unsignedPol = gpdevidence.UnsignedDigestAsUnsigned
		}
		s := gpdevidence.NewStore()
		if err := gprec.IngestAllIntoStore(s, paths, unsignedPol); err != nil {
			out["error"] = err.Error()
			emit(out, 2)
		}
		var genAt time.Time
		if ts := strings.TrimSpace(*recommendationGeneratedAt); ts != "" {
			t, err := time.Parse(time.RFC3339, ts)
			if err != nil {
				out["error"] = fmt.Sprintf("recommendation-generated-at: %v", err)
				emit(out, 2)
			}
			genAt = t.UTC()
		}
		rep := gprec.BuildRecommendationReport(s, gprec.EngineOptions{
			Tenant:        strings.TrimSpace(*recommendTenant),
			DryRunOnly:    *recommendDryRunOnly,
			UnsignedLabel: unsignedPol,
			GeneratedAt:   genAt,
		})
		b, err := json.MarshalIndent(rep, "", "  ")
		if err != nil {
			out["error"] = err.Error()
			emit(out, 2)
		}
		b = append(b, '\n')
		os.Stdout.Write(b)
		if p := strings.TrimSpace(*recommendationOutputPath); p != "" {
			if err := os.WriteFile(p, b, 0o600); err != nil {
				out["error"] = err.Error()
				emit(out, 2)
			}
		}
		os.Exit(0)

	default:
		out["error"] = fmt.Sprintf("unknown mode %q", *mode)
		emit(out, 2)
	}
}

func emit(v map[string]interface{}, code int) {
	b, _ := json.MarshalIndent(v, "", "  ")
	os.Stdout.Write(b)
	os.Stdout.Write([]byte("\n"))
	os.Exit(code)
}

func boolExit(hasErr bool) int {
	if hasErr {
		return 2
	}
	return 0
}

// stringList supports repeated -ingest-request-input flags (flag.Value).
type stringList []string

func (s *stringList) String() string {
	if s == nil || len(*s) == 0 {
		return ""
	}
	return strings.Join(*s, ",")
}

func (s *stringList) Set(v string) error {
	v = strings.TrimSpace(v)
	if v == "" {
		return fmt.Errorf("empty -ingest-request-input")
	}
	*s = append(*s, v)
	return nil
}
