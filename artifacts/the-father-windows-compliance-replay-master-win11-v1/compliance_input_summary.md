# Compliance Input Summary

Replay fixture: `master-win11-real-evidence.ready.json`

- VM: `master-win11`
- namespace: `karl`
- cluster proof lineage: `karl-metal-01@ovh`
- node: `karl-lab-metal-01`
- pod UID: `7b6a904a-1c9a-4a44-9b37-1dc737304773`
- QEMU PID/start: `96` / `Thu May 7 18:58:03 2026`
- Windows boot evidence: preserved (guest same-boot=true)
- machine identity evidence: preserved (hash-only, no raw GUID)
- QMP evidence present: true
- fluidShell evidence present: true
- RAM balloon evidence present: true
- CPU actuator capability present: true
- policy ready annotation: `hyperdensity.karl.io/windows-ready=true`
- pool context: standalone (`isPoolChild=false`, provisioning-only=true)
- migration/recreate/rollout constraints: preserved by source evidence
- rollback/return-to-floor: verified by source evidence
