# Test Summary

Actuator test coverage includes:

1. dry-run accepted
2. apply accepted
3. rollback accepted
4. return-to-floor accepted
5. stale request rejected
6. kill-switch blocked
7. pod UID mismatch rejected
8. qemu PID mismatch rejected
9. qemu start mismatch rejected
10. cgroup path mismatch rejected
11. symlink traversal rejected
12. parent cgroup write rejected
13. non-`cpu.max` controller rejected
14. out-of-bounds requested cpu.max rejected
15. missing rollback target rejected

Bundle append test coverage includes:

16. append ready run
17. append blocked run
18. broken chain rejected
19. duplicate run rejected
20. deterministic append with fixed time
