# Blocker Report

Primary blockers:

1. `cpu_scale_up_guest_not_confirmed`
   - QMP observed 7 vCPUs, guest remained at 6 logical processors.

2. `cpu_scale_down_live_unplug_timeout`
   - live return to floor (`7 -> 6`) timed out.

3. `WINDOWS_HYPERDENSITY_PARITY_BLOCKED_BY_MEMORY_RETURN`
   - RAM branch not safely executable after CPU branch failure.

4. `hyperdensity_parity_partial_success_not_total_feasibility`
   - Four-proof matrix not completed.
