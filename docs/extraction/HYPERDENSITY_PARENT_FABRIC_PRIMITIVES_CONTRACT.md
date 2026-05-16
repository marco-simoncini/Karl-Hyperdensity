# Hyperdensity Parent Fabric — primitives contract (Sprint 49)

## Package

`pkg/hyperdensity/parentfabric/primitives` — **stdlib only**; no `k8s.io/apimachinery/resource.Quantity`.

## Nested map helpers

| Function | Behavior |
|----------|----------|
| `StringAt` | Walk path; coerce final value to string; `ok=false` if missing |
| `Int64At` | Numeric coercion; `ok=false` if not coercible |
| `Float64At` | Float coercion |
| `MapAt` | `map[string]interface{}` at path |
| `SliceAt` | `[]interface{}` at path |

## Quantity helpers

| Function | Supported inputs (minimal) |
|----------|---------------------------|
| `NormalizeCPUQuantity` | `100m`, whole cores `1`, `2` |
| `NormalizeMemoryQuantity` | `128Mi`, `1Gi`, `512Ki`, plain bytes `1000` |

Unknown input → `ok=false` (no panic).

## Golden

`pkg/hyperdensity/parentfabric/primitives/testdata/primitives_contract.golden.json`

## Limitations vs Dashboard

Dashboard `hyperdensityParseCPUQuantity` / `hyperdensityParseMemoryQuantity` use **Kubernetes** `resource.ParseQuantity`. Hyperdensity primitives implement a **narrow contract** for golden tests and future copy — **not** full K8s quantity semantics.

## Related

- `HYPERDENSITY_PARENT_FABRIC_WORKLOAD_PREREQUISITES.md`
- Dashboard `HYPERDENSITY_PARENT_FABRIC_PRIMITIVES_AUDIT_M34.md`
