# Clean Baseline After Recovery

Recovery method used:

- controlled **lab reset** (stop/start), outside Hyperdensity runtime criteria.

Before lab reset (incoherent):

- pod: `virt-launcher-master-win11-w7mxv`
- VMI UID: `a565fc95-d79f-4cbd-a852-e0467e5f2110`
- QEMU PID/start: `92` / `Thu May 7 18:16:51 2026`
- QMP/libvirt CPU: `7`
- guest CPU: `6`
- Windows last boot: `/Date(1778177812500)/`

After lab reset (coherent):

- pod: `virt-launcher-master-win11-kmwgg`
- pod UID: `7b6a904a-1c9a-4a44-9b37-1dc737304773`
- VMI UID: `b31dcab6-5a99-432e-a6fe-21607a9e3403`
- node: `karl-lab-metal-01` (same node)
- QEMU PID/start: `96` / `Thu May 7 18:58:03 2026`
- QMP/libvirt CPU: `6`
- guest CPU: `6`
- Windows last boot: `/Date(1778180311500)/` (changed)
- machineGuidHash: `d1fdf2ce69932d7ac2f9e3497be7f845d4a3ef9140c45bd0d7b0ea47c195508e` (unchanged)
- pending reboot: false
- VMIM objects: none
- migration evidence: none

Classification:

- `baseline_lab_reset_recovered`
