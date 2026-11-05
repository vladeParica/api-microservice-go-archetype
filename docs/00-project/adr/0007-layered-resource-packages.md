# ADR-0007 — Layered resource packages

> **AI-DLC phase:** `02-design`  ·  **Status:** accepted
> **Date:** 2026-11-05  ·  **Owner:** vladeParica

## Context
The original handlers opened transactions, called generated queries and wrote responses in one function,
leaking transactions on error and making testing impossible without a database.

## Decision
Package by resource (`internal/user`). Inside it: `Handler` (HTTP) → `Service` (rules, validation) →
`Repository` interface (persistence), with a PostgreSQL implementation that maps driver errors to domain
errors (`ErrNotFound`, `ErrEmailTaken`).

## Alternatives considered
- **Package by layer (`handlers/`, `services/`)** — rejected: scatters one feature across packages.
- **Full hexagonal / clean architecture** — rejected: too much ceremony for a template.

## Consequences
- **Positive:** features are self-contained; handlers and services are unit-tested with a fake repository.
- **Negative / costs:** cross-resource transactions need an explicit unit-of-work when they appear.
