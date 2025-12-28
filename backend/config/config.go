package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBURL     string
	Port      string
	JWTSecret string
	RedisURL  string
	KafkaURL  string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, reading from environment variables")
	}

	return &Config{
		DBURL:     getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/qcsystem?sslmode=disable"),
		Port:      getEnv("PORT", "8081"),
		JWTSecret: getEnv("JWT_SECRET", "super-secret-key"),
		RedisURL:  getEnv("REDIS_URL", "localhost:6379"),
		KafkaURL:  getEnv("KAFKA_URL", "localhost:9092"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
