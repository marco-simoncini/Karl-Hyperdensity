package host

import (
	"encoding/json"
	"os"
	"time"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/crdv1alpha1"
)

const hostAPIVersion = "runtime.karl.io/v1alpha1"

// BuildHostStatus produces a Host CR-shaped JSON document (local only; no kube apply).
func BuildHostStatus(cfg *Config, nodeName string, observedPorts []crdv1alpha1.ObjectRef) crdv1alpha1.Host {
	if cfg == nil {
		cfg = &Config{}
	}
	if nodeName == "" {
		nodeName, _ = os.Hostname()
	}
	hostID := cfg.Spec.HostID
	if hostID == "" {
		hostID = "karl-host-unknown"
	}
	runtimeMode := "sandbox"
	if !cfg.Spec.SandboxMode {
		runtimeMode = "disabled"
	} else if cfg.Spec.SandboxApplyEnabled {
		runtimeMode = "preview"
	}
	caps := ReportCapabilities(cfg)
	capsRaw, _ := json.Marshal(caps)
	now := time.Now().UTC()
	safetyMode := "sandbox"
	if !cfg.Spec.SandboxMode {
		safetyMode = "production-blocked"
	}
	return crdv1alpha1.Host{
		APIVersion: hostAPIVersion,
		Kind:       "Host",
		Metadata: crdv1alpha1.ObjectMeta{
			Name: hostID,
			Labels: map[string]string{
				"khr.karl.io/sandbox": "true",
			},
		},
		Spec: crdv1alpha1.HostSpec{
			HostID:      hostID,
			NodeName:    nodeName,
			Provider:    "khr.native",
			RuntimeMode: runtimeMode,
			Labels:      copyLabels(cfg.Spec.AllowedLabels),
		},
		Status: crdv1alpha1.HostStatus{
			Phase: "Observed",
			Conditions: []crdv1alpha1.HostCondition{{
				Type:               "Ready",
				Status:             "False",
				Reason:             "NoController",
				Message:            "KHR-I status from karl-host-runtime only; no reconciler loop",
				LastTransitionTime: now.Format(time.RFC3339),
			}, {
				Type:               "SandboxOnly",
				Status:             "True",
				Reason:             "ProductionUnsupported",
				Message:            "Host registration is sandbox/read-only in KHR-I",
				LastTransitionTime: now.Format(time.RFC3339),
			}},
			Capabilities:          capsRaw,
			ObservedResourcePorts: observedPorts,
			LastHeartbeatTime:     now.Format(time.RFC3339),
			RuntimeVersion:        RuntimeVersion,
			SafetyMode:            safetyMode,
		},
	}
}
