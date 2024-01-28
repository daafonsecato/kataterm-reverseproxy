package config

import (
	"os"
	"strconv"
)

type Config struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	// Add other configuration fields as needed
}

func LoadConfig() (*Config, error) {
	return &Config{
		DBHost:     os.Getenv("postgres"),
		DBPort:     5432,
		DBUser:     os.Getenv("your_username"),
		DBPassword: os.Getenv("your_password"),
		DBName:     os.Getenv("your_database_name"),
		// Load other configurations as needed
	}, nil
}
