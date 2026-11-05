# Security policy

## Supported versions

Only the latest release and the `develop` branch receive security fixes.

## Reporting a vulnerability

Please **do not open a public issue**. Report it privately through
[GitHub Security Advisories](https://github.com/vladeParica/api-microservice-go-archetype/security/advisories/new).

You can expect an acknowledgement within 5 business days. Once a fix is available, the advisory is
published with credit to the reporter unless anonymity is requested.

## Scope

The archetype intentionally ships without authentication or rate limiting (see the
[threat model](docs/02-design/threat-model.md)). Services generated from it must add both before
going to production.
