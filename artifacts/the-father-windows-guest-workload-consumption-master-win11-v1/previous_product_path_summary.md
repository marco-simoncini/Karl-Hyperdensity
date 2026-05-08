# Previous Product Path Summary

Read and reused productized patterns from:

- `artifacts/the-father-windows-product-path-definitive-master-win11-v1/`
- `artifacts/the-father-windows-node-fluid-actuator-mvp-v1/`
- `artifacts/the-father-windows-fluidvirt-product-model-v1/`

Extracted controls reused:

- CPU floor/ceiling: `300000 100000` / `600000 100000`
- RAM floor/ceiling request: `12884901888` / `13958643712`
- actuator request+allowlist identity pinning
- QMP balloon pattern (`query-balloon`, `balloon`)
- continuity proofs (same QEMU PID/start, same Windows boot/machine hash)
- compliance replay gate and deterministic evidence chain approach
