# CLI Summary

CLI: `cmd/karl-fluid-windows-executor`

Flags:

- `-target`, `-lease`, `-gate`, optional `-approval`
- `-mode plan|dry-run|apply-plan-only`
- `-out`, `-evaluation-time`, `-pretty`

Safety statement in help:

- controlled apply planning only
- no autonomous apply
- no runtime mutation
- no vCPU hotplug
- no logical CPU scaling
