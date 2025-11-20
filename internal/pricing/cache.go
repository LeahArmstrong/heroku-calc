package pricing

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	cacheTTL      = 24 * time.Hour
	cacheDir      = ".heroku-calc"
	cacheFileName = "pricing_cache.json"
)

type cacheEntry struct {
	Data      Data      `json:"data"`
	Timestamp time.Time `json:"timestamp"`
}

// GetCachePath returns the path to the pricing cache file
func GetCachePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	cachePath := filepath.Join(home, cacheDir, cacheFileName)
	return cachePath, nil
}

// LoadFromCache loads pricing data from cache if available and not expired
func LoadFromCache() (*Data, error) {
	cachePath, err := GetCachePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(cachePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("cache not found")
		}
		return nil, fmt.Errorf("failed to read cache: %w", err)
	}

	var entry cacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, fmt.Errorf("failed to parse cache: %w", err)
	}

	// Check if cache is expired
	if time.Since(entry.Timestamp) > cacheTTL {
		return nil, fmt.Errorf("cache expired")
	}

	return &entry.Data, nil
}

// SaveToCache saves pricing data to cache
func SaveToCache(data *Data) error {
	cachePath, err := GetCachePath()
	if err != nil {
		return err
	}

	// Ensure cache directory exists
	cacheDir := filepath.Dir(cachePath)
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	entry := cacheEntry{
		Data:      *data,
		Timestamp: time.Now(),
	}

	jsonData, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal cache data: %w", err)
	}

	if err := os.WriteFile(cachePath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write cache: %w", err)
	}

	return nil
}

// IsCacheExpired checks if the cache is expired without loading it
func IsCacheExpired() bool {
	cachePath, err := GetCachePath()
	if err != nil {
		return true
	}

	info, err := os.Stat(cachePath)
	if err != nil {
		return true
	}

	return time.Since(info.ModTime()) > cacheTTL
}
