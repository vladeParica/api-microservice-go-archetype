# Contributing

Thanks for your interest in improving the archetype.

## Workflow

1. Branch from `develop`: `feature/<name>`, `fix/<name>`, `docs/<name>` or `chore/<name>`.
2. Keep commits atomic and use [Conventional Commits](https://www.conventionalcommits.org/) in English:
   `feat(user): add list endpoint`.
3. Run the checks before pushing:
   ```bash
   make fmt vet lint test
   make sqlc   # if you changed db/queries or db/migrations
   ```
4. Open a pull request against `develop`. `main` only receives releases from `develop`.

## Design changes

Non-trivial decisions need an ADR in [`docs/00-project/adr/`](docs/00-project/adr/README.md), and the
affected [lifecycle docs](docs/README.md) must be updated in the same pull request.

## Code of conduct

By participating you agree to the [Code of Conduct](CODE_OF_CONDUCT.md).
