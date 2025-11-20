package pricing

import (
	"fmt"
)

// Fetch attempts to fetch fresh pricing data from Heroku
// Currently returns bundled data as Heroku doesn't expose a public pricing API
// In the future, this could scrape the Heroku pricing page or use other methods
func Fetch() (*Data, error) {
	// For now, we'll use bundled data as there's no official pricing API
	// In production, this could:
	// 1. Scrape https://www.heroku.com/pricing
	// 2. Use Elements Marketplace API if available
	// 3. Maintain an updated pricing service

	// For this implementation, we'll just return bundled data
	return LoadBundled()
}

// Get returns pricing data using the hybrid approach:
// 1. Try to load from cache
// 2. If cache miss or expired, try to fetch fresh data
// 3. Fall back to bundled data if fetch fails
func Get() (*Data, error) {
	// Try cache first
	if data, err := LoadFromCache(); err == nil {
		return data, nil
	}

	// Try to fetch fresh data
	data, err := Fetch()
	if err != nil {
		// Fall back to bundled data
		return LoadBundled()
	}

	// Save to cache for next time
	_ = SaveToCache(data) // Ignore cache save errors

	return data, nil
}

// GetDynoPrice looks up pricing for a specific dyno type
func (d *Data) GetDynoPrice(dynoType string) (*DynoPrice, error) {
	// Normalize dyno type to lowercase
	dynoType = normalizeKey(dynoType)

	if price, ok := d.Dynos[dynoType]; ok {
		return &price, nil
	}
	return nil, fmt.Errorf("dyno type not found: %s", dynoType)
}

// GetPostgresPrice looks up pricing for a specific Postgres plan
func (d *Data) GetPostgresPrice(plan string) (*PostgresPrice, error) {
	// Normalize plan name to lowercase
	plan = normalizeKey(plan)

	if price, ok := d.Postgres[plan]; ok {
		return &price, nil
	}
	return nil, fmt.Errorf("postgres plan not found: %s", plan)
}

// GetRedisPrice looks up pricing for a specific Redis plan
func (d *Data) GetRedisPrice(plan string) (*RedisPrice, error) {
	// Normalize plan name to lowercase
	plan = normalizeKey(plan)

	if price, ok := d.Redis[plan]; ok {
		return &price, nil
	}
	return nil, fmt.Errorf("redis plan not found: %s", plan)
}

// normalizeKey converts a key to lowercase for consistent lookups
func normalizeKey(key string) string {
	// Convert to lowercase and replace underscores with hyphens
	normalized := ""
	for _, ch := range key {
		if ch == '_' {
			normalized += "-"
		} else if ch >= 'A' && ch <= 'Z' {
			normalized += string(ch + 32)
		} else {
			normalized += string(ch)
		}
	}
	return normalized
}
