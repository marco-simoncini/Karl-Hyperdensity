package resourceport

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/crdv1alpha1"
	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/host"
)

const (
	AnnotationSource             = "karl.io/source"
	AnnotationRuntimeVersion       = "karl.io/runtime-version"
	AnnotationObservedAt           = "karl.io/observed-at"
	AnnotationSafetyMode           = "karl.io/safety-mode"
	AnnotationEmissionMode         = "karl.io/emission-mode"
	LabelManagedBy                 = "karl.io/managed-by"
	LabelSandboxNamespace          = "karl.io/sandbox-namespace"
	ManagedByValue                 = "karl-host-runtime"
)

// CRDocumentMeta is owner/source metadata stamped on preview and applied CRs.
type CRDocumentMeta struct {
	Source          string
	RuntimeVersion  string
	ObservedAt      string
	SafetyMode      string
	EmissionMode    string
	SandboxNamespace string
}

// CandidateToCR converts a loop candidate to a stable cluster-scoped ResourcePort CR.
func CandidateToCR(c Candidate, meta CRDocumentMeta) crdv1alpha1.ResourcePort {
	name := clusterResourcePortName(meta.SandboxNamespace, c.Metadata.Name)
	labels := map[string]string{
		"khr.karl.io/sandbox": "true",
		LabelManagedBy:        ManagedByValue,
		LabelSandboxNamespace: meta.SandboxNamespace,
	}
	for k, v := range c.Metadata.Labels {
		if k != "" && v != "" {
			labels[k] = v
		}
	}
	annotations := map[string]string{
		AnnotationSource:         meta.Source,
		AnnotationRuntimeVersion: meta.RuntimeVersion,
		AnnotationObservedAt:     meta.ObservedAt,
		AnnotationSafetyMode:     meta.SafetyMode,
		AnnotationEmissionMode:     meta.EmissionMode,
	}
	status, _ := json.Marshal(map[string]any{
		"phase":      c.Status.Phase,
		"observedAt": c.Status.ObservedAt,
	})
	return crdv1alpha1.ResourcePort{
		APIVersion: "runtime.karl.io/v1alpha1",
		Kind:       "ResourcePort",
		Metadata: crdv1alpha1.ObjectMeta{
			Name:        name,
			Labels:      labels,
			Annotations: annotations,
		},
		Spec: crdv1alpha1.ResourcePortSpec{
			Provider:     c.Spec.Provider,
			ShellRef:     c.Spec.ShellRef,
			CellRef:      c.Spec.CellRef,
			Capabilities: append([]string(nil), c.Spec.Capabilities...),
			Hotplug:      &c.Spec.Hotplug,
			Ports:        c.Spec.Ports,
			Notes:        "KHR-K sandbox preview; no autonomous reconcile",
		},
		Status: status,
	}
}

func clusterResourcePortName(namespace, portName string) string {
	base := sanitizeDNS1123(namespace + "-" + portName)
	if !strings.HasPrefix(base, "khr-") {
		base = "khr-" + base
	}
	if len(base) > 253 {
		base = base[:253]
	}
	return strings.Trim(base, "-.")
}

func sanitizeDNS1123(s string) string {
	s = strings.ToLower(s)
	var b strings.Builder
	lastDash := false
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
			lastDash = false
		case r == '-' || r == '.' || r == '/':
			if !lastDash && b.Len() > 0 {
				b.WriteByte('-')
				lastDash = true
			}
		default:
			if !lastDash && b.Len() > 0 {
				b.WriteByte('-')
				lastDash = true
			}
		}
	}
	return strings.Trim(b.String(), "-")
}

func metaFromConfig(cfg *host.Config, observedAt, emissionMode, namespace string) CRDocumentMeta {
	return CRDocumentMeta{
		Source:           ManagedByValue,
		RuntimeVersion:   host.RuntimeVersion,
		ObservedAt:       observedAt,
		SafetyMode:       "sandbox",
		EmissionMode:     emissionMode,
		SandboxNamespace: namespace,
	}
}

// RenderCRJSON returns stable indented JSON for a ResourcePort CR.
func RenderCRJSON(rp crdv1alpha1.ResourcePort) ([]byte, error) {
	return json.MarshalIndent(rp, "", "  ")
}

// RenderCRFiles writes stable JSON CR documents under dir.
func RenderCRFiles(dir string, ports []crdv1alpha1.ResourcePort) ([]string, error) {
	if dir == "" {
		return nil, fmt.Errorf("output directory is required")
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	var paths []string
	for _, rp := range ports {
		path := filepath.Join(dir, "resourceport-"+rp.Metadata.Name+".json")
		raw, err := RenderCRJSON(rp)
		if err != nil {
			return nil, err
		}
		if err := os.WriteFile(path, raw, 0o644); err != nil {
			return nil, err
		}
		paths = append(paths, path)
	}
	return paths, nil
}
