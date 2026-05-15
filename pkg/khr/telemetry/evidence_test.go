package telemetry

import (
	"testing"
	"time"
)

func TestParseKVInt64Lines(t *testing.T) {
	s := "nr_periods 10\nusage_usec 99\n"
	m, w := parseKVInt64Lines(s)
	if len(w) != 0 || m["nr_periods"] != 10 || m["usage_usec"] != 99 {
		t.Fatalf("m=%v w=%v", m, w)
	}
}

func TestParseIOStat(t *testing.T) {
	s := "259:0 rbytes=1 wbytes=2\n\n259:1 rbytes=3\n"
	out := parseIOStat(s)
	if len(out) < 1 {
		t.Fatal(out)
	}
}

func TestConfidenceFor(t *testing.T) {
	if confidenceFor(MetricsBundle{CPUStat: map[string]int64{"a": 1}, MemoryCurrent: "5"}) != "high" {
		t.Fatal()
	}
	if confidenceFor(MetricsBundle{CPUStat: map[string]int64{"a": 1}}) != "medium" {
		t.Fatal()
	}
}

func TestBuildEvidenceUsesFixedClockInTest(t *testing.T) {
	t.Setenv("KHR_TEST_TELEMETRY_NOW", "2026-05-15T15:00:00Z")
	ev := BuildEvidence([]string{"w1"}, nil, MetricsBundle{CPUStat: map[string]int64{"x": 1}, MemoryCurrent: "10"})
	if ev.ObservedAt != "2026-05-15T15:00:00Z" {
		t.Fatalf("got %q", ev.ObservedAt)
	}
	_ = time.Now()
}
