# Do Not Merge List

Items explicitly excluded from direct merge/porting in this milestone:

1. `Karl-Dashboard` `The-Father-Windows` frontend (`React/TSX/CSS/UI`) and any stale UI lineage.
2. Historical artifact bulk trees (`artifacts/the-father-windows-*`) as raw import payload.
3. Superseded or duplicate contract variants not aligned to Parent Fabric naming and Linux baseline semantics.
4. Heavy raw runtime proof dumps/log archives that inflate repository size without contract value.
5. `Karl-OS-ISO` changes (including packaging of Node Fluid Actuator) for this phase.
6. Local runtime files, secrets, tokens, cookies, kubeconfig, and environment-specific credentials.
7. Build outputs and generated binaries/obj artifacts.
8. Any change introducing raw QMP/libvirt/QGA/QOM/K8s patch controls as product surface.
9. Any branch-level direct merge from `The-Father-Windows` into `main`/`The-Father`.
