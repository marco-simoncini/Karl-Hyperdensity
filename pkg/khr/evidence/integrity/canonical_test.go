package integrity

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestCanonicalJSONStable(t *testing.T) {
	v := map[string]any{
		"zebra": 1,
		"alpha": "x",
		"beta":  []string{"a", "b"},
	}
	b1, err := CanonicalJSON(v)
	if err != nil {
		t.Fatal(err)
	}
	b2, err := CanonicalJSON(v)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(b1, b2) {
		t.Fatalf("canonical mismatch:\n%s\n%s", b1, b2)
	}
	// Map keys are sorted in encoding/json.
	want := `{"alpha":"x","beta":["a","b"],"zebra":1}`
	if string(b1) != want {
		t.Fatalf("got %q want %q", b1, want)
	}
	var enc bytes.Buffer
	e := json.NewEncoder(&enc)
	e.SetEscapeHTML(false)
	_ = e.Encode(v)
	raw := enc.Bytes()
	if len(raw) > 0 && raw[len(raw)-1] == '\n' {
		raw = raw[:len(raw)-1]
	}
	if string(raw) != want {
		t.Fatalf("encoder baseline %q", raw)
	}
}
