package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rakibulbh/ai-finance-manager/internal/models"
	"github.com/rakibulbh/ai-finance-manager/internal/rest/mocks"
)

// Tests based on Ruby test/specifications from maybe/test/controllers/api/v1/transactions_controller_test.rb

// Test "should create transaction with valid parameters"
func TestTransactionHandler_Create_Success(t *testing.T) {
	store := mocks.NewTransactionStore()
	handler := NewTransactionHandler(store)

	accountID := uuid.New()
	categoryID := uuid.New()
	reqBody := CreateTransactionRequest{
		AccountID:    accountID,
		Amount:       25.00,
		Date:         time.Now(),
		Name:         "Test Transaction",
		CategoryID:   &categoryID,
		MerchantName: "Test Merchant",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/transactions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	familyID := uuid.New()
	ctx := context.WithValue(req.Context(), "family_id", familyID)
	handler.Create(w, req.WithContext(ctx))

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d. Body: %s", w.Code, w.Body.String())
	}

	var response models.Entry
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.Name != "Test Transaction" {
		t.Errorf("Expected name 'Test Transaction', got %s", response.Name)
	}
	if response.Amount != 25.00 {
		t.Errorf("Expected amount 25.00, got %f", response.Amount)
	}
}

// Test "should create transaction without merchant"
func TestTransactionHandler_Create_NoMerchant(t *testing.T) {
	store := mocks.NewTransactionStore()
	handler := NewTransactionHandler(store)

	accountID := uuid.New()
	reqBody := CreateTransactionRequest{
		AccountID: accountID,
		Amount:    25.00,
		Date:      time.Now(),
		Name:      "Test Transaction",
		// No merchant name
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/transactions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	familyID := uuid.New()
	ctx := context.WithValue(req.Context(), "family_id", familyID)
	handler.Create(w, req.WithContext(ctx))

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}
}

// Test "should default to current date when not provided"
func TestTransactionHandler_Create_DefaultDate(t *testing.T) {
	store := mocks.NewTransactionStore()
	handler := NewTransactionHandler(store)

	accountID := uuid.New()
	reqBody := CreateTransactionRequest{
		AccountID: accountID,
		Amount:    25.00,
		Name:      "Test Transaction",
		// No date - should default to now
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/transactions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	familyID := uuid.New()
	ctx := context.WithValue(req.Context(), "family_id", familyID)
	handler.Create(w, req.WithContext(ctx))

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	var response models.Entry
	json.NewDecoder(w.Body).Decode(&response)

	// Check that date was set (should be recent)
	if time.Since(response.Date) > time.Second {
		t.Error("Expected date to be set to current time")
	}
}

// Test "should handle invalid JSON"
func TestTransactionHandler_Create_InvalidJSON(t *testing.T) {
	store := mocks.NewTransactionStore()
	handler := NewTransactionHandler(store)

	req := httptest.NewRequest("POST", "/transactions", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	familyID := uuid.New()
	ctx := context.WithValue(req.Context(), "family_id", familyID)
	handler.Create(w, req.WithContext(ctx))

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	var response map[string]interface{}
	json.NewDecoder(w.Body).Decode(&response)
	if response["error"] != "Invalid request body" {
		t.Errorf("Expected 'Invalid request body' error, got %v", response["error"])
	}
}

// Test "should handle merchant creation failure"
func TestTransactionHandler_Create_MerchantError(t *testing.T) {
	store := mocks.NewTransactionStore()
	store.MerchantError = &MerchantError{Message: "Failed to create merchant"}
	handler := NewTransactionHandler(store)

	accountID := uuid.New()
	reqBody := CreateTransactionRequest{
		AccountID:    accountID,
		Amount:       25.00,
		Name:         "Test Transaction",
		MerchantName: "Test Merchant",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/transactions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	familyID := uuid.New()
	ctx := context.WithValue(req.Context(), "family_id", familyID)
	handler.Create(w, req.WithContext(ctx))

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}

	var response map[string]interface{}
	json.NewDecoder(w.Body).Decode(&response)
	if response["error"] != "Failed to handle merchant" {
		t.Errorf("Expected 'Failed to handle merchant' error, got %v", response["error"])
	}
}

// Test "should handle transaction creation failure"
func TestTransactionHandler_Create_TransactionError(t *testing.T) {
	store := mocks.NewTransactionStore()
	store.CreateError = &TransactionError{Message: "Database error"}
	handler := NewTransactionHandler(store)

	accountID := uuid.New()
	reqBody := CreateTransactionRequest{
		AccountID: accountID,
		Amount:    25.00,
		Name:      "Test Transaction",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/transactions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	familyID := uuid.New()
	ctx := context.WithValue(req.Context(), "family_id", familyID)
	handler.Create(w, req.WithContext(ctx))

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}

	var response map[string]interface{}
	json.NewDecoder(w.Body).Decode(&response)
	if response["error"] != "Failed to create transaction" {
		t.Errorf("Expected 'Failed to create transaction' error, got %v", response["error"])
	}
}

// Test "should create transfer successfully"
func TestTransactionHandler_CreateTransfer_Success(t *testing.T) {
	store := mocks.NewTransactionStore()
	handler := NewTransactionHandler(store)

	fromAccountID := uuid.New()
	toAccountID := uuid.New()
	reqBody := CreateTransferRequest{
		FromAccountID: fromAccountID,
		ToAccountID:   toAccountID,
		Amount:        100.00,
		Date:          time.Now(),
		Name:          "Transfer to Savings",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/transfers", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.CreateTransfer(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d. Body: %s", w.Code, w.Body.String())
	}

	var response map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["message"] != "Transfer successful" {
		t.Errorf("Expected 'Transfer successful' message, got %v", response["message"])
	}
}

// Test "should default date when creating transfer"
func TestTransactionHandler_CreateTransfer_DefaultDate(t *testing.T) {
	store := mocks.NewTransactionStore()
	handler := NewTransactionHandler(store)

	fromAccountID := uuid.New()
	toAccountID := uuid.New()
	reqBody := CreateTransferRequest{
		FromAccountID: fromAccountID,
		ToAccountID:   toAccountID,
		Amount:        100.00,
		Name:          "Transfer",
		// No date - should default to now
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/transfers", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.CreateTransfer(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}
}

// Test "should handle invalid JSON for transfer"
func TestTransactionHandler_CreateTransfer_InvalidJSON(t *testing.T) {
	store := mocks.NewTransactionStore()
	handler := NewTransactionHandler(store)

	req := httptest.NewRequest("POST", "/transfers", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.CreateTransfer(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	var response map[string]interface{}
	json.NewDecoder(w.Body).Decode(&response)
	if response["error"] != "Invalid request body" {
		t.Errorf("Expected 'Invalid request body' error, got %v", response["error"])
	}
}

// Test "should handle transfer creation failure"
func TestTransactionHandler_CreateTransfer_TransferError(t *testing.T) {
	store := mocks.NewTransactionStore()
	store.TransferError = &TransferError{Message: "Insufficient funds"}
	handler := NewTransactionHandler(store)

	fromAccountID := uuid.New()
	toAccountID := uuid.New()
	reqBody := CreateTransferRequest{
		FromAccountID: fromAccountID,
		ToAccountID:   toAccountID,
		Amount:        100.00,
		Name:          "Transfer",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/transfers", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.CreateTransfer(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}

	var response map[string]interface{}
	json.NewDecoder(w.Body).Decode(&response)
	if response["error"] != "Failed to create transfer" {
		t.Errorf("Expected 'Failed to create transfer' error, got %v", response["error"])
	}
}

// Test "should handle merchant reuse"
func TestTransactionHandler_Create_ReuseMerchant(t *testing.T) {
	store := mocks.NewTransactionStore()
	handler := NewTransactionHandler(store)

	accountID := uuid.New()

	// First transaction with merchant
	reqBody1 := CreateTransactionRequest{
		AccountID:    accountID,
		Amount:       25.00,
		Name:         "First Transaction",
		MerchantName: "Coffee Shop",
	}
	body1, _ := json.Marshal(reqBody1)

	req1 := httptest.NewRequest("POST", "/transactions", bytes.NewBuffer(body1))
	req1.Header.Set("Content-Type", "application/json")

	w1 := httptest.NewRecorder()
	familyID := uuid.New()
	ctx1 := context.WithValue(req1.Context(), "family_id", familyID)
	handler.Create(w1, req1.WithContext(ctx1))

	// Second transaction with same merchant
	reqBody2 := CreateTransactionRequest{
		AccountID:    accountID,
		Amount:       15.00,
		Name:         "Second Transaction",
		MerchantName: "Coffee Shop",
	}
	body2, _ := json.Marshal(reqBody2)

	req2 := httptest.NewRequest("POST", "/transactions", bytes.NewBuffer(body2))
	req2.Header.Set("Content-Type", "application/json")

	w2 := httptest.NewRecorder()
	// Same familyID
	ctx2 := context.WithValue(req2.Context(), "family_id", familyID)
	handler.Create(w2, req2.WithContext(ctx2))

	if w2.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w2.Code)
	}

	// Should only have one merchant ID
	if len(store.Merchants) != 1 {
		t.Errorf("Expected 1 merchant, got %d", len(store.Merchants))
	}
}

// Test "should handle different merchants"
func TestTransactionHandler_Create_DifferentMerchants(t *testing.T) {
	store := mocks.NewTransactionStore()
	handler := NewTransactionHandler(store)

	accountID := uuid.New()

	// First transaction
	reqBody1 := CreateTransactionRequest{
		AccountID:    accountID,
		Amount:       25.00,
		Name:         "Coffee",
		MerchantName: "Coffee Shop",
	}
	body1, _ := json.Marshal(reqBody1)

	req1 := httptest.NewRequest("POST", "/transactions", bytes.NewBuffer(body1))
	req1.Header.Set("Content-Type", "application/json")

	w1 := httptest.NewRecorder()
	familyID := uuid.New()
	ctx1 := context.WithValue(req1.Context(), "family_id", familyID)
	handler.Create(w1, req1.WithContext(ctx1))

	// Second transaction with different merchant
	reqBody2 := CreateTransactionRequest{
		AccountID:    accountID,
		Amount:       50.00,
		Name:         "Groceries",
		MerchantName: "Grocery Store",
	}
	body2, _ := json.Marshal(reqBody2)

	req2 := httptest.NewRequest("POST", "/transactions", bytes.NewBuffer(body2))
	req2.Header.Set("Content-Type", "application/json")

	w2 := httptest.NewRecorder()
	ctx2 := context.WithValue(req2.Context(), "family_id", familyID)
	handler.Create(w2, req2.WithContext(ctx2))

	if w2.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w2.Code)
	}

	// Should have two merchant IDs
	if len(store.Merchants) != 2 {
		t.Errorf("Expected 2 merchants, got %d", len(store.Merchants))
	}
}

// Custom error types
type MerchantError struct {
	Message string
}

func (e *MerchantError) Error() string {
	return e.Message
}

type TransactionError struct {
	Message string
}

func (e *TransactionError) Error() string {
	return e.Message
}

type TransferError struct {
	Message string
}

func (e *TransferError) Error() string {
	return e.Message
}
