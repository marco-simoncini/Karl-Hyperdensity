package crdv1alpha1

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func testdata(t *testing.T, parts ...string) string {
	t.Helper()
	root := filepath.Join("..", "..", "..")
	return filepath.Join(append([]string{root}, parts...)...)
}

func TestParseResourceLeaseFullFixture(t *testing.T) {
	raw, err := os.ReadFile(testdata(t, "examples", "khr", "resourcelease-linux-envelope-full.json"))
	if err != nil {
		t.Fatal(err)
	}
	var rl ResourceLease
	if err := json.Unmarshal(raw, &rl); err != nil {
		t.Fatal(err)
	}
	if rl.APIVersion != "hyperdensity.karl.io/v1alpha1" || rl.Kind != "ResourceLease" {
		t.Fatalf("unexpected type meta %+v", rl)
	}
	if rl.Metadata.Name != "linux-envelope-full-demo" {
		t.Fatalf("metadata: %q", rl.Metadata.Name)
	}
	if rl.Spec.Resource != "cpu" || rl.Spec.Mode != "envelope" {
		t.Fatalf("spec: %+v", rl.Spec)
	}
	if rl.Spec.Donor.Name != "lease-demo-donor-cell" || rl.Spec.Receiver.Name != "lease-demo-receiver-cell" {
		t.Fatalf("refs: %+v %+v", rl.Spec.Donor, rl.Spec.Receiver)
	}
	if rl.Spec.DurationSeconds == nil || *rl.Spec.DurationSeconds != 120 {
		t.Fatal("expected durationSeconds 120")
	}
	if rl.Spec.DryRunOnly == nil || !*rl.Spec.DryRunOnly {
		t.Fatal("expected dryRunOnly true")
	}
}

func TestParseResourcePortFullFixture(t *testing.T) {
	raw, err := os.ReadFile(testdata(t, "examples", "khr", "resourceport-linux-envelope-full.json"))
	if err != nil {
		t.Fatal(err)
	}
	var rp ResourcePort
	if err := json.Unmarshal(raw, &rp); err != nil {
		t.Fatal(err)
	}
	if len(rp.Spec.AppliesToRuntimeProviderIDs) != 2 {
		t.Fatalf("appliesTo: %v", rp.Spec.AppliesToRuntimeProviderIDs)
	}
	if len(rp.Spec.Ports.CPU.Modes) < 2 {
		t.Fatal("expected cpu modes")
	}
	if rp.Spec.Ports.Disk == nil || len(rp.Spec.Ports.Disk.Modes) == 0 {
		t.Fatal("expected disk modes")
	}
}

func TestParseCellFullFixture(t *testing.T) {
	raw, err := os.ReadFile(testdata(t, "examples", "khr", "cell-linux-envelope-full.json"))
	if err != nil {
		t.Fatal(err)
	}
	var c Cell
	if err := json.Unmarshal(raw, &c); err != nil {
		t.Fatal(err)
	}
	if c.Spec.ShellRef.Name == "" || c.Spec.RuntimeProviderRef.Name == "" {
		t.Fatal("expected shell and runtime refs")
	}
	if c.Spec.HostRef == nil || c.Spec.HostRef.Name == "" {
		t.Fatal("expected hostRef")
	}
}

func TestParseRuntimeProviderFullFixture(t *testing.T) {
	raw, err := os.ReadFile(testdata(t, "examples", "khr", "runtimeprovider-linux-cgroup-full.json"))
	if err != nil {
		t.Fatal(err)
	}
	var rp RuntimeProvider
	if err := json.Unmarshal(raw, &rp); err != nil {
		t.Fatal(err)
	}
	if rp.Spec.Driver != "linux-systemd" {
		t.Fatalf("driver %q", rp.Spec.Driver)
	}
	if rp.Spec.ID != "linux.systemd.v1" {
		t.Fatalf("id %q", rp.Spec.ID)
	}
}
