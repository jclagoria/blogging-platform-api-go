# Blogging Platform API

A Go RESTful API for a personal blogging platform with CRUD operations on blog posts. Design-first approach using OpenAPI 3.1.0 spec as the source of truth, with oapi-codegen for type-safe server stubs and MongoDB Atlas for persistence.

Project from [roadmap.sh](https://roadmap.sh/projects/blogging-platform-api).

## Tech Stack

| Category | Technology |
|----------|------------|
| Language | Go 1.26.3 |
| HTTP Router | Chi (via oapi-codegen) |
| Code Generation | oapi-codegen v2.8.0 |
| Database | MongoDB Atlas |
| Testing | testify |
| Linting | golangci-lint |
| CORS | rs/cors |
| Docs | Swagger UI + Redoc (CDN, Go embed) |

## Installation

```bash
git clone https://github.com/juanka/blogging-platform-api.git
cd blogging-platform-api
go build -o blogging-platform-api ./cmd/blogging-platform-api
```

## Configuration

Copy `.env.example` to `.env` and configure:

```bash
cp .env.example .env
```

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `MONGODB_URI` | No | — | MongoDB Atlas connection string. Unset uses in-memory store |
| `PORT` | No | `8080` | Server listen port |

## Usage

### Health check

```bash
curl http://localhost:8080/health
# 200 OK
```

### Create a blog post

```bash
curl -X POST http://localhost:8080/posts \
  -H "Content-Type: application/json" \
  -d '{
    "title": "My First Blog Post",
    "content": "This is the content of my first blog post.",
    "category": "Technology",
    "tags": ["Tech", "Programming"]
  }'
# 201 Created
```

### Get all blog posts

```bash
curl http://localhost:8080/posts
# 200 OK — returns array of posts
```

### Filter blog posts

```bash
curl "http://localhost:8080/posts?term=tech"
# 200 OK — wildcard search on title, content, category
```

### Get a single blog post

```bash
curl http://localhost:8080/posts/{id}
# 200 OK
```

### Update a blog post

```bash
curl -X PUT http://localhost:8080/posts/{id} \
  -H "Content-Type: application/json" \
  -d '{
    "title": "My Updated Blog Post",
    "content": "Updated content.",
    "category": "Technology",
    "tags": ["Tech"]
  }'
# 200 OK
```

### Delete a blog post

```bash
curl -X DELETE http://localhost:8080/posts/{id}
# 204 No Content
```

## API Documentation

Interactive documentation is served at runtime:

- **Swagger UI**: [http://localhost:8080/docs/](http://localhost:8080/docs/)
- **Redoc**: [http://localhost:8080/redoc/](http://localhost:8080/redoc/)
- **OpenAPI Spec**: [http://localhost:8080/docs/openapi.yaml](http://localhost:8080/docs/openapi.yaml)

## Error Responses

All errors follow RFC 7807 Problem Details format:

```json
{
  "type": "/problems/not-found",
  "title": "Not Found",
  "status": 404,
  "detail": "Post not found"
}
```

Problem types:
- `/problems/not-found` — resource does not exist
- `/problems/validation-error` — request body validation failed
- `/problems/internal-error` — unexpected server error

## Project Structure

```
├── cmd/blogging-platform-api/   # Entry point
├── docs/api/                    # OpenAPI spec (source of truth)
├── internal/
│   ├── generated/               # oapi-codegen output (do not edit)
│   ├── handlers/                # CRUD handlers + PostStore interface
│   ├── middleware/               # CORS middleware
│   └── docs/                    # Embedded Swagger UI + Redoc
├── adr/                         # Architecture Decision Records
```

## Development

```bash
# Run tests
go test ./...

# Lint
golangci-lint run

# Build
go build -o blogging-platform-api ./cmd/blogging-platform-api
```

## Database

- **Development**: In-memory store (default, no `MONGODB_URI` set)
- **Production**: MongoDB Atlas (set `MONGODB_URI` in `.env`)

Post IDs are MongoDB ObjectId hex strings (e.g. `507f1f77bcf86cd799439011`).

## Architecture Decisions

See `adr/` for all architecture decisions. Key ADRs:

- ADR-0001: OpenAPI Version Selection (3.1.0)
- ADR-0002: Go OpenAPI Codegen Tool (oapi-codegen)
- ADR-0005: Error Response Schema (RFC 7807)
- ADR-0011: MongoDB Database Selection (deviation from digest recommendation)

## License

MIT
