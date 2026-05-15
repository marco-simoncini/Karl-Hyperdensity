package evidence

import (
	"sort"
	"sync"
)

// Store is an in-memory evidence index (no persistence, no API).
type Store struct {
	mu               sync.RWMutex
	indices          []EvidenceIndex
	shaFirstPosition map[string]int // first seen index for sha (for duplicate counting)
	duplicateTotal   int
}

// NewStore returns an empty store.
func NewStore() *Store {
	return &Store{
		indices:          nil,
		shaFirstPosition: map[string]int{},
	}
}

// StoreBundle appends an index row. Duplicate bundleSha256 increments duplicateTotal.
func (s *Store) StoreBundle(idx EvidenceIndex) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if idx.BundleSha256 != "" {
		if _, ok := s.shaFirstPosition[idx.BundleSha256]; ok {
			s.duplicateTotal++
		} else {
			s.shaFirstPosition[idx.BundleSha256] = len(s.indices)
		}
	}
	s.indices = append(s.indices, idx)
}

// DuplicateTotal returns ingest-time duplicate SHA-256 observations (before dedupe pass).
func (s *Store) DuplicateTotal() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.duplicateTotal
}

// DeduplicateBySha256 collapses indices to one row per bundleSha256, keeping the latest IndexedAt.
// Returns the number of removed rows.
func (s *Store) DeduplicateBySha256() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.indices) == 0 {
		return 0
	}
	best := map[string]EvidenceIndex{}
	order := []string{}
	for _, idx := range s.indices {
		sha := idx.BundleSha256
		if sha == "" {
			sha = "__empty__" + idx.ArtifactID + idx.IndexedAt
		}
		if _, seen := best[sha]; !seen {
			order = append(order, sha)
		}
		prev := best[sha]
		if prev.IndexedAt == "" || idx.IndexedAt >= prev.IndexedAt {
			best[sha] = idx
		}
	}
	out := make([]EvidenceIndex, 0, len(order))
	for _, sha := range order {
		out = append(out, best[sha])
	}
	removed := len(s.indices) - len(out)
	s.indices = out
	s.shaFirstPosition = map[string]int{}
	for i := range s.indices {
		if s.indices[i].BundleSha256 != "" {
			s.shaFirstPosition[s.indices[i].BundleSha256] = i
		}
	}
	return removed
}

// Snapshot returns a copy of all indices (sorted by IndexedAt for stability).
func (s *Store) Snapshot() []EvidenceIndex {
	s.mu.RLock()
	defer s.mu.RUnlock()
	cp := append([]EvidenceIndex(nil), s.indices...)
	sort.Slice(cp, func(i, j int) bool {
		if cp[i].IndexedAt == cp[j].IndexedAt {
			return cp[i].ArtifactID < cp[j].ArtifactID
		}
		return cp[i].IndexedAt < cp[j].IndexedAt
	})
	return cp
}

// Len returns the number of index rows.
func (s *Store) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.indices)
}
