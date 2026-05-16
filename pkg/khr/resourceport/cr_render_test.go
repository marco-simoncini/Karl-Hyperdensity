package resourceport

import (
	"testing"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/crdv1alpha1"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/host"
)

func TestCandidateToCRStableRender(t *testing.T) {
	cfg := loopConfig(true)
	c := ReportCandidate(cfg, "khr-runtime-sandbox/Shell/target-1", "khr-runtime-sandbox/Cell/target-1", "khr-runtime-sandbox", "target-1-port")
	meta := CRDocumentMeta{
		Source:           ManagedByValue,
		RuntimeVersion:   host.RuntimeVersion,
		ObservedAt:       "2026-05-16T12:00:00Z",
		SafetyMode:       "sandbox",
		EmissionMode:     EmissionModeCRPreview,
		SandboxNamespace: "khr-runtime-sandbox",
	}
	rp1 := CandidateToCR(c, meta)
	rp2 := CandidateToCR(c, meta)
	b1, err := RenderCRJSON(rp1)
	if err != nil {
		t.Fatal(err)
	}
	b2, err := RenderCRJSON(rp2)
	if err != nil {
		t.Fatal(err)
	}
	if string(b1) != string(b2) {
		t.Fatalf("render not stable:\n%s\n---\n%s", b1, b2)
	}
	if rp1.Metadata.Name != "khr-runtime-sandbox-target-1-port" {
		t.Fatalf("name=%q", rp1.Metadata.Name)
	}
	if rp1.Metadata.Annotations[AnnotationEmissionMode] != EmissionModeCRPreview {
		t.Fatal("missing emission mode annotation")
	}
}

func TestRenderCRFiles(t *testing.T) {
	dir := t.TempDir()
	cfg := loopConfig(true)
	c := ReportCandidate(cfg, "khr-runtime-sandbox/Shell/t", "khr-runtime-sandbox/Cell/t", "khr-runtime-sandbox", "t-port")
	meta := metaFromConfig(cfg, "2026-05-16T12:00:00Z", EmissionModeCRPreview, "khr-runtime-sandbox")
	paths, err := RenderCRFiles(dir, []crdv1alpha1.ResourcePort{CandidateToCR(c, meta)})
	if err != nil {
		t.Fatal(err)
	}
	if len(paths) != 1 {
		t.Fatalf("paths=%v", paths)
	}
}
