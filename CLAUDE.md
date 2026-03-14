# Project: magento2-catalog-graphql-go

High-performance Go drop-in replacement for Magento 2's `products()` GraphQL query using gqlgen.

## Architecture

- **Schema-first GraphQL** via gqlgen — edit `graph/schema.graphqls`, then `go run github.com/99designs/gqlgen generate`
- **Never edit** `graph/generated.go` or `graph/model/models_gen.go` — they are auto-generated
- **Magento Enterprise Edition** — all EAV JOINs use `row_id`, not `entity_id`
- **Read-only** — the service never writes to MySQL

## Project Structure

```
cmd/server/           Entry point
graph/                GraphQL schema, resolvers, generated code
internal/
  app/                HTTP server bootstrap
  cache/              Redis client
  config/             Config loader (Viper: env vars > YAML > defaults)
  database/           MySQL connection (DSN, pooling, UTC timezone)
  middleware/         CORS, caching, logging, panic recovery, store resolution
  repository/         Data access layer — one file per domain (SQL queries)
  service/            Business logic — query orchestration, field detection, type mapping
tests/                Integration + comparison tests (HTTP-based, no internal imports)
```

## Build & Test

```bash
go build -o server ./cmd/server/                    # build
go vet ./...                                        # lint
go test ./tests/ -v -timeout 60s                    # integration tests (needs MySQL)
go test ./tests/ -run TestCompare -v -timeout 300s  # comparison tests (needs Go + Magento running)
```

## Key Conventions

- **Go 1.25** — use current language features
- **Error handling**: wrap with `fmt.Errorf("context: %w", err)`, use `errors.Is`/`errors.As`
- **Naming**: `CamelCase` exported, `camelCase` unexported, no stutter (`repository.NewProductRepository` not `repository.NewRepositoryProduct`)
- **Config**: all settings via env vars (`DB_HOST`, `DB_PORT`, etc.) with sensible defaults; config file optional
- **Logging**: zerolog structured JSON logging; use `log.Info().Str("key", val).Msg("message")`
- **Testing**: table-driven tests, `-race` flag, `t.Helper()` for test helpers
- **Concurrency**: sender closes channels, tie goroutines to `context.Context`
- **Context**: always first parameter `ctx context.Context`
- **Interfaces**: small (1-3 methods), defined at consumer; return concrete types
- **Dependencies**: favor stdlib; justify new dependencies
- **Field-selective loading**: `CollectRequestedFields()` inspects GraphQL AST — only query data the client requested
- **EAV attributes**: loaded once at startup from `eav_attribute` table, cached in `AttributeRepository`

## Common Patterns

### Adding an EAV attribute
1. Add field to `graph/schema.graphqls` → regenerate → add to `ProductEAVValues` struct in `internal/repository/product.go` → map in `internal/service/products.go` `mapProductToModel()`

### Adding a filter
1. Add to `ProductAttributeFilterInput` in schema → regenerate → handle in `FindProducts()` in `internal/repository/product.go`

### Adding a repository
1. Create `internal/repository/your_domain.go` → wire in `graph/resolver.go` → call from service layer
