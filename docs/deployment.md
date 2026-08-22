# Deployment — Blogging Platform API

## Local Development

```bash
# Run with in-memory store (no MongoDB required)
go run ./cmd/blogging-platform-api/

# Run with MongoDB Atlas
# Ensure .env has MONGODB_URI and MONGODB_DATABASE
go run ./cmd/blogging-platform-api/
```

## Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `MONGODB_URI` | No | MongoDB Atlas connection string. Falls back to in-memory store if unset. |
| `MONGODB_DATABASE` | No | MongoDB database name. Defaults to `blogging-platform` if unset. |

## Build

```bash
go build -o blogging-platform-api ./cmd/blogging-platform-api/
```

## API Documentation

Once running:
- Swagger UI: `http://localhost:8080/docs/swagger-ui/`
- Redoc: `http://localhost:8080/docs/redoc/`
- OpenAPI spec: `http://localhost:8080/docs/openapi.yaml`

## Production

Not yet defined. Static Go binary can be containerized or deployed directly.
