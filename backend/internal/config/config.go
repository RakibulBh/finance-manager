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

	var cfg Config
	err = viper.Unmarshal(&cfg)

	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
