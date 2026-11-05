# PRD — User resource (reference flow)

> **AI-DLC phase:** `01-requirements`  ·  **Status:** in review
> **Author:** vladeParica  ·  **Date:** 2026-11-05

## 1. Problem
The archetype needs one complete, minimal resource that shows every layer working end to end, so a
Developer can copy its shape when adding new resources.

## 2. Users and actors
- **API client** — creates and reads users.
- **Developer** — uses the flow as the reference implementation.

## 3. Main flow (happy path)
1. API client sends `POST /api/v1/users` with `name` and `email`.
2. Service normalizes input (trim, lower-case email) and validates it.
3. Repository inserts the row; the API returns `201`, the user and a `Location` header.
4. API client sends `GET /api/v1/users/{id}` and receives `200` with the user.

## 4. Business rules
- `name` is required after trimming whitespace.
- `email` must be a valid address and is unique (case-insensitive by normalization).
- `id` is a positive 64-bit integer assigned by the database.
- Unknown JSON fields are rejected; bodies are limited to 1 MiB.

## 5. Edge, abuse and risk scenarios
| Scenario | Expected system behavior |
|---|---|
| Duplicate email (any case) | `409 Conflict`, no row inserted |
| Malformed JSON, trailing data or unknown field (e.g. `password`) | `400 Bad Request` |
| Body larger than 1 MiB | `400 Bad Request`, body not fully read |
| Non-numeric, zero or negative `id` | `400 Bad Request` |
| Non-existent `id` | `404 Not Found` |
| Database unavailable | `500` with generic message; `/readyz` returns `503` |
| Handler panic | Recovered, `500`, stack trace only in logs |
| Enumeration of ids by an anonymous client | Possible — no auth in the archetype (see 8) |

## 6. Acceptance criteria
- [x] Each scenario in section 5 is covered by a unit test in `internal/user`, `internal/httpx` or `internal/server`.
- [x] Error responses use `{"error": "<message>"}` and never leak internal errors.
- [ ] Flow verified against a real PostgreSQL instance (integration test).

## 7. Out of scope
- Update, delete and list operations; pagination.
- Passwords and authentication.

## 8. Open decisions (Human-in-the-Loop)
- `<TODO — Human-in-the-Loop: default authentication mechanism to ship (none, JWT middleware, API key)>`
