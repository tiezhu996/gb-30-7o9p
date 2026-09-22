package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config aggregates runtime configuration from environment variables.
type Config struct {
	ServerPort    string
	CORSOrigins   string
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	JWTSecret     string
	JWTExpire     time.Duration
	RateLimitReq  int
	RateLimitWin  time.Duration
	RedisAddr     string
	RedisPassword string
	RedisDB       int
	MinIOEndpoint string
	MinIOUser     string
	MinIOPassword string
	MinIOBucket   string
	MinIOUseSSL   bool
}

// Load reads configuration from environment variables with defaults.
func Load() *Config {
	return &Config{
		ServerPort:    getEnv("SERVER_PORT", "8080"),
		CORSOrigins:   getEnv("APP_CORS_ORIGINS", "http://localhost:8101"),
		DBHost:        getEnv("DB_HOST", "db"),
		DBPort:        getEnv("DB_PORT", "5601"),
		DBUser:        getEnv("DB_USER", "gbadopt_user"),
		DBPassword:    getEnv("DB_PASSWORD", "gbadopt_pwd"),
		DBName:        getEnv("DB_NAME", "gbadopt_db"),
		JWTSecret:     getEnv("JWT_SECRET", "change_me_to_a_long_random_string"),
		JWTExpire:     time.Duration(getEnvInt("JWT_EXPIRE_HOURS", 72)) * time.Hour,
		RateLimitReq:  getEnvInt("RATE_LIMIT_REQUESTS", 60),
		RateLimitWin:  time.Duration(getEnvInt("RATE_LIMIT_WINDOW_SECONDS", 60)) * time.Second,
		RedisAddr:     getEnv("REDIS_ADDR", "redis:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnvInt("REDIS_DB", 0),
		MinIOEndpoint: getEnv("MINIO_ENDPOINT", "minio:9000"),
		MinIOUser:     getEnv("MINIO_ROOT_USER", "minioadmin"),
		MinIOPassword: getEnv("MINIO_ROOT_PASSWORD", "minioadmin"),
		MinIOBucket:   getEnv("MINIO_BUCKET", "gbadopt"),
		MinIOUseSSL:   getEnv("MINIO_USE_SSL", "false") == "true",
	}
}

// DSN returns the PostgreSQL connection string.
func (c *Config) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName)
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
