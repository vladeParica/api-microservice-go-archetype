# api-microservice-go-archetype

[![CI](https://github.com/vladeParica/api-microservice-go-archetype/actions/workflows/ci.yml/badge.svg)](https://github.com/vladeParica/api-microservice-go-archetype/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/vladeParica/api-microservice-go-archetype.svg)](https://pkg.go.dev/github.com/vladeParica/api-microservice-go-archetype)
[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)

Production-ready template for Go HTTP microservices: standard library first, PostgreSQL with sqlc,
structured logging, health probes and graceful shutdown.

## Quick start

```bash
go install golang.org/x/tools/cmd/gonew@latest
gonew github.com/vladeParica/api-microservice-go-archetype github.com/your-org/your-service

make up                      # postgres + migrations + api on :8080
curl localhost:8080/readyz
```

Run `make help` for all tasks.

## Documentation

| | |
|---|---|
| [Architecture](docs/02-design/architecture.md) | Components, stack and C4 diagrams |
| [Decisions](docs/00-project/adr/README.md) | Architecture Decision Records |
| [API](docs/02-design/api-contracts.md) | Endpoints and error format |
| [Deployment](docs/05-deployment/deploy-guide.md) | Configuration and runbook |

Documentation follows the [AI-DLC](docs/) lifecycle (`docs/00-project` → `06-monitoring`).

## License

[Apache 2.0](LICENSE)
