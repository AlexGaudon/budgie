package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost         string
	DBUserName     string
	DBUserPassword string
	DBName         string
	DBPort         string
	ServerPort     string

	AccessTokenExpiresIn  time.Duration
	RefreshTokenExpiresIn time.Duration
	AccessTokenMaxAge     int
	RefreshTokenMaxAge    int
}

var cfg Config

func GetConfig() *Config {
	if (cfg == Config{}) {
		LoadConfig()
	}

	return &cfg
}

func LoadConfig() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Failed to load .env file: %s", err)
	}

	c := &Config{}

	c.DBHost = os.Getenv("DB_HOST")
	c.ServerPort = os.Getenv("SERVER_PORT")
	c.DBUserName = os.Getenv("DB_USER")
	c.DBUserPassword = os.Getenv("DB_PASS")
	c.DBName = os.Getenv("DB_NAME")
	c.DBPort = os.Getenv("DB_PORT")

	c.AccessTokenExpiresIn = time.Minute * 15
	c.RefreshTokenExpiresIn = time.Hour * 1

	c.AccessTokenMaxAge = int(time.Minute * 15)
	c.RefreshTokenMaxAge = int(time.Hour * 1)

	cfg = *c
}
