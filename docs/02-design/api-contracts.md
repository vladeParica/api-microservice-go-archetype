# API contracts — API Microservice Go Archetype

> **AI-DLC phase:** `02-design`
> Source of truth: the route registrations in `internal/*/handler.go`. Samples in [`api/requests.http`](../../api/requests.http).

## Conventions
- Base path: `/api/v1` (probes at the root).
- Format: JSON (`Content-Type: application/json`). Auth: none (see threat model).
- Error envelope: `{"error": "<message>"}`. Every response carries `X-Request-ID`.

## Endpoints
| Method + path | Auth | Purpose | Request | Response |
|---|---|---|---|---|
| `GET /healthz` | none | Liveness | — | `200 {"status":"ok","version":"…"}` |
| `GET /readyz` | none | Readiness | — | `200` or `503` with per-dependency `checks` |
| `POST /api/v1/users` | none | Create user | `{"name","email"}` | `201` user + `Location`; `400`; `409` |
| `GET /api/v1/users/{id}` | none | Get user | — | `200` user; `400`; `404` |

User representation: `{"id": 1, "name": "John", "email": "john@example.com", "createdAt": "2026-11-05T10:00:00Z"}`.

## Business rules reflected in the API
- Email is normalized to lower case and unique → `409` on duplicates ([PRD](../01-requirements/prd-user-resource.md)).
- Unknown fields are rejected → `400`.
