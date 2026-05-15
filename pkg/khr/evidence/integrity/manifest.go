package integrity

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

// ChainOfCustody records non-authoritative local process context (not a security boundary).
type ChainOfCustody struct {
	Hostname       string `json:"hostname,omitempty"`
	User           string `json:"user,omitempty"`
	ProcessID      int    `json:"processId"`
	ExecutablePath string `json:"executablePath,omitempty"`
	WorkingDir     string `json:"workingDirectory,omitempty"`
}

// ArtifactManifest is the sidecar JSON for an evidence bundle (local integrity, Sprint 10).
type ArtifactManifest struct {
	ArtifactID         string         `json:"artifactId"`
	AgentID            string         `json:"agentId"`
	GeneratedAt        string         `json:"generatedAt"`
	BundleSha256       string         `json:"bundleSha256"`
	BundleBytes        int            `json:"bundleBytes"`
	SigningMode        string         `json:"signingMode"`
	SignaturePresent   bool           `json:"signaturePresent"`
	SourceMode         string         `json:"sourceMode"`
	MutationsForbidden bool           `json:"mutationsForbidden"`
	ChainOfCustody     ChainOfCustody `json:"chainOfCustody"`
	SignatureAlgorithm string         `json:"signatureAlgorithm,omitempty"`
	SignatureBase64    string         `json:"signature,omitempty"`
	IntegrityNotes     []string       `json:"integrityNotes,omitempty"`
}

// BuildChainOfCustody captures best-effort local metadata (never enables apply or transport).
func BuildChainOfCustody() ChainOfCustody {
	if strings.TrimSpace(os.Getenv("KHR_TEST_INTEGRITY_CHAIN_STUB")) == "1" {
		return ChainOfCustody{
			Hostname:       "test-host.example",
			User:           "test-user",
			ProcessID:      4242,
			ExecutablePath: "/tmp/khr-linux-agent.test",
			WorkingDir:     "/tmp",
		}
	}
	c := ChainOfCustody{ProcessID: os.Getpid()}
	if h, err := os.Hostname(); err == nil {
		c.Hostname = h
	}
	if u, err := user.Current(); err == nil {
		c.User = u.Username
	}
	if len(os.Args) > 0 {
		c.ExecutablePath = os.Args[0]
	}
	if wd, err := os.Getwd(); err == nil {
		c.WorkingDir = wd
	}
	return c
}

func generatedAtRFC3339() time.Time {
	s := strings.TrimSpace(os.Getenv("KHR_TEST_INTEGRITY_NOW"))
	if s == "" {
		return time.Now().UTC()
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		t2, err2 := time.Parse(time.RFC3339Nano, s)
		if err2 != nil {
			return time.Now().UTC()
		}
		return t2.UTC()
	}
	return t.UTC()
}

// BuildManifest assembles the manifest for a canonical bundle payload.
func BuildManifest(agentID, artifactID, signingMode string, canonical []byte, bundleSHA string, sigB64, sigAlg string) ArtifactManifest {
	gen := generatedAtRFC3339().Format(time.RFC3339)
	m := ArtifactManifest{
		ArtifactID:         strings.TrimSpace(artifactID),
		AgentID:            agentID,
		GeneratedAt:        gen,
		BundleSha256:       bundleSHA,
		BundleBytes:        len(canonical),
		SigningMode:        signingMode,
		SignaturePresent:   sigB64 != "",
		SourceMode:         "collect-evidence",
		MutationsForbidden: true,
		ChainOfCustody:     BuildChainOfCustody(),
		IntegrityNotes: []string{
			"Local integrity metadata only; not a production signature or admission decision.",
			"No transport, no ingest API, no apply: bundle remains read-only evidence.",
		},
	}
	if sigB64 != "" {
		m.SignatureBase64 = sigB64
		m.SignatureAlgorithm = sigAlg
	}
	return m
}

// MarshalManifestJSON returns indented JSON with trailing newline for file emission.
func MarshalManifestJSON(m ArtifactManifest) ([]byte, error) {
	return json.MarshalIndent(m, "", "  ")
}

// ValidateSigningMode returns an error for unknown modes.
func ValidateSigningMode(mode string) error {
	switch strings.TrimSpace(mode) {
	case "", "none", "local-dev":
		return nil
	default:
		return fmt.Errorf("unknown signing-mode %q (allowed: none, local-dev)", strings.TrimSpace(mode))
	}
}

// NormalizeSigningMode returns "none" when empty.
func NormalizeSigningMode(mode string) string {
	if strings.TrimSpace(mode) == "" {
		return "none"
	}
	return strings.TrimSpace(mode)
}

// WriteDigestFile writes bundleSha256 as a single line plus newline.
func WriteDigestFile(path, shaHex string) error {
	if strings.TrimSpace(path) == "" {
		return nil
	}
	dir := filepath.Dir(path)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	return os.WriteFile(path, []byte(strings.TrimSpace(shaHex)+"\n"), 0o600)
}

// WriteManifestFile writes manifest JSON with mode 0o600.
func WriteManifestFile(path string, m ArtifactManifest) error {
	if strings.TrimSpace(path) == "" {
		return nil
	}
	b, err := MarshalManifestJSON(m)
	if err != nil {
		return err
	}
	b = append(b, '\n')
	dir := filepath.Dir(path)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	return os.WriteFile(path, b, 0o600)
}
