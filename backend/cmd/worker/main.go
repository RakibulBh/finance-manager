package main

import (
	"context"
	"log"

	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rakibulbh/ai-finance-manager/internal/config"
	"github.com/rakibulbh/ai-finance-manager/internal/jobs"
	"github.com/rakibulbh/ai-finance-manager/internal/repository/postgres"
	"github.com/rakibulbh/ai-finance-manager/internal/services"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Could not load config: ", err)
	}

	// 1. Database
	dbPool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
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

	log.Printf("Worker server started on %s", cfg.RedisAddr)
	if err := srv.Run(mux); err != nil {
		log.Fatal(err)
	}
}
