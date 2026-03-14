---
name: add-attribute
description: Add a new EAV product attribute to the GraphQL schema. Use when adding custom Magento attributes to the API.
argument-hint: <attribute_code>
---

Add the Magento EAV attribute `$ARGUMENTS` to the GraphQL API. Follow these steps exactly:

## 1. Schema

Add the field to `ProductInterface` in `graph/schema.graphqls`. Use the appropriate GraphQL type:
- `varchar` / `text` backend → `String`
- `int` backend → `Int`
- `decimal` backend → `Float`
- `datetime` backend → `String`
- `select` backend → `Int` (option_id) or `String` (label), depending on use case

## 2. Regenerate

```bash
go run github.com/99designs/gqlgen generate
```

## 3. Repository

In `internal/repository/product.go`:
- Add the field to the `ProductEAVValues` struct
- Add the EAV JOIN in `FindProducts()` using the attribute's `backend_type` to determine which value table (`catalog_product_entity_varchar`, `_int`, `_text`, `_decimal`, `_datetime`)
- Use `COALESCE(store_value.value, default_value.value)` for store scoping
- JOINs use `row_id` (Magento Enterprise Edition)

## 4. Service mapping

In `internal/service/products.go`, in `mapProductToModel()`, map the new field from `ProductEAVValues` to the generated model type.

## 5. Verify

```bash
go build ./...
go vet ./...
```

## Important

- The attribute must exist in Magento's `eav_attribute` table with `entity_type_code = 'catalog_product'`
- The `AttributeRepository` cache (loaded at startup) provides `attribute_id` and `backend_type`
- Check existing attributes in `ProductEAVValues` for patterns to follow
