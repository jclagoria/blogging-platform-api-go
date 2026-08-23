# ADR-0001: Use `github.com/rs/cors` for CORS Middleware

- Status: Accepted
- Date: 2026-08-23

## Context

The API needs CORS headers for local frontend development. A CORS library is needed that integrates cleanly with Chi router middleware.

## Decision

Use `github.com/rs/cors` to handle CORS. The library provides a `cors.Handler()` middleware that wraps the Chi router and adds appropriate headers to all responses.

## Consequences

**Positive:**
- Standard, well-maintained library
- Integrates with any `http.Handler` (Chi compatible)
- Handles preflight automatically
- Permissive defaults suitable for local dev

**Negative:**
- Permissive config (`AllowAll`) is not production-ready — origin restrictions needed later
- Adds a dependency (already in project)

## Alternatives Considered

- **Manual header setting**: Rejected — error-prone, misses edge cases (preflight, credentials)
- **`rs/cors` alternatives (`go-chi/cors`, etc.)**: Not considered — `rs/cors` is the de facto standard
