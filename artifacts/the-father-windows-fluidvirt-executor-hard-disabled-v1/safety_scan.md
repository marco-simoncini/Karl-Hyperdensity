# Safety Scan

Pattern scan executed on changed files.

- No secrets/tokens/kubeconfig/private keys found.
- No deploy commands (`kubectl apply/patch`, `helm upgrade`) found.
- No frontend/dashboard scope files changed.
- Forbidden QMP command strings appear only in policy docs as rejected/non-executed references.
