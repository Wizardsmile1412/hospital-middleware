package config

import (
	"os"
	"strconv"
)

type Config struct {
	AppEnv string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	JWTSecret        string
	JWTAccessExpiry  int
	JWTRefreshExpiry int

	HospitalAAPIURL string
	ServerPort      string
}

func Load() *Config {
	return &Config{
		AppEnv: getEnv("APP_ENV", "development"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "hospital_middleware"),

		JWTSecret:        getEnv("JWT_SECRET", "supersecretkey"),
		JWTAccessExpiry:  getEnvInt("JWT_ACCESS_EXPIRY", 900),
		JWTRefreshExpiry: getEnvInt("JWT_REFRESH_EXPIRY", 604800),

		HospitalAAPIURL: getEnv("HOSPITAL_A_API_URL", "https://hospital-a.api.co.th"),
		ServerPort:      getEnv("SERVER_PORT", "8080"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}
