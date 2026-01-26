package utils

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func GetString(key string) (string, error) {
	value, exists := os.LookupEnv(key)
	if !exists {
		return "", nil
	}
	return value, nil
}

func GetInt(key string) (int, error) {
	value := os.Getenv(key)
	if value == "" {
		return 0, nil
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("environment variable %s must be a valid integer", key)
	}
	return intValue, nil
}

func GetBool(key string) (bool, error) {
	value := os.Getenv(key)
	if value == "" {
		return false, nil
	}
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return false, fmt.Errorf("environment variable %s must be a valid boolean", key)
	}

	return boolValue, nil
}

func GetDuration(key string) (time.Duration, error) {
	return time.ParseDuration(os.Getenv(key))
}
