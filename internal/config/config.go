package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env         string
	Storage     string
	Address     string
	HelpAPI     string
	Timeout     time.Duration
	IdleTimeout time.Duration
}

func InitConfig() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	// check if file exist
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	err := godotenv.Load(configPath)
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	timeout, err := time.ParseDuration(os.Getenv("TIMEOUT"))
	if err != nil {
		log.Fatalf("Error parsing TIMEOUT: %v", err)
	}

	idleTimeout, err := time.ParseDuration(os.Getenv("IDLE_TIMEOUT"))
	if err != nil {
		log.Fatalf("Error parsing IDLE_TIMEOUT: %v", err)
	}

	return &Config{
		Env:         os.Getenv("ENV"),
		Storage:     os.Getenv("STORAGE"),
		Address:     os.Getenv("ADDRESS"),
		HelpAPI:     os.Getenv("HELP_API"),
		Timeout:     timeout,
		IdleTimeout: idleTimeout,
	}
}
