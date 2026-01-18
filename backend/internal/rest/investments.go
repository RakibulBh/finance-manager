package rest

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rakibulbh/ai-finance-manager/internal/models"
)

type InvestmentStore interface {
	GetOrCreateSecurity(ctx context.Context, ticker, name string) (uuid.UUID, error)
	CreateTrade(ctx context.Context, entry *models.Entry, trade *models.Trade) error
	GetActiveTickers(ctx context.Context) ([]string, error)
	UpdateSecurityPrice(ctx context.Context, ticker string, price float64) error
}

type MarketDataProvider interface {
	GetQuote(ticker string) (float64, error)
}

type InvestmentHandler struct {
	repo InvestmentStore
}

func NewInvestmentHandler(repo InvestmentStore) *InvestmentHandler {
	return &InvestmentHandler{repo: repo}
}

type CreateTradeRequest struct {
	AccountID    uuid.UUID `json:"account_id"`
	Ticker       string    `json:"ticker"`
	SecurityName string    `json:"security_name"`
	Qty          float64   `json:"qty"`
	Price        float64   `json:"price"`
	Currency     string    `json:"currency"`
	Date         time.Time `json:"date"`
	Kind         string    `json:"kind"` // "buy", "sell"
}

func (h *InvestmentHandler) CreateTrade(w http.ResponseWriter, r *http.Request) {
	var req CreateTradeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// 1. Get/Create Security
	secID, err := h.repo.GetOrCreateSecurity(r.Context(), req.Ticker, req.SecurityName)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to handle security")
		return
	}

	// 2. Calculate Amount
	// Buy: -$1500 (money leaves account)
	// Sell: +$1500 (money enters account)
	amount := req.Qty * req.Price
	if req.Kind == "buy" || req.Kind == "" {
		amount = -amount
		req.Kind = "buy"
	}

	// 3. Prepare Entry
	entry := &models.Entry{
		AccountID: req.AccountID,
		Amount:    amount,
		Date:      req.Date,
		Name:      req.Kind + " " + req.Ticker,
		Currency:  req.Currency,
	}
	if entry.Date.IsZero() {
		entry.Date = time.Now()
	}
	if entry.Currency == "" {
		entry.Currency = "USD"
	}

	// 4. Prepare Trade
	trade := &models.Trade{
		SecurityID: secID,
		Qty:        req.Qty,
		Price:      req.Price,
		Kind:       req.Kind,
	}

	// 5. Save in DB
	if err := h.repo.CreateTrade(r.Context(), entry, trade); err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to record trade")
		return
	}

	sendJSON(w, http.StatusCreated, entry)
}
