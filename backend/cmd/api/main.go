package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rakibulbh/ai-finance-manager/internal/config"
	"github.com/rakibulbh/ai-finance-manager/internal/repository/postgres"
	"github.com/rakibulbh/ai-finance-manager/internal/rest"
	"github.com/rakibulbh/ai-finance-manager/internal/services"
)


func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Could not load config: ", err)
	}

	// connect to the db
	dbPool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}
	defer dbPool.Close()

	// 1. Repositories
	userRepo := postgres.NewUserRepository(dbPool)
	accountRepo := postgres.NewAccountRepository(dbPool)
	ledgerRepo := postgres.NewLedgerRepository(dbPool)
	investmentRepo := postgres.NewInvestmentRepository(dbPool)
	plaidRepo := postgres.NewPlaidRepository(dbPool)

	// 2. Services
	plaidService := services.NewPlaidService(cfg.PlaidClientID, cfg.PlaidSecret, cfg.PlaidEnv, cfg.EncryptionKey)
	asynqClient := asynq.NewClient(asynq.RedisClientOpt{Addr: cfg.RedisAddr})
	defer asynqClient.Close()

	// 3. Handlers
	authHandler := rest.NewAuthHandler(userRepo, cfg.JWTSecret)
	accountHandler := rest.NewAccountHandler(accountRepo)
	transactionHandler := rest.NewTransactionHandler(ledgerRepo)
	investmentHandler := rest.NewInvestmentHandler(investmentRepo)
	plaidHandler := rest.NewPlaidHandler(plaidService, plaidRepo, asynqClient)

	// 4. Router Setup
	r := rest.NewRouter(rest.RouterConfig{
		AuthHandler:        authHandler,
		AccountHandler:     accountHandler,
		TransactionHandler: transactionHandler,
		InvestmentHandler:  investmentHandler,
		PlaidHandler:       plaidHandler,
		JWTSecret:         cfg.JWTSecret,
	})




	fmt.Printf("Starting application in %s mode...\n", cfg.Environment)
	fmt.Printf("Server running on port :%s\n", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}
