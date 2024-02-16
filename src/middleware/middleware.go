package middleware

import (
	"blog-server/src/auth"
	"context"
	"net/http"
	"strings"
)

// AuthMiddleware - Middleware для проверки JWT-токена
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

		claims, err := auth.VerifyToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Добавляем данные о пользователе в контекст запроса
		ctx := r.Context()
		ctx = context.WithValue(ctx, "username", claims.Username)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
