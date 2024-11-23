package configs

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost string
	Port       string

	DBUser                          string
	DBPassword                      string
	DBAddress                       string
	DBName                          string
	JWTSecret                       string
	JWTExpirationInSeconds          int64
	RefreshTokenExpirationInSeconds int64
	SslMode                         string
	CACertPath                      string
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()

	return Config{
		PublicHost: getEnv("PUBLIC_HOST", "http://localhost"),
		Port:       getEnv("PORT", "8080"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASS", ""),
		DBAddress:  getEnv("DB_HOST", "127.0.0.1"),
		DBName:     getEnv("DB_NAME", "educations"),
		JWTSecret: getEnv("JWT_SECRET",
			"not-secret-secret-anymore"),
		JWTExpirationInSeconds:          getEnvAsInt("JWT_EXP", 600),
		RefreshTokenExpirationInSeconds: getEnvAsInt("REFRESH_TOKEN_EXP", 86400),
		SslMode:                         getEnv("SSL_MODE", "disable"),
		CACertPath:                      getEnv("CA_CERT_PATH", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}
		return i
	}

	return fallback
}
