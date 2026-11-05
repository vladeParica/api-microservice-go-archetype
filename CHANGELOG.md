# Changelog

All notable changes to this project are documented here. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/) and the project uses
[Semantic Versioning](https://semver.org/).

## [Unreleased]

## [1.0.0] - 2026-11-05

### Added
- AI-DLC lifecycle documentation, ADR-0001…0008, C4 diagrams and threat model.
- Health probes (`/healthz`, `/readyz`), request id and panic recovery middleware.
- CI (lint, race tests, sqlc diff, docker build), release workflow, Dependabot.
- Contributing guide, security policy, code of conduct and issue/PR templates.

### Changed
- Standard Go layout (`cmd/`, `internal/`), stdlib routing, `log/slog`, validated config with env overrides.
- REST routes: `POST /api/v1/users`, `GET /api/v1/users/{id}`; default port `8080`.
- golang-migrate migrations with `BIGINT` ids, unique email and `TIMESTAMPTZ`.
- Distroless non-root container image.

### Removed
- zap, lumberjack, gorilla/mux, gorm and lib/pq dependencies.

### Fixed
- Data race in the global request logger.
- Transactions left open on errors; panics in handlers.
- Broken Dockerfile and untrappable `os.Kill` signal handling.

[Unreleased]: https://github.com/vladeParica/api-microservice-go-archetype/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/vladeParica/api-microservice-go-archetype/releases/tag/v1.0.0
