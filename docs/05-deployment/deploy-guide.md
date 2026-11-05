# Deployment guide — API Microservice Go Archetype

> **AI-DLC phase:** `05-deployment`
> Secrets live in the environment, never in the repo.

## 1. Services to provision
| Service | For | Notes |
|---|---|---|
| Container runtime (Kubernetes, ECS, Cloud Run…) | API | Runs the image built from `Dockerfile` |
| PostgreSQL 16 | Data | Managed service with automated backups / PITR |
| Ingress / load balancer | TLS, routing | Terminates HTTPS; rate limiting |

## 2. Environment variables
| Variable | Use | Secret |
|---|---|---|
| `APP_NAME`, `APP_ENVIRONMENT` | Log metadata | no |
| `HTTP_PORT` | Listening port (default `8080`) | no |
| `HTTP_SHUTDOWN_TIMEOUT` | Grace period on `SIGTERM` (default `10s`) | no |
| `LOG_LEVEL`, `LOG_FORMAT` | `info` / `json` in production | no |
| `DATABASE_HOST`, `DATABASE_PORT`, `DATABASE_NAME` | Connection | no |
| `DATABASE_USER`, `DATABASE_PASSWORD` | Credentials | **yes** |
| `DATABASE_SSLMODE` | `require` or `verify-full` in production | no |
| `DATABASE_MAX_CONNS` | Pool size per replica | no |

**Secret generation:** `<TODO — Human-in-the-Loop: secret manager per platform>`; e.g. `openssl rand -base64 32` for the DB password.

## 3. Steps
1. Build and tag the image: `make docker-build VERSION=vX.Y.Z`.
2. Apply migrations: `make migrate-up DATABASE_URL=…` (or run the `migrate/migrate` image as a job).
3. Deploy the image with the variables above.
4. Configure probes: liveness `GET /healthz`, readiness `GET /readyz`.
5. Smoke test: `GET /readyz` returns `200`, then create and fetch a user.

## 4. Runbook / rollback
- **Rollback (app):** redeploy the previous image tag; the API is stateless.
- **Rollback (data):** `make migrate-down` for the last migration, or restore from backup/PITR.
- **Migrations:** additive and backward compatible so the previous image keeps working.
- **Shutdown:** on `SIGTERM` the server stops accepting connections and drains for `HTTP_SHUTDOWN_TIMEOUT`.

## 5. Go-live checklist
- [ ] Secrets set  · [ ] Migrations applied, no seed  · [ ] Backups + restore tested
- [ ] HTTPS at ingress  · [ ] Authentication added  · [ ] Smoke tests pass  · [ ] Monitoring on
