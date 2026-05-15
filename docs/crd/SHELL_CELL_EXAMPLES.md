# Shell / Cell Example Set (v1alpha1)

**Location:** `examples/crds/`  
**Namespace convention:** examples use `karl-sandbox` for namespaced objects. Adjust before apply.

## How to read these files

Each file is **contract illustration** only. Controllers shipped in later sprints will enforce referential integrity (e.g., that `ShellClass` exists, that `ResourceLease.mode` is allowed by `ResourcePort`).

## File map

| File | Illustrates |
|------|-------------|
| `shell-windows-desktop-kubevirt-legacy.yaml` | Windows desktop Shell backed by **legacy KubeVirt** hints (`kubeVirtLegacy` block is optional metadata for Sprint 3 mapper). |
| `shell-windows-app-remoteapp.yaml` | Windows **App Shell** with `WindowsApp` reference and gateway route refs. |
| `shell-linux-cell-systemd.yaml` | Linux workspace Shell preferring **systemd** provider with cgroup-friendly profile. |
| `cell-linux-cgroup-envelope.yaml` | `Cell` bound to a Linux Shell with **inline** cgroup envelope ResourcePort snapshot. |
| `cell-native-vmlike-fluidvirt-candidate.yaml` | VM-like Cell using **vm-native-candidate** driver (FluidVirt R&D alignment — not a GA claim). |
| `resourceport-linux-envelope.yaml` | Cluster `ResourcePort` for **Linux cgroup envelope** CPU/RAM modes. |
| `resourceport-vmlike-hotplug-prewired.yaml` | VM-like profile with **prewired** hot-add / balloon modes (capability truth, not universal promise). |
| `resourcelease-cpu-burst.yaml` | `ResourceLease` between two Cells for **CPU burst** envelope transfer. |
| `resourcefuture-predictive-memory.yaml` | `ResourceFuture` with **memory pressure** predicate and lease template. |
| `gatewayroute-rdp-desktop.yaml` | `GatewayRoute` for **full desktop** RDP via `karl.rdpgw.v1`. |
| `gatewayroute-remoteapp.yaml` | RemoteApp-oriented route with **remoteapplication\*** fields. |
| `windowsapp-example.yaml` | Catalog `WindowsApp` with RemoteApp mapping hints. |

## Multi-document bundles

Some files contain multiple resources separated by `---` so a reviewer can `kubectl apply -f <file>` after CRDs are installed (still no runtime side-effects from KARL controllers in Sprint 2).

## KubeVirt legacy provider (Sprint 3)

Contracts and provider-scoped examples:

| Path | Content |
|------|---------|
| `docs/providers/KUBEVIRT_*.md` | Legacy provider contract, Shell/Cell mapping, labels/handles, migration safety |
| `api/providers/kubevirt/*.yaml` | Declarative `providers.karl.io/v1alpha1` bundles (**not** `apiextensions.k8s.io` CRDs) |
| `examples/providers/kubevirt/*.yaml` | `Shell` / `Cell` / `RuntimeProvider` / `ResourcePort` / `ResourceLease` aligned with **kubevirt-legacy** |

**Apply note:** `examples/crds/shell-windows-desktop-kubevirt-legacy.yaml` and `examples/providers/kubevirt/runtimeprovider-kubevirt-legacy.yaml` both define cluster `RuntimeProvider` metadata `name: kubevirt-legacy-v1`. Apply **one** bundle or remove the duplicate document to avoid `AlreadyExists`.

**Prerequisites:** `examples/providers/kubevirt/shell-linux-vm-kubevirt-legacy.yaml` and related Cells reference `ShellPool` `pool-linux-dev-example` and `Host` `worker-01-example` from the `examples/crds/` seed bundles (`shell-linux-cell-systemd.yaml`, `shell-windows-desktop-kubevirt-legacy.yaml`).

## Legacy KubeVirt mapping

`Shell.spec.kubeVirtLegacy` is an **optional hint surface** for mappers (Sprint 3+). It does **not** change KubeVirt VM behavior by itself. See `docs/providers/KUBEVIRT_SHELL_CELL_MAPPING.md`.

## FluidVirt alignment

`cell-native-vmlike-fluidvirt-candidate.yaml` references `RuntimeProvider` id `vm.native.fluidvirt.candidate.v1` as a **placeholder** for the fork/research repository `marco-simoncini/FluidVirt`. Hardening and rename are deferred.
