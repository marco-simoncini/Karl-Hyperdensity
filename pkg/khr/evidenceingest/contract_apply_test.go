package evidenceingest

import (
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestPreparedYAMLHasNoApplyAuthorizationTokens(t *testing.T) {
	dir := writeAlignedFixture(t)
	opts := DefaultPrepareOptions()
	out, err := PrepareIngestRequest(
		filepath.Join(dir, "bundle.json"),
		filepath.Join(dir, "manifest.json"),
		filepath.Join(dir, "digest.txt"),
		opts,
	)
	if err != nil {
		t.Fatal(err)
	}
	low := strings.ToLower(string(out))
	for _, tok := range []string{
		"allowapply",
		"applyenabled",
		"unsafeapply",
		"wouldapply",
		"applyauthorized",
	} {
		if strings.Contains(low, tok) {
			t.Fatalf("unexpected token %q in output", tok)
		}
	}
	var doc map[string]interface{}
	if err := yaml.Unmarshal(out, &doc); err != nil {
		t.Fatal(err)
	}
	sp := doc["spec"].(map[string]interface{})
	if sp["dryRunOnly"] != false {
		t.Fatalf("expected dryRunOnly false by default")
	}
}
