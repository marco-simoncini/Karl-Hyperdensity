# Topology Analysis

Libvirt domain (`dumpxml`) highlights:

- `<vcpu placement='static' current='7'>24</vcpu>`
- vcpu 0..5 non-hotpluggable baseline
- vcpu 6 hotpluggable and enabled
- vcpu 7..23 hotpluggable and disabled
- CPU topology:
  - sockets=4
  - cores=6
  - threads=1
  - max vCPU = 24
- Machine type: `pc-q35-rhel9.8.0`

QMP mapping:

- Active hotplugged vCPU path: `/machine/peripheral/vcpu6`
- Baseline CPUs on `/machine/unattached/device[...]`

Guest capability indicators:

- QGA supports `guest-get-vcpus`
- QGA reports `guest-set-vcpus` disabled
- Guest reports all 6 visible CPUs as `can-offline=false`

Hotplug model implications:

- CPU add is visible to QMP/libvirt.
- Windows guest does not auto-online the new hotplugged CPU in this configuration.
- Live unplug path is likely blocked by guest-side offlining limitations.
