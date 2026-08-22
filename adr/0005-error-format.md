# ADR-0005: Error Response Schema

- Status: Accepted
- Date: 2026-08-22
- Related: LAG-382

## Context

All error responses need a consistent, structured format.

## Decision

Use RFC 7807 Problem Details for all error responses.

## Consequences

- Fields: `type`, `title`, `status`, `detail`, `errors[]` (validation only)
- Problem types: `about:blank`, `/problems/not-found`, `/problems/validation-error`, `/problems/internal-error`
- Validation error codes: `REQUIRED`, `INVALID_FORMAT`, `TOO_SHORT`, `TOO_LONG`, `TOO_SMALL`, `TOO_LARGE`
- Each validation error has: `field`, `message`, `code`

## Alternatives Considered

- Custom error format: less interoperable, reinventing a standard
- Plain string errors: no structured parsing possible
