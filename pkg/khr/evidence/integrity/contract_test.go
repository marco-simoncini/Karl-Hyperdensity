package integrity

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestArtifactManifestJSONHasNoApplyTriggerFields(t *testing.T) {
	t.Setenv("KHR_TEST_INTEGRITY_NOW", "2026-07-01T12:00:00Z")
	t.Setenv("KHR_TEST_INTEGRITY_CHAIN_STUB", "1")
	m := BuildManifest("ag", "art", "none", []byte(`{}`), SHA256Hex([]byte(`{}`)), "", "")
	b, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	s := strings.ToLower(string(b))
	for _, forbidden := range []string{`"allowapply"`, `"unsafeapply"`, `"applyenabled"`, `"wouldapply"`} {
		if strings.Contains(s, forbidden) {
			t.Fatalf("unexpected token %s in %s", forbidden, b)
		}
	}
}
