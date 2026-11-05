# Test strategy — API Microservice Go Archetype

> **AI-DLC phase:** `04-testing`

## Levels
| Level | Scope | Tooling | What it covers |
|---|---|---|---|
| Unit | Handlers, service rules, middleware, config, logger | `testing`, `httptest`, in-memory fake repository | Happy path and PRD edge cases |
| Integration | Repository + PostgreSQL | **Not yet** — planned with testcontainers | Constraint mapping (`23505` → `409`), `ErrNoRows` → `404` |
| End-to-end | Running stack | `make up` + [`api/requests.http`](../../api/requests.http) | Manual smoke test |
| Static | Whole codebase | `go vet`, golangci-lint, `sqlc diff` | Bugs, style, generated code drift |

All automated levels run in CI on every push and pull request; unit tests run with `-race`.

## Critical / safety paths (must be covered)
- Error mapping never leaks internal errors (`500` → generic message).
- Panics are recovered; readiness reports `503` when the database is down.
- Strict JSON decoding (unknown fields, trailing data).

## Fixtures / seed data
- Unit tests use the in-memory `fakeRepository` in `internal/user/handler_test.go`.
- No seed data is shipped; production databases only receive migrations.

## Known gaps
| Gap | Risk | Decision (fix now / accept) |
|---|---|---|
| No integration tests for `PostgresRepository` | Driver error mapping untested | Fix next (backlog) |
| 1 MiB body limit not asserted by a test | Regression unnoticed | Fix next |
| No load tests | Unknown capacity | Accept; each service sizes itself |
