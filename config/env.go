package config

import (
	"fmt"
	"os"
)

var (
	PORT = getEnv("PORT", "3000")

	DB = getEnv("DB", "test.db")

	JWT_TOKEN_SECRET = getEnv("JWT_TOKEN_SECRET", "secret")
	JWT_TOKEN_EXP    = getEnv("JWT_TOKEN_EXP", "10h")
)

func getEnv(name string, fallback string) string {
	if value, exists := os.LookupEnv(name); exists {
		return value
	}

	if fallback != "" {
		return fallback
	}

	panic(fmt.Sprintf("Environment variable not found :: %v", name))
}
