package lanediscovery

import (
	"strings"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/windowslane"
)

// WorkloadHint carries discovery inputs for classification.
type WorkloadHint struct {
	Name       string
	Namespace  string
	OSFamily   string // linux | windows | unknown
	VMType     string // container | vm | unknown
	Running    bool
	SandboxPod bool
	NodeName   string
}

// ClassifyWorkload returns lane, provider, classification, liveScale observed, and optional block.
func ClassifyWorkload(h WorkloadHint) (lane, provider, classification string, liveScale bool, block *BlockedState) {
	os := strings.ToLower(strings.TrimSpace(h.OSFamily))
	vmType := strings.ToLower(strings.TrimSpace(h.VMType))

	switch {
	case h.SandboxPod && os != "windows":
		lane = LaneLinuxContainerCgroup
		provider = "khr.native"
		classification = ClassificationLiveInPlaceCapable
		liveScale = h.Running
		return lane, provider, classification, liveScale, nil

	case vmType == "vm" && os == "linux":
		lane = LaneLinuxVMCompatibility
		provider = windowslane.ProviderKubevirtCompatibility
		classification = ClassificationCompatibilityFallback
		liveScale = false
		if !h.Running {
			block = &BlockedState{
				State: "observation-only", Reason: "VM stopped; live scale not observable",
				TargetRef: ref(h.Namespace, h.Name), Lane: lane,
			}
			classification = ClassificationObservationOnly
		}
		return lane, provider, classification, liveScale, block

	case vmType == "vm" && os == "windows":
		lane = LaneWindowsVMSession
		provider = windowslane.ProviderKubevirtCompatibility
		classification = ClassificationCompatibilityFallback
		liveScale = false
		block = &BlockedState{
			State:     windowslane.BlockedRequiresRestart,
			Reason:    "kubevirt.compatibility: Windows VM memory/CPU change may require restart; live-in-place not asserted",
			TargetRef: ref(h.Namespace, h.Name),
			Lane:      lane,
		}
		if h.Running {
			classification = ClassificationObservationOnly
		} else {
			classification = ClassificationUnsupported
			block.Reason = "Windows VM stopped; lane unsupported for live scale observation"
		}
		return lane, provider, classification, liveScale, block

	case vmType == "vm":
		lane = LaneKubevirtCompatibility
		provider = windowslane.ProviderKubevirtCompatibility
		classification = ClassificationCompatibilityFallback
		block = &BlockedState{
			State:     windowslane.BlockedProviderUnsupported,
			Reason:    "unknown guest OS on kubevirt VM; compatibility-fallback only",
			TargetRef: ref(h.Namespace, h.Name),
			Lane:      lane,
		}
		return lane, provider, classification, false, block

	default:
		lane = LaneLinuxContainerCgroup
		provider = "parent-fabric.observed"
		classification = ClassificationObservationOnly
		return lane, provider, classification, false, nil
	}
}

func ref(ns, name string) string {
	if ns == "" {
		return name
	}
	return ns + "/Cell/" + name
}

// InferOSFamily from VM/pod name and labels.
func InferOSFamily(name string, labels map[string]string) string {
	if labels != nil {
		for _, k := range []string{"kubevirt.io/os", "karl.io/os-family", "khr.karl.io/os-family"} {
			if v := strings.ToLower(labels[k]); v != "" {
				if strings.Contains(v, "win") {
					return "windows"
				}
				if strings.Contains(v, "linux") {
					return "linux"
				}
			}
		}
	}
	lower := strings.ToLower(name)
	if strings.Contains(lower, "win") || strings.Contains(lower, "windows") {
		return "windows"
	}
	return "linux"
}

// InferSessionType for Windows workloads.
func InferSessionType(osFamily string, running bool) string {
	if osFamily != "windows" {
		return ""
	}
	if running {
		return "rdp-session-candidate"
	}
	return "session-offline"
}
