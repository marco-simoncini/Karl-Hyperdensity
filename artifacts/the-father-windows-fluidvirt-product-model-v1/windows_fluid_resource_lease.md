# WindowsFluidResourceLease

Lease model supports:

- `cpu-entitlement`
- `ram-balloon`
- `combined-envelope`

Safety:

- rollback target required
- return-to-floor target required
- guest ACK evidence required
- VM spec patch and vCPU hotplug requests rejected
