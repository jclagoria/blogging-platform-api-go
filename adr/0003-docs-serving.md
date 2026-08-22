# ADR-0003: Serving Swagger UI from Go API

- Status: Accepted
- Date: 2026-08-22
- Related: LAG-380

## Context

The API needs interactive documentation (Swagger UI) and reference documentation (Redoc) accessible via HTTP.

## Decision

Use Go `embed` package to bundle Swagger UI and Redoc static files into the binary.

## Consequences

- Single binary deployment, no external file dependencies
- Static files must be present at build time in `docs/api/`
- Swagger UI at `/docs/`, Redoc at `/redoc/`
- No runtime file serving overhead

## Alternatives Considered

- Separate static file server: adds deployment complexity
- CDN-hosted Swagger UI: adds external dependency
