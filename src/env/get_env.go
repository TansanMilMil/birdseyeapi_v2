package env

import (
	"os"
	"strconv"
)

func GetEnv(key string, default_val string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return default_val
	}
	return val
}

func GetEnvInt(key string, default_val int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		return default_val
	}
	intVal, err := strconv.Atoi(val)
	if err != nil {
		return default_val
	}
	return intVal
}
