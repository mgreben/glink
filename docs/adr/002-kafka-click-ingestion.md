# ADR-002: Ingest redirect clicks through Kafka

## Status

Accepted

## Context

Writing a PostgreSQL row per redirect makes click tracking the limiting path under high read traffic. Small analytics loss during failures is acceptable.

## Decision

Publish `link.clicked` events asynchronously to Kafka using franz-go. The producer has a bounded buffer and never delays the redirect. A consumer group batches events into daily PostgreSQL aggregates. Statistics read those aggregates rather than raw click events.

## Consequences

- Redirect availability and latency do not depend on Kafka acknowledgements.
- Kafka producer overflow or broker failures can lose events.
- Consumer retries may duplicate a batch after a crash, so analytics are approximate.
- Kafka needs production-grade replication, monitoring, retention, and authentication outside local Compose.
