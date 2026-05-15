package resourcelease

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func testRepo(t *testing.T) string {
	t.Helper()
	return filepath.Join("..", "..", "..")
}

func TestCellContextParseAllowedFixture(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join(testRepo(t), "examples", "khr", "golden", "inputs", "cell-context-allowed.json"))
	if err != nil {
		t.Fatal(err)
	}
	var ctx CellContext
	if err := json.Unmarshal(raw, &ctx); err != nil {
		t.Fatal(err)
	}
	if ctx.DonorPlatform != "linux" || ctx.ReceiverPlatform != "linux" {
		t.Fatalf("platforms %+v", ctx)
	}
	if ctx.DonorCell == nil || ctx.ReceiverCell == nil {
		t.Fatal("expected donor and receiver cells")
	}
	if ctx.DonorCell.Spec.RuntimeProviderRef.Name != "linux-systemd-v1" {
		t.Fatalf("donor provider %q", ctx.DonorCell.Spec.RuntimeProviderRef.Name)
	}
	if ctx.DonorRuntimeProvider == nil || ctx.DonorRuntimeProvider.Spec.Driver != "linux-systemd" {
		t.Fatal("expected donor runtime provider snapshot")
	}
}
