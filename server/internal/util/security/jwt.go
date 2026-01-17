package security

import (
	"errors"
	"server/internal/util/env"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(userID int32) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(env.GetEnv("JWT_SECRET", "")))
}

func ParseJWT(tokenStr string) (int32, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		return []byte([]byte(env.GetEnv("JWT_SECRET", ""))), nil
	})
	if err != nil || !token.Valid {
		return 0, errors.New("invalid token")
	}

	claims := token.Claims.(jwt.MapClaims)
	return int32(claims["sub"].(float64)), nil
}
