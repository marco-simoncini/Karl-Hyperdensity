package evidenceingest

// PrepareOptions configures local EvidenceIngestRequest generation (no network).
type PrepareOptions struct {
	Format                 string
	DryRunOnly             bool
	Namespace              string
	Name                   string
	NodeName               string
	HostID                 string
	Tenant                 string
	RequireDigestMatch     bool
	AllowUnsigned          bool
	AllowLocalDevSignature bool
}
