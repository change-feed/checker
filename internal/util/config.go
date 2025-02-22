package util

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
)

var (
	MaxSnapshots = getEnvInt("MAX_SNAPSHOTS", 3)
)

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		slog.Warn(fmt.Sprintf("Invalid integer value for %s: %s. Using default: %d", key, valueStr, defaultValue))
		return defaultValue
	}
	return value
}
