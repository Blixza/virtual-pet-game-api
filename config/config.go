package config

import (
	"fmt"
	"virtual_pet_game/pkg/logger"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type Config struct {
	HttpPort string
	Level    string
	Secret   string
	Log      *zap.Logger
}

func Load(path string) Config {
	cfg := Config{}
	log := logger.New(cfg.Level)
	cfg.Log = log

	err := godotenv.Load(path)
	if err != nil {
		fmt.Printf("failed to load config by path %s. Using default: %v\n", path, err)
		return Config{
			HttpPort: "8080",
			Level:    "level",
			Secret: "secret-terces-secret",
			Log:      log,
		}
	}

	return cfg
}
