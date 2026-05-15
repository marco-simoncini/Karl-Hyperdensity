// Package crdv1alpha1 holds JSON/YAML-friendly structs aligned to Karl CRDs (Sprint 2 v1alpha1 subset).
package crdv1alpha1

// ObjectMeta mirrors metadata fields used by KHR dry-run fixtures.
type ObjectMeta struct {
	Name            string            `json:"name,omitempty"`
	Namespace       string            `json:"namespace,omitempty"`
	Labels          map[string]string `json:"labels,omitempty"`
	Annotations     map[string]string `json:"annotations,omitempty"`
	UID             string            `json:"uid,omitempty"`
	ResourceVersion string            `json:"resourceVersion,omitempty"`
}
