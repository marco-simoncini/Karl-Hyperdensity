package shellpassport_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/shellpassport"
)

func repoRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}

func TestSprint2SchemaFilesExist(t *testing.T) {
	root := repoRoot(t)
	for _, name := range shellpassport.SchemaFilesRequiredSprint2() {
		if _, err := os.Stat(filepath.Join(root, "schemas", name)); err != nil {
			t.Fatalf("missing schema %s: %v", name, err)
		}
	}
}

func TestSprint2ExamplesValidate(t *testing.T) {
	if err := shellpassport.ValidateSprint2Examples(repoRoot(t)); err != nil {
		t.Fatalf("ValidateSprint2Examples: %v", err)
	}
}

func TestBlockedShellCannotBeEligible(t *testing.T) {
	sh := map[string]interface{}{
		"shellId": "x", "supportProfile": "p",
		"claimBoundary": []interface{}{"ok"},
		"shellKind": "vm_linux", "targetRef": "ns/VirtualMachine/vm",
		"blocked": true, "donorEligibility": "eligible", "receiverEligibility": "blocked",
	}
	if err := shellpassport.ValidateShellRegistrySurface(map[string]interface{}{
		"registryId": shellpassport.RegistryID,
		"shells":     []interface{}{sh},
	}); err == nil {
		t.Fatal("expected blocked+eligible rejection")
	}
}

func TestExcludedShellKindRejected(t *testing.T) {
	sh := map[string]interface{}{
		"shellId": "x", "supportProfile": "p",
		"claimBoundary": []interface{}{"ok"},
		"shellKind": "container_windows", "targetRef": "ns/Deployment/x",
		"blocked": false, "donorEligibility": "eligible", "receiverEligibility": "blocked",
	}
	if err := shellpassport.ValidateShellRegistrySurface(map[string]interface{}{
		"registryId": shellpassport.RegistryID,
		"shells":     []interface{}{sh},
	}); err == nil {
		t.Fatal("expected excluded kind rejection")
	}
}

func TestCapabilityEvidenceFluidVirtOnly(t *testing.T) {
	doc := map[string]interface{}{
		"source": "Dashboard", "actuator": "Dashboard",
		"rebootRequired": false, "recreateRequired": false, "migrationRequired": false,
		"runtimeMutationSupportedFamilies": []interface{}{"cpu_entitlement"},
		"unsupportedFamilies":              []interface{}{},
	}
	if err := shellpassport.ValidateShellCapabilityEvidence(doc); err == nil {
		t.Fatal("expected non-FluidVirt rejection")
	}
}
