# Hyperdensity / KHR — ResourceLease examples (Sprint 89)

## Summary

Three **documental** examples for minimal ResourceLease contract. Not applied to cluster.

---

## Example 1 — Windows DaaS KHR-native desired lease

```yaml
apiVersion: hyperdensity.karl.io/v1alpha1
kind: ResourceLease
metadata:
  name: winsrv-daas-001
spec:
  shell:
    kind: windowsDesktop
    ref: shell/winsrv-daas-001
  cell:
    ref: cell/winsrv-daas-001
  provider: khr.native
  resources:
    cpu: { request: "4", limit: "8" }
    memory: { request: "8Gi", limit: "16Gi" }
  storage:
    defaultDiscardPolicy: deleteOnStop
    promoteToImage: { enabled: true }
    disks:
      - name: os
        role: os
        mode: ephemeralOverlay
        source: { type: goldenImage, ref: golden/winsrv-2026-q1 }
        discardPolicy: deleteOnStop
      - name: profile
        role: profile
        mode: persistent
        source: { type: pvc, ref: profile-pvc-user-001 }
      - name: scratch
        role: scratch
        mode: scratch
        discardPolicy: deleteOnStop
  network:
    attachments:
      - name: tenant
        networkRef: karlnetwork/tenant-acme
        role: primary
        isolation: strict
    exposure:
      ingress: rdp-GW
      egress: controlled
      directPublic: false
  policy: { tenantIsolation: strict }
  expiration: { ttl: 8h }
status:
  phase: Pending
```

---

## Example 2 — Public cloud fallback lease

```yaml
apiVersion: hyperdensity.karl.io/v1alpha1
kind: ResourceLease
metadata:
  name: vmlike-aws-fallback-001
spec:
  shell:
    kind: vmLike
    ref: shell/vm-aws-001
  cell:
    ref: cell/vm-aws-001
  provider: kubevirt.public-cloud-fallback
  resources:
    cpu: { request: "2", limit: "4" }
    memory: { request: "4Gi", limit: "8Gi" }
  storage:
    disks:
      - name: root
        role: os
        mode: ephemeralOverlay
        source: { type: pvc, ref: karl-os-nfs }
      - name: data
        role: data
        mode: persistent
        source: { type: image, ref: ubuntu-22.04 }
  network:
    providerNetwork:
      provider: kubevirt.legacy.ovn
    attachments:
      - name: primary
        networkRef: segment/tenant-aws
        providerBinding: ovn.logicalPort
  policy: {}
status:
  phase: Pending
```

VM/VMPool compatibility via KubeVirt; Shell/Cell product model preserved.

---

## Example 3 — Baremetal native lease

```yaml
apiVersion: hyperdensity.karl.io/v1alpha1
kind: ResourceLease
metadata:
  name: baremetal-shell-001
spec:
  shell:
    kind: linuxShell
    ref: shell/bm-001
  cell:
    ref: cell/bm-001
  provider: khr.native
  storage:
    disks:
      - name: root
        role: os
        mode: ephemeralClone
        source: { type: snapshot, ref: snap/golden-linux-v1 }
  network:
    providerNetwork:
      provider: baremetal.bridge
    attachments:
      - name: lan
        networkRef: karlnetwork/bm-vlan-100
        providerBinding: baremetal.vlan
status:
  phase: Pending
```

Modes: `ephemeralOverlay` or `ephemeralClone`; network: `baremetal.bridge` or `baremetal.vlan`.

---

## Related

- `HYPERDENSITY_KHR_RESOURCELEASE_MINIMAL_CONTRACT.md`


---

## Sprint 90 (inventory facts)

Sprint 90 adds read-only KubeVirt and OVN/SDN capability inventory mapped to ResourceLease contract. No CRD, no controller, no runtime. See HYPERDENSITY_KHR_KUBEVIRT_CAPABILITY_INVENTORY.md and related Sprint 90 docs.
