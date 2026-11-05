# Architecture — API Microservice Go Archetype

> **AI-DLC phase:** `02-design`  ·  **Status:** accepted
> **Related ADRs:** ADR-0001 … ADR-0008

## 1. Architectural style
A single stateless HTTP service backed by PostgreSQL, packaged by resource with handler → service →
repository layers (ADR-0007). Dependencies are wired explicitly in `cmd/api/main.go`; there is no global state.

| Layer | Technology | ADR |
|---|---|---|
| Layout | `cmd/` + `internal/` | ADR-0001 |
| HTTP | `net/http` ServeMux (Go 1.22 patterns) | ADR-0002 |
| Logging | `log/slog` to stdout | ADR-0003 |
| Data access | PostgreSQL, pgx v5, sqlc | ADR-0004 |
| Migrations | golang-migrate | ADR-0005 |
| Configuration | Viper + environment | ADR-0006 |
| Runtime image | distroless, non-root | ADR-0008 |

## 2. Main components
- **`cmd/api`** — composition root: config → logger → pool → resources → server.
- **`internal/server`** — `http.Server` with timeouts and graceful shutdown; `RequestLogger` and `Recoverer` middleware.
- **`internal/user`** — reference resource (handler, service, repository).
- **`internal/health`** — liveness and readiness probes.
- **`internal/httpx`** — JSON decode/encode and error envelope.
- **`internal/config`, `internal/logger`, `internal/postgres`** — infrastructure.
- **`internal/database`** — sqlc generated code (do not edit).

All state lives in PostgreSQL; any number of replicas can run behind a load balancer.

## 3. Data model (high level)
- **users** — `id BIGINT IDENTITY PK`, `name TEXT NOT NULL`, `email TEXT NOT NULL UNIQUE`, `created_at TIMESTAMPTZ NOT NULL DEFAULT now()`.

## 4. Access control
None in the archetype (see PRD open decision). Services must add authentication middleware before
exposing endpoints publicly.

## 5. Quality attributes
- **Security:** input limits and strict decoding, generic error messages, non-root image, no secrets in repo. See [threat-model.md](threat-model.md).
- **Resilience:** read/write/idle timeouts, panic recovery, readiness probe, graceful shutdown (`HTTP_SHUTDOWN_TIMEOUT`).
- **Cost / operations:** one static binary on a distroless base image, logs to stdout.
- **Availability:** horizontally scalable; database pool size via `DATABASE_MAX_CONNS`.

See the C4 diagrams in [`docs/architecture/`](../architecture/).
