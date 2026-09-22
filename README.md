# glink

A URL shortener written in Go. It creates a short code for an HTTP(S) URL and redirects users to the original address.

## Run

Docker and Docker Compose are required.

```sh
cp .env.example .env
docker compose up --build
```

Docker Compose starts PostgreSQL, Redis, Kafka, and ClickHouse, applies migrations, and runs the application at `http://localhost:8080`.

Redis caches resolved short links. Configure it with `REDIS_ADDR`, optional `REDIS_PASSWORD`, and `LINK_CACHE_TTL` (defaults to `1h`); see `.env.example`.

Successful redirects publish an asynchronous `link.clicked` event to Kafka. A consumer inserts events in ClickHouse batches; the redirect never waits for Kafka. Configure Kafka with `KAFKA_BROKERS`, `KAFKA_CLICK_TOPIC`, `KAFKA_CLICK_CONSUMER_GROUP`, and `KAFKA_MAX_BUFFERED_RECORDS`, and ClickHouse with `CLICKHOUSE_ADDR`.

## API

Create a short link:

```sh
curl -X POST http://localhost:8080/links/ \
  -H 'Content-Type: application/json' \
  -d '{"original_url":"https://example.com"}'
```

`201 Created` response:

```json
{
  "id": 1,
  "original_url": "https://example.com",
  "code": "a1B2c3",
  "created_at": "2026-09-21T12:00:00Z"
}
```

Redirect a user:

```sh
curl -i http://localhost:8080/r/a1B2c3
```

Get the number of redirects for a link. `from` and `to` are optional RFC 3339 timestamps and both bounds are inclusive:

```sh
curl 'http://localhost:8080/links/a1B2c3/stats?from=2026-09-01T00:00:00Z&to=2026-09-30T23:59:59Z'
```

## Development

```sh
go test ./...
go run ./cmd/app
```

Running without Docker requires available PostgreSQL, Redis, and Kafka instances. Set `POSTGRES_DSN`, `REDIS_ADDR`, and `KAFKA_BROKERS`; `HTTP_ADDR` defaults to `:8080`.
