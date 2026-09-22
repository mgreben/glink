# ADR-001: Use Redis as the shared link-resolution cache

## Status

Accepted

## Context

Resolving `GET /r/{code}` currently reads PostgreSQL for every redirect. This endpoint is expected to be read-heavy and the application may run with multiple replicas.

## Decision

Use Redis as a cache-aside cache keyed by `link:{code}` with a configurable positive TTL (default: one hour). PostgreSQL remains the source of truth. Cache read and write failures are ignored by the link service so that Redis outages do not prevent redirects. Successfully created links are written through to the cache.

## Alternatives considered

- Process-local cache: simpler, but cache entries are not shared between replicas.
- Redis as source of truth: lower database reads, but weakens persistence and recovery guarantees.

## Consequences

- Positive: shared low-latency cache and a clear path to horizontal scaling.
- Negative: an additional service to operate; a cache miss during a Redis outage increases PostgreSQL load.

## Failure modes

Redis failures fall back to PostgreSQL. Cache corruption is treated as a miss and refreshed on the next successful database read.
