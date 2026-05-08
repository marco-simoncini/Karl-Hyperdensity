# Test Summary

Added test suite `pkg/windowsfluidvirt/controlled_apply_plan_test.go` covering:

- default gate deny
- awaiting approval and apply-ready
- autonomous apply rejected
- dry-run required
- kill switch blocked
- missing audit/rollback/return/workload verification blocked
- master-win11 action composition
- pool-child allowed as individual target
- pool scaling blocked
- vcpu/logical/vm-spec rejected
- CLI deterministic output + fixture phase checks
