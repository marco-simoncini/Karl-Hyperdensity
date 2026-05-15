package recommendation

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	gpevidence "github.com/marco-simoncini/Karl-Hyperdensity/pkg/grandepadre/evidence"
)

// CollectIngestPaths merges explicit files and optional directory globs (*.yaml, *.yml only).
func CollectIngestPaths(files []string, dir string) ([]string, error) {
	seen := map[string]bool{}
	var out []string
	for _, f := range files {
		f = strings.TrimSpace(f)
		if f == "" || seen[f] {
			continue
		}
		seen[f] = true
		out = append(out, f)
	}
	d := strings.TrimSpace(dir)
	if d != "" {
		for _, pat := range []string{"*.yaml", "*.yml"} {
			matches, err := filepath.Glob(filepath.Join(d, pat))
			if err != nil {
				return nil, err
			}
			for _, m := range matches {
				if seen[m] {
					continue
				}
				seen[m] = true
				out = append(out, m)
			}
		}
	}
	sort.Strings(out)
	if len(out) == 0 {
		return nil, fmt.Errorf("no ingest paths: pass -ingest-request-input one or more times and/or -ingest-request-dir")
	}
	return out, nil
}

// IngestAllIntoStore reads each path and ingests into the store (no dedupe until BuildActionSlate).
func IngestAllIntoStore(s *gpevidence.Store, paths []string, pol gpevidence.UnsignedDigestTrustPolicy) error {
	for _, p := range paths {
		b, err := os.ReadFile(p)
		if err != nil {
			return fmt.Errorf("read %s: %w", p, err)
		}
		if _, err := s.Ingest(b, pol); err != nil {
			return fmt.Errorf("ingest %s: %w", p, err)
		}
	}
	return nil
}
