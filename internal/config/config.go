package config

import (
	"encoding/json"
	"flag"
	"os"
	"time"

	"github.com/aga-absolut/Vault-System/internal/models"
	"github.com/caarlos0/env"
)

type Config struct {
	Address       string        `env:"ADDRESS" envDefault:":3200"`
	DatabaseDSN   string        `env:"DATABASE_DSN"`
	TokenTTL      time.Duration `json:"token_ttl"`
	JWTSecret     string        `json:"jwt_secret"`
	EncryptionKey string        `json:"encryption_key"`
}

func MustLoadConfig() *Config {
	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		panic(err)
	}

	flag.StringVar(&cfg.Address, "a", ":3200", "address for strarting server")
	flag.StringVar(&cfg.DatabaseDSN, "d", "postgres://postgres:absolute_1@localhost:5432/keeper", "value for connecting to database")
	flag.Parse()

	if err := ParseConfigFromFile("secret.json", cfg); err != nil {
		panic(err)
	}

	return cfg
}

func ParseConfigFromFile(name string, cfg *Config) error {
	file, err := os.Open(name)
	if err != nil {
		return err
	}
	defer file.Close()

	var configFile models.SecretConfig
	if err := json.NewDecoder(file).Decode(&configFile); err != nil {
		return err
	}

	cfg.JWTSecret = configFile.JWTSecret
	cfg.EncryptionKey = configFile.EncryptionKey

	ttl, err := time.ParseDuration(configFile.TokenTTL)
	if err != nil {
		return err
	}

	cfg.TokenTTL = ttl

	return nil
}
