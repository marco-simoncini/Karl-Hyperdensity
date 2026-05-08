# CPU Ceiling Workload

Actuator flow:

- dry-run decision: `accepted`
- apply decision: `applied`
- applied cpu.max: `600000 100000`

Workload outcome:

- runs: 3
- median iterationsPerSecond: **163259.57**
- improvement vs floor: **76.58%**
- cpu.stat usage delta (usec): 291768550
- cpu.stat throttled delta (usec): 722584
- cpu.stat nr_throttled delta: 260

Interpretation:

- throughput increased materially at ceiling
- throttling dropped sharply vs floor
- CPU_WORKLOAD_UP_CONFIRMED = true
