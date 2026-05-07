# GateSet Evaluator

`EvaluateWindowsFluidUnlockGateSet` evaluates Gate 0/1/2 and returns aggregate status:

- `GATE_SET_PASSED`
- `GATE_SET_BLOCKED`
- `GATE_SET_QUARANTINED`
- `GATE_SET_FAILED`

Aggregate pass does not enable executor and does not imply unlock.
