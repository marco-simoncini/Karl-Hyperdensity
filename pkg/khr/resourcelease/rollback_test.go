package resourcelease

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRollbackBaseline(t *testing.T) {
	dir := t.TempDir()
	marker := filepath.Join(dir, "apply-marker.txt")
	if err := os.WriteFile(marker, []byte("v1"), 0o644); err != nil {
		t.Fatal(err)
	}
	bl, err := CaptureBaseline("b1", dir)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(marker, []byte("v2"), 0o644); err != nil {
		t.Fatal(err)
	}
	res := RollbackBaseline(bl)
	if !res.RolledBack {
		t.Fatalf("res=%+v", res)
	}
	raw, _ := os.ReadFile(marker)
	if string(raw) != "v1" {
		t.Fatalf("marker=%q", raw)
	}
}
