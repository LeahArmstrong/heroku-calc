package config

import "time"

// Config represents the .heroku-calc.yml configuration file
type Config struct {
	// AppName is the Heroku application name
	AppName string `yaml:"app_name"`

	// GitRemote is the git remote name (e.g., "heroku", "production")
	GitRemote string `yaml:"git_remote,omitempty"`

	// SafeEnvVars are the environment variables safe to query
	SafeEnvVars []string `yaml:"safe_env_vars"`

	// ExcludedEnvVars are vars explicitly excluded by the user
	ExcludedEnvVars []string `yaml:"excluded_env_vars,omitempty"`

	// LastUpdated timestamp
	LastUpdated time.Time `yaml:"last_updated"`

	// ProjectPath is the path to the Rails project (if using --project flag)
	ProjectPath string `yaml:"project_path,omitempty"`
}

// HerokuEnvVar represents a single environment variable
type HerokuEnvVar struct {
	Name  string
	Value string
}

// DynoFormation represents the dyno configuration
type DynoFormation struct {
	Type     string // "web", "worker", etc.
	Quantity int
	Size     string // "Standard-1X", "Performance-M", etc.
}

// Addon represents a Heroku addon
type Addon struct {
	Name    string
	Plan    string
	Price   string
	AddedAt time.Time
}

// AnalysisResult represents the output of configuration analysis
type AnalysisResult struct {
	DatabaseAnalysis *DatabaseAnalysis
	RedisAnalysis    *RedisAnalysis
	WebTierAnalysis  *WebTierAnalysis
	Recommendations  []Recommendation
}

// DatabaseAnalysis contains database connection analysis
type DatabaseAnalysis struct {
	DatabaseURL      string
	PostgresPlan     string
	MaxConnections   int
	CurrentUsage     int
	WebDynos         int
	WorkersPerDyno   int
	ThreadsPerWorker int
	SidekiqDynos     int
	SidekiqThreads   int
	TotalRequired    int
	BufferPercent    float64
	Status           AnalysisStatus
	Issues           []string
}

// RedisAnalysis contains Redis/cache analysis
type RedisAnalysis struct {
	RedisURL           string
	RedisPlan          string
	MaxConnections     int
	SidekiqConcurrency int
	RedisPoolSize      int
	EstimatedUsage     int
	Status             AnalysisStatus
	Issues             []string
}

// WebTierAnalysis contains web tier concurrency analysis
type WebTierAnalysis struct {
	DynoType         string
	DynoMemoryMB     int
	WebConcurrency   int
	RailsMaxThreads  int
	TotalThreads     int
	MemoryPerThread  int
	Status           AnalysisStatus
	Issues           []string
}

// Recommendation represents a suggested configuration change
type Recommendation struct {
	Category    string // "database", "redis", "web", "cost"
	Severity    RecommendationSeverity
	Title       string
	Description string
	Current     string
	Suggested   string
	EnvVarName  string // If applicable
	Impact      string // Cost or performance impact
	AutoApply   bool   // Whether this can be auto-applied
}

// AnalysisStatus represents the health status of a component
type AnalysisStatus string

const (
	StatusCritical AnalysisStatus = "critical"
	StatusWarning  AnalysisStatus = "warning"
	StatusOptimal  AnalysisStatus = "optimal"
	StatusUnknown  AnalysisStatus = "unknown"
)

// RecommendationSeverity indicates priority level
type RecommendationSeverity string

const (
	SeverityCritical RecommendationSeverity = "critical"
	SeverityHigh     RecommendationSeverity = "high"
	SeverityMedium   RecommendationSeverity = "medium"
	SeverityLow      RecommendationSeverity = "low"
	SeverityInfo     RecommendationSeverity = "info"
)
