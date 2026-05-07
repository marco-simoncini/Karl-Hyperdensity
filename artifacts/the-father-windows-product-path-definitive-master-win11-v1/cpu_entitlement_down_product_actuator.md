# CPU Entitlement Down (Return-to-Floor)

Mechanism:

- product actuator `return-to-floor` on `cpu.max`.

Evidence:

- `raw_logs_sanitized/actuator_cpu_down_output.json`

Resulting state:

- `cpu.max` changed from `600000 100000` to `300000 100000`.

Workload metrics after down:

- floor run2:
  - elapsed: `28057 ms`
  - usage delta: `81354187 usec`
  - throttled delta: `59507913 usec`
  - nr_throttled delta: `215`

Interpretation:

- floor behavior is restored and coherent with constrained CPU entitlement.

Result:

- `CPU_ENTITLEMENT_DOWN_CONFIRMED=true`
