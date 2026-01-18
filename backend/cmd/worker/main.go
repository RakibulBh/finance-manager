package main

import (
	"context"
	"fmt"
	"os"

	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rakibulbh/ai-finance-manager/internal/config"
	"github.com/rakibulbh/ai-finance-manager/internal/jobs"
	"github.com/rakibulbh/ai-finance-manager/internal/logger"
	"github.com/rakibulbh/ai-finance-manager/internal/repository/postgres"
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

	// 1. Database
	dbPool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		logger.Error("Unable to create connection pool", zap.Error(err))
		os.Exit(1)
	}
	defer dbPool.Close()

	// 2. Services
	plaidService := services.NewPlaidService(cfg.PlaidClientID, cfg.PlaidSecret, cfg.PlaidEnv, cfg.EncryptionKey)

	// 3. Repositories
	plaidRepo := postgres.NewPlaidRepository(dbPool)
	ledgerRepo := postgres.NewLedgerRepository(dbPool)
	accountRepo := postgres.NewAccountRepository(dbPool)

	// 4. Worker Setup
	svc := &jobs.WorkerServices{
		Plaid:    plaidService,
		DB:       plaidRepo,
		Ledger:   ledgerRepo,
		Accounts: accountRepo,
	}

	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: cfg.RedisAddr},
		asynq.Config{Concurrency: 10},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc(jobs.TypeSyncAccount, func(ctx context.Context, t *asynq.Task) error {
		return jobs.HandleSyncAccountTask(ctx, t, svc)
	})

	fmt.Printf("Worker server started on %s\n", cfg.RedisAddr)
	if err := srv.Run(mux); err != nil {
		logger.Error("Worker failed", zap.Error(err))
		os.Exit(1)
	}
}
