package utils

import (
	"os"
	"strconv"
)

// GetString retrieves the value of the environment variable named by the key.
// If the variable is not present, it returns the provided fallback value.
func GetString(key string, fallback string) string {
	val, ok := os.LookupEnv(key)

	if !ok {
		return fallback
	}

	return val
}

// GetInt retrieves the value of the environment variable named by the key and converts it to an int.
// If the variable is not present or cannot be converted, it returns the provided fallback value.
func GetInt(key string, fallback int) int {
	val, ok := os.LookupEnv(key)

	if !ok {
		return fallback
	}

	intVal, err := strconv.Atoi(val)

	if err != nil {
		return fallback
	}

	return intVal
}

// GetBool retrieves the value of the environment variable named by the key and converts it to a bool.
// This function is case-insensitive
// If the variable is not present or cannot be converted, it returns the provided fallback value.
func GetBool(key string, fallback bool) bool {
	val, ok := os.LookupEnv(key)

	if !ok {
		return fallback
	}

	boolVal, err := strconv.ParseBool(val)

	if err != nil {
		return fallback
	}

	return boolVal
}
