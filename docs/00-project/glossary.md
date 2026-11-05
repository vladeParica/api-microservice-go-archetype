# Glossary — API Microservice Go Archetype

> **AI-DLC phase:** `00-project`
> The ubiquitous language of the project. Every other document must use these terms consistently.

| Term | Definition |
|---|---|
| Archetype | This repository: the template new services are generated from with `gonew`. |
| Service | A microservice generated from the archetype. |
| Resource | A domain entity exposed over HTTP, implemented as one package under `internal/` (e.g. `user`). |
| Handler | HTTP transport layer of a resource: decodes requests, maps errors to status codes. |
| Service layer | Business rules and validation of a resource (`Service` type). Avoid "controller". |
| Repository | Persistence port of a resource; the PostgreSQL implementation wraps sqlc queries. |
| Query | SQL statement in `db/queries/`, compiled to Go by sqlc. |
| Migration | Versioned, reversible schema change in `db/migrations/` (`NNNNNN_name.up/down.sql`). |
| Liveness | `GET /healthz` — the process is running. Never checks dependencies. |
| Readiness | `GET /readyz` — the service can serve traffic (database reachable). |
| Request ID | Correlation identifier per request, read from or written to `X-Request-ID` and logged. |
| Developer | Actor who generates and extends a service from the archetype. |
| Operator | Actor who deploys and runs a service. |
| API client | Any system consuming the HTTP API. |
