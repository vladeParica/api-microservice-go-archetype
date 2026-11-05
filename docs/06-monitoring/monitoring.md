# Monitoring & observability — API Microservice Go Archetype

> **AI-DLC phase:** `06-monitoring`

## Health
- `GET /healthz` — process alive; never checks dependencies.
- `GET /readyz` — `200` when PostgreSQL answers a ping within 2s, otherwise `503`.

## Key metrics
Derived today from the access log (`msg="http request"`); a metrics endpoint is in the backlog.

| Metric | Why it matters | Source |
|---|---|---|
| Request latency p95 | User experience | `duration` field |
| 5xx rate | Failures | `status` field |
| Readiness failures | Database availability | `readiness check failed` log |
| Recovered panics | Bugs in handlers | `panic recovered` log |

## Alerts
| Condition | Severity | Who / action |
|---|---|---|
| 5xx rate > 1% for 5 min | high | On-call: check logs by `request_id` |
| `/readyz` failing for 2 min | high | On-call: check database |
| Any `panic recovered` | medium | Owning team: fix bug |

## SLOs
- `<TODO — Human-in-the-Loop: availability and latency targets per service>`

## Dashboards / logs
- JSON logs on stdout with `service`, `version`, `env` and `request_id`; correlate with the `X-Request-ID` response header.
