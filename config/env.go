package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var (
	PORT = getEnv("PORT", "3000")

	DB_USER          = getEnv("DB_USER", "root")
	DB_ROOT_PASSWORD = getEnv("DB_ROOT_PASSWORD", "root")
	DB_NAME          = getEnv("DB_NAME", "go_api")
	DB_PORT          = getEnv("DB_LOCAL_PORT", "3306")

	JWT_TOKEN_SECRET = getEnv("JWT_TOKEN_SECRET", "secret")
	JWT_TOKEN_EXP    = getEnv("JWT_TOKEN_EXP", "10h")
)

func getEnv(name string, fallback string) string {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Print("Error loading .env file")
	}

	if value, exists := os.LookupEnv(name); exists {
		return value
	}

	if fallback != "" {
		return fallback
	}

	panic(fmt.Sprintf("Environment variable not found :: %v", name))
}
