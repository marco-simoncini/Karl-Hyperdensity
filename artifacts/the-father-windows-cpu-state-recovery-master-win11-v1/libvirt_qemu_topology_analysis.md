# Libvirt/QEMU Topology Analysis

Domain: `karl_master-win11`

From `dumpxml` / `dominfo` / `vcpucount` / `domstats` / QMP:

- current vCPU (runtime): `7`
- config vCPU baseline: `6`
- max vCPU: `24`
- topology in QEMU/libvirt:
  - sockets=`4`
  - cores=`6`
  - threads=`1`
  - `-smp 6,maxcpus=24,sockets=4,cores=6,threads=1`
- machine type: `pc-q35-rhel9.8.0`
- cpu model: `EPYC-Genoa`
- ACPI enabled in domain features

Hotplug state:

- static non-hotpluggable baseline vCPU: ids `0..5`
- hotpluggable enabled runtime vCPU: id `6`
- disabled hotpluggable slots: ids `7..23`
- QMP hotplugged CPU exists and is realized/hotplugged:
  - path `/machine/peripheral/vcpu6`
  - `qom-get realized=true`
  - `qom-get hotplugged=true`
- QMP supports `device_del`
- QOM `qom-list /machine/peripheral` exposes child `vcpu6`

Key divergence signal:

- runtime vCPU object `vcpu6` exists and is active in QEMU thread accounting (`vcpu.6.state=1`) but guest still reports 6 logical CPUs.

Best technical interpretation:

- CPU object is present in QEMU/libvirt, but guest-side CPU online/consumption path is not completing.
- Live unplug likely blocks because guest never fully accepted the extra CPU for safe offlining/unplug workflow.
