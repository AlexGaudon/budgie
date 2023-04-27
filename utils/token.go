package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/alexgaudon/budgie/config"
	"github.com/golang-jwt/jwt/v5"
)

func CreateToken(userId string, ttl time.Duration) (string, error) {
	secret := config.GetConfig().JWTSecret

	now := time.Now().UTC()

	expiresAt := now.Add(ttl).Unix()

	claims := &jwt.MapClaims{
		"iat": now.Unix(),
		"nbf": now.Unix(),
		"sub": userId,
		"exp": expiresAt,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func GetTokenClaims(token string) (jwt.MapClaims, error) {
	validatedToken, err := ValidateToken(token)

	if err != nil {
		return nil, err
	}

	if validatedToken.Valid {
		claims := validatedToken.Claims.(jwt.MapClaims)
		return claims, nil
	}
	return nil, fmt.Errorf("token is invalid")
}

func ValidateToken(token string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")

	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
}
