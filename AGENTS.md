# AGENTS.md

## Project

`glink` is a URL-shortening HTTP service built with Go 1.25, PostgreSQL, chi, and fx.

## Commands

```sh
go test ./...
gofmt -w <changed-go-files>
docker compose up --build
```

## Structure

- `cmd/app` — application entry point.
- `internal/app` — dependency injection and application lifecycle.
- `internal/handler` — HTTP routes, DTOs, and response serialization.
- `internal/modules/links` — link domain model and service.
- `internal/modules/links/postgres` — PostgreSQL repository.
- `pkg` — independent reusable packages.
- `db/postgres` — SQL migrations.

## Conventions

- Represent expected domain states with sentinel errors and check them using `errors.Is`.
- Translate PostgreSQL driver errors to domain errors in the repository; wrap infrastructure errors with `%w`.
- Define HTTP status codes and public messages in handlers. Use `writeJSON` and `writeError` for JSON responses.
- Validate input at the HTTP boundary. Do not add application-specific code to `pkg` unless it is independent and reusable.
- New logic needs tests. Run `gofmt` and `go test ./...` before completing changes.
