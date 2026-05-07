# Previous Product Path Summary

Source proof:

- commit: `ce05534`
- artifact: `artifacts/the-father-windows-product-path-definitive-master-win11-v1/`
- verdict: `WINDOWS_FLUIDVIRT_PRODUCT_PATH_CONFIRMED`

Extracted baseline:

- CPU floor/ceiling: `300000 100000` / `600000 100000`
- RAM floor/ceiling: `12884901888` / `13958643712`
- target cgroup path (proof run): `/host-sys-fs-cgroup/.../cri-containerd-...scope`
- same QEMU PID/start and same Windows boot were preserved
- guest ACK true, rollback true, return-to-floor true
