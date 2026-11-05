# ADR-0002 — Standard library HTTP routing

> **AI-DLC phase:** `02-design`  ·  **Status:** accepted
> **Date:** 2026-11-05  ·  **Owner:** vladeParica

## Context
Go 1.22 added method and path-parameter patterns to `net/http.ServeMux` (`GET /users/{id}`), the main
reason to depend on gorilla/mux.

## Decision
Use `http.ServeMux` with Go 1.22 patterns and small `func(http.Handler) http.Handler` middleware.

## Alternatives considered
- **gorilla/mux** — rejected: no longer needed; extra dependency.
- **chi / gin / echo** — rejected: useful, but each service can adopt one if it truly needs it.

## Consequences
- **Positive:** zero routing dependencies; standard `http.Handler` everywhere.
- **Negative / costs:** no built-in route groups or regex constraints; middleware chaining is manual.
