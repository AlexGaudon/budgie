package config

import "time"

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
	c := &Config{}

	c.AccessTokenExpiresIn = time.Minute * 15
	c.RefreshTokenExpiresIn = time.Hour * 1

	c.AccessTokenMaxAge = int(time.Minute * 15)
	c.RefreshTokenMaxAge = int(time.Hour * 1)

	cfg = *c
}
