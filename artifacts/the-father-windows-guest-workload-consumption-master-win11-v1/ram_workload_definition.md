# RAM Workload Definition

Guest memory workload executed via QGA `guest-exec` + PowerShell/.NET.

- bounded allocation request: 512 MiB (`536870912` bytes)
- chunked allocations: 64 MiB arrays
- touch pages to enforce real commitment
- short hold, then explicit release + GC
- output JSON includes:
  - visibleMemoryBefore
  - availableMemoryBefore
  - allocationBytesRequested
  - allocationBytesSucceeded
  - availableMemoryDuring
  - availableMemoryAfterRelease
  - success/errors
