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

## Legacy KubeVirt mapping

`Shell.spec.kubeVirtLegacy` is an **optional hint surface** for Sprint 3 controllers. It must not be interpreted as changing KubeVirt behavior by itself.

## FluidVirt alignment

`cell-native-vmlike-fluidvirt-candidate.yaml` references `RuntimeProvider` id `vm.native.fluidvirt.candidate.v1` as a **placeholder** for the fork/research repository `marco-simoncini/FluidVirt`. Hardening and rename are deferred.
