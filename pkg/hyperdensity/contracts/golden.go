package contracts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

// CanonicalSummaryJSON returns stable JSON for ParentFabricSummary for golden tests
// and extraction aids. Format: UTF-8, two-space indent, no HTML escaping in strings
// (encoding/json Encoder with SetEscapeHTML(false)), no trailing newline.
//
// Field order follows the struct declaration order in summary.go (encoding/json default).
func CanonicalSummaryJSON(summary ParentFabricSummary) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(summary); err != nil {
		return nil, fmt.Errorf("canonical summary json: %w", err)
	}
	out := buf.Bytes()
	// json.Encoder.Encode always appends a final newline; strip for stable byte comparison.
	if len(out) > 0 && out[len(out)-1] == '\n' {
		out = out[:len(out)-1]
	}
	return out, nil
}

// WriteCanonicalSummary writes CanonicalSummaryJSON(summary) to path, then a single
// trailing newline (POSIX text file convention).
func WriteCanonicalSummary(path string, summary ParentFabricSummary) error {
	b, err := CanonicalSummaryJSON(summary)
	if err != nil {
		return err
	}
	data := append(append([]byte(nil), b...), '\n')
	return os.WriteFile(path, data, 0o644)
}

// CompareSummaryGolden compares actual to golden JSON bytes by canonicalizing both
// summaries (re-marshal through CanonicalSummaryJSON). Extra whitespace or key order
// in golden is normalized as long as it unmarshals to ParentFabricSummary.
func CompareSummaryGolden(actual ParentFabricSummary, golden []byte) error {
	want, err := ParseParentFabricSummary(golden)
	if err != nil {
		return fmt.Errorf("parse golden: %w", err)
	}
	gotCanon, err := CanonicalSummaryJSON(actual)
	if err != nil {
		return fmt.Errorf("canonicalize actual: %w", err)
	}
	wantCanon, err := CanonicalSummaryJSON(want)
	if err != nil {
		return fmt.Errorf("canonicalize golden: %w", err)
	}
	if !bytes.Equal(gotCanon, wantCanon) {
		return fmt.Errorf("golden mismatch:\n--- got (canonical)\n%s\n--- want (canonical)\n%s", gotCanon, wantCanon)
	}
	return nil
}
