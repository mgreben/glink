# glink

A URL shortener written in Go. It creates a short code for an HTTP(S) URL and redirects users to the original address.

## Run

Docker and Docker Compose are required.

```sh
cp .env.example .env
docker compose up --build
```

Docker Compose starts PostgreSQL and Redis, applies migrations, and runs the application at `http://localhost:8080`.

Redis caches resolved short links. Configure it with `REDIS_ADDR`, optional `REDIS_PASSWORD`, and `LINK_CACHE_TTL` (defaults to `1h`); see `.env.example`.

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

## Development

```sh
go test ./...
go run ./cmd/app
```

Running without Docker requires available PostgreSQL and Redis instances. Set `POSTGRES_DSN` and `REDIS_ADDR`; `HTTP_ADDR` defaults to `:8080`.
