package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"
)

// StoreConfig holds cached store configuration values.
type StoreConfig struct {
	WebsiteID            int
	BaseCurrency         string
	ProductURLSuffix     string
	CategoryURLSuffix    string
	MediaBaseURL         string  // e.g. "https://example.com/media/catalog/product"
	StockThresholdQty    float64 // cataloginventory/options/stock_threshold_qty; 0 = disabled
	ProductCanonicalTag  bool    // catalog/seo/product_canonical_tag; false = disabled
}

// StoreConfigRepository caches per-store configuration.
type StoreConfigRepository struct {
	db    *sql.DB
	cache map[int]*StoreConfig
	mu    sync.RWMutex
}

func NewStoreConfigRepository(db *sql.DB) *StoreConfigRepository {
	return &StoreConfigRepository{
		db:    db,
		cache: make(map[int]*StoreConfig),
	}
}

// Get returns the cached store config, loading it on first access.
func (r *StoreConfigRepository) Get(ctx context.Context, storeID int) *StoreConfig {
	r.mu.RLock()
	if cfg, ok := r.cache[storeID]; ok {
		r.mu.RUnlock()
		return cfg
	}
	r.mu.RUnlock()

	// Load website_id for this store
	var websiteID int
	if err := r.db.QueryRowContext(ctx, "SELECT website_id FROM store WHERE store_id = ?", storeID).Scan(&websiteID); err != nil {
		websiteID = 1
	}

	// Batch-load all config values in a single query
	configPaths := []string{
		"currency/options/base",
		"catalog/seo/product_url_suffix",
		"catalog/seo/category_url_suffix",
		"web/secure/base_media_url",
		"web/secure/base_url",
		"cataloginventory/options/stock_threshold_qty",
		"catalog/seo/product_canonical_tag",
	}

	placeholders := make([]string, len(configPaths))
	args := make([]interface{}, len(configPaths)+1)
	for i, p := range configPaths {
		placeholders[i] = "?"
		args[i] = p
	}
	args[len(configPaths)] = storeID

	query := `SELECT path, value, scope, scope_id
		FROM core_config_data
		WHERE path IN (` + strings.Join(placeholders, ",") + `)
		AND (scope = 'default' OR (scope = 'stores' AND scope_id = ?))
		ORDER BY scope_id ASC`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		// Fall back to defaults
		cfg := &StoreConfig{
			WebsiteID:        websiteID,
			BaseCurrency:     "USD",
			ProductURLSuffix: ".html",
			CategoryURLSuffix: ".html",
			MediaBaseURL:     "http://localhost/media/catalog/product",
		}
		r.mu.Lock()
		r.cache[storeID] = cfg
		r.mu.Unlock()
		return cfg
	}
	defer rows.Close()

	// Collect values; store-scoped overrides default because of ORDER BY scope_id ASC
	configValues := make(map[string]string)
	for rows.Next() {
		var path, value, scope string
		var scopeID int
		if err := rows.Scan(&path, &value, &scope, &scopeID); err != nil {
			continue
		}
		// Store-scoped values override default values (they come later due to ORDER BY)
		configValues[path] = value
	}

	// Build media base URL
	mediaBaseURL := ""
	if v, ok := configValues["web/secure/base_media_url"]; ok && v != "" {
		mediaBaseURL = strings.TrimRight(v, "/")
	} else {
		baseURL := "http://localhost/"
		if v, ok := configValues["web/secure/base_url"]; ok && v != "" {
			baseURL = v
		}
		mediaBaseURL = strings.TrimRight(baseURL, "/") + "/media"
	}
	mediaBaseURL += "/catalog/product"

	// Parse float config values
	stockThreshold := 0.0
	if v, ok := configValues["cataloginventory/options/stock_threshold_qty"]; ok {
		fmt.Sscanf(v, "%f", &stockThreshold)
	}
	canonicalTag := false
	if v, ok := configValues["catalog/seo/product_canonical_tag"]; ok {
		canonicalTag = v == "1"
	}

	cfg := &StoreConfig{
		WebsiteID:           websiteID,
		BaseCurrency:        getConfigOrDefault(configValues, "currency/options/base", "USD"),
		ProductURLSuffix:    getConfigOrDefault(configValues, "catalog/seo/product_url_suffix", ".html"),
		CategoryURLSuffix:   getConfigOrDefault(configValues, "catalog/seo/category_url_suffix", ".html"),
		MediaBaseURL:        mediaBaseURL,
		StockThresholdQty:   stockThreshold,
		ProductCanonicalTag: canonicalTag,
	}

	r.mu.Lock()
	r.cache[storeID] = cfg
	r.mu.Unlock()

	return cfg
}

func getConfigOrDefault(m map[string]string, key, defaultVal string) string {
	if v, ok := m[key]; ok && v != "" {
		return v
	}
	return defaultVal
}

// getMediaBaseURL builds the product media base URL from core_config_data.
// Checks web/secure/base_media_url first; if NULL, uses web/secure/base_url + "media/".
// Always appends "catalog/product".
func (r *StoreConfigRepository) getMediaBaseURL(ctx context.Context, storeID int) string {
	var mediaURL sql.NullString
	// Try explicit media URL first (store-scoped, then default)
	_ = r.db.QueryRowContext(ctx,
		"SELECT value FROM core_config_data WHERE path = 'web/secure/base_media_url' AND scope_id = ? AND scope = 'stores'", storeID,
	).Scan(&mediaURL)
	if !mediaURL.Valid {
		_ = r.db.QueryRowContext(ctx,
			"SELECT value FROM core_config_data WHERE path = 'web/secure/base_media_url' AND scope = 'default'",
		).Scan(&mediaURL)
	}

	base := ""
	if mediaURL.Valid && mediaURL.String != "" {
		base = strings.TrimRight(mediaURL.String, "/")
	} else {
		// Fall back to base_url + /media
		var baseURL string
		err := r.db.QueryRowContext(ctx,
			"SELECT value FROM core_config_data WHERE path = 'web/secure/base_url' AND scope_id = ? AND scope = 'stores'", storeID,
		).Scan(&baseURL)
		if err != nil {
			_ = r.db.QueryRowContext(ctx,
				"SELECT value FROM core_config_data WHERE path = 'web/secure/base_url' AND scope = 'default'",
			).Scan(&baseURL)
		}
		if baseURL == "" {
			baseURL = "http://localhost/"
		}
		base = strings.TrimRight(baseURL, "/") + "/media"
	}

	return base + "/catalog/product"
}

// getFloatConfig reads a float config value from core_config_data, returning 0 if not found.
func (r *StoreConfigRepository) getFloatConfig(ctx context.Context, path string, storeID int) float64 {
	var val float64
	// Try store-scoped, then default
	err := r.db.QueryRowContext(ctx, "SELECT value FROM core_config_data WHERE path = ? AND scope_id = ? AND scope = 'stores'", path, storeID).Scan(&val)
	if err == nil {
		return val
	}
	err = r.db.QueryRowContext(ctx, "SELECT value FROM core_config_data WHERE path = ? AND scope = 'default' AND scope_id = 0", path).Scan(&val)
	if err == nil {
		return val
	}
	return 0
}
