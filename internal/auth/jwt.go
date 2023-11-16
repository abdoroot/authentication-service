package auth

import (
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

func IsUserAuthorizedWithClaim(tokenString string) (jwt.MapClaims, bool) {
	if tokenString != "" {
		tkn, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil {
			log.Println("error parse token")
			return nil, false
		}
		//get jwt claim
		claims, ok := tkn.Claims.(jwt.MapClaims)
		if !ok {
			log.Println("claims error: token claims are not of type jwt.MapClaims")
			return nil, false
		}
		return claims, true
	}
	log.Println("Token not found")
	return nil, false
}
