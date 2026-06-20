package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                  string
	DBURL                 string
	JWTSecret             string
	PasswordEncryptionKey string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found")
	}

	return &Config{
		Port:                  os.Getenv("HTTP_PORT"),
		DBURL:                 os.Getenv("DATABASE_URL"),
		JWTSecret:             os.Getenv("JWT_SECRET"),
		PasswordEncryptionKey: os.Getenv("PASSWORD_ENCRYPTION_KEY"),
	}
}
