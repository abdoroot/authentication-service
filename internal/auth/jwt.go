package auth

import (
	"context"
	"log"
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

func IsUserAuthorizedWithClaim(ctx context.Context) (jwt.MapClaims, bool) {
	token := ctx.Value("token")
	if tokenString, ok := token.(string); ok {
		tkn, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil {
			log.Println("error parse token")
			return nil, false
		}
		//get jwt claim
		if claims, ok := tkn.Claims.(jwt.MapClaims); ok {
			return claims, true
		}
		log.Println("claims error not type of jwt.MapClaims")
		return nil, false
	}
	log.Println("token not string", token)
	return nil, false
}
