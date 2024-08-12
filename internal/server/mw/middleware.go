package mw

import (
	"Marketplace/internal/utils"
	"context"
	"net/http"
	"strings"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization Header", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		//http.Error(w, parts[0]+" is not a Bearer", http.StatusUnauthorized)
		//http.Error(w, "Invalid Authorization Header", http.StatusUnauthorized)
		//return

		tokenString := parts[0]
		userID, err := utils.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid Token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// Сохраняем userID в контексте для последующего использования в других обработчиках
		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
