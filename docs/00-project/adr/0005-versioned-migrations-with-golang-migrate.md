# ADR-0005 — Versioned migrations with golang-migrate

> **AI-DLC phase:** `02-design`  ·  **Status:** accepted
> **Date:** 2026-11-05  ·  **Owner:** vladeParica

## Context
Schema changes must be repeatable across environments and reversible. sqlc reads the same files to
infer the schema.

## Decision
Store migrations as `db/migrations/NNNNNN_name.{up,down}.sql` and apply them with golang-migrate, run as a
separate step (Compose `migrate` service, `make migrate-up`), never on application startup.

## Alternatives considered
- **goose** — viable; rejected only to keep plain SQL files without annotations.
- **Migrate on startup** — rejected: races between replicas and couples deploys to schema changes.

## Consequences
- **Positive:** explicit, auditable schema history; sqlc and migrate share one source.
- **Negative / costs:** deploy pipelines need an extra migration step.
