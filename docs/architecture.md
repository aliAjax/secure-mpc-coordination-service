# Architecture and operations

The domain layer has no transport or storage dependency. Application commands enforce transitions and idempotency before invoking repository/crypto ports. The snapshot adapter writes a temporary file then atomically renames it. A production deployment should use a transactional SQL adapter with unique constraints on computation and participant identifiers.

## State machine

`draft -> committed -> running -> waiting_shares -> reconstructing -> succeeded`; `aborted` and `expired` are terminal. Round state is `open -> collecting -> complete`, with timeout/abort terminal branches.

## SLO and recovery

Target p95 API latency is 100ms for metadata operations. State snapshots are copied to durable storage; restore by setting `MPC_STATE_FILE`. Evidence digest uses SHA-256 over computation/round/result and is independently verifiable.
