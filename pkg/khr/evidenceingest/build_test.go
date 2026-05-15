package evidenceingest

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/crdv1alpha1"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/discovery"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/evidence"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/evidence/integrity"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/telemetry"
	"gopkg.in/yaml.v3"
)

func cellDemo() *crdv1alpha1.Cell {
	return &crdv1alpha1.Cell{
		APIVersion: "runtime.karl.io/v1alpha1",
		Kind:       "Cell",
		Metadata: crdv1alpha1.ObjectMeta{
			Name:      "demo-cell",
			Namespace: "karl-sandbox",
		},
	}
}

func writeAlignedFixture(t *testing.T) string {
	t.Helper()
	t.Setenv("KHR_TEST_COLLECTED_AT", "2026-05-15T16:00:00Z")
	t.Setenv("KHR_TEST_TELEMETRY_NOW", "2026-05-15T15:00:00Z")
	dir := t.TempDir()
	disc := &discovery.CgroupDiscoveryOutput{
		AgentID:            "khr-agent-test",
		CgroupVersion:      "v2",
		DiscoveryMode:      "read-only",
		ScannedRoot:        "/cg",
		SelectedPath:       "/cg/scope",
		BlockedReasons:     []string{},
		Warnings:           []string{},
		MutationsForbidden: true,
	}
	tel := evidence.TelemetrySnapshot{
		Skipped:       false,
		TelemetryMode: "read-only",
		CgroupPath:    "/cg/scope",
		Metrics: telemetry.MetricsBundle{
			CPUStat:       map[string]int64{"usage_usec": 1},
			MemoryCurrent: "4096",
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
	dry := evidence.DryRunSkippedPayload("skip")
	bundle := evidence.BuildCollectEvidenceBundle("0.0.1-sprint12", "khr-agent-test", cellDemo(), disc, tel, dry, "")
	bj, err := json.MarshalIndent(bundle, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	bundlePath := filepath.Join(dir, "bundle.json")
	if err := os.WriteFile(bundlePath, bj, 0o600); err != nil {
		t.Fatal(err)
	}
	canonical, err := integrity.CanonicalJSON(bundle)
	if err != nil {
		t.Fatal(err)
	}
	sha := integrity.SHA256Hex(canonical)
	man := integrity.BuildManifest(bundle.AgentID, "fixture-art", "none", canonical, sha, "", "")
	mj, err := json.MarshalIndent(man, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(dir, "manifest.json")
	if err := os.WriteFile(manifestPath, mj, 0o600); err != nil {
		t.Fatal(err)
	}
	digestPath := filepath.Join(dir, "digest.txt")
	if err := os.WriteFile(digestPath, []byte(sha+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	return dir
}

func TestPrepareIngestRequestDigestAligned(t *testing.T) {
	dir := writeAlignedFixture(t)
	opts := DefaultPrepareOptions()
	out, err := PrepareIngestRequest(
		filepath.Join(dir, "bundle.json"),
		filepath.Join(dir, "manifest.json"),
		filepath.Join(dir, "digest.txt"),
		opts,
	)
	if err != nil {
		t.Fatal(err)
	}
	var doc map[string]interface{}
	if err := yaml.Unmarshal(out, &doc); err != nil {
		t.Fatal(err)
	}
	st := doc["status"].(map[string]interface{})
	if st["digestMatch"] != true {
		t.Fatalf("digestMatch: %v", st["digestMatch"])
	}
	if st["signatureStatus"] != "None" {
		t.Fatalf("signatureStatus: %v", st["signatureStatus"])
	}
	meta := doc["metadata"].(map[string]interface{})
	if _, ok := meta["annotations"]; ok {
		t.Fatalf("unexpected annotations: %v", meta["annotations"])
	}
}

func TestPrepareIngestRequestMismatchWarning(t *testing.T) {
	dir := writeAlignedFixture(t)
	if err := os.WriteFile(filepath.Join(dir, "digest.txt"), []byte("0000000000000000000000000000000000000000000000000000000000000000\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	opts := DefaultPrepareOptions()
	out, err := PrepareIngestRequest(
		filepath.Join(dir, "bundle.json"),
		filepath.Join(dir, "manifest.json"),
		filepath.Join(dir, "digest.txt"),
		opts,
	)
	if err != nil {
		t.Fatal(err)
	}
	var doc map[string]interface{}
	if err := yaml.Unmarshal(out, &doc); err != nil {
		t.Fatal(err)
	}
	meta := doc["metadata"].(map[string]interface{})
	ann := meta["annotations"].(map[string]interface{})
	raw := ann[annotationPreparationWarnings].(string)
	if !strings.Contains(raw, "digest file") {
		t.Fatalf("expected digest mismatch warning, got %q", raw)
	}
	st := doc["status"].(map[string]interface{})
	if st["digestMatch"] == true {
		t.Fatal("expected digestMatch false")
	}
}

func TestPrepareIngestRequestLocalDevAnnotation(t *testing.T) {
	dir := writeAlignedFixture(t)
	manifestPath := filepath.Join(dir, "manifest.json")
	raw, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	var man integrity.ArtifactManifest
	if err := json.Unmarshal(raw, &man); err != nil {
		t.Fatal(err)
	}
	man.SigningMode = "local-dev"
	man.SignaturePresent = true
	man.SignatureAlgorithm = integrity.SignatureAlgorithmLocalDev
	man.SignatureBase64 = "AA==" // placeholder; annotation tier is what we test
	mb, err := json.MarshalIndent(man, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(manifestPath, mb, 0o600); err != nil {
		t.Fatal(err)
	}
	opts := DefaultPrepareOptions()
	out, err := PrepareIngestRequest(
		filepath.Join(dir, "bundle.json"),
		manifestPath,
		filepath.Join(dir, "digest.txt"),
		opts,
	)
	if err != nil {
		t.Fatal(err)
	}
	var doc map[string]interface{}
	if err := yaml.Unmarshal(out, &doc); err != nil {
		t.Fatal(err)
	}
	meta := doc["metadata"].(map[string]interface{})
	ann := meta["annotations"].(map[string]interface{})
	if ann[annotationSignatureTrustTier] != trustTierDevOnly {
		t.Fatalf("annotations: %v", ann)
	}
	st := doc["status"].(map[string]interface{})
	if st["signatureStatus"] != "DevOnly" {
		t.Fatalf("signatureStatus: %v", st["signatureStatus"])
	}
	sp := doc["spec"].(map[string]interface{})["policy"].(map[string]interface{})
	if sp["allowLocalDevSignature"] != true {
		t.Fatalf("policy: %v", sp)
	}
}
