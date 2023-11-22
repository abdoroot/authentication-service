package auth

import (
	"fmt"
	"log"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

const (
	secret               string = "jwts-screct-code@123.."
	tokenLifespan        int    = 3  //h
	refreshTokenLifespan int    = 24 //h
)

func GenerateToken(userId, userEmail string) (map[string]string, error) {
	//access token
	claims := jwt.MapClaims{}
	claims["user_id"] = userId
	claims["user_email"] = userEmail
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(tokenLifespan)).Unix()
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := accessToken.SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}

	//refresh token
	rtClaims := jwt.MapClaims{}
	rtClaims["user_id"] = userId
	rtClaims["exp"] = time.Now().Add(time.Hour * time.Duration(refreshTokenLifespan)).Unix()
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	rt, err := refreshToken.SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"access_token":  t,
		"refresh_token": rt,
	}, nil
}

func IsUserAuthorizedWithClaim(accessToken string) (jwt.MapClaims, bool) {
	if accessToken != "" {
		tkn, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
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

		if _, ok := claims["user_email"]; !ok {
			return nil, false
		}

		return claims, true
	}
	log.Println("Token not found")
	return nil, false
}

func RefreshAccessToken(refreshToken string) (string, error) {
	if refreshToken != "" {
		tkn, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil {
			return "", fmt.Errorf("error parse accessToken")
		}

		//get jwt claim
		claims, ok := tkn.Claims.(jwt.MapClaims)
		if !ok {
			return "", fmt.Errorf("claims error: token claims are not of type jwt.MapClaims")
		}

		//parse user_id
		if userId, ok := claims["user_id"]; ok {
			//get user data
			user, err := FindUserById(userId.(string))
			if err != nil {
				return "", err
			}
			//generate new token and refresh token
			mp, err := GenerateToken(user.UserId, user.Email)
			if err != nil {
				return "", err
			}
			if rt, ok := mp["refresh_token"]; ok {
				return rt, nil
			}
			return "", err
		}

		return "", fmt.Errorf("unknowow error")
	}
	return "", fmt.Errorf("refreshToken empty")
}

func CheckTokenExpiry(tokenString string) bool {
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}

	expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)
	// 20 minutes threshold before expiration
	return time.Until(expirationTime) < 20*time.Minute
}
