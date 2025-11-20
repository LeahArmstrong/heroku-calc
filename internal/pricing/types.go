package pricing

// Data represents the complete pricing information
type Data struct {
	Version  string                  `json:"version"`
	Dynos    map[string]DynoPrice    `json:"dynos"`
	Postgres map[string]PostgresPrice `json:"postgres"`
	Redis    map[string]RedisPrice   `json:"redis"`
}

// DynoPrice represents pricing for a dyno type
type DynoPrice struct {
	Name         string  `json:"name"`
	MemoryMB     int     `json:"memory_mb"`
	PriceMonthly float64 `json:"price_monthly"`
	PriceHourly  float64 `json:"price_hourly"`
	Compute      string  `json:"compute"`
}

// PostgresPrice represents pricing for a Postgres plan
type PostgresPrice struct {
	Name            string  `json:"name"`
	MaxConnections  int     `json:"max_connections"`
	PriceMonthly    float64 `json:"price_monthly"`
	StorageGB       int     `json:"storage_gb"`
	RAMMB           int     `json:"ram_mb"`
	HighAvailability bool   `json:"high_availability,omitempty"`
}

// RedisPrice represents pricing for a Redis plan
type RedisPrice struct {
	Name            string  `json:"name"`
	MaxConnections  int     `json:"max_connections"`
	PriceMonthly    float64 `json:"price_monthly"`
	MaxMemoryMB     int     `json:"max_memory_mb"`
	EvictionPolicy  string  `json:"eviction_policy,omitempty"`
	HighAvailability bool   `json:"high_availability,omitempty"`
}
