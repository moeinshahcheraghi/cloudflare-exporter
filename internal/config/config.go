package config

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	APIToken       string
	ZoneID         string
	Port           string
	ScrapeInterval time.Duration
}

func LoadFromEnv() (*Config, error) {
	apiToken := os.Getenv("CLOUDFLARE_API_TOKEN")
	zoneID := os.Getenv("CLOUDFLARE_ZONE_ID")

	if apiToken == "" {
		return nil, fmt.Errorf("CLOUDFLARE_API_TOKEN environment variable is required")
	}

	if zoneID == "" {
		return nil, fmt.Errorf("CLOUDFLARE_ZONE_ID environment variable is required")
	}

	return &Config{
		APIToken:       apiToken,
		ZoneID:         zoneID,
		Port:           getEnvOrDefault("EXPORTER_PORT", "9199"),
		ScrapeInterval: 60 * time.Second,
	}, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}