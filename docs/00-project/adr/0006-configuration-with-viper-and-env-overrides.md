# ADR-0006 — Configuration with Viper and env overrides

> **AI-DLC phase:** `02-design`  ·  **Status:** accepted
> **Date:** 2026-11-05  ·  **Owner:** vladeParica

## Context
Services run locally, in Compose and in orchestrators. Configuration must come from the environment in
production while keeping a readable file for development.

## Decision
Load an optional YAML file (`configs/config.yaml`) with Viper, let every key be overridden by an
environment variable (`DATABASE_HOST`, `HTTP_PORT`, …), and validate the result at startup. The config is
returned as a value, never stored globally.

## Alternatives considered
- **Environment only (envconfig / caarlos0/env)** — lighter; rejected to keep a documented YAML file.
- **Global `viper` instance** — rejected: hidden coupling and hard to test.

## Consequences
- **Positive:** 12-factor friendly; misconfiguration fails fast with a clear error.
- **Negative / costs:** Viper brings several transitive dependencies.
