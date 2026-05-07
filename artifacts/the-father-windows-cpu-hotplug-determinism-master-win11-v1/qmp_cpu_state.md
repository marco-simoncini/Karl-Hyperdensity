# QMP CPU State

Read-only QMP evidence:

- `query-status`: running
- `query-cpus-fast`: 7 vCPUs (cpu-index 0..6)
- `query-hotpluggable-cpus`: baseline topology sockets=4, cores=6, threads=1; hotplugged slot exposed at `/machine/peripheral/vcpu6`
- `query-memory-size-summary`: base memory only, plugged memory 0 (read-only check only)

Libvirt runtime:

- `dominfo`: `CPU(s): 7`
- `vcpucount --live`: `7`
- `domstats --vcpu`: `vcpu.current=7`, `vcpu.maximum=24`

Interpretation:

- QEMU/libvirt sees the extra hotplugged vCPU active.
- Guest does not consume it as an online logical processor.
