# Project Status Report — 2026-03-13

## Overview

Go drop-in replacement for Magento 2's `products()` GraphQL query. Reads directly from Magento's MySQL database (EAV schema). Built with gqlgen (schema-first).

- **Build**: Compiles cleanly
- **Tests**: 19/19 passing (0.9s)
- **No TODOs or placeholder code**

---

## Architecture

```
HTTP Client
    │
    ▼
RecoveryMiddleware → CORSMiddleware → LoggingMiddleware → StoreMiddleware → CacheMiddleware
    │
    ▼
GraphQL Handler (gqlgen)
    │
    ▼
ProductService.GetProducts()
    │
    ├── ProductRepository     (EAV JOINs, filters, sorting, pagination)
    ├── PriceRepository       (price_range, tier_prices)
    ├── MediaRepository       (media_gallery)
    ├── InventoryRepository   (stock_status, qty, min/max sale qty)
    ├── CategoryRepository    (product→category mapping)
    ├── URLRepository         (url_rewrites)
    ├── ConfigurableRepository (variants, options, swatches)
    ├── BundleRepository      (bundle items, options, dynamic pricing)
    ├── ProductLinkRepository (related/upsell/crosssell)
    ├── AggregationRepository (faceted filters)
    ├── ReviewRepository      (ratings, reviews)
    ├── SearchRepository      (search suggestions)
    └── StoreConfigRepository (currency, base URL, thresholds)
```

---

## Implementation Summary

| Category | Implemented | Skipped | Open |
|----------|-------------|---------|------|
| Core fields (#1-24) | 23 | 1 | 0 |
| Image fields (#25-30) | 5 | 1 | 0 |
| Price fields (#31-34) | 2 | 2 | 0 |
| Categories (#35) | 1 | 0 | 0 |
| Inventory (#36-40) | 5 | 0 | 0 |
| URL fields (#41-44) | 3 | 1 | 0 |
| Related products (#45-48) | 3 | 1 | 0 |
| Reviews (#49-51) | 3 | 0 | 0 |
| Custom attributes (#52) | 0 | 1 | 0 |
| Websites (#53) | 0 | 1 | 0 |
| Product types (5) | 5 | 0 | 0 |
| Configurable sub-fields | 2 | 1 | 0 |
| Physical/Customizable/Routable | 1 | 4 | 0 |
| Response fields (6) | 5 | 1 | 0 |
| Filters (7+dynamic) | 6 | 1 | 0 |
| Sorts (4+dynamic) | 4 | 1 | 0 |
| **Totals** | **68** | **16** | **0** |

---

## Skipped Fields (with justification)

| Field | Reason |
|-------|--------|
| `tier_price` | Deprecated → use `price_tiers` |
| `price` (ProductPrices) | Deprecated → use `price_range` |
| `tier_prices` | Deprecated → use `price_tiers` |
| `media_gallery_entries` | Deprecated → use `media_gallery` |
| `url_path` | Deprecated → use `url_key` + `url_suffix` |
| `product_links` | Deprecated → use `related_products`/`upsell`/`crosssell` |
| `filters` (LayerFilter) | Deprecated → use `aggregations` |
| `websites` | Deprecated |
| `custom_attributesV2` | Requires runtime schema introspection |
| `configurable_product_options_selection` | Input-driven variant resolution, rarely used |
| `CustomizableProductInterface.options` | No customizable options in catalog |
| `RoutableInterface` (3 fields) | Used by `route()` query, not `products()` |
| Dynamic EAV filters/sorts | Requires dynamic schema generation |

---

## Infrastructure

| Component | Status |
|-----------|--------|
| MySQL connection pooling | 25 open / 10 idle |
| Redis response caching | 5min TTL, optional |
| CORS middleware | Allow all origins |
| Store-scoped multi-tenancy | Via `Store` header |
| Request logging (zerolog) | Structured JSON |
| Panic recovery with stack trace | Yes |
| GraphQL complexity limiting | 1000 / depth 15 |
| Health check (`/health`) | DB ping |
| Graceful shutdown | 10s timeout |
| Field-selective batch loading | Skip unused data |

---

## Key Files

| File | Lines | Purpose |
|------|-------|---------|
| `graph/schema.graphqls` | 1220 | GraphQL schema |
| `graph/model/models_gen.go` | 2968 | Generated types |
| `internal/service/products.go` | 1096 | Query orchestration |
| `internal/repository/product.go` | 370 | EAV product queries |
| `internal/repository/configurable.go` | 411 | Configurable product data |
| `internal/repository/aggregation.go` | 281 | Faceted filters |
| `internal/repository/bundle.go` | 213 | Bundle product data |
| `internal/app/app.go` | 137 | Application bootstrap |
| `integration_test.go` | ~600 | 19 integration tests |

---

## Open Questions

1. **RoutableInterface** — `relative_url`/`redirect_code`/`type` skipped. Schema declares `SimpleProduct implements RoutableInterface`. Does the frontend query these on product results?
2. **CustomizableProductInterface.options** — Returns nil. Does the catalog have customizable options?
3. **`configurable_product_options_selection`** — Frontend using this or just `variants` + `configurable_options`?
4. **Search** — MySQL LIKE on 4 fields. Sufficient for catalog size, or need OpenSearch?
5. **Cache invalidation** — 5min TTL only. Need active invalidation on product updates?
6. **`custom_attributesV2`** — Only non-deprecated skipped field. Frontend querying it?

---

## Database

- **Host**: localhost:3306
- **Schema**: `m24-ploom-ae`
- **Products**: 85 (53 simple, 11 configurable, 21 bundle)
- **Currency**: AED
- **Store**: website_id=1, store_id=2
- **Edition**: Enterprise (uses `row_id` in EAV value tables)
