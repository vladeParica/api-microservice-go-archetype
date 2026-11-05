# Backlog & implementation notes — API Microservice Go Archetype

> **AI-DLC phase:** `03-implementation`

## Status
Core archetype implemented and green in CI; hardening items below remain.

## Done
- [x] Standard layout, module path matching the repository — ADR-0001
- [x] Config loader with env overrides and validation — ADR-0006
- [x] slog logger with request-scoped context — ADR-0003
- [x] HTTP server, middleware, graceful shutdown, health probes — ADR-0002
- [x] Reference `user` resource with sqlc repository — ADR-0004, ADR-0007
- [x] Migrations, Compose stack, distroless image — ADR-0005, ADR-0008
- [x] CI (lint, race tests, sqlc diff, docker build), Dependabot

## In progress / remaining
- [ ] Integration tests against PostgreSQL (testcontainers)
- [ ] Optional authentication middleware — blocked by PRD open decision
- [ ] Metrics endpoint (Prometheus / OpenTelemetry)
- [ ] OpenAPI specification
- [ ] Release workflow publishing the image on tags

## Implementation notes
- Add a resource: `make migrate-new name=…` → write queries → `make sqlc` → copy `internal/user` → register routes in `cmd/api/main.go`.
- `internal/database` is generated; edit `db/queries` and regenerate instead.
- The version is injected at build time: `-ldflags "-X main.version=…"` (Makefile and Dockerfile do it).
- `.gitattributes` forces LF so gofmt and sqlc diff behave the same on Windows.

## Deviations from design
None.
