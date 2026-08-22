# Backend Architecture — Blogging Platform API

## Architecture Style

Monolithic Go application with layered architecture:

```
cmd/blogging-platform-api/    — entry point (main.go)
internal/
  generated/                  — oapi-codegen output (server interfaces, types, Chi router)
  handlers/                   — HTTP handlers, store interface, error helpers
  middleware/                  — CORS middleware
  docs/                       — embedded static files (Swagger UI, Redoc, OpenAPI spec)
docs/api/                     — OpenAPI 3.1.0 spec (source of truth)
```

## Request Flow

```
HTTP Request
  → Chi Router (generated)
    → Middleware (CORS)
      → Handler (implements ServerInterface)
        → Store (PostStore interface)
          → MongoDB / InMemory
```

## Key Patterns

- **Design-First**: OpenAPI spec is the source of truth; code is generated from it
- **Interface-Based Stores**: `PostStore` interface allows swapping MongoDB for in-memory (dev/test)
- **RFC 7807 Errors**: All error responses use Problem Details format
- **Embedded Docs**: Swagger UI and Redoc served via Go `embed` package

## Data Flow

1. Client sends HTTP request
2. Chi router dispatches to generated handler stub
3. Handler implementation processes request
4. Store interface abstracts persistence
5. Response serialized as JSON per OpenAPI spec
