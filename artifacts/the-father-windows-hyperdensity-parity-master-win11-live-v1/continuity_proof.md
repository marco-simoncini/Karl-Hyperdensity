# Continuity Proof

- same VM UID: true
- same namespace: true
- same VMI UID: true
- same node: true
- same virt-launcher pod UID: true
- same QEMU PID: true (`92`)
- same QEMU start time: true (`Thu May 7 18:16:51 2026`)
- same Windows boot: true (`/Date(1778177812500)/`)
- same machine identity hash: true (`d1fdf2ce69932d7ac2f9e3497be7f845d4a3ef9140c45bd0d7b0ea47c195508e`)
- no VMIM: true
- no LiveMigration objects observed: true

Continuity holds, but parity still fails because mutation proofs are incomplete.
