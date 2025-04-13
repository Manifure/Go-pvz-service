package auth

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const RoleKey = contextKey("role")

func AuthMiddleware(next http.HandlerFunc, allowedRoles ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenStr := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		claims, err := ParseJWT(tokenStr)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		role := claims.Role
		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				ctx := context.WithValue(r.Context(), "role", role)
				next(w, r.WithContext(ctx))
				return
			}
		}
		http.Error(w, "forbidden", http.StatusForbidden)

	}
}
