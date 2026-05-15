package evidenceingest

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCRDFilesContainOpenAPISchema(t *testing.T) {
	root := filepath.Join("..", "..", "..", "api", "crds", "hyperdensity.karl.io")
	for _, name := range []string{"evidencebundle.yaml", "evidenceingestrequest.yaml"} {
		p := filepath.Join(root, name)
		b, err := os.ReadFile(p)
		if err != nil {
			t.Fatal(err)
		}
		s := string(b)
		if !strings.Contains(s, "openAPIV3Schema") {
			t.Fatalf("%s: missing openAPIV3Schema", name)
		}
		if !strings.Contains(s, "hyperdensity.karl.io") {
			t.Fatalf("%s: missing group", name)
		}
	}
}
