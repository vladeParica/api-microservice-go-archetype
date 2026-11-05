# C4 — Context diagram

> **AI-DLC phase:** `02-design`  ·  Level 1 (the system in its environment).

```mermaid
C4Context
    title System context — API Microservice

    Person(client, "API client", "Consumes the HTTP/JSON API")
    Person(operator, "Operator", "Deploys and monitors the service")

    System(api, "API Microservice", "Generated from the archetype; exposes resources over HTTP")

    System_Ext(db, "PostgreSQL", "Persistent storage")
    System_Ext(logs, "Log platform", "Collects stdout logs")

    Rel(client, api, "Calls", "HTTPS/JSON")
    Rel(operator, api, "Probes /healthz, /readyz", "HTTP")
    Rel(api, db, "Reads/writes", "PostgreSQL protocol")
    Rel(api, logs, "Writes structured logs", "stdout")
```

**Notes**
- TLS termination is expected at the ingress / load balancer in front of the service.
