# Design Implications

CPU path is currently non-deterministic for this Windows profile.

Actionable directions:

- Prearm topology for Windows CPU lifecycle:
  - keep target sockets/cores pre-defined;
  - validate whether guest can see/offline planned CPU slots without runtime divergence.
- Evaluate machine/type + CPU model pairing:
  - confirm ACPI CPU hotplug signaling behavior with `pc-q35-rhel9.8.0` + `EPYC-Genoa`.
- Add explicit Windows-side CPU online/offline evidence adapter:
  - correlate OS processor groups and per-processor state after hotplug/unplug requests.
- Consider QMP device-id managed flow (`device_add`/`device_del` against known vCPU IDs) only with strict guardrails and before/after proof blocks.
- If unplug remains unsupported in this guest/profile, classify Windows CPU parity as blocked for dynamic downscale until image/platform changes are validated.
- Cold reset may be required to restore coherence from this exact mismatch state, but it is outside Hyperdensity in-place success criteria.

Explicit non-directions for this phase:

- no RAM testing
- no replica/pool scaling
- no migration/recreate workaround
