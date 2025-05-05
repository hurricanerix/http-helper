/*
Package config provides functions for reading config values
from the environment.  These functions are designed to not
fail.  If an env variable is not set, or parsing fails, the
provided default value is used.
*/
package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

func IntEnv(key string, defaultValue int) int {
	value := StringEnv(key, strconv.FormatInt(int64(defaultValue), 10))
	i, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return defaultValue
	}
	return int(i)
}

func StringSliceEnv(key string, defaultValue string) []string {
	value := StringEnv(key, defaultValue)
	return strings.Split(value, ",")
}

func BoolEnv(key string, defaultValue bool) bool {
	value := StringEnv(key, strconv.FormatBool(defaultValue))
	b, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return b
}

func DurationEnv(key string, defaultValue time.Duration) time.Duration {
	value := StringEnv(key, defaultValue.String())
	d, err := time.ParseDuration(value)
	if err != nil {
		return defaultValue
	}
	return d
}

func StringEnv(key string, defaultVal string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultVal
	}
	return val
}
