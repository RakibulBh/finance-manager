package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/rakibulbh/ai-finance-manager/internal/models"
	"github.com/rakibulbh/ai-finance-manager/internal/rest/mocks"
)

// Tests based on Ruby test/specifications from maybe/test/controllers/api/v1/accounts_controller_test.rb

// Test "should create account with valid parameters"
func TestAccountHandler_Create_Success(t *testing.T) {
	store := mocks.NewAccountStore()
	handler := NewAccountHandler(store)

	familyID := uuid.New()
	reqBody := models.Account{
		Name:           "Test Checking Account",
		Balance:        1000.50,
		Currency:       "USD",
		Type:           "depository",
		Subtype:        "checking",
		Classification: "asset",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/accounts", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Set family_id in context (simulating authenticated request)
	ctx := context.WithValue(req.Context(), "family_id", familyID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.Create(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d. Body: %s", w.Code, w.Body.String())
	}

	var response models.Account
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.Name != "Test Checking Account" {
		t.Errorf("Expected name 'Test Checking Account', got %s", response.Name)
	}
	if response.Balance != 1000.50 {
		t.Errorf("Expected balance 1000.50, got %f", response.Balance)
	}
	if response.ID == (uuid.UUID{}) {
		t.Error("Expected account ID to be set")
	}
}

// Test "should auto-set classification to asset for depository accounts"
func TestAccountHandler_Create_AutoClassificationAsset(t *testing.T) {
	store := mocks.NewAccountStore()
	handler := NewAccountHandler(store)

	familyID := uuid.New()
	reqBody := models.Account{
		Name:    "Savings Account",
		Balance: 5000.00,
		Currency: "USD",
		Type:    "depository",
		Subtype: "savings",
		// Classification not set - should default to asset
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/accounts", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), "family_id", familyID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.Create(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	var response models.Account
	json.NewDecoder(w.Body).Decode(&response)

	if response.Classification != "asset" {
		t.Errorf("Expected auto-classification 'asset', got '%s'", response.Classification)
	}
}

// Test "should auto-set classification to liability for credit card accounts"
func TestAccountHandler_Create_AutoClassificationLiability(t *testing.T) {
	store := mocks.NewAccountStore()
	handler := NewAccountHandler(store)

	familyID := uuid.New()
	reqBody := models.Account{
		Name:    "Credit Card",
		Balance: -250.00,
		Currency: "USD",
		Type:    "credit_card",
		// Classification not set - should default to liability
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/accounts", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), "family_id", familyID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.Create(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	var response models.Account
	json.NewDecoder(w.Body).Decode(&response)

	if response.Classification != "liability" {
		t.Errorf("Expected auto-classification 'liability', got '%s'", response.Classification)
	}
}

// Test "should require authentication (family_id in context)"
func TestAccountHandler_Create_Unauthorized(t *testing.T) {
	store := mocks.NewAccountStore()
	handler := NewAccountHandler(store)

	reqBody := models.Account{
		Name:    "Test Account",
		Balance: 100,
		Currency: "USD",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/accounts", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	// No family_id in context

	w := httptest.NewRecorder()
	handler.Create(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}

	var response map[string]interface{}
	json.NewDecoder(w.Body).Decode(&response)
	if response["error"] != "Family ID missing from context" {
		t.Errorf("Expected 'Family ID missing from context' error, got %v", response["error"])
	}
}

// Test "should require name and currency"
func TestAccountHandler_Create_MissingRequiredFields(t *testing.T) {
	store := mocks.NewAccountStore()
	handler := NewAccountHandler(store)

	familyID := uuid.New()
	reqBody := models.Account{
		Balance: 100,
		// Missing Name and Currency
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/accounts", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), "family_id", familyID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	var response map[string]interface{}
	json.NewDecoder(w.Body).Decode(&response)
	if response["error"] != "Name and Currency are required" {
		t.Errorf("Expected 'Name and Currency are required' error, got %v", response["error"])
	}
}

// Test "should list accounts successfully"
func TestAccountHandler_List_Success(t *testing.T) {
	store := mocks.NewAccountStore()
	handler := NewAccountHandler(store)

	familyID := uuid.New()

	// Add some test accounts
	store.AddAccount(familyID, models.Account{
		ID:             uuid.New(),
		Name:           "Checking",
		Balance:        1000,
		Currency:       "USD",
		Classification: "asset",
	})
	store.AddAccount(familyID, models.Account{
		ID:             uuid.New(),
		Name:           "Credit Card",
		Balance:        100,
		Currency:       "USD",
		Classification: "liability",
	})

	req := httptest.NewRequest("GET", "/accounts", nil)
	ctx := context.WithValue(req.Context(), "family_id", familyID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.List(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	var response map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	data := response["data"].(map[string]interface{})
	accounts := data["accounts"].([]interface{})

	if len(accounts) != 2 {
		t.Errorf("Expected 2 accounts, got %d", len(accounts))
	}

	netWorth := data["net_worth"].(float64)
	if netWorth != 900 {
		t.Errorf("Expected net worth 900, got %f", netWorth)
	}
}

// Test "should return empty list for family with no accounts"
func TestAccountHandler_List_NoAccounts(t *testing.T) {
	store := mocks.NewAccountStore()
	handler := NewAccountHandler(store)

	familyID := uuid.New()

	req := httptest.NewRequest("GET", "/accounts", nil)
	ctx := context.WithValue(req.Context(), "family_id", familyID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.List(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.NewDecoder(w.Body).Decode(&response)

	data := response["data"].(map[string]interface{})
	accounts := data["accounts"].([]interface{})

	if len(accounts) != 0 {
		t.Errorf("Expected 0 accounts, got %d", len(accounts))
	}
}

// Test "should require authentication for list"
func TestAccountHandler_List_Unauthorized(t *testing.T) {
	store := mocks.NewAccountStore()
	handler := NewAccountHandler(store)

	req := httptest.NewRequest("GET", "/accounts", nil)
	// No family_id in context

	w := httptest.NewRecorder()
	handler.List(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}

	var response map[string]interface{}
	json.NewDecoder(w.Body).Decode(&response)
	if response["error"] != "Family ID missing from context" {
		t.Errorf("Expected 'Family ID missing from context' error, got %v", response["error"])
	}
}

// Test "should handle invalid JSON"
func TestAccountHandler_Create_InvalidJSON(t *testing.T) {
	store := mocks.NewAccountStore()
	handler := NewAccountHandler(store)

	familyID := uuid.New()

	req := httptest.NewRequest("POST", "/accounts", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), "family_id", familyID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	var response map[string]interface{}
	json.NewDecoder(w.Body).Decode(&response)
	if response["error"] != "Invalid request body" {
		t.Errorf("Expected 'Invalid request body' error, got %v", response["error"])
	}
}

// Test "should handle loan type accounts"
func TestAccountHandler_Create_LoanType(t *testing.T) {
	store := mocks.NewAccountStore()
	handler := NewAccountHandler(store)

	familyID := uuid.New()
	reqBody := models.Account{
		Name:     "Mortgage",
		Balance:  250000,
		Currency: "USD",
		Type:     "loan",
		Subtype:  "mortgage",
		LoanDetails: &models.LoanDetails{
			InterestRate: 4.5,
			TermMonths:   360,
		},
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/accounts", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), "family_id", familyID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.Create(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	var response models.Account
	json.NewDecoder(w.Body).Decode(&response)

	if response.Classification != "liability" {
		t.Errorf("Expected auto-classification 'liability' for loan, got '%s'", response.Classification)
	}
}

// Test "should handle property type accounts"
func TestAccountHandler_Create_PropertyType(t *testing.T) {
	store := mocks.NewAccountStore()
	handler := NewAccountHandler(store)

	familyID := uuid.New()
	reqBody := models.Account{
		Name:     "Family Home",
		Balance:  450000,
		Currency: "USD",
		Type:     "property",
		PropertyDetails: &models.PropertyDetails{
			Address: "123 Main St",
			Sqft:    2000,
		},
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/accounts", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), "family_id", familyID)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.Create(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	var response models.Account
	json.NewDecoder(w.Body).Decode(&response)

	if response.Classification != "asset" {
		t.Errorf("Expected classification 'asset' for property, got '%s'", response.Classification)
	}
	if response.PropertyDetails == nil {
		t.Error("Expected PropertyDetails to be set")
	}
}
