# policy_pack

Default policy pack: `windows-fluid-admission-policy-v1`.

Conservative defaults:

- allow only `in-place-qmp` runtime mode;
- require certified fluid shell conditions;
- require no migration, no reboot, no recreate;
- require same node/pod/qemu/boot/machine continuity;
- require QMP ACK and guest ACK;
- require rollback and return-to-floor readiness;
- deny pool replica model and generic Windows VM;
- disable mutation/apply in this phase;
- enforce evidence freshness and score threshold for future-apply admission.
