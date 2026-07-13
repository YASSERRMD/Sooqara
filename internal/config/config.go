package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all application settings loaded from environment variables.
type Config struct {
	AgnesAPIKey   string
	AgnesBaseURL  string
	AgnesPollURL  string
	RPM           int
	DBPath        string
	StorageDir    string
	Addr          string
	Workers       int
	LogLevel      string
}

// Load reads configuration from environment variables.
func Load() Config {
	return Config{
		AgnesAPIKey:  os.Getenv("AGNES_API_KEY"),
		AgnesBaseURL: getEnvOrDefault("AGNES_BASE_URL", "https://apihub.agnes-ai.com/v1"),
		AgnesPollURL: getEnvOrDefault("AGNES_POLL_URL", "https://apihub.agnes-ai.com/agnesapi"),
		RPM:          getEnvOrDefaultInt("SOOQARA_RPM", 18),
		DBPath:       getEnvOrDefault("SOOQARA_DB", "./sooqara.db"),
		StorageDir:   getEnvOrDefault("SOOQARA_STORAGE", "./storage"),
		Addr:         getEnvOrDefault("SOOQARA_ADDR", ":8080"),
		Workers:      getEnvOrDefaultInt("SOOQARA_WORKERS", 3),
		LogLevel:     getEnvOrDefault("SOOQARA_LOG_LEVEL", "info"),
	}
}

// getEnvOrDefault returns the environment variable named key or defaultValue if empty.
func getEnvOrDefault(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

// getEnvOrDefaultInt parses the environment variable as integer or returns defaultValue.
func getEnvOrDefaultInt(key string, defaultValue int) int {
	v := os.Getenv(key)
	if v == "" {
		return defaultValue
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return defaultValue
	}
	return n
}

// Validate checks all required fields and returns a joined error for every problem.
func (c Config) Validate() error {
	var errs []error
	if c.AgnesAPIKey == "" {
		errs = append(errs, fmt.Errorf("AGNES_API_KEY is required"))
	}
	if c.RPM <= 0 {
		errs = append(errs, fmt.Errorf("SOOQARA_RPM must be positive, got %d", c.RPM))
	}
	if c.Workers <= 0 {
		errs = append(errs, fmt.Errorf("SOOQARA_WORKERS must be positive, got %d", c.Workers))
	}
	if c.LogLevel == "" {
		errs = append(errs, fmt.Errorf("SOOQARA_LOG_LEVEL is required"))
	}
	if _, err := time.ParseDuration("0s"); err != nil {
		errs = append(errs, fmt.Errorf("invalid log level %q", c.LogLevel))
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}
