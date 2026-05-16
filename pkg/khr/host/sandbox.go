package host

import "strings"

// ApplyGate explains whether sandbox apply is permitted.
type ApplyGate struct {
	Allowed bool   `json:"allowed"`
	Reason  string `json:"reason,omitempty"`
}

// SandboxApplyAllowed enforces Linux-only sandbox policy with namespace/label allowlist.
func SandboxApplyAllowed(cfg *Config, namespace string, labels map[string]string) ApplyGate {
	if cfg == nil {
		return ApplyGate{Reason: "config is nil"}
	}
	if !cfg.Spec.LinuxOnly {
		return ApplyGate{Reason: "linuxOnly must be true"}
	}
	if !cfg.Spec.SandboxMode {
		return ApplyGate{Reason: "sandboxMode must be true"}
	}
	if !cfg.Spec.SandboxApplyEnabled {
		return ApplyGate{Reason: "sandboxApplyEnabled is false (default)"}
	}
	if len(cfg.Spec.AllowedNamespaces) > 0 {
		ok := false
		for _, ns := range cfg.Spec.AllowedNamespaces {
			if ns == namespace {
				ok = true
				break
			}
		}
		if !ok {
			return ApplyGate{Reason: "namespace not in allowedNamespaces allowlist"}
		}
	}
	for k, want := range cfg.Spec.AllowedLabels {
		if labels[k] != want {
			return ApplyGate{Reason: "required label " + k + " mismatch"}
		}
	}
	return ApplyGate{Allowed: true, Reason: "sandbox apply permitted"}
}

// NamespaceAllowed checks read-only namespace policy.
func NamespaceAllowed(cfg *Config, namespace string) bool {
	if cfg == nil || len(cfg.Spec.AllowedNamespaces) == 0 {
		return true
	}
	for _, ns := range cfg.Spec.AllowedNamespaces {
		if ns == namespace {
			return true
		}
	}
	return false
}

// LabelKeyAllowed returns true if key is in allowlist (value not checked).
func LabelKeyAllowed(cfg *Config, key string) bool {
	if cfg == nil || len(cfg.Spec.AllowedLabels) == 0 {
		return true
	}
	_, ok := cfg.Spec.AllowedLabels[key]
	return ok
}

// NormalizeNamespace trims whitespace.
func NormalizeNamespace(ns string) string {
	return strings.TrimSpace(ns)
}

// ProductionNamespaces must never be used by sandbox loops.
var ProductionNamespaces = []string{
	"karl-system", "kube-system", "kube-public", "kube-node-lease",
	"default", "ingress", "longhorn-system", "kubevirt", "cdi",
}

// ProductionNamespaceBlocked returns true for production namespaces.
func ProductionNamespaceBlocked(namespace string) bool {
	ns := NormalizeNamespace(namespace)
	for _, blocked := range ProductionNamespaces {
		if ns == blocked {
			return true
		}
	}
	return false
}

// LabelsAllowlistMatch requires every allowlist key to match labels.
func LabelsAllowlistMatch(cfg *Config, labels map[string]string) bool {
	if cfg == nil || len(cfg.Spec.AllowedLabels) == 0 {
		return true
	}
	for k, want := range cfg.Spec.AllowedLabels {
		if labels[k] != want {
			return false
		}
	}
	return true
}
