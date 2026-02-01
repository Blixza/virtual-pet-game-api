package db_config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	Name     string
	Port     string
	Host     string
	User     string
	Password string
	SSLMode  string
}

func Load(path string) DBConfig {
	err := godotenv.Load(path)
	if err != nil {
		fmt.Printf("failed to load config by path %s. Using default: %v\n", path, err)
		return DBConfig{
			Name: "pet_game",
			Port: "5432",
			Host: "localhost",
			User: "postgres",
			Password: "5432",
			SSLMode: "disable",
		}
	}

	cfg := DBConfig{}
	
	cfg.Name = os.Getenv("DB_NAME")
	cfg.Host = os.Getenv("DB_HOST")
	cfg.User = os.Getenv("DB_USER")
	cfg.Password = os.Getenv("DB_PASSWORD")
	cfg.SSLMode = os.Getenv("DB_SSLMODE")

	return cfg
}
