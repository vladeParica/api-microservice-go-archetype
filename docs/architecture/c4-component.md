# C4 — Component diagram: API

> **AI-DLC phase:** `02-design`  ·  Level 3 (components inside a container).

```mermaid
C4Component
    title Components — API

    Container_Boundary(api, "API") {
        Component(mw, "Middleware", "internal/server", "Request ID, access log, panic recovery")
        Component(health, "Health handler", "internal/health", "Liveness and readiness")
        Component(handler, "User handler", "internal/user", "Decode, map errors to status codes")
        Component(svc, "User service", "internal/user", "Normalize and validate input")
        Component(repo, "User repository", "internal/user + internal/database", "sqlc queries, driver error mapping")
    }

    ContainerDb(db, "Database", "PostgreSQL", "users")

    Rel(mw, handler, "Routes to")
    Rel(mw, health, "Routes to")
    Rel(handler, svc, "Calls")
    Rel(svc, repo, "Uses")
    Rel(repo, db, "Reads/writes")
    Rel(health, db, "Pings")
```

## Rules / behavior
1. Only the handler knows HTTP; only the repository knows PostgreSQL.
2. Domain errors (`ErrNotFound`, `ErrEmailTaken`, `ValidationError`) are the contract between layers.
