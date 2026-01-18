package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rakibulbh/ai-finance-manager/internal/config"
	"github.com/rakibulbh/ai-finance-manager/internal/logger"
	"github.com/rakibulbh/ai-finance-manager/internal/repository/postgres"
	"github.com/rakibulbh/ai-finance-manager/internal/rest"
	"github.com/rakibulbh/ai-finance-manager/internal/services"
	"go.uber.org/zap"
)


func main() {
	if err := logger.InitLogger(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("Could not load config", zap.Error(err))
		os.Exit(1)
	}

	// connect to the db
	dbPool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		logger.Error("Unable to create connection pool", zap.Error(err))
		os.Exit(1)
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
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		logger.Error("Server failed to start", zap.Error(err))
		os.Exit(1)
	}
}
