# Fake Runtime Boundary

- `fakeRuntimeOnly=true`
- `usesTemporaryFilesOnly=true`
- `touchesRealCgroup=false`
- `touchesRealQmp=false`
- `touchesRealQga=false`
- `touchesHostRuntime=false`
- `requiresNoPrivileges=true`
- `deterministicReplay=true`
- `safeForCi=true`

No privileged runtime operations are allowed in this milestone.
