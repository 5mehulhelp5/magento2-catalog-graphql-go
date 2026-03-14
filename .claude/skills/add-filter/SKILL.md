---
name: add-filter
description: Add a new product filter to the GraphQL schema. Use when adding filterable attributes like brand, color, or custom attributes.
argument-hint: <filter_name>
---

Add a new product filter `$ARGUMENTS` to the GraphQL API. Follow these steps:

## 1. Schema

Add the filter to `ProductAttributeFilterInput` in `graph/schema.graphqls`:

```graphql
input ProductAttributeFilterInput {
    # ... existing filters ...
    $ARGUMENTS: FilterEqualTypeInput    # or FilterRangeInput for numeric ranges
}
```

Choose the filter type:
- `FilterEqualTypeInput` — for `eq` and `in` operations (SKU, brand, color)
- `FilterRangeInput` — for `from`/`to` ranges (price, weight)
- `FilterMatchTypeInput` — for `match` text search (name)

## 2. Regenerate

```bash
go run github.com/99designs/gqlgen generate
```

## 3. Repository

In `internal/repository/product.go` `FindProducts()`, add filter handling. Follow the existing patterns for `filter.Sku`, `filter.Name`, etc.

For EAV attribute filters:
- Look up the attribute via `r.attrRepo.GetByCode("$ARGUMENTS")`
- JOIN the appropriate EAV value table based on `BackendType`
- Use `row_id` for JOINs (Enterprise Edition)
- Include `store_id IN (0, ?)` for store scoping

For flat table filters:
- Add a WHERE clause directly on `catalog_product_entity` columns

## 4. Verify

```bash
go build ./...
go vet ./...
```

Test with a GraphQL query:
```graphql
{ products(filter: { $ARGUMENTS: { eq: "test" } }) { items { sku name } total_count } }
```
