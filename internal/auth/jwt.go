package auth

import (
	"context"
	"fmt"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

const (
	secret        string = "jwts-screct-code@123.."
	tokenLifespan int    = 3
)

func GenerateToken(userId, userEmail string) (string, error) {
	claims := jwt.MapClaims{}
	claims["user_id"] = userId
	claims["user_email"] = userEmail
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(tokenLifespan)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func IsUserAuthorized(ctx context.Context) bool {
	token := ctx.Value("token")
	if tokenString, ok := token.(string); ok {
		_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		})
		return err != nil
	}
	return false
}
