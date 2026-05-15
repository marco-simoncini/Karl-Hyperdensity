# Evidence ingest security boundaries

## What this contract does **not** grant

- **Not authorization:** digest match and signatures prove **integrity** (or dev placeholders), not identity of a trusted operator or service account.
- **Not admission:** `Accepted` phases (future controller semantics) are **workflow** states, not Kubernetes `ValidatingAdmissionPolicy` replacements unless explicitly wired elsewhere.
- **Not transport security:** TLS, mTLS, and OIDC are out of scope for Sprint 11 contract files.

## local-dev signatures

`signingMode: local-dev` uses Ed25519 keys held locally (see Sprint 10). They are labeled **DevOnly** in annotations and `signatureStatus`. They **must not** be treated as production code-signing or PKI.

## Policy knobs

`EvidenceIngestRequest.spec.policy` allows a controller to:

- Reject unsigned bundles when `allowUnsigned: false`.
- Reject digest mismatches when `requireDigestMatch: true`.
- Gate `local-dev` with `allowLocalDevSignature`.

Defaults for the CLI stub lean toward **strict digest** with **unsigned allowed** when no production signature exists—operators tighten policy in-cluster.

## blast-radius

Evidence may include paths and host metadata (`chainOfCustody` in manifests). Store and replicate according to tenant data classification; this contract does not define retention or encryption at rest.
