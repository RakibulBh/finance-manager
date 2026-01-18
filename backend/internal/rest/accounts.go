package rest

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/rakibulbh/ai-finance-manager/internal/models"
)

type AccountStore interface {
	Create(ctx context.Context, acc *models.Account) error
	ListByFamilyID(ctx context.Context, familyID uuid.UUID) ([]models.Account, error)
	GetNetWorth(ctx context.Context, familyID uuid.UUID) (float64, error)
}

type AccountHandler struct {
	repo AccountStore
}

func NewAccountHandler(repo AccountStore) *AccountHandler {
	return &AccountHandler{repo: repo}
}

func (h *AccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	var acc models.Account
	if err := json.NewDecoder(r.Body).Decode(&acc); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	fID, ok := r.Context().Value("family_id").(uuid.UUID)
	if !ok {
		sendError(w, http.StatusUnauthorized, "Family ID missing from context")
		return
	}
	acc.FamilyID = fID

	// 2. Validate
	if acc.Name == "" || acc.Currency == "" {
		sendError(w, http.StatusBadRequest, "Name and Currency are required")
		return
	}

	// 3. Determine Classification
	if acc.Classification == "" {
		if acc.Type == "credit_card" || acc.Type == "loan" {
			acc.Classification = "liability"
		} else {
			acc.Classification = "asset"
		}
	}

	// 4. Create
	if err := h.repo.Create(r.Context(), &acc); err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to create account")
		return
	}

	sendJSON(w, http.StatusCreated, acc)
}

func (h *AccountHandler) List(w http.ResponseWriter, r *http.Request) {
	familyID, ok := r.Context().Value("family_id").(uuid.UUID)
	if !ok {
		sendError(w, http.StatusUnauthorized, "Family ID missing from context")
		return
	}

	accounts, err := h.repo.ListByFamilyID(r.Context(), familyID)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to fetch accounts")
		return
	}
	if accounts == nil {
		accounts = []models.Account{}
	}

	netWorth, err := h.repo.GetNetWorth(r.Context(), familyID)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to calculate net worth")
		return
	}

	response := map[string]interface{}{
		"data": map[string]interface{}{
			"accounts":  accounts,
			"net_worth": netWorth,
		},
	}

	sendJSON(w, http.StatusOK, response)
}
