# ADR-0002: Go OpenAPI Codegen Tool

- Status: Accepted
- Date: 2026-08-22
- Related: LAG-379

## Context

Need a Go code generator that produces server interfaces and router stubs from an OpenAPI 3.1.0 spec.

## Decision

Use oapi-codegen v2.8.0.

## Consequences

- Go-native, generates Chi router and server interfaces
- Generated code goes to `internal/generated/`
- Must not edit generated code directly
- Requires Go 1.25+

## Alternatives Considered

- go-swagger: older, less active, OpenAPI 3.0 only
- libopenapi: library for reading specs, not code generation
