package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rakibulbh/ai-finance-manager/internal/config"
	"github.com/rakibulbh/ai-finance-manager/internal/repository/postgres"
	"github.com/rakibulbh/ai-finance-manager/internal/rest"
	authMW "github.com/rakibulbh/ai-finance-manager/internal/rest/middleware"
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

	// 2. Handlers
	authHandler := rest.NewAuthHandler(userRepo, cfg.JWTSecret)
	accountHandler := rest.NewAccountHandler(accountRepo)

	// 3. Router Setup
	r := chi.NewRouter()

	// Base Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// API Routes
	r.Route("/api", func(r chi.Router) {
		// Public Routes
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)

		// Protected Routes
		r.Group(func(r chi.Router) {
			r.Use(authMW.AuthMiddleware([]byte(cfg.JWTSecret)))

			r.Route("/accounts", func(r chi.Router) {
				r.Post("/", accountHandler.Create)
				r.Get("/", accountHandler.List)
			})
		})
	})

	fmt.Printf("Starting application in %s mode...\n", cfg.Environment)
	fmt.Printf("Server running on port :%s\n", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}
