package integrity

import "testing"

func TestSHA256HexDeterministic(t *testing.T) {
	in := []byte(`{"alpha":"x","beta":["a","b"],"zebra":1}`)
	got := SHA256Hex(in)
	want := SHA256Hex(in)
	if got != want {
		t.Fatalf("non-deterministic")
	}
	if len(got) != 64 {
		t.Fatalf("hex length: %d", len(got))
	}
}
