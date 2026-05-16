package crdv1alpha1

import "encoding/json"

// Host is runtime.karl.io/v1alpha1 Host (KHR-I registration foundation).
type Host struct {
	APIVersion string     `json:"apiVersion"`
	Kind       string     `json:"kind"`
	Metadata   ObjectMeta `json:"metadata"`
	Spec       HostSpec   `json:"spec"`
	Status     HostStatus `json:"status,omitempty"`
}

// HostSpec is the KHR-I Host spec subset.
type HostSpec struct {
	HostID              string            `json:"hostId"`
	NodeName            string            `json:"nodeName"`
	Provider            string            `json:"provider"`
	RuntimeMode         string            `json:"runtimeMode"`
	Labels              map[string]string `json:"labels,omitempty"`
	Taints              []HostTaint       `json:"taints,omitempty"`
	KubernetesNodeName  string            `json:"kubernetesNodeName,omitempty"`
}

// HostTaint is a scheduling taint on the host record.
type HostTaint struct {
	Key    string `json:"key,omitempty"`
	Value  string `json:"value,omitempty"`
	Effect string `json:"effect,omitempty"`
}

// HostStatus is the KHR-I Host status subset.
type HostStatus struct {
	Phase                 string              `json:"phase,omitempty"`
	Conditions            []HostCondition     `json:"conditions,omitempty"`
	Capabilities          json.RawMessage     `json:"capabilities,omitempty"`
	ObservedResourcePorts []ObjectRef         `json:"observedResourcePorts,omitempty"`
	LastHeartbeatTime     string              `json:"lastHeartbeatTime,omitempty"`
	RuntimeVersion        string              `json:"runtimeVersion,omitempty"`
	SafetyMode            string              `json:"safetyMode,omitempty"`
}

// HostCondition is a standard condition entry.
type HostCondition struct {
	Type               string `json:"type"`
	Status             string `json:"status"`
	Reason             string `json:"reason,omitempty"`
	Message            string `json:"message,omitempty"`
	LastTransitionTime string `json:"lastTransitionTime,omitempty"`
}

// ObjectRef is a namespaced object reference.
type ObjectRef struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace,omitempty"`
}
