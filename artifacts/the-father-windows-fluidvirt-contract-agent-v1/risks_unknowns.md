# Risks and Unknowns

## Risks

- Inventory test execution is blocked by missing local .NET toolchain.
- Guest critical event scan is basic and may need tighter filtering per certified signal set.
- Driver truth is conservative placeholder; `memoryAdapterVerified` remains false by default.
- Contract parity between guest evidence hash semantics and runtime continuity consumers must be validated in integration lab.

## Unknowns

- exact backend endpoint contract for `modules.fluidShell` evidence ingest URL in production topology
- final policy for machine GUID hashing/redaction compatibility across controllers
- thresholds for critical event classification (event IDs and windows)

## Mitigation

- keep verdict validation-blocked until inventory tests run in environment with `dotnet`
- execute integration/lab prompt with controlled Windows shell and evidence replay
