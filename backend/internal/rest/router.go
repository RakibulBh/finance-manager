package rest

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	authMW "github.com/rakibulbh/ai-finance-manager/internal/rest/middleware"
)

type RouterConfig struct {
	AuthHandler        *AuthHandler
	AccountHandler     *AccountHandler
	TransactionHandler *TransactionHandler
	InvestmentHandler  *InvestmentHandler
	PlaidHandler       *PlaidHandler
	JWTSecret         string
}

func NewRouter(cfg RouterConfig) *chi.Mux {
	r := chi.NewRouter()

	r.Use(authMW.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Route("/api", func(r chi.Router) {
		// Public Routes
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", cfg.AuthHandler.Register)
			r.Post("/login", cfg.AuthHandler.Login)
		})

		// Protected Routes
		r.Group(func(r chi.Router) {
			r.Use(authMW.AuthMiddleware([]byte(cfg.JWTSecret)))

			r.Route("/accounts", func(r chi.Router) {
				r.Post("/", cfg.AccountHandler.Create)
				r.Get("/", cfg.AccountHandler.List)
			})

			r.Route("/transactions", func(r chi.Router) {
				r.Post("/", cfg.TransactionHandler.Create)
			})

			r.Route("/transfers", func(r chi.Router) {
				r.Post("/", cfg.TransactionHandler.CreateTransfer)
			})

			r.Route("/investments", func(r chi.Router) {
				r.Post("/trade", cfg.InvestmentHandler.CreateTrade)
			})

			r.Route("/plaid", func(r chi.Router) {
				r.Post("/create_link_token", cfg.PlaidHandler.CreateLinkToken)
				r.Post("/exchange_public_token", cfg.PlaidHandler.ExchangePublicToken)
			})
		})
	})

	return r
}
