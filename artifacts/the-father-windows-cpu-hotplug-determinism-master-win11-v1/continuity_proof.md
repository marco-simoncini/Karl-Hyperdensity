# Continuity Proof

Checked before and after recovery attempt:

- VM name/namespace unchanged: `master-win11` / `karl`
- VMI UID unchanged: `a565fc95-d79f-4cbd-a852-e0467e5f2110`
- virt-launcher pod unchanged: `virt-launcher-master-win11-w7mxv`
- pod UID unchanged: `6e0a597f-b508-450d-a4ac-d5e4262c8615`
- node unchanged: `karl-lab-metal-01`
- pod restart count unchanged: `0`
- QEMU PID unchanged: `92`
- QEMU start time unchanged: `Thu May 7 18:16:51 2026`
- Windows last boot unchanged: `/Date(1778177812500)/`
- machineGuid hash unchanged: `d1fdf2ce69932d7ac2f9e3497be7f845d4a3ef9140c45bd0d7b0ea47c195508e`
- pending reboot: `false`
- VMIM objects in namespace: none
- migration/recreate evidence during test window: none observed
