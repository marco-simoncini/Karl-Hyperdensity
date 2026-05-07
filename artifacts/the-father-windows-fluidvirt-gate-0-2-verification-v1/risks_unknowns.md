# Risks and Unknowns

- Gate 2 replay detection currently relies on deterministic replay flags and freshness windows; full signed replay prevention remains future milestone work.
- Gate 1 still depends on quality of collected evidence payloads; low-quality data ingestion can produce conservative blocking.
- GateSet pass semantics remain intentionally non-executable and require explicit communication to avoid misinterpretation.
