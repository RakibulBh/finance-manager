package rest

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rakibulbh/ai-finance-manager/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type UserStore interface {
    CreateUser(ctx context.Context, email, password string, familyID uuid.UUID) (*models.User, error)
    FindByEmail(ctx context.Context, email string) (*models.User, string, error)
}

type RegisterRequest struct {
    Email      string `json:"email"`
    Password   string `json:"password"`
    FamilyName string `json:"family_name"`
}

type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type AuthHandler struct {
    repo      UserStore
    jwtSecret []byte
}

func NewAuthHandler(repo UserStore, jwtSecret string) *AuthHandler {
    return &AuthHandler{
        repo:      repo,
        jwtSecret: []byte(jwtSecret),
    }
}



// POST /register
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
    var req RegisterRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        sendError(w, http.StatusBadRequest, "Invalid request body")
        return
    }

    // Simple validation
    if req.Email == "" || req.Password == "" {
        sendError(w, http.StatusBadRequest, "Email and password are required")
        return
    }

    // Generate a Family ID based on the family_name input (or just random)
    // The prompt implies "family_name" comes in, but "family_id" goes out.
    // We will generate a deterministic UUID based on the name, or just a new one.
    // Let's create a new UUID for the family ID for this registration.
    familyID := uuid.New()

    user, err := h.repo.CreateUser(r.Context(), req.Email, req.Password, familyID)
    if err != nil {
        if strings.Contains(err.Error(), "duplicate key value") {
            sendError(w, http.StatusConflict, "Email already registered")
            return
        }
        log.Printf("DB Error: %v", err)
        sendError(w, http.StatusInternalServerError, "Failed to create user")
        return
    }

    // Response format requested: [{ "data": { ... } }]
    response := []map[string]interface{}{
        {
            "data": map[string]interface{}{
                "user_id":    user.ID,
                "email":      user.Email,
                "family_id":  user.FamilyID,
                "message":    "Registration successful",
            },
        },
    }

    sendJSON(w, http.StatusCreated, response)
}



func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    var req LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        sendError(w, http.StatusBadRequest, "Invalid request body")
        return
    }

    user, hashedPassword, err := h.repo.FindByEmail(r.Context(), req.Email)
    if err != nil {
        if err == pgx.ErrNoRows {
            sendError(w, http.StatusUnauthorized, "Invalid email or password")
        } else {
            log.Printf("DB Error: %v", err)
            sendError(w, http.StatusInternalServerError, "Database error")
        }
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)); err != nil {
        sendError(w, http.StatusUnauthorized, "Invalid email or password")
        return
    }

    // Generate JWT
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id":   user.ID.String(),
        "email":     user.Email,
        "family_id": user.FamilyID.String(),
        "exp":       time.Now().Add(time.Hour * 24 * 7).Unix(), // 1 week
    })

    tokenString, err := token.SignedString(h.jwtSecret)
    if err != nil {
        sendError(w, http.StatusInternalServerError, "Failed to generate token")
        return
    }

    response := map[string]interface{}{
        "data": map[string]interface{}{
            "token": tokenString,
            "user": map[string]interface{}{
                "id":        user.ID,
                "email":     user.Email,
                "family_id": user.FamilyID,
            },
        },
    }

    sendJSON(w, http.StatusOK, response)
}
