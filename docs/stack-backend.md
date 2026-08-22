# Backend Stack — Blogging Platform API

## Decision Summary

| Category | Decision | Rationale | Trade-offs |
|----------|----------|-----------|------------|
| Backend Runtime | Go 1.26.3 | Strong stdlib, fast compilation, static binary | Less ergonomic than higher-level languages |
| HTTP Router | Chi (via oapi-codegen) | Lightweight, stdlib-compatible, generated from spec | Less feature-rich than Gin/Echo |
| Code Generation | oapi-codegen v2.8.0 | Type-safe server stubs from OpenAPI 3.1.0 spec | Requires spec-first workflow |
| Database | MongoDB Atlas | Flexible document schema, managed service | No ACID transactions for simple CRUD |
| MongoDB Driver | go.mongodb.org/mongo-driver/v2 v2.8.0 | Official driver, v2 API | — |
| Env Loading | godotenv | Loads `.env` files automatically | — |
| CORS | rs/cors | Simple, permissive for local dev | — |
| Testing | testify | Popular, simple assertions | Additional dependency |
| Docs Serving | Go embed + Swagger UI/Redoc CDN | No runtime dependencies, built into binary | CDN dependency for UI |

## Dependencies

```
github.com/getkin/kin-openapi v0.132.0
github.com/go-chi/chi/v5 v5.2.3
github.com/joho/godotenv v1.5.1
github.com/rs/cors v1.11.0
github.com/stretchr/testify v1.11.1
go.mongodb.org/mongo-driver/v2 v2.8.0
github.com/oapi-codegen/runtime v1.1.1
```
