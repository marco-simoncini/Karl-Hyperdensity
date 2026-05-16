package cgroup

import "testing"

func TestComputeMemoryTargetScaleUpDown(t *testing.T) {
	got, err := ComputeMemoryTarget(128*1024*1024, false, 64*1024*1024, "scaleUp")
	if err != nil || got != 192*1024*1024 {
		t.Fatalf("scaleUp got=%d err=%v", got, err)
	}
	got, err = ComputeMemoryTarget(128*1024*1024, false, 200*1024*1024, "scaleDown")
	if err != nil || got != minSandboxMemoryBytes {
		t.Fatalf("scaleDown floor got=%d err=%v", got, err)
	}
	got, err = ComputeMemoryTarget(0, true, 256*1024*1024, "envelope")
	if err != nil || got != 256*1024*1024 {
		t.Fatalf("envelope got=%d err=%v", got, err)
	}
}

func TestParseFormatMemoryValue(t *testing.T) {
	b, u, err := ParseMemoryValue("max")
	if err != nil || !u || b != 0 {
		t.Fatalf("max parse=%d unlimited=%v err=%v", b, u, err)
	}
	b, u, err = ParseMemoryValue("134217728")
	if err != nil || u || b != 134217728 {
		t.Fatalf("bytes parse=%d unlimited=%v err=%v", b, u, err)
	}
	if FormatMemoryValue(134217728) != "134217728" {
		t.Fatal("format bytes")
	}
}
