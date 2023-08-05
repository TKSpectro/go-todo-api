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
	ROOT_PATH = getEnv("GTA_ROOT_PATH", "")
	PORT      = getEnv("PORT", "3000")

	DB_USER          = ifIsTestElse(getEnv("TEST_DB_USER", "root"), getEnv("DB_USER", "root"))
	DB_ROOT_PASSWORD = ifIsTestElse(getEnv("TEST_DB_ROOT_PASSWORD", "root"), getEnv("DB_ROOT_PASSWORD", "root"))
	DB_NAME          = ifIsTestElse(getEnv("TEST_DB_NAME", "go_api_test"), getEnv("DB_NAME", "go_api"))
	DB_PORT          = ifIsTestElse(getEnv("TEST_DB_LOCAL_PORT", "3306"), getEnv("DB_LOCAL_PORT", "3306"))

	// getEnv returns a string that we have to time.ParseDuration
	JWT_TOKEN_EXP   = getEnvTimeDurationParse("JWT_TOKEN_EXP", "1h")
	JWT_REFRESH_EXP = getEnvTimeDurationParse("JWT_REFRESH_EXP", "10m")

	ALLOWED_IPS = getEnvList("ALLOWED_IPS")

	// Testing environment variables
	IS_TEST = getEnvBool("IS_TEST", "false")
)

func ifThenElse(condition bool, a string, b string) string {
	if condition {
		return a
	}
	return b
}

func ifIsTestElse(a string, b string) string {
	return ifThenElse(getEnvBool("IS_TEST", "false"), a, b)
}

func getEnv(name string, fallback string) string {
	if !isEnvLoaded {
		path, _ := os.LookupEnv("GTA_ROOT_PATH")
		path2, _ := os.Getwd()
		println(path2)

		envPath := ".env"
		if path != "" {
			envPath = path + "/.env"
		}

		if err := godotenv.Load(envPath); err != nil {
			fmt.Println("Error loading .env file")
			fmt.Println(err)
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

func getEnvBool(name string, fallback string) bool {
	value := getEnv(name, fallback)

	if value == "" {
		return false
	}

	return value == "true" || value == "1"
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
