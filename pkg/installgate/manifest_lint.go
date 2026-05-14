package installgate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"sigs.k8s.io/yaml"
)

type LintResult struct {
	Passed   bool     `json:"passed"`
	Blockers []string `json:"blockers"`
	Warnings []string `json:"warnings,omitempty"`
}

func LintManifestDir(dir string) (LintResult, error) {
	required := []string{
		"kustomization.yaml", "namespace.yaml", "serviceaccount.yaml", "role.yaml", "rolebinding.yaml",
		"deployment.yaml", "service-metrics.yaml", "leader-election-role.yaml", "leader-election-rolebinding.yaml",
		"configmap.yaml", "crd.yaml", "rollback-plan.yaml", "install-gate.yaml",
	}
	var blockers []string
	for _, name := range required {
		if _, err := os.Stat(filepath.Join(dir, name)); err != nil {
			blockers = append(blockers, "missing_"+strings.TrimSuffix(name, ".yaml"))
		}
	}

	docs, err := readYAMLDocs(dir)
	if err != nil {
		return LintResult{}, err
	}

	for _, doc := range docs {
		kind := strMap(doc, "kind")
		if kind == "ClusterRoleBinding" && strNested(doc, "roleRef", "name") == "cluster-admin" {
			blockers = append(blockers, "cluster_admin")
		}

		if kind == "Role" || kind == "ClusterRole" {
			for _, rule := range mapSlice(doc["rules"]) {
				resources := stringSlice(rule["resources"])
				verbs := stringSlice(rule["verbs"])
				if has(resources, "*") {
					blockers = append(blockers, "wildcard_resources")
				}
				if has(verbs, "*") {
					blockers = append(blockers, "wildcard_verbs")
				}
				if has(resources, "pods/exec") {
					blockers = append(blockers, "pods_exec_allowed")
				}
				if has(resources, "nodes") && (has(verbs, "update") || has(verbs, "patch") || has(verbs, "delete")) {
					blockers = append(blockers, "nodes_write_allowed")
				}
			}
		}

		if kind == "Deployment" {
			spec := mapAny(doc["spec"])
			tpl := mapAny(spec["template"])
			podSpec := mapAny(tpl["spec"])
			labels := mapAny(mapAny(doc["metadata"])["labels"])
			if labels["hyperdensity.karl.io/general-production-auto-allowed"] != "false" {
				blockers = append(blockers, "missing_general_production_auto_label_false")
			}
			if labels["hyperdensity.karl.io/production-auto-with-policy"] != "false" {
				blockers = append(blockers, "missing_production_auto_with_policy_label_false")
			}
			if strAny(podSpec["serviceAccountName"]) == "" {
				blockers = append(blockers, "missing_service_account")
			}
			if boolAny(podSpec["hostPID"]) || boolAny(podSpec["hostNetwork"]) || boolAny(podSpec["hostIPC"]) {
				blockers = append(blockers, "host_namespace_unsafe")
			}
			for _, c := range mapSlice(podSpec["containers"]) {
				if mapAny(c["livenessProbe"]) == nil {
					blockers = append(blockers, "missing_liveness_probe")
				}
				if mapAny(c["readinessProbe"]) == nil {
					blockers = append(blockers, "missing_readiness_probe")
				}
				if mapAny(c["startupProbe"]) == nil {
					blockers = append(blockers, "missing_startup_probe")
				}
				if mapAny(c["resources"]) == nil || mapAny(mapAny(c["resources"])["requests"]) == nil || mapAny(mapAny(c["resources"])["limits"]) == nil {
					blockers = append(blockers, "missing_resource_limits")
				}
				sc := mapAny(c["securityContext"])
				if boolAny(sc["privileged"]) {
					blockers = append(blockers, "privileged_container")
				}
				if boolAny(sc["allowPrivilegeEscalation"]) {
					blockers = append(blockers, "allow_privilege_escalation")
				}
				for _, env := range mapSlice(c["env"]) {
					name := strAny(env["name"])
					value := strings.ToLower(strAny(env["value"]))
					if name == "GENERAL_PRODUCTION_AUTO_ALLOWED" && value != "false" {
						blockers = append(blockers, "production_auto_enabled")
					}
					if name == "PRODUCTION_AUTO_WITH_POLICY" && value != "false" {
						blockers = append(blockers, "production_auto_with_policy_enabled")
					}
				}
			}
		}
	}

	blockers = dedupe(blockers)
	return LintResult{Passed: len(blockers) == 0, Blockers: blockers}, nil
}

func readYAMLDocs(dir string) ([]map[string]interface{}, error) {
	patterns := []string{"*.yaml", "*.yml"}
	var out []map[string]interface{}
	for _, p := range patterns {
		files, err := filepath.Glob(filepath.Join(dir, p))
		if err != nil {
			return nil, err
		}
		for _, f := range files {
			raw, err := os.ReadFile(f)
			if err != nil {
				return nil, err
			}
			parts := strings.Split(string(raw), "\n---")
			for _, part := range parts {
				part = strings.TrimSpace(part)
				if part == "" {
					continue
				}
				var doc map[string]interface{}
				if err := yaml.Unmarshal([]byte(part), &doc); err != nil {
					return nil, fmt.Errorf("%s: %w", f, err)
				}
				if len(doc) > 0 {
					out = append(out, doc)
				}
			}
		}
	}
	return out, nil
}

func mapAny(v interface{}) map[string]interface{} {
	if m, ok := v.(map[string]interface{}); ok {
		return m
	}
	return nil
}
func mapSlice(v interface{}) []map[string]interface{} {
	raw, ok := v.([]interface{})
	if !ok {
		return nil
	}
	out := make([]map[string]interface{}, 0, len(raw))
	for _, item := range raw {
		if m, ok := item.(map[string]interface{}); ok {
			out = append(out, m)
		}
	}
	return out
}
func stringSlice(v interface{}) []string {
	raw, ok := v.([]interface{})
	if !ok {
		return nil
	}
	out := make([]string, 0, len(raw))
	for _, item := range raw {
		if s, ok := item.(string); ok {
			out = append(out, s)
		}
	}
	return out
}
func has(list []string, needle string) bool {
	for _, item := range list {
		if item == needle {
			return true
		}
	}
	return false
}
func dedupe(in []string) []string {
	if len(in) == 0 {
		return nil
	}
	seen := map[string]struct{}{}
	out := make([]string, 0, len(in))
	for _, v := range in {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}
	return out
}
func strMap(m map[string]interface{}, key string) string {
	if m == nil {
		return ""
	}
	return strAny(m[key])
}
func strNested(m map[string]interface{}, key1, key2 string) string {
	return strAny(mapAny(m[key1])[key2])
}
func strAny(v interface{}) string {
	s, _ := v.(string)
	return s
}
func boolAny(v interface{}) bool {
	b, _ := v.(bool)
	return b
}
