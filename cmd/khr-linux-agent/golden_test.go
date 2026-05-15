package main

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"
)

func moduleRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	return filepath.Clean(filepath.Join(wd, "..", ".."))
}

func TestGoldenDryRunStdout(t *testing.T) {
	root := moduleRoot(t)
	cfg := filepath.Join(root, "examples", "khr", "khr-linux-agent-config.yaml")
	cmdDir := filepath.Join(root, "cmd", "khr-linux-agent")

	cases := []struct {
		name   string
		lease  string
		port   string
		cell   string
		golden string
		extra  []string
	}{
		{
			name:   "allowed",
			lease:  filepath.Join(root, "examples", "khr", "resourcelease-linux-envelope-full.json"),
			port:   filepath.Join(root, "examples", "khr", "resourceport-linux-envelope-full.json"),
			cell:   filepath.Join(root, "examples", "khr", "golden", "inputs", "cell-context-allowed.json"),
			golden: filepath.Join(root, "examples", "khr", "golden", "dryrun-allowed.json"),
		},
		{
			name:   "blocked_nonlinux",
			lease:  filepath.Join(root, "examples", "khr", "resourcelease-linux-envelope-full.json"),
			port:   filepath.Join(root, "examples", "khr", "resourceport-linux-envelope-full.json"),
			cell:   filepath.Join(root, "examples", "khr", "golden", "inputs", "cell-context-nonlinux.json"),
			golden: filepath.Join(root, "examples", "khr", "golden", "dryrun-blocked-nonlinux.json"),
		},
		{
			name:   "blocked_mode",
			lease:  filepath.Join(root, "examples", "khr", "golden", "inputs", "lease-blocked-mode.json"),
			port:   filepath.Join(root, "examples", "khr", "resourceport-linux-envelope-full.json"),
			golden: filepath.Join(root, "examples", "khr", "golden", "dryrun-blocked-mode.json"),
		},
		{
			name:   "blocked_resourceport",
			lease:  filepath.Join(root, "examples", "khr", "resourcelease-linux-envelope-full.json"),
			port:   filepath.Join(root, "examples", "khr", "golden", "inputs", "resourceport-static-cpu-only.json"),
			golden: filepath.Join(root, "examples", "khr", "golden", "dryrun-blocked-resourceport.json"),
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("KHR_TEST_CGROUP_VERSION", "v2")
			args := []string{
				"run", ".",
				"-mode", "dry-run",
				"-config", cfg,
				"-lease-input", tc.lease,
				"-resource-port-input", tc.port,
			}
			if tc.cell != "" {
				args = append(args, "-cell-context", tc.cell)
			}
			args = append(args, tc.extra...)
			cmd := exec.Command("go", args...)
			cmd.Dir = cmdDir
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			if err := cmd.Run(); err != nil {
				t.Fatalf("go run: %v stderr=%s", err, stderr.String())
			}
			wantRaw, err := os.ReadFile(tc.golden)
			if err != nil {
				t.Fatal(err)
			}
			var want, got any
			if err := json.Unmarshal(wantRaw, &want); err != nil {
				t.Fatal(err)
			}
			if err := json.Unmarshal(stdout.Bytes(), &got); err != nil {
				t.Fatalf("stdout json: %v body=%q", err, stdout.String())
			}
			if !reflect.DeepEqual(want, got) {
				t.Fatalf("golden mismatch for %s\n--- want ---\n%s\n--- got ---\n%s", tc.name, string(wantRaw), stdout.String())
			}
		})
	}
}
