package rest

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/rakibulbh/ai-finance-manager/internal/jobs"
	"github.com/rakibulbh/ai-finance-manager/internal/models"
)

type PlaidManager interface {
	CreateLinkToken(ctx context.Context, userID, clientName string) (string, error)
	ExchangePublicToken(ctx context.Context, publicToken string) (string, string, error)
	EncryptToken(token string) (string, error)
}

type PlaidDB interface {
	SaveItem(ctx context.Context, item *models.PlaidItem) error
}

type PlaidHandler struct {
	manager PlaidManager
	repo    PlaidDB
	queue   *asynq.Client
}

func NewPlaidHandler(manager PlaidManager, repo PlaidDB, queue *asynq.Client) *PlaidHandler {
	return &PlaidHandler{
		manager: manager,
		repo:    repo,
		queue:   queue,
	}
}


func (h *PlaidHandler) CreateLinkToken(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		sendError(w, http.StatusUnauthorized, "User ID missing from context")
		return
	}

	token, err := h.manager.CreateLinkToken(r.Context(), userID.String(), "Maybe Finance")
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to create link token")
		return
	}

	sendJSON(w, http.StatusOK, map[string]string{"link_token": token})
}

type exchangeRequest struct {
	PublicToken     string `json:"public_token"`
	InstitutionID   string `json:"institution_id"`
	InstitutionName string `json:"institution_name"`
}

func (h *PlaidHandler) ExchangePublicToken(w http.ResponseWriter, r *http.Request) {
	familyID, ok := r.Context().Value("family_id").(uuid.UUID)
	if !ok {
		sendError(w, http.StatusUnauthorized, "Family ID missing from context")
		return
	}

	var req exchangeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	accessToken, itemID, err := h.manager.ExchangePublicToken(r.Context(), req.PublicToken)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to exchange token")
		return
	}

	encryptedToken, err := h.manager.EncryptToken(accessToken)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to encrypt token")
		return
	}

	item := &models.PlaidItem{
		FamilyID:        familyID,
		AccessToken:     encryptedToken,
		ItemID:          itemID,
		InstitutionID:   req.InstitutionID,
		InstitutionName: req.InstitutionName,
		Status:          "active",
	}

	if err := h.repo.SaveItem(r.Context(), item); err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to save plaid item")
		return
	}

	// 6. Enqueue initial sync
	task, err := jobs.NewSyncAccountTask(familyID, itemID)
	if err == nil {
		_, err = h.queue.Enqueue(task)
	}
	if err != nil {
		// Just log, the account is linked anyway
		http.Error(w, "", 0) // dummy to avoid unused variable if need be, but let's just log
		println("Failed to enqueue initial sync:", err.Error())
	}

	sendJSON(w, http.StatusOK, map[string]string{"message": "Account linked successfully", "item_id": itemID})
}
