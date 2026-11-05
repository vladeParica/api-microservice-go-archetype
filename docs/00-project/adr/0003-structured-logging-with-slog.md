# ADR-0003 — Structured logging with slog

> **AI-DLC phase:** `02-design`  ·  **Status:** accepted
> **Date:** 2026-11-05  ·  **Owner:** vladeParica

## Context
The original logger used zap plus lumberjack file rotation and a package-level logger mutated on every
request, which caused a data race between concurrent requests.

## Decision
Use `log/slog` writing to stdout (JSON in deployed environments, text locally). A request-scoped logger
carrying `request_id` travels in the `context.Context`.

## Alternatives considered
- **zap / zerolog** — rejected: faster in benchmarks, but the gain is irrelevant for typical services.
- **File rotation in-process** — rejected: log collection belongs to the platform (12-factor).

## Consequences
- **Positive:** stdlib, no global state, logs ready for any collector.
- **Negative / costs:** services that need sampling or very high throughput must add a custom handler.
