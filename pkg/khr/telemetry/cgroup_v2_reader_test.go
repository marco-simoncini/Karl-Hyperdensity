package telemetry

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/cgroup"
)

func TestReadCgroupV2MetricsFull(t *testing.T) {
	root := t.TempDir()
	write(t, filepath.Join(root, "cpu.stat"), "usage_usec 42\nnr_periods 1\n")
	write(t, filepath.Join(root, "memory.current"), "4096\n")
	write(t, filepath.Join(root, "memory.max"), "max\n")
	write(t, filepath.Join(root, "memory.events"), "pgfault 3\n")
	write(t, filepath.Join(root, "io.stat"), "8:0 rbytes=1\n\n8:1 wbytes=2\n")
	res, _, b := cgroup.ValidateCgroupPathForTelemetry(root, "")
	if len(b) != 0 {
		t.Fatal(b)
	}
	m, w, bl := ReadCgroupV2Metrics(res)
	if len(bl) != 0 {
		t.Fatal(bl)
	}
	if m.MemoryCurrent != "4096" || m.CPUStat["usage_usec"] != 42 {
		t.Fatalf("m=%+v w=%v", m, w)
	}
	if len(m.IOStat) == 0 {
		t.Fatal("expected io.stat")
	}
}

func TestReadCgroupV2MetricsMissingBothBlocked(t *testing.T) {
	root := t.TempDir()
	write(t, filepath.Join(root, "memory.max"), "max\n")
	res, _, b := cgroup.ValidateCgroupPathForTelemetry(root, "")
	if len(b) != 0 {
		t.Fatal(b)
	}
	_, _, bl := ReadCgroupV2Metrics(res)
	if len(bl) == 0 {
		t.Fatal("expected blocked when cpu.stat and memory.current missing")
	}
}

func TestReadCgroupV2MetricsMissingFileWarning(t *testing.T) {
	root := t.TempDir()
	write(t, filepath.Join(root, "cpu.stat"), "usage_usec 1\n")
	res, _, b := cgroup.ValidateCgroupPathForTelemetry(root, "")
	if len(b) != 0 {
		t.Fatal(b)
	}
	_, w, bl := ReadCgroupV2Metrics(res)
	if len(bl) != 0 {
		t.Fatal(bl)
	}
	if len(w) == 0 {
		t.Fatal("expected warning for missing memory.current")
	}
}

func write(t *testing.T, p, content string) {
	t.Helper()
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
