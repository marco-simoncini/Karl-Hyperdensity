//go:build ignore

// Run: go run scripts/gen_grandepadre_evidence_fixtures.go
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/crdv1alpha1"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/discovery"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/evidence"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/evidence/integrity"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/telemetry"
	"gopkg.in/yaml.v3"
)

func cell() *crdv1alpha1.Cell {
	return &crdv1alpha1.Cell{
		APIVersion: "runtime.karl.io/v1alpha1",
		Kind:       "Cell",
		Metadata: crdv1alpha1.ObjectMeta{Name: "demo-cell", Namespace: "karl-sandbox"},
	}
}

func main() {
	os.Setenv("KHR_TEST_COLLECTED_AT", "2026-05-15T16:00:00Z")
	os.Setenv("KHR_TEST_TELEMETRY_NOW", "2026-05-15T15:00:00Z")
	ready := buildBundle(true, nil)
	blocked := buildBundle(false, []string{"cgroup path not resolved for telemetry"})
	lowConf := buildBundle(true, nil)
	lowConf.EvidenceSummary.Confidence = "low"
	lowConf.EvidenceSummary.ReadyForGrandePadre = true

	readyDoc := ingestDoc("gp-ready", "fixture-ready", ready)
	blockedDoc := ingestDoc("gp-blocked", "fixture-blocked", blocked)
	lowDoc := ingestDoc("gp-low", "fixture-low", lowConf)

	rdCanon, err := integrity.CanonicalJSON(ready)
	if err != nil {
		panic(err)
	}
	rdSHA := integrity.SHA256Hex(rdCanon)
	manDev := integrity.BuildManifest(ready.AgentID, "fixture-devonly", "local-dev", rdCanon, rdSHA, "AQIDBAUGBw==", "ed25519")
	devOnlyDoc := ingestDocWithManifest("gp-devonly", "fixture-devonly", ready, manDev)

	for name, doc := range map[string]map[string]interface{}{
		"examples/grandepadre/evidence-store/ingest-request-ready.yaml":   readyDoc,
		"examples/grandepadre/evidence-store/ingest-request-blocked.yaml": blockedDoc,
	} {
		writeYAMLFixture(name, doc)
	}
	recDir := "examples/grandepadre/recommendation"
	if err := os.MkdirAll(recDir, 0o755); err != nil {
		panic(err)
	}
	writeYAMLFixture(filepath.Join(recDir, "recommendation-input-ready.yaml"), readyDoc)
	writeYAMLFixture(filepath.Join(recDir, "recommendation-input-blocked.yaml"), blockedDoc)
	writeYAMLFixture(filepath.Join(recDir, "recommendation-input-low-confidence.yaml"), lowDoc)
	writeYAMLFixture(filepath.Join(recDir, "recommendation-input-devonly.yaml"), devOnlyDoc)
	failDoc := copyYAMLDoc(readyDoc)
	failDoc["spec"].(map[string]interface{})["digest"] = "0000000000000000000000000000000000000000000000000000000000000000"
	writeYAMLFixture(filepath.Join(recDir, "recommendation-input-integrity-failed.yaml"), failDoc)
	fmt.Println("wrote recommendation fixtures under", recDir)
}

func copyYAMLDoc(m map[string]interface{}) map[string]interface{} {
	b, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(b, &out); err != nil {
		panic(err)
	}
	return out
}

func writeYAMLFixture(name string, doc map[string]interface{}) {
	b, err := yaml.Marshal(doc)
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile(name, b, 0o644); err != nil {
		panic(err)
	}
	fmt.Println("wrote", name)
}

func buildBundle(ready bool, blocked []string) *evidence.CollectEvidenceBundle {
	disc := &discovery.CgroupDiscoveryOutput{
		AgentID:            "khr-agent-fixture",
		CgroupVersion:      "v2",
		DiscoveryMode:      "read-only",
		ScannedRoot:        "/tmp/khr-evidence-example-root",
		SelectedPath:       "/tmp/khr-evidence-example-root/karl.slice/karl-shell-dev-linux-systemd-001.scope",
		BlockedReasons:     []string{},
		Warnings:           []string{},
		MutationsForbidden: true,
	}
	tel := evidence.TelemetrySnapshot{
		TelemetryMode: "read-only",
		CgroupPath:    "/tmp/khr-evidence-example-root/karl.slice/karl-shell-dev-linux-systemd-001.scope",
		Metrics: telemetry.MetricsBundle{
			CPUStat:       map[string]int64{"usage_usec": 100},
			MemoryCurrent: "8192",
		},
		Evidence: telemetry.Evidence{
			ObservedAt:     "2026-05-15T15:00:00Z",
			Source:         "cgroup-v2",
			Confidence:     "high",
			Warnings:       []string{},
			BlockedReasons: []string{},
		},
		MutationsForbidden: true,
	}
	dry := evidence.DryRunSkippedPayload("no lease or resource port inputs provided for optional dry-run")
	b := evidence.BuildCollectEvidenceBundle("0.0.1-sprint13", "khr-agent-fixture", cell(), disc, tel, dry, "")
	if !ready {
		b.EvidenceSummary.ReadyForGrandePadre = false
		b.EvidenceSummary.BlockedReasons = blocked
		b.EvidenceSummary.RecommendedNextAction = "Resolve blocking conditions before Grande Padre promotion."
	}
	return b
}

func ingestDoc(name, artifact string, bundle *evidence.CollectEvidenceBundle) map[string]interface{} {
	canonical, err := integrity.CanonicalJSON(bundle)
	if err != nil {
		panic(err)
	}
	sha := integrity.SHA256Hex(canonical)
	man := integrity.BuildManifest(bundle.AgentID, artifact, "none", canonical, sha, "", "")
	var bundleObj, manObj interface{}
	bj, _ := json.Marshal(bundle)
	json.Unmarshal(bj, &bundleObj)
	mj, _ := json.Marshal(man)
	json.Unmarshal(mj, &manObj)
	return map[string]interface{}{
		"apiVersion": "hyperdensity.karl.io/v1alpha1",
		"kind":       "EvidenceIngestRequest",
		"metadata": map[string]interface{}{
			"name":      name,
			"namespace": "karl-sandbox",
		},
		"spec": map[string]interface{}{
			"artifactId": artifact,
			"bundle":     bundleObj,
			"manifest":   manObj,
			"digest":     sha,
			"source": map[string]interface{}{
				"agentId": bundle.AgentID,
			},
			"policy": map[string]interface{}{
				"requireDigestMatch":     true,
				"allowUnsigned":          true,
				"allowLocalDevSignature": false,
				"maxBundleBytes":         int64(1048576),
			},
			"dryRunOnly": true,
			"ttlSeconds": int64(86400),
		},
	}
}

func ingestDocWithManifest(name, artifact string, bundle *evidence.CollectEvidenceBundle, man integrity.ArtifactManifest) map[string]interface{} {
	canonical, err := integrity.CanonicalJSON(bundle)
	if err != nil {
		panic(err)
	}
	sha := integrity.SHA256Hex(canonical)
	if man.BundleSha256 != sha {
		panic("manifest bundleSha256 mismatch canonical digest")
	}
	var bundleObj, manObj interface{}
	bj, _ := json.Marshal(bundle)
	json.Unmarshal(bj, &bundleObj)
	mj, _ := json.Marshal(man)
	json.Unmarshal(mj, &manObj)
	return map[string]interface{}{
		"apiVersion": "hyperdensity.karl.io/v1alpha1",
		"kind":       "EvidenceIngestRequest",
		"metadata": map[string]interface{}{
			"name":      name,
			"namespace": "karl-sandbox",
		},
		"spec": map[string]interface{}{
			"artifactId": artifact,
			"bundle":     bundleObj,
			"manifest":   manObj,
			"digest":     sha,
			"source": map[string]interface{}{
				"agentId": bundle.AgentID,
			},
			"policy": map[string]interface{}{
				"requireDigestMatch":     true,
				"allowUnsigned":          true,
				"allowLocalDevSignature": true,
				"maxBundleBytes":         int64(1048576),
			},
			"dryRunOnly": true,
			"ttlSeconds": int64(86400),
		},
	}
}
