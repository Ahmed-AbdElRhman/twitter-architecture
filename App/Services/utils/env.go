package utils

import (
	"fmt"
	"os"
	"strconv"
)

const (
	hostKey      = "DB_HOST"
	portKey      = "DB_PORT"
	userKey      = "DB_USER"
	passwordKey  = "DB_PASSWORD"
	dbnameKey    = "DB_NAME"
	jwtSecretKey = "JWT_SECRET"
)

var (
	Host       = getEnvVar(hostKey, "localhost")
	Port       = getEnvVarInt(portKey, 5432)
	User       = getEnvVar(userKey, "postgres")
	Password   = getEnvVar(passwordKey, "postgres")
	DBName     = getEnvVar(dbnameKey, "auth_app")
	JWT_SECRET = getEnvVar(jwtSecretKey, "7c3d9f82a1e5b409f6e7c8d2a4b0c9d3e5f7a1b2c8d9e0f3a4b5c6d7e8f9a0b")
)

func getEnvVar(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

func getEnvVarInt(key string, defaultValue int) int {
	strValue := getEnvVar(key, "")
	if strValue == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(strValue)
	if err != nil {
		fmt.Printf("Invalid integer value for %s, using default %d\n", key, defaultValue)
		return defaultValue
	}

	return value
}
