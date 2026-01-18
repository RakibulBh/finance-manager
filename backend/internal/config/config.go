package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	DatabaseURL string `mapstructure:"DATABASE_URL"`
	JWTSecret   string `mapstructure:"JWT_SECRET"`
	Port        string `mapstructure:"PORT"`
	Environment string `mapstructure:"ENVIRONMENT"`

	// Plaid
	PlaidClientID string `mapstructure:"PLAID_CLIENT_ID"`
	PlaidSecret   string `mapstructure:"PLAID_SECRET"`
	PlaidEnv      string `mapstructure:"PLAID_ENV"`

	// Encryption
	EncryptionKey string `mapstructure:"ENCRYPTION_KEY"`

	// Redis
	RedisAddr string `mapstructure:"REDIS_ADDR"`
}


func LoadConfig() (*Config, error) {
	// 1. Load .env file if it exists
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}

	// 2. Set up Viper
	viper.AutomaticEnv()

	// Bind environment variables to struct fields
	viper.BindEnv("DATABASE_URL")
	viper.BindEnv("JWT_SECRET")
	viper.BindEnv("PORT")
	viper.BindEnv("ENVIRONMENT")
	viper.BindEnv("PLAID_CLIENT_ID")
	viper.BindEnv("PLAID_SECRET")
	viper.BindEnv("PLAID_ENV")
	viper.BindEnv("ENCRYPTION_KEY")
	viper.BindEnv("REDIS_ADDR")


	var cfg Config
	err = viper.Unmarshal(&cfg)

	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
