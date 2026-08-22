# ADR-0006: Test Strategy

- Status: Accepted
- Date: 2026-08-22
- Related: LAG-383

## Context

Need a testing approach that covers handlers, middleware, and error paths without testing generated code.

## Decision

- testify (`assert`/`require`, no suite) for assertions
- Interface-based PostStore fake for mocking
- Table-driven subtests with `t.Run()`
- httptest for all HTTP handler tests
- Test files alongside source (`*_test.go`)

## Consequences

- Full coverage of handlers, middleware, error paths
- No external mock library needed (5-6 methods, trivial to maintain)
- Tests are co-located with code they exercise
- oapi-codegen generated code is skipped (library output)

## Alternatives Considered

- stdlib testing only: more verbose assertions
- gomock/mockgen: overkill for a 5-method interface
- Suite-based tests: adds abstraction layer without benefit
