# Project charter — API Microservice Go Archetype

> **AI-DLC phase:** `00-project`  ·  **Status:** draft
> **Owner:** `<TODO — Human-in-the-Loop>`  ·  **Date:** 2026-11-05

## Vision
Start a new Go HTTP microservice in minutes with production defaults already in place, so teams spend
their time on business logic instead of boilerplate.

## Context
Every new service re-implemented configuration, logging, database access, health checks, shutdown and
CI, each slightly differently. The original version of this archetype was built for that purpose inside
a company and has been refactored to follow idiomatic Go conventions.

## Problem
Developers creating Go services need a consistent, reviewable starting point that is small enough to
understand in an afternoon and complete enough to deploy.

## Scope
**In scope**
- HTTP/JSON API skeleton with one example resource (`user`).
- PostgreSQL access (pgx + sqlc) and versioned migrations.
- Configuration, structured logging, health probes, graceful shutdown.
- Tests, linting, CI, container image and local Compose stack.

**Out of scope (explicit, not omissions)**
- Authentication / authorization — depends on each organization's identity provider.
- Metrics and tracing exporters — vendor specific (see `06-monitoring`).
- gRPC, messaging, caching — added per service when needed.

## Success metrics
- A new service is created with `gonew` and passes CI without code changes.
- `make up` yields a running API with a ready database in under 2 minutes.
- Zero lint findings and green tests on `develop`.

## Principles / non-negotiables
- Standard library first; a dependency must justify itself in an ADR.
- No global mutable state; dependencies are wired explicitly in `cmd/api`.
- Configuration via environment; no secrets committed to the repository.
- Every generated artifact (sqlc code) is reproducible and checked in CI.

## Key risks
| Risk | Impact | Mitigation / owner |
|---|---|---|
| Template drifts from Go ecosystem practices | Services start outdated | Dependabot + periodic ADR review |
| Example code copied without understanding | Insecure services (no auth) | Charter and README state auth is out of scope |
| Toolchain versions diverge (Go, sqlc, migrate) | Inconsistent builds | Versions pinned in CI, Dockerfile and Makefile |

## Open decisions (Human-in-the-Loop)
- `<TODO — Human-in-the-Loop: confirm owner/maintainers of the archetype>`
- `<TODO — Human-in-the-Loop: minimum supported Go version policy>`
