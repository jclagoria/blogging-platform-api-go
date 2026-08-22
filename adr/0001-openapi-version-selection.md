# ADR-0001: OpenAPI Version Selection

- Status: Accepted
- Date: 2026-08-22
- Related: LAG-378

## Context

The project needs an OpenAPI specification version that oapi-codegen v2.8.0 supports fully without migration concerns.

## Decision

Use OpenAPI 3.1.0.

## Consequences

- Full support in oapi-codegen v2.8.0, no version mismatch
- Newer spec with JSON Schema 2020-12 alignment
- Fewer third-party validation tools than 3.0, but sufficient for our needs

## Alternatives Considered

- OpenAPI 3.0.x: wider tooling support but oapi-codegen would need version handling
