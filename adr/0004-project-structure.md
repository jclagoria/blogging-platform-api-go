# ADR-0004: Go Project Directory Structure

- Status: Accepted
- Date: 2026-08-22
- Related: LAG-381

## Context

The project needs a standard Go layout that other Go developers will recognize.

## Decision

Use standard Go project layout:

- `cmd/blogging-platform-api/` — application entry point
- `docs/api/` — OpenAPI spec and static doc files
- `internal/generated/` — oapi-codegen output
- `internal/handlers/` — handler implementations
- `internal/middleware/` — middleware (CORS)

## Consequences

- Familiar to Go developers
- `internal/` prevents external imports
- Clear separation between generated and hand-written code

## Alternatives Considered

- Flat layout: less scalable, harder to navigate
