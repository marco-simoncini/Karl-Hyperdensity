package windowsfluidvirt

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestMinimalFixturesAreValidJSONObjects(t *testing.T) {
	root := filepath.Join("..", "..", "examples", "windows-fluid-product-fixtures")
	fixtures := []string{
		"product_model_minimal.json",
		"action_slate_minimal.json",
		"blockers_minimal.json",
	}
	for _, fixture := range fixtures {
		data, err := os.ReadFile(filepath.Join(root, fixture))
		if err != nil {
			t.Fatalf("read %s: %v", fixture, err)
		}
		var payload map[string]any
		if err := json.Unmarshal(data, &payload); err != nil {
			t.Fatalf("decode %s: %v", fixture, err)
		}
		if len(payload) == 0 {
			t.Fatalf("fixture %s must not be empty", fixture)
		}
	}
}
