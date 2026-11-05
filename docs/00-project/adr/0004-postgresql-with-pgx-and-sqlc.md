# ADR-0004 — PostgreSQL with pgx and sqlc

> **AI-DLC phase:** `02-design`  ·  **Status:** accepted
> **Date:** 2026-11-05  ·  **Owner:** vladeParica

## Context
Services need a relational store with type-safe access and no runtime reflection surprises. The original
module also listed gorm and lib/pq without using them.

## Decision
Use PostgreSQL through `pgx/v5` (`pgxpool`) and generate the data-access code with sqlc from SQL files.
Generated code lives in `internal/database` and CI fails if it is out of date (`sqlc diff`).

## Alternatives considered
- **gorm** — rejected: hides SQL, reflection-based, harder to reason about performance.
- **database/sql + lib/pq** — rejected: lib/pq is in maintenance mode; pgx is faster and richer.

## Consequences
- **Positive:** compile-time checked queries; plain SQL is reviewable.
- **Negative / costs:** developers must run `make sqlc` after editing queries; dynamic queries need hand-written code.
