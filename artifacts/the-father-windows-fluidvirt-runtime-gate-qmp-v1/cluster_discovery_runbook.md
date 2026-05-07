# Cluster Discovery Runbook

Runbook location:

- `docs/runbooks/windows-fluidvirt-readonly-cluster-discovery-v1.md`

Highlights:

- read-only kubectl commands only
- strict no-deploy/no-mutation policy
- classification outcomes for:
  - existing Windows VM candidate
  - replica pool context-only
  - blocked missing annotations/guest/qmp/migration

Sanitization policy:

- no token dump
- no kubeconfig dump
- no secret/bearer disclosure
