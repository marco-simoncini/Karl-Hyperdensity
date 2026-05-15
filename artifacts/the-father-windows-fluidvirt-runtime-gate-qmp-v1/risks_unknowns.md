# Risks and Unknowns

- Sidecar socket transport is a skeleton and depends on runtime QMP framing behavior in target pod.
- Cluster probe is read-only and does not assert annotation gates on the candidate VM yet.
- Runtime gate uses canonical blockers without introducing new annotation-specific blocker IDs.
- Full end-to-end proof still requires future runtime controller and sidecar deployment in lab.
