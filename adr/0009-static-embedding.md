# ADR-0001: Static File Embedding for Documentation

- Status: Accepted
- Date: 2026-08-22

## Context

The API needs to serve Swagger UI and Redoc for interactive documentation. Static files (HTML, CSS, JS) must be available at runtime without external file dependencies.

## Decision

Use Go's `embed` package to embed static files into the binary at build time. HTML files load Swagger UI / Redoc from CDN (unpkg.com) rather than bundling the full libraries.

## Consequences

- **Positive**: No runtime file I/O, single binary deployment, no external file dependencies
- **Positive**: Standard library solution, no new dependencies
- **Negative**: Binary size increases marginally
- **Negative**: CDN dependency for JS/CSS (requires internet on first load)

## Alternatives Considered

1. **Serve from filesystem**: Rejected — requires files to exist on disk at runtime
2. **Bundle Swagger UI / Redoc**: Rejected — increases binary size significantly, CDN is simpler
3. **Use a third-party Go library**: Rejected — `embed` is stdlib and sufficient
