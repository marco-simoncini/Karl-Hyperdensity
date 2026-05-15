# CPU Scale-Down Retest (-1)

Execution status: **NOT EXECUTED**

Reason:

- Deterministic `-CPU` retest is valid only after a successful `+CPU` guest-confirmed retest.
- `+CPU` retest was intentionally skipped because the environment stayed in mismatch quarantine (`QMP=7`, `guest=6`) after recovery timeout.

Guardrail applied:

- No chained mutations were executed once return-to-floor failed.
