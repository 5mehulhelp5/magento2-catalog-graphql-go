---
name: compare
description: Compare Go service responses with Magento PHP for a specific GraphQL query. Use when debugging field-level differences.
argument-hint: <graphql-query>
disable-model-invocation: true
---

Compare the Go service and Magento PHP responses for a GraphQL query.

## Prerequisites

Both services must be running:
- Go service at `GO_GRAPHQL_URL` (default: http://localhost:8080/graphql)
- Magento at `MAGE_GRAPHQL_URL` (default: http://localhost/graphql)

## Steps

### 1. Query both services

Send the same GraphQL query to both endpoints with the `Store: default` header. Time both requests.

```bash
# Go service
curl -s -w "\nTime: %{time_total}s" \
  -H 'Content-Type: application/json' \
  -H 'Store: default' \
  -d '{"query":"$ARGUMENTS"}' \
  ${GO_GRAPHQL_URL:-http://localhost:8080/graphql}

# Magento
curl -s -w "\nTime: %{time_total}s" \
  -H 'Content-Type: application/json' \
  -H 'Store: default' \
  -d '{"query":"$ARGUMENTS"}' \
  ${MAGE_GRAPHQL_URL:-http://localhost/graphql}
```

### 2. Compare responses

Parse both JSON responses and compare field-by-field:
- Report matching fields
- Report mismatches with both values
- Report fields present in one response but not the other
- Use float tolerance of 0.01 for price comparisons
- Sort arrays by `sku` or `value` before comparing for stable results

### 3. Report

Show:
- Response time comparison (Go vs Magento)
- Total fields compared
- Mismatches with path, Go value, and Magento value
- Overall verdict: IDENTICAL or list of differences
