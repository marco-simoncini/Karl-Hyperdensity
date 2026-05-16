package resourceport

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

func runKubectlJSON(context, namespace, labelSelector string) ([]byte, error) {
	args := []string{
		"--context", context,
		"-n", namespace,
		"get", "pods",
		"-l", labelSelector,
		"-o", "json",
	}
	cmd := exec.Command("kubectl", args...)
	return cmd.Output()
}

func parsePodList(raw []byte, opts LoopOptions) []SandboxTarget {
	var list struct {
		Items []struct {
			Metadata struct {
				Name   string            `json:"name"`
				Labels map[string]string `json:"labels"`
			} `json:"metadata"`
			Spec struct {
				Containers []struct {
					Name string `json:"name"`
				} `json:"containers"`
			} `json:"spec"`
		} `json:"items"`
	}
	if err := json.Unmarshal(raw, &list); err != nil {
		return nil
	}
	out := make([]SandboxTarget, 0, len(list.Items))
	for _, item := range list.Items {
		container := ""
		if len(item.Spec.Containers) > 0 {
			container = item.Spec.Containers[0].Name
		}
		t := SandboxTarget{
			Namespace: opts.Namespace,
			PodName:   item.Metadata.Name,
			Container: container,
			Labels:    item.Metadata.Labels,
			HostID:    opts.Config.Spec.HostID,
			CgroupPath: ObserveCgroupPath(opts.Config, SandboxTarget{Namespace: opts.Namespace, PodName: item.Metadata.Name}),
		}
		out = append(out, t)
	}
	return out
}

// CurrentKubeContext returns kubectl current-context (empty if unavailable).
func CurrentKubeContext() string {
	out, err := exec.Command("kubectl", "config", "current-context").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// FormatKubectlError includes stderr when kubectl fails.
func FormatKubectlError(err error, out []byte) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("kubectl: %w (%s)", err, strings.TrimSpace(string(out)))
}
