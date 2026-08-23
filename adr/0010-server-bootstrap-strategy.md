# ADR-0001: Server Bootstrap Strategy

## Status

Accepted

## Context

The blogging platform needs to wire together generated code, handlers, middleware, MongoDB, and embedded docs into a running server. The entry point must:
- Initialize the Chi router from oapi-codegen output
- Register all handlers (health, posts)
- Apply CORS middleware
- Connect to MongoDB Atlas (or fall back to in-memory store)
- Serve embedded documentation
- Handle graceful shutdown

## Decision

Use a single `main.go` file with the following structure:

1. Load environment variables with godotenv
2. Initialize MongoDB connection (if MONGODB_URI is set)
3. Create PostStore (MongoPostStore or InMemoryPostStore based on env)
4. Create handler struct implementing ServerInterface
5. Initialize Chi router from generated code
6. Register handlers on router
7. Apply CORS middleware
8. Start HTTP server
9. Listen for SIGINT/SIGTERM via signal.NotifyContext
10. On shutdown: close server, close MongoDB connection

## Consequences

### Positive
- Simple, no unnecessary abstractions
- All wiring visible in one place
- Easy to understand and modify
- Environment-based store selection works for dev and production

### Negative
- All initialization logic in one function
- No dependency injection framework
- Harder to unit test the wiring itself (but handlers/stores are tested separately)

## Alternatives Considered

1. **Dependency injection container**: More testable but adds complexity for a simple app
2. **Config struct with validation**: More structured but godotenv + os.Getenv is sufficient
3. **Separate init functions**: More modular but unnecessary indirection
