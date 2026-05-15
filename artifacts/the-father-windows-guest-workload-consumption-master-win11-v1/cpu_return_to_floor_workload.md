# CPU Return-to-Floor Workload

A first return-to-floor attempt with `-mode return-to-floor` kept `cpu.max` at baseline (`600000 100000`).
A corrective bounded rerun used actuator `-mode apply` to restore floor deterministically.

Corrective return evidence:

- down apply decision: `applied`
- observedAfterCpuMax: `300000 100000`
- post-return runs: 3
- post-return median iterationsPerSecond: **91982.86**
- floor median reference: **92456.54**
- post-return cpu.stat throttled delta (usec): 114318814
- post-return nr_throttled delta: 501

CPU_WORKLOAD_DOWN_CONFIRMED = true (throughput returned to floor-consistent range and throttling returned to floor-like behavior).
