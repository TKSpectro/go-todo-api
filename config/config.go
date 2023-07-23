package config

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

var (
	isEnvLoaded = false
)

var (
	PORT = getEnv("PORT", "3000")

	DB_USER          = getEnv("DB_USER", "root")
	DB_ROOT_PASSWORD = getEnv("DB_ROOT_PASSWORD", "root")
	DB_NAME          = getEnv("DB_NAME", "go_api")
	DB_PORT          = getEnv("DB_LOCAL_PORT", "3306")

	// getEnv returns a string that we have to time.ParseDuration
	JWT_TOKEN_EXP   = getEnvTimeDurationParse("JWT_TOKEN_EXP", "1h")
	JWT_REFRESH_EXP = getEnvTimeDurationParse("JWT_REFRESH_EXP", "10m")

	ALLOWED_IPS = getEnvList("ALLOWED_IPS")
)

func getEnv(name string, fallback string) string {
	if !isEnvLoaded {
		if err := godotenv.Load(".env"); err != nil {
			fmt.Println("Error loading .env file")
		} else {
			isEnvLoaded = true
		}
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

func getEnvList(name string) []string {
	value := getEnv(name, "")

	if value == "" {
		return []string{}
	}

	return strings.Split(value, ",")
}

func getEnvTimeDurationParse(name string, fallback string) time.Duration {
	value := getEnv(name, fallback)

	if value == "" {
		return 0
	}

	parsed, err := time.ParseDuration(value)
	if err != nil {
		log.Println("Error parsing duration for environment variable: "+name, value, err)
		return 0
	}

	return parsed
}
