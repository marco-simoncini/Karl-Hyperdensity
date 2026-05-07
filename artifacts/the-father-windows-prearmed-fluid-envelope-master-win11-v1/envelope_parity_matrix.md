# Envelope Parity Matrix

- CPU entitlement up runtime: **FAIL**
- CPU entitlement down runtime: **FAIL**
- RAM entitlement up runtime: **PASS**
- RAM entitlement down runtime: **PASS**
- same QEMU process (proof phase): **PASS**
- same Windows boot (proof phase): **PASS**
- same VM: **PASS**
- same pod/node (proof phase): **PASS**
- no VMIM: **PASS**
- no migration: **PASS**
- no rollout/recreate during proof: **PASS**
- guest ACK throughout proof: **PASS**
- rollback / return-to-floor:
  - CPU: **NOT VERIFIED**
  - RAM: **VERIFIED**

Overall envelope verdict:

- `WINDOWS_PREARMED_FLUID_ENVELOPE_BLOCKED_BY_CPU_CONTROL`
