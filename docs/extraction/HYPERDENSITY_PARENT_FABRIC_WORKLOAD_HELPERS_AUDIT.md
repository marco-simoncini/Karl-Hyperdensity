# Hyperdensity Parent Fabric — workload helpers audit (Sprint 48)

## Source file

| | |
|--|--|
| **Dashboard path** | `pkg/server/hyperdensity_parent_fabric_workload_helpers.go` |
| **Package** | `server` |
| **Lines** | ~729 |
| **Type definitions** | **0** |
| **Functions** | **46** |
| **Const block** | Execution mode / mechanism string literals (in-file) |

## Import set (source)

```go
import (
	"fmt"
	"strings"
	"time"
)
```

Stdlib only in the **import block**. No `k8s.io/*`, `net/http`, or `client-go` imports.

## Forbidden / coupling signals (in source body)

| Signal | Examples |
|--------|----------|
| **KubeVirt API path literals** | `hyperdensityVirtualMachineAPIPath`, `…InstanceAPIPath`, `guestosinfo` subresource |
| **Kubernetes API path literals** | `hyperdensityAppsWorkloadAPIPath` (`/apis/apps/v1/…`), pod/resize paths |
| **Same-package runtime symbols** | `hyperdensityExecutionActionBurst`, `hyperdensityNestedSlice`, `HyperdensityPilotObservedState`, `hyperdensityPilotObservedStateFromDeployment` (defined in other `.go` files) |
| **Untyped workload maps** | `map[string]interface{}` observation helpers tied to live object shapes |
| **Apply / execution coupling** | Mode selection, ready reasons, mechanism mapping for live apply |

## Candidate functions (stdlib-shaped only — not copied in Sprint 48)

| Function | Note |
|----------|------|
| `hyperdensityAppsWorkloadResource` | Kind → apps resource string |
| `hyperdensityPilotWorkloadTerm` | Kind → display term |
| `hyperdensityExecutionSupportsLiveApplyKind` | Kind allowlist |

Too small/isolated to justify a package copy without the rest of the file; copying would misrepresent the module boundary.

## Excluded functions (deferred — all others)

| Category | Functions (sample) |
|----------|-------------------|
| **API path builders** | `hyperdensityAppsWorkloadAPIPath`, `hyperdensityPodAPIPath`, `hyperdensityPodResizeAPIPath`, `hyperdensityVirtualMachine*APIPath` |
| **Execution mode / apply** | `hyperdensityExecutionModeForKind`, `hyperdensityNoRolloutExecutionModeForKind`, `hyperdensityExecutionReadyReason*`, `hyperdensityExecutionMechanismForMode`, … |
| **Mode capability predicates** | `hyperdensityExecutionModeSupportsWorkloadPatch`, `…PodResize`, `…VMLiveUpdate`, … |
| **Observed state from API objects** | `hyperdensityObservedWorkload*`, `hyperdensityObservedPod*`, `hyperdensityPilotObservedStateFrom*` |

## Verdict

**`copy-deferred`**

### Motivation

The file is **not** extractable as a stdlib-only Hyperdensity package without:

1. Duplicating or importing large **server-package** helper surfaces (`hyperdensityNestedSlice`, quantity parsers, execution action constants).
2. Embedding **KubeVirt/Kubernetes API path** construction (forbidden for pure-core).
3. Carrying **runtime/apply** semantics (execution modes, pilot observed state, live object maps).

Sprint 48 delivers **audit + placeholder only** (`pkg/hyperdensity/parentfabric/workload/doc.go`). See **`HYPERDENSITY_PARENT_FABRIC_WORKLOAD_HELPERS_DEFERRED.md`**.

## Rollback

Remove `workload/doc.go` and audit/deferred docs; no Dashboard source changes required.

## Sprint 49 follow-up

Stdlib **`parentfabric/primitives`** added as prerequisite — does **not** change this file's **`copy-deferred`** verdict. See **`HYPERDENSITY_PARENT_FABRIC_WORKLOAD_PREREQUISITES.md`**.

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_PREREQUISITES.md`
- `HYPERDENSITY_PARENT_FABRIC_PRIMITIVES_CONTRACT.md`
- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_HELPERS_DEFERRED.md`
- `HYPERDENSITY_PARENT_FABRIC_EXECUTION_TYPES_DRIFT_GUARD.md`
- Dashboard `docs/hyperdensity/HYPERDENSITY_PARENT_FABRIC_WORKLOAD_HELPERS_AUDIT_M33.md`
