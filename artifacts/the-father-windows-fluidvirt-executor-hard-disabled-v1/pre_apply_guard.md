# Pre-Apply Guard

`WindowsFluidPreApplyGuard` enforces:

- `executorEnabled=false`
- `mutationWindowOpen=false`
- `qmpMutationAllowed=false`
- `clusterMutationAllowed=false`

Guard phase can be ready/blocked/quarantined/needs revalidation, but executor remains disabled in every case.
