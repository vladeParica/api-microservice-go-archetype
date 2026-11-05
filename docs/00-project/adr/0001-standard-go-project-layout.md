# ADR-0001 — Standard Go project layout

> **AI-DLC phase:** `02-design`  ·  **Status:** accepted
> **Date:** 2026-11-05  ·  **Owner:** vladeParica

## Context
The original archetype placed everything under `pkg/` with Java-style folders (`controllers`, `models`,
`components`). `pkg/` signals importable public APIs, which a service does not offer.

## Decision
Follow the official [Organizing a Go module](https://go.dev/doc/modules/layout) guide: `cmd/api` for the
entry point (wiring only), `internal/` for all application code, `db/` for SQL assets.

## Alternatives considered
- **`pkg/` layout** — rejected: exposes internals to other modules and is not idiomatic for services.
- **golang-standards/project-layout** — rejected: not official and heavier than needed.

## Consequences
- **Positive:** compiler-enforced encapsulation; familiar to any Go developer.
- **Negative / costs:** code meant to be shared must be extracted into its own module.
