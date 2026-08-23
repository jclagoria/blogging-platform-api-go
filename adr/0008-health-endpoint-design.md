# ADR-0001: Health Endpoint Design

- **Status**: Accepted
- **Date**: 2026-08-22
- **Deciders**: Lagoria team

## Context

The health endpoint is the first vertical slice implementing the API contract. It proves the entire pipeline works: OpenAPI spec → oapi-codegen → handler → test. It needs to be simple, fast, and require no external dependencies.

## Decision

Implement the health endpoint as a static JSON response handler that returns `{"status": "ok"}` without any store dependencies or business logic.

### Alternatives Considered

1. **Database connectivity check**: Rejected — adds complexity for a first vertical slice; can be added later if needed.
2. **Dependency health checks**: Rejected — premature for initial implementation.

## Consequences

- ✅ Validates the entire pipeline works end-to-end
- ✅ No external dependencies required
- ✅ Fast, predictable response
- ✅ No authentication needed
- ⚠️ Does not report database or service health (acceptable for initial implementation)
