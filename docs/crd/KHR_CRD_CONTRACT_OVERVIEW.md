# KARL CRD Contract Overview (v1alpha1)

**Track:** Sprint 2 — contract / schema only.  
**Branch:** `KHR`  
**Repository:** `marco-simoncini/Karl-Hyperdensity`

## Purpose

This sprint delivers **installable CustomResourceDefinition (CRD) manifests** and **example Custom Resources** for the KARL execution and market model:

- **Shell** is the product-facing workload identity.
- **Cell** is the node-materialized execution unit (future KHR apply target).
- **Hyperdensity** objects (`ResourceLease`, `ResourceFuture`) express **market contracts** between donors and receivers.
- **Gateway** objects model **ingress** toward RDP / RemoteApp (`rdp-GW` alignment) without changing gateway code in this sprint.

No controllers, daemons, ISO changes, or Dashboard runtime changes are included.

## API layout

| API group | Scope | Kinds |
|-----------|-------|-------|
| `runtime.karl.io/v1alpha1` | Mixed | `Shell`, `ShellLease` (Namespaced); `ShellClass`, `ShellPool`, `Host`, `HostPool`, `RuntimeProvider`, `ResourcePort` (Cluster) |
| `hyperdensity.karl.io/v1alpha1` | Namespaced | `ResourceLease`, `ResourceFuture` |
| `gateway.karl.io/v1alpha1` | Namespaced | `GatewayRoute`, `WindowsApp` |

## File layout

```text
api/crds/runtime.karl.io/*.yaml
api/crds/hyperdensity.karl.io/*.yaml
api/crds/gateway.karl.io/*.yaml
examples/crds/*.yaml
docs/crd/*.md
```

## Installation (cluster admins)

Apply CRDs in dependency order (CRDs have no inter-CRD install order, any order is fine):

```bash
kubectl apply -f api/crds/runtime.karl.io/
kubectl apply -f api/crds/hyperdensity.karl.io/
kubectl apply -f api/crds/gateway.karl.io/
```

**Client-side validation** (no cluster required beyond kube-apiserver not needed for dry-run=client):

```bash
find api/crds -name '*.yaml' -print0 | xargs -0 -I{} kubectl apply --dry-run=client -f {}
```

If `kubectl` is not available, rely on CI / review; the repo `scripts/validate.sh` invokes `scripts/validate_crds.sh` when `kubectl` is present.

## Design rules

1. **Kubernetes remains the declarative substrate** — these are CRDs, not a parallel API server.
2. **KHR is not shipped here** — `Host` / `Cell` are contracts for future agent reconciliation.
3. **KubeVirt is not modified** — legacy mapping may appear as **opaque** `providerHandle` or optional `spec.kubeVirtLegacy` hints on `Shell`.
4. **ResourcePort is capability truth** — Hyperdensity must reject impossible `ResourceLease` modes when wired (Sprint 6+).
5. **v1alpha1** — breaking changes are expected; promotion to `v1beta1` requires conversion webhooks and compatibility gates.

## Cross-references

- Architecture: `docs/architecture/KARL_HOST_RUNTIME_VISION.md`, `KARL_SHELL_CELL_MODEL.md`, `KARL_HYPERDENSITY_KHR_FUSION.md`
- ADRs: `docs/adr/ADR-0001-khr-shell-cell-runtime-model.md` through `ADR-0003-*.md`

## Example apply order (optional sandbox)

Example manifests are **syntax-checked** in CI; **semantic** validation against OpenAPI requires CRDs to be installed on a cluster first.

Recommended order when applying everything to a dev cluster:

1. Install all CRDs (`api/crds/**`).
2. `examples/crds/resourceport-vmlike-hotplug-prewired.yaml`
3. `examples/crds/resourceport-linux-envelope.yaml`
4. `examples/crds/shell-windows-desktop-kubevirt-legacy.yaml` (Host / HostPool / ShellPool / Windows desktop Shell)
5. `examples/crds/shell-linux-cell-systemd.yaml`
6. `examples/crds/cell-linux-cgroup-envelope.yaml`
7. `examples/crds/resourcelease-cpu-burst.yaml`
8. `examples/crds/resourcefuture-predictive-memory.yaml`
9. `examples/crds/windowsapp-example.yaml`
10. `examples/crds/shell-windows-app-remoteapp.yaml`
11. `examples/crds/gatewayroute-*.yaml`
12. `examples/crds/cell-native-vmlike-fluidvirt-candidate.yaml` (depends on `pool-linux-dev-example` from step 5)

## Forward pointers (Sprint 3+)

- Sprint 3: KubeVirt legacy wrapper controller reads `Shell` + emits/labels existing VM objects.
- Sprint 5: KHR MVP consumes `Cell`, `Host`, and approved `ResourceLease` documents.
