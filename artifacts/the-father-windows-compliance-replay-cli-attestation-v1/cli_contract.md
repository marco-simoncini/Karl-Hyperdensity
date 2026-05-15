# CLI Contract

CLI: `karl-fluid-compliance-replay`

Flags:

- `-input` (required)
- `-evaluation-time` (RFC3339, deterministic mode)
- `-emit-attestation`
- `-attestation-mode` (`unsigned-dev` | `future-signable`)
- `-pretty`

Output includes:

- replay identity (`replayId`, `inputRef`, `evaluationTime`)
- compliance decision (`compliancePhase`, `hyperdensityReady`, `risk`)
- VM context (`vmRef`, `namespace`, `shellRef`, `poolContext`)
- blocker/remediation vectors
- audit hashes (`evidenceHash`, `replayHash`) and `auditRefs`
- `mutationFlags` all false in read-only mode
