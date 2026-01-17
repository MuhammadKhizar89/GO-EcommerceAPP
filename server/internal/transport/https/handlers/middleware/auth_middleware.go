package middleware

import (
	"net/http"
	"server/internal/util/env"
	"server/internal/util/response"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			response.WriteJson(w, http.StatusUnauthorized,
				response.GernalResponse{Success: false, Message: "unauthorized", Data: nil})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(env.GetEnv("JWT_SECRET", "")), nil
		})

		if err != nil || !token.Valid {
			response.WriteJson(w, http.StatusUnauthorized,
				response.GernalResponse{Success: false, Message: "unauthorized", Data: nil})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			response.WriteJson(w, http.StatusUnauthorized,
				response.GernalResponse{Success: false, Message: "invalid token", Data: nil})
			return
		}

		sub, ok := claims["sub"].(float64)
		if !ok {
			response.WriteJson(w, http.StatusUnauthorized,
				response.GernalResponse{Success: false, Message: "invalid token payload", Data: nil})
			return
		}

		userID := int32(sub)

		ctx := SetUserID(r.Context(), userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
