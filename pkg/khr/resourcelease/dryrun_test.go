package resourcelease

import (
	"encoding/json"
	"testing"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/crdv1alpha1"
)

func leaseFixture(mode, resource string) *crdv1alpha1.ResourceLease {
	l := &crdv1alpha1.ResourceLease{}
	l.APIVersion = "hyperdensity.karl.io/v1alpha1"
	l.Kind = "ResourceLease"
	l.Metadata.Name = "fixture"
	l.Metadata.Namespace = "karl-sandbox"
	l.Spec.Mode = mode
	l.Spec.Resource = resource
	l.Spec.Donor.Kind = "Cell"
	l.Spec.Donor.Name = "a"
	l.Spec.Donor.Namespace = "ns"
	l.Spec.Receiver.Kind = "Cell"
	l.Spec.Receiver.Name = "b"
	l.Spec.Receiver.Namespace = "ns"
	return l
}

func portEnvelope() *crdv1alpha1.ResourcePort {
	p := &crdv1alpha1.ResourcePort{}
	p.APIVersion = "runtime.karl.io/v1alpha1"
	p.Kind = "ResourcePort"
	p.Spec.Ports.CPU.Modes = []string{"static", "envelope"}
	p.Spec.Ports.Memory.Modes = []string{"static", "envelope"}
	return p
}

func TestDryRunAllowedEnvelopeCPU(t *testing.T) {
	res := DryRun(leaseFixture("envelope", "cpu"), portEnvelope(), &CellContext{DonorPlatform: "linux", ReceiverPlatform: "linux"})
	if res.Blocked || !res.Allowed {
		t.Fatalf("expected allowed, got %+v", res)
	}
	if res.ResourcePortCheck != "compatible" {
		t.Fatalf("expected compatible port check")
	}
}

func TestDryRunBlockedNonEnvelopeMode(t *testing.T) {
	res := DryRun(leaseFixture("hotAdd", "cpu"), portEnvelope(), &CellContext{DonorPlatform: "linux", ReceiverPlatform: "linux"})
	if !res.Blocked || res.Allowed {
		t.Fatalf("expected blocked")
	}
}

func TestDryRunBlockedResourcePort(t *testing.T) {
	p := &crdv1alpha1.ResourcePort{}
	p.Spec.Ports.CPU.Modes = []string{"static"}
	p.Spec.Ports.Memory.Modes = []string{"static"}
	res := DryRun(leaseFixture("envelope", "cpu"), p, &CellContext{DonorPlatform: "linux", ReceiverPlatform: "linux"})
	if !res.Blocked {
		t.Fatal("expected blocked for missing envelope mode")
	}
}

func TestDryRunBlockedMissingPort(t *testing.T) {
	res := DryRun(leaseFixture("envelope", "cpu"), nil, &CellContext{DonorPlatform: "linux", ReceiverPlatform: "linux"})
	if !res.Blocked {
		t.Fatal("expected blocked without port")
	}
}

func TestDryRunBlockedNonLinux(t *testing.T) {
	res := DryRun(leaseFixture("envelope", "cpu"), portEnvelope(), &CellContext{DonorPlatform: "windows", ReceiverPlatform: "linux"})
	if !res.Blocked {
		t.Fatal("expected blocked for non-linux")
	}
}

func TestLeaseJSONUnmarshal(t *testing.T) {
	raw := []byte(`{"apiVersion":"hyperdensity.karl.io/v1alpha1","kind":"ResourceLease","metadata":{"name":"x"},"spec":{"donor":{"kind":"Cell","name":"a","namespace":"ns"},"receiver":{"kind":"Cell","name":"b","namespace":"ns"},"resource":"memory","mode":"envelope"}}`)
	l := &crdv1alpha1.ResourceLease{}
	if err := json.Unmarshal(raw, l); err != nil {
		t.Fatal(err)
	}
	res := DryRun(l, portEnvelope(), &CellContext{DonorPlatform: "linux", ReceiverPlatform: "linux"})
	if !res.Allowed {
		t.Fatalf("expected allowed %+v", res)
	}
}
