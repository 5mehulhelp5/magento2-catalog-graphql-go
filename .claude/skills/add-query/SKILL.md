---
name: add-query
description: Add a new root-level GraphQL query (e.g., categoryList, cmsPage). Use when extending the API beyond products().
argument-hint: <query_name>
---

Add a new root-level GraphQL query `$ARGUMENTS`. Follow these steps:

## 1. Schema

In `graph/schema.graphqls`, add the query to the `Query` type and define all needed types:

```graphql
type Query {
    products(...): Products
    $ARGUMENTS(...): YourReturnType
}
```

Define input types, return types, and any nested types needed.

## 2. Regenerate

```bash
go run github.com/99designs/gqlgen generate
```

This creates a resolver stub in `graph/schema.resolvers.go`.

## 3. Repository

Create `internal/repository/<domain>.go`:
- Struct with `*sql.DB` field (and `*AttributeRepository` if using EAV)
- Constructor: `NewXxxRepository(db *sql.DB) *XxxRepository`
- Query methods that accept `ctx context.Context` as first parameter
- Use `row_id` for EAV value table JOINs (Enterprise Edition)
- Use `COALESCE(store_value, default_value)` for store-scoped attributes

## 4. Service (optional)

If business logic is needed beyond simple data fetching, create `internal/service/<domain>.go` following the `ProductService` pattern.

## 5. Wiring

In `graph/resolver.go`:
- Add the new repository/service to the `Resolver` struct
- Instantiate it in `NewResolver()`

## 6. Resolver

In `graph/schema.resolvers.go`, implement the generated stub by calling your service/repository.

## 7. Verify

```bash
go build ./...
go vet ./...
```

## Patterns to follow

- Look at `graph/resolver.go` for dependency wiring
- Look at `graph/schema.resolvers.go` for how `Products()` delegates to `ProductService`
- Look at any file in `internal/repository/` for SQL query patterns
