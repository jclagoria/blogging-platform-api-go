# ADR-0002: MongoDB Database Selection

- Status: Accepted
- Date: 2026-08-23
- Related: LAG-392

## Context

The tech research digest (derived from repomix outputs) explicitly recommends PostgreSQL as the default database and rejects MongoDB:

> "PostgreSQL default. Full ACID, JSONB, extensions prevent premature polyglot."
> "MongoDB default — No ACID across documents, migration is expensive full rewrite." (Alternatives Rejected)

The blogging platform project uses MongoDB Atlas instead. This ADR documents the justification for this deviation.

## Decision

Use MongoDB Atlas as the primary database, deviating from the digest's PostgreSQL recommendation.

## Rationale

### Why MongoDB fits this project

1. **Schema flexibility**: Blog posts have variable structure (title, content, category, optional tags). MongoDB's document model maps directly to the `Post` struct without ORM mapping or JOIN tables.

2. **Managed service**: MongoDB Atlas provides free tier, automatic scaling, built-in backups, and zero operational overhead — ideal for a personal project.

3. **Simpler data model**: The application has a single entity (`Post`) with no relationships, no transactions, and no complex queries. ACID transactions across documents are unnecessary.

4. **Development speed**: No schema migrations needed. Adding fields to `postDoc` struct is sufficient — MongoDB is schema-less.

5. **Go driver quality**: `go.mongodb.org/mongo-driver/v2` is the official driver with mature BSON support, connection pooling, and context-based timeouts.

### Why PostgreSQL would be overkill

1. **No relational data**: The blog has one entity type. No foreign keys, no JOINs, no referential integrity constraints needed.

2. **No complex queries**: Search is a simple regex across three fields. No window functions, CTEs, or advanced SQL features required.

3. **No transactional requirements**: Creating, updating, or deleting a single post is an atomic operation. No multi-document transactions needed.

4. **Migration cost is low**: If the project needs PostgreSQL later, the data model is simple enough to migrate with a straightforward script.

## Consequences

### Positive

- Zero schema migration overhead
- Free tier on MongoDB Atlas for development
- Direct struct-to-document mapping (no ORM)
- Flexible field additions without ALTER TABLE

### Negative

- No ACID transactions across documents (acceptable for single-document operations)
- Regex search is less efficient than PostgreSQL full-text search at scale
- Potential migration cost if relational features are needed later
- No SQL ecosystem (pgAdmin, psql, standard SQL queries)

## Mitigations

| Concern | Mitigation |
|---------|------------|
| Search efficiency | MongoDB Atlas provides text indexes if regex becomes slow |
| Migration risk | Data model is simple (one collection, flat documents) — export to JSON, import to PostgreSQL |
| No ACID | Single-document operations are atomic in MongoDB |
| Vendor lock-in | MongoDB driver is open-source; Atlas is optional (can self-host) |

## Alternatives Considered

1. **PostgreSQL**: Rejected — adds complexity (schema, migrations, ORM) without benefit for a single-entity CRUD app
2. **SQLite**: Rejected — no managed service, poor concurrent write performance
3. **DynamoDB**: Rejected — AWS lock-in, complex data modeling (per digest)

## Digest Deviation Acknowledgment

This decision deviates from the tech research digest's recommendation. The deviation is justified by:

- Project scope (simple CRUD, single entity)
- Development speed priority
- No current or planned need for ACID transactions
- Low migration cost if requirements change

The deviation should be revisited if:
- The project adds multiple related entities requiring JOINs
- Complex queries exceed regex search capabilities
- Multi-document transactions become necessary
