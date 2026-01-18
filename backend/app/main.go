package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rakibulbh/ai-finance-manager/internal/config"
	"github.com/rakibulbh/ai-finance-manager/internal/repository/postgres"
	"github.com/rakibulbh/ai-finance-manager/internal/rest"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Could not load config: ", err)
	}

	fmt.Printf("Starting application in %s mode...\n", cfg.Environment)
	fmt.Printf("Server will run on port %s\n", cfg.Port)
	fmt.Printf("Database URL: %s\n", cfg.DatabaseURL)

	// connect to the db
	dbPool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
    if err != nil {
        log.Fatalf("Unable to create connection pool: %v\n", err)
    }
    defer dbPool.Close()

	// Auth
	userRepository := postgres.NewUserRepository(dbPool)
	_ = rest.NewAuthHandler(userRepository, cfg.JWTSecret)
}
