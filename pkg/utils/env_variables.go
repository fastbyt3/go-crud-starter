package utils

import "os"

func GetEnvVarOrDefault(variableName string, defaultValue string) string {
	value, exists := os.LookupEnv(variableName)
	if exists {
		return value
	}
	return defaultValue
}
