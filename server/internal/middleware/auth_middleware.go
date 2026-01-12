package middleware

import (
	"context"
	"net/http"
	"os"
	"server/internal/response"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			response.WriteJson(w, http.StatusOK, response.GernalResponse{Success: false, Message: "unauthorized", Data: nil})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			response.WriteJson(w, http.StatusOK, response.GernalResponse{Success: false, Message: "unauthorized", Data: nil})
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		userID := int32(claims["sub"].(float64))

		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
