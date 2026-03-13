package repository

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
)

// AttributeMetadata holds the metadata for a single EAV attribute.
type AttributeMetadata struct {
	AttributeID  int
	Code         string
	BackendType  string // varchar, int, text, decimal, datetime, static
	FrontendInput string
}

// AttributeRepository manages EAV attribute metadata with in-memory caching.
type AttributeRepository struct {
	db    *sql.DB
	cache map[string]*AttributeMetadata // keyed by attribute_code
	mu    sync.RWMutex
	loaded bool
}

func NewAttributeRepository(db *sql.DB) *AttributeRepository {
	return &AttributeRepository{
		db:    db,
		cache: make(map[string]*AttributeMetadata),
	}
}

// LoadProductAttributes loads all catalog_product attribute metadata into cache.
func (r *AttributeRepository) LoadProductAttributes(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.loaded {
		return nil
	}

	query := `
		SELECT ea.attribute_id, ea.attribute_code, ea.backend_type, COALESCE(ea.frontend_input, '')
		FROM eav_attribute ea
		JOIN eav_entity_type et ON ea.entity_type_id = et.entity_type_id
		WHERE et.entity_type_code = 'catalog_product'
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return fmt.Errorf("load attributes failed: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		attr := &AttributeMetadata{}
		if err := rows.Scan(&attr.AttributeID, &attr.Code, &attr.BackendType, &attr.FrontendInput); err != nil {
			return fmt.Errorf("scan attribute failed: %w", err)
		}
		r.cache[attr.Code] = attr
	}

	r.loaded = true
	return nil
}

// Get returns the attribute metadata for the given attribute code.
func (r *AttributeRepository) Get(code string) *AttributeMetadata {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.cache[code]
}

// GetID returns the attribute ID for the given attribute code.
func (r *AttributeRepository) GetID(code string) int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if attr, ok := r.cache[code]; ok {
		return attr.AttributeID
	}
	return 0
}

// TableForType returns the EAV value table name for a given backend_type.
func TableForType(backendType string) string {
	switch backendType {
	case "varchar":
		return "catalog_product_entity_varchar"
	case "int":
		return "catalog_product_entity_int"
	case "text":
		return "catalog_product_entity_text"
	case "decimal":
		return "catalog_product_entity_decimal"
	case "datetime":
		return "catalog_product_entity_datetime"
	default:
		return ""
	}
}
