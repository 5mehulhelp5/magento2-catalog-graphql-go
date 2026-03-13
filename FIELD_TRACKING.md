# ProductInterface Field Implementation Tracking

## Legend
- [ ] Not started
- [~] In progress
- [x] Implemented
- [-] Skipped (deprecated/not applicable)

---

## ProductInterface Core Fields

| # | Field | Type | Status | Phase | Notes |
|---|-------|------|--------|-------|-------|
| 1 | `uid` | `ID!` | [x] | 2 | Base64-encoded entity_id |
| 2 | `id` | `Int` | [x] | 2 | Deprecated, alias for entity_id |
| 3 | `name` | `String` | [x] | 1 | EAV varchar, store-scoped |
| 4 | `sku` | `String` | [x] | 1 | catalog_product_entity.sku |
| 5 | `type_id` | `String` | [x] | 2 | Deprecated, use __typename |
| 6 | `description` | `ComplexTextValue` | [x] | 2 | EAV text, store-scoped |
| 7 | `short_description` | `ComplexTextValue` | [x] | 2 | EAV text, store-scoped |
| 8 | `special_price` | `Float` | [x] | 2 | EAV decimal |
| 9 | `special_from_date` | `String` | [x] | 2 | EAV datetime, deprecated |
| 10 | `special_to_date` | `String` | [x] | 2 | EAV datetime |
| 11 | `attribute_set_id` | `Int` | [x] | 2 | Deprecated |
| 12 | `meta_title` | `String` | [x] | 2 | EAV varchar |
| 13 | `meta_keyword` | `String` | [x] | 2 | EAV varchar |
| 14 | `meta_description` | `String` | [x] | 2 | EAV varchar |
| 15 | `new_from_date` | `String` | [x] | 2 | EAV datetime |
| 16 | `new_to_date` | `String` | [x] | 2 | EAV datetime |
| 17 | `tier_price` | `Float` | [-] | 2 | Deprecated, use price_tiers |
| 18 | `options_container` | `String` | [x] | 2 | EAV varchar |
| 19 | `created_at` | `String` | [x] | 2 | Deprecated, entity table |
| 20 | `updated_at` | `String` | [x] | 2 | Deprecated, entity table |
| 21 | `country_of_manufacture` | `String` | [x] | 2 | EAV varchar |
| 22 | `manufacturer` | `Int` | [x] | 2 | EAV int |
| 23 | `gift_message_available` | `String` | [x] | 2 | EAV varchar |
| 24 | `canonical_url` | `String` | [x] | 2 | Computed |

## Image Fields

| # | Field | Type | Status | Phase | Notes |
|---|-------|------|--------|-------|-------|
| 25 | `image` | `ProductImage` | [x] | 2 | EAV varchar → full URL |
| 26 | `small_image` | `ProductImage` | [x] | 2 | EAV varchar → full URL |
| 27 | `thumbnail` | `ProductImage` | [x] | 2 | EAV varchar → full URL |
| 28 | `swatch_image` | `String` | [x] | 2 | EAV varchar |
| 29 | `media_gallery` | `[MediaGalleryInterface]` | [x] | 2 | ProductImage + ProductVideo |
| 30 | `media_gallery_entries` | `[MediaGalleryEntry]` | [-] | 2 | Deprecated, use media_gallery |

## Price Fields

| # | Field | Type | Status | Phase | Notes |
|---|-------|------|--------|-------|-------|
| 31 | `price_range` | `PriceRange!` | [x] | 2 | catalog_product_index_price |
| 32 | `price` | `ProductPrices` | [-] | 2 | Deprecated, use price_range |
| 33 | `price_tiers` | `[TierPrice]` | [x] | 2 | catalog_product_entity_tier_price |
| 34 | `tier_prices` | `[ProductTierPrices]` | [-] | 2 | Deprecated, use price_tiers |

## Category Fields

| # | Field | Type | Status | Phase | Notes |
|---|-------|------|--------|-------|-------|
| 35 | `categories` | `[CategoryInterface]` | [x] | 2 | catalog_category_product JOIN |

## Inventory Fields

| # | Field | Type | Status | Phase | Notes |
|---|-------|------|--------|-------|-------|
| 36 | `stock_status` | `ProductStockStatus` | [x] | 2 | cataloginventory_stock_status |
| 37 | `only_x_left_in_stock` | `Float` | [x] | 2 | stock item qty, threshold-based |
| 38 | `quantity` | `Float` | [x] | 2 | stock item qty |
| 39 | `min_sale_qty` | `Float` | [x] | 2 | cataloginventory_stock_item |
| 40 | `max_sale_qty` | `Float` | [x] | 2 | cataloginventory_stock_item |

## URL Fields

| # | Field | Type | Status | Phase | Notes |
|---|-------|------|--------|-------|-------|
| 41 | `url_key` | `String` | [x] | 2 | EAV varchar |
| 42 | `url_suffix` | `String` | [x] | 2 | From core_config_data |
| 43 | `url_path` | `String` | [-] | 2 | Deprecated, use url_key + url_suffix |
| 44 | `url_rewrites` | `[UrlRewrite]` | [x] | 2 | url_rewrite table |

## Related Products

| # | Field | Type | Status | Phase | Notes |
|---|-------|------|--------|-------|-------|
| 45 | `related_products` | `[ProductInterface]` | [x] | 4 | catalog_product_link (type=1) |
| 46 | `upsell_products` | `[ProductInterface]` | [x] | 4 | catalog_product_link (type=4) |
| 47 | `crosssell_products` | `[ProductInterface]` | [x] | 4 | catalog_product_link (type=5) |
| 48 | `product_links` | `[ProductLinksInterface]` | [-] | 4 | Deprecated, use related/upsell/crosssell |

## Review Fields

| # | Field | Type | Status | Phase | Notes |
|---|-------|------|--------|-------|-------|
| 49 | `rating_summary` | `Float!` | [x] | 4 | review_entity_summary |
| 50 | `review_count` | `Int!` | [x] | 4 | review_entity_summary |
| 51 | `reviews` | `ProductReviews!` | [x] | 4 | review / review_detail tables |

## Custom Attributes

| # | Field | Type | Status | Phase | Notes |
|---|-------|------|--------|-------|-------|
| 52 | `custom_attributesV2` | `ProductCustomAttributes` | [-] | 4 | Dynamic EAV attributes, requires runtime schema |

## Websites (deprecated)

| # | Field | Type | Status | Phase | Notes |
|---|-------|------|--------|-------|-------|
| 53 | `websites` | `[Website]` | [-] | 4 | Deprecated |

---

## Product Types

| Type | Status | Phase | Notes |
|------|--------|-------|-------|
| `SimpleProduct` | [x] | 1-2 | ProductInterface + PhysicalProductInterface (weight) + CustomizableProductInterface |
| `ConfigurableProduct` | [x] | 4 | + variants, configurable_options, swatch_data |
| `VirtualProduct` | [x] | 4 | ProductInterface + CustomizableProductInterface (no weight) |
| `GroupedProduct` | [x] | 4 | + items: [GroupedProductItem] (stub, no data in DB) |
| `BundleProduct` | [x] | 4 | + bundle_items, options, dynamic pricing, ship_bundle_items |

## ConfigurableProduct Sub-Fields

| Field | Type | Status | Phase |
|-------|------|--------|-------|
| `variants` | `[ConfigurableVariant]` | [x] | 4 |
| `configurable_options` | `[ConfigurableProductOptions]` | [x] | 4 |
| `configurable_product_options_selection` | `ConfigurableProductOptionsSelection` | [-] | 4 |

## PhysicalProductInterface

| Field | Type | Status | Phase |
|-------|------|--------|-------|
| `weight` | `Float` | [x] | 2 |

## CustomizableProductInterface

| Field | Type | Status | Phase |
|-------|------|--------|-------|
| `options` | `[CustomizableOptionInterface]` | [-] | 4 |

## RoutableInterface

| Field | Type | Status | Phase |
|-------|------|--------|-------|
| `relative_url` | `String` | [-] | 2 |
| `redirect_code` | `Int!` | [-] | 2 |
| `type` | `UrlRewriteEntityTypeEnum` | [-] | 2 |

---

## Products Response Fields

| Field | Type | Status | Phase |
|-------|------|--------|-------|
| `items` | `[ProductInterface]` | [x] | 1 |
| `total_count` | `Int` | [x] | 3 |
| `page_info` | `SearchResultPageInfo` | [x] | 3 |
| `aggregations` | `[Aggregation]` | [x] | 5 |
| `sort_fields` | `SortFields` | [x] | 5 |
| `suggestions` | `[SearchSuggestion]` | [x] | 5 |
| `filters` | `[LayerFilter]` | [-] | 5 |

---

## Filter Input Fields (ProductAttributeFilterInput)

| Field | Type | Status | Phase |
|-------|------|--------|-------|
| `sku` | `FilterEqualTypeInput` | [x] | 3 |
| `name` | `FilterMatchTypeInput` | [x] | 3 |
| `category_uid` | `FilterEqualTypeInput` | [x] | 3 |
| `category_id` | `FilterEqualTypeInput` | [x] | 3 |
| `category_url_path` | `FilterEqualTypeInput` | [x] | 3 |
| `url_key` | `FilterEqualTypeInput` | [x] | 3 |
| `price` | `FilterRangeTypeInput` | [x] | 3 |
| Dynamic EAV filters | Various | [-] | 3 |

## Sort Input Fields (ProductAttributeSortInput)

| Field | Type | Status | Phase |
|-------|------|--------|-------|
| `relevance` | `SortEnum` | [x] | 3 |
| `position` | `SortEnum` | [x] | 3 |
| `name` | `SortEnum` | [x] | 3 |
| `price` | `SortEnum` | [x] | 3 |
| Dynamic EAV sorts | `SortEnum` | [-] | 3 |
