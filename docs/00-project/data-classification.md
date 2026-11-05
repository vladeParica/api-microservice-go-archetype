# Data classification — API Microservice Go Archetype

> **AI-DLC phase:** `00-project`
> Inventory every kind of data by sensitivity so design is driven by it, not guessed.

| Data | Category | Sensitivity | Handling requirements |
|---|---|---|---|
| User name and email (`users` table) | PII | Confidential | TLS to the database in production (`DATABASE_SSLMODE`), never logged, least-privilege DB user |
| Database credentials | Secret | Restricted | Injected via environment / secret manager, never committed or logged |
| Request logs (method, path, status, request id) | Operational | Internal | No bodies or PII; retention set by the log platform |
| Configuration defaults (`configs/config.yaml`) | Configuration | Internal | Local development values only |

**Sensitivity levels:** `Public` < `Internal` < `Confidential` < `Restricted`.

## Cross-cutting requirements
- Request bodies are never logged; error responses never expose internal errors.
- Services generated from the archetype must re-classify their own data before going to production.
- `<TODO — Human-in-the-Loop: applicable regulation (e.g. GDPR) and retention policy per service>`
