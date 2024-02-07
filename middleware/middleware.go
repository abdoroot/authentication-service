package middleware

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/abdoroot/authentication-service/internal/auth"
)

func HttpLoginMiddleware(next http.HandlerFunc, srv *auth.Auth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-Api-Token")
		if token == "" {
			Unauthorized(w)
			return
		}
		claims, ok := auth.IsUserAuthorizedWithClaim(token)
		if !ok {
			Unauthorized(w)
			return
		}
		userId := claims["user_id"].(string)
		user, err := srv.Store.GetUserById(r.Context(), userId)
		if !ok {
			log.Println("GetUserById err:", err)
			return
		}
		ctxWithUser := context.WithValue(context.Background(), "user", user)
		request := r.WithContext(ctxWithUser)
		next.ServeHTTP(w, request)
	}
}

func Unauthorized(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	data, _ := json.Marshal(map[string]any{
		"err": "Unauthorized",
	})
	w.Write(data)
}
