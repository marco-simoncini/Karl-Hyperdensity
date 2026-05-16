package lanediscovery

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type kubeVM struct {
	Namespace string
	Name      string
	Labels    map[string]string
	Running   bool
	Status    string
	NodeName  string
}

type kubePod struct {
	Namespace  string
	Name       string
	Labels     map[string]string
	NodeName   string
	Running    bool
	Sandbox    bool
	NativeLive bool
}

type kubeNode struct {
	Name   string
	Labels map[string]string
	Ready  bool
}

type kubeResourcePort struct {
	Name      string
	Namespace string
	ShellRef  string
	CellRef   string
	Provider  string
}

func kubectlJSON(context string, args ...string) ([]byte, error) {
	base := []string{"--context", context}
	base = append(base, args...)
	cmd := exec.Command("kubectl", base...)
	out, err := cmd.Output()
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("kubectl %v: %w (%s)", args, err, strings.TrimSpace(string(ee.Stderr)))
		}
		return nil, fmt.Errorf("kubectl %v: %w", args, err)
	}
	return out, nil
}

func discoverNodes(context string, observedAt string) ([]DiscoveredHost, error) {
	raw, err := kubectlJSON(context, "get", "nodes", "-o", "json")
	if err != nil {
		return nil, err
	}
	var list struct {
		Items []struct {
			Metadata struct {
				Name   string            `json:"name"`
				Labels map[string]string `json:"labels"`
			} `json:"metadata"`
			Status struct {
				Conditions []struct {
					Type   string `json:"type"`
					Status string `json:"status"`
				} `json:"conditions"`
			} `json:"status"`
		} `json:"items"`
	}
	if err := json.Unmarshal(raw, &list); err != nil {
		return nil, err
	}
	out := make([]DiscoveredHost, 0, len(list.Items))
	for _, n := range list.Items {
		ready := false
		for _, c := range n.Status.Conditions {
			if c.Type == "Ready" && c.Status == "True" {
				ready = true
				break
			}
		}
		out = append(out, DiscoveredHost{
			HostID:      "karl-host-" + n.Metadata.Name,
			NodeName:    n.Metadata.Name,
			Provider:    "khr.native",
			RuntimeMode: "cluster-observed",
			Labels:      n.Metadata.Labels,
			Ready:       ready,
			ObservedAt:  observedAt,
		})
	}
	return out, nil
}

func discoverVMs(context string) ([]kubeVM, error) {
	raw, err := kubectlJSON(context, "get", "virtualmachines.kubevirt.io", "-A", "-o", "json")
	if err != nil {
		return nil, err
	}
	var list struct {
		Items []struct {
			Metadata struct {
				Name      string            `json:"name"`
				Namespace string            `json:"namespace"`
				Labels    map[string]string `json:"labels"`
			} `json:"metadata"`
			Status struct {
				PrintableStatus string `json:"printableStatus"`
				Ready           bool   `json:"ready"`
			} `json:"status"`
		} `json:"items"`
	}
	if err := json.Unmarshal(raw, &list); err != nil {
		return nil, err
	}
	out := make([]kubeVM, 0, len(list.Items))
	for _, vm := range list.Items {
		st := strings.ToLower(vm.Status.PrintableStatus)
		running := vm.Status.Ready || st == "running"
		out = append(out, kubeVM{
			Namespace: vm.Metadata.Namespace,
			Name:      vm.Metadata.Name,
			Labels:    vm.Metadata.Labels,
			Running:   running,
			Status:    vm.Status.PrintableStatus,
		})
	}
	return out, nil
}

func discoverPods(context string) ([]kubePod, error) {
	raw, err := kubectlJSON(context, "get", "pods", "-A", "-o", "json")
	if err != nil {
		return nil, err
	}
	var list struct {
		Items []struct {
			Metadata struct {
				Name      string            `json:"name"`
				Namespace string            `json:"namespace"`
				Labels    map[string]string `json:"labels"`
			} `json:"metadata"`
			Spec struct {
				NodeName string `json:"nodeName"`
			} `json:"spec"`
			Status struct {
				Phase string `json:"phase"`
			} `json:"status"`
		} `json:"items"`
	}
	if err := json.Unmarshal(raw, &list); err != nil {
		return nil, err
	}
	out := make([]kubePod, 0)
	for _, p := range list.Items {
		if strings.HasPrefix(p.Metadata.Name, "virt-launcher-") {
			continue
		}
		sandbox := p.Metadata.Labels["khr.karl.io/sandbox"] == "true" ||
			p.Metadata.Namespace == "khr-runtime-sandbox" ||
			strings.HasPrefix(p.Metadata.Name, "khr-runtime-")
		if !sandbox {
			continue
		}
		running := p.Status.Phase == "Running"
		out = append(out, kubePod{
			Namespace:  p.Metadata.Namespace,
			Name:       p.Metadata.Name,
			Labels:     p.Metadata.Labels,
			NodeName:   p.Spec.NodeName,
			Running:    running,
			Sandbox:    sandbox,
			NativeLive: IsNativeLiveWorkload(p.Metadata.Name, p.Metadata.Labels, p.Metadata.Namespace),
		})
	}
	return out, nil
}

func discoverClusterResourcePorts(context string) ([]kubeResourcePort, error) {
	raw, err := kubectlJSON(context, "get", "resourceports", "-A", "-o", "json")
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") || strings.Contains(err.Error(), "not found") {
			return nil, nil
		}
		return nil, err
	}
	var list struct {
		Items []struct {
			Metadata struct {
				Name      string `json:"name"`
				Namespace string `json:"namespace"`
			} `json:"metadata"`
			Spec struct {
				ShellRef string `json:"shellRef"`
				CellRef  string `json:"cellRef"`
				Provider string `json:"provider"`
			} `json:"spec"`
		} `json:"items"`
	}
	if err := json.Unmarshal(raw, &list); err != nil {
		return nil, err
	}
	out := make([]kubeResourcePort, 0, len(list.Items))
	for _, rp := range list.Items {
		out = append(out, kubeResourcePort{
			Name:      rp.Metadata.Name,
			Namespace: rp.Metadata.Namespace,
			ShellRef:  rp.Spec.ShellRef,
			CellRef:   rp.Spec.CellRef,
			Provider:  rp.Spec.Provider,
		})
	}
	return out, nil
}

func enrichVMNodes(context string, vms []kubeVM) {
	raw, err := kubectlJSON(context, "get", "pods", "-A", "-l", "kubevirt.io=virt-launcher", "-o", "json")
	if err != nil {
		return
	}
	var list struct {
		Items []struct {
			Metadata struct {
				Namespace string            `json:"namespace"`
				Labels    map[string]string `json:"labels"`
			} `json:"metadata"`
			Spec struct {
				NodeName string `json:"nodeName"`
			} `json:"spec"`
			Status struct {
				Phase string `json:"phase"`
			} `json:"status"`
		} `json:"items"`
	}
	if json.Unmarshal(raw, &list) != nil {
		return
	}
	vmNode := map[string]string{}
	for _, p := range list.Items {
		vmName := p.Metadata.Labels["kubevirt.io/vm"]
		if vmName == "" {
			continue
		}
		key := p.Metadata.Namespace + "/" + vmName
		if p.Status.Phase == "Running" {
			vmNode[key] = p.Spec.NodeName
		}
	}
	for i := range vms {
		key := vms[i].Namespace + "/" + vms[i].Name
		if n, ok := vmNode[key]; ok {
			vms[i].NodeName = n
			vms[i].Running = true
		}
	}
}

func buildFromCluster(context string) (Result, error) {
	now := time.Now().UTC()
	observedAt := now.Format(time.RFC3339)
	res := Result{
		Mode:           ModeLaneDiscovery,
		ClusterContext: context,
		ObservedAt:     observedAt,
		Safety:         DefaultSafetyPolicy(),
		Summary:        map[string]int{},
	}

	hosts, err := discoverNodes(context, observedAt)
	if err != nil {
		return res, err
	}
	res.DiscoveredHosts = hosts

	vms, err := discoverVMs(context)
	if err != nil {
		return res, err
	}
	enrichVMNodes(context, vms)

	pods, err := discoverPods(context)
	if err != nil {
		return res, err
	}

	seenShell := map[string]bool{}
	seenCell := map[string]bool{}
	laneCounts := map[string]int{}

	addWorkload := func(h WorkloadHint) {
		lane, provider, class, liveScale, block := ClassifyWorkload(h)
		laneCounts[lane]++
		shellRef := fmt.Sprintf("%s/Shell/%s", h.Namespace, h.Name)
		cellRef := fmt.Sprintf("%s/Cell/%s", h.Namespace, h.Name)
		if !seenShell[shellRef] {
			seenShell[shellRef] = true
			rc := "container"
			if h.VMType == "vm" {
				rc = "kubevirt-vm"
			}
			res.DiscoveredShells = append(res.DiscoveredShells, DiscoveredShell{
				Ref: shellRef, Namespace: h.Namespace, Name: h.Name,
				RuntimeClass: rc, OSFamily: h.OSFamily, VMType: h.VMType,
				ProviderBinding: provider,
			})
		}
		if !seenCell[cellRef] {
			seenCell[cellRef] = true
			res.DiscoveredCells = append(res.DiscoveredCells, DiscoveredCell{
				Ref: cellRef, ShellRef: shellRef, Namespace: h.Namespace, Name: h.Name,
				VMType: h.VMType, OSFamily: h.OSFamily,
				SessionType: InferSessionType(h.OSFamily, h.Running),
				NodeName: h.NodeName, Running: h.Running, ProviderBinding: provider,
			})
		}
		portRef := fmt.Sprintf("%s/ResourcePort/%s-port", h.Namespace, h.Name)
		res.DiscoveredResourcePorts = append(res.DiscoveredResourcePorts, DiscoveredResourcePort{
			Ref: portRef, ShellRef: shellRef, CellRef: cellRef,
			Lane: lane, ProviderBinding: provider, Classification: class,
			LiveScaleCapabilityObserved: liveScale,
		})
		if block != nil {
			res.BlockedStates = append(res.BlockedStates, *block)
		}
	}

	for _, vm := range vms {
		osFamily := InferOSFamily(vm.Name, vm.Labels)
		addWorkload(WorkloadHint{
			Name: vm.Name, Namespace: vm.Namespace, OSFamily: osFamily,
			VMType: "vm", Running: vm.Running, NodeName: vm.NodeName,
		})
	}

	for _, p := range pods {
		addWorkload(WorkloadHint{
			Name: p.Name, Namespace: p.Namespace, OSFamily: "linux",
			VMType: "container", Running: p.Running, SandboxPod: p.Sandbox,
			NativeLive: p.NativeLive, NodeName: p.NodeName,
		})
	}

	clusterPorts, _ := discoverClusterResourcePorts(context)
	for _, rp := range clusterPorts {
		res.DiscoveredResourcePorts = append(res.DiscoveredResourcePorts, DiscoveredResourcePort{
			Ref: fmt.Sprintf("%s/ResourcePort/%s", rp.Namespace, rp.Name),
			ShellRef: rp.ShellRef, CellRef: rp.CellRef,
			Lane: LaneLinuxContainerCgroup, ProviderBinding: rp.Provider,
			Classification: ClassificationObservationOnly,
			LiveScaleCapabilityObserved: true, ClusterObserved: true,
		})
		laneCounts[LaneLinuxContainerCgroup]++
	}

	res.LaneCapabilities = summarizeLaneCapabilities(laneCounts, res.DiscoveredResourcePorts)
	res.Summary["hosts"] = len(res.DiscoveredHosts)
	res.Summary["shells"] = len(res.DiscoveredShells)
	res.Summary["cells"] = len(res.DiscoveredCells)
	res.Summary["resourcePorts"] = len(res.DiscoveredResourcePorts)
	res.Summary["blockedStates"] = len(res.BlockedStates)
	return res, nil
}

func summarizeLaneCapabilities(counts map[string]int, ports []DiscoveredResourcePort) []LaneCapability {
	providerByLane := map[string]string{}
	classByLane := map[string]string{}
	liveByLane := map[string]bool{}
	for _, p := range ports {
		if providerByLane[p.Lane] == "" {
			providerByLane[p.Lane] = p.ProviderBinding
			classByLane[p.Lane] = p.Classification
			liveByLane[p.Lane] = p.LiveScaleCapabilityObserved
		}
	}
	out := make([]LaneCapability, 0, len(counts))
	for lane, n := range counts {
		out = append(out, LaneCapability{
			Lane:                        lane,
			Classification:              classByLane[lane],
			ProviderBinding:             providerByLane[lane],
			LiveScaleCapabilityObserved: liveByLane[lane],
			WorkloadCount:               n,
		})
	}
	return out
}
