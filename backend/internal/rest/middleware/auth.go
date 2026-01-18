package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func AuthMiddleware(jwtSecret []byte) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            authHeader := r.Header.Get("Authorization")
            if authHeader == "" {
                sendError(w, http.StatusUnauthorized, "Authorization header missing")
                return
            }

            // Expecting format: Bearer <token>
            tokenString := ""
            if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
                tokenString = authHeader[7:]
            } else {
                sendError(w, http.StatusUnauthorized, "Authorization header missing")
                return
            }

            token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
                if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                    return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
                }
                return jwtSecret, nil
            })

            if err != nil {
                sendError(w, http.StatusUnauthorized, "Invalid token")
                return
            }

            if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
                // Extract user_id
                userIDStr, ok := claims["user_id"].(string)
                if !ok {
                    sendError(w, http.StatusUnauthorized, "Invalid token claims")
                    return
                }

                userID, err := uuid.Parse(userIDStr)
                if err != nil {
                    sendError(w, http.StatusUnauthorized, "Invalid user ID in token")
                    return
                }

                familyIDStr, ok := claims["family_id"].(string)
                if !ok {
                    sendError(w, http.StatusUnauthorized, "Invalid token claims")
                    return
                }

                familyID, err := uuid.Parse(familyIDStr)
                if err != nil {
                    sendError(w, http.StatusUnauthorized, "Invalid family ID in token")
                    return
                }

                // Pass user_id and family_id down via context
                ctx := context.WithValue(r.Context(), "user_id", userID)
                ctx = context.WithValue(ctx, "family_id", familyID)
                next.ServeHTTP(w, r.WithContext(ctx))
            } else {
                sendError(w, http.StatusUnauthorized, "Invalid token")
            }
        })
    }
}
