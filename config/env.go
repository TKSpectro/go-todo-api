package config

import (
	"fmt"
	"log"
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
	// getEnv returns a string that we have to time.ParseDuration
	JWT_TOKEN_EXP = getEnv("JWT_TOKEN_EXP", "15m")
)

func getEnv(name string, fallback string) string {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Print("Error loading .env file")
	}

	if value, exists := os.LookupEnv(name); exists {
		return value
	}

	if fallback != "" {
		return fallback
	}

	log.Println("Environment variable not found: " + name)

	return ""
}
