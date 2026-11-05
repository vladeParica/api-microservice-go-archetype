# ADR-0008 — Distroless container image

> **AI-DLC phase:** `05-deployment`  ·  **Status:** accepted
> **Date:** 2026-11-05  ·  **Owner:** vladeParica

## Context
The service is a static Go binary. The original image used `alpine:latest` (unpinned, with a shell and
package manager) and ran as root.

## Decision
Multi-stage build: compile with `CGO_ENABLED=0` in `golang:<version>-alpine`, run on
`gcr.io/distroless/static-debian12:nonroot` as a non-root user. Version is injected with `-ldflags`.

## Alternatives considered
- **scratch** — rejected: no CA certificates, timezone data or non-root user out of the box.
- **alpine** — rejected: larger attack surface not needed at runtime.

## Consequences
- **Positive:** small image, minimal CVE surface, non-root by default.
- **Negative / costs:** no shell for debugging; use ephemeral debug containers instead. Container
  healthchecks must be done by the orchestrator over HTTP.
