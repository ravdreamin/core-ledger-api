package middleware

import (
	"net/http"
)

func RequireRole(next http.HandlerFunc, allowedRoles ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userRole := r.Context().Value(RoleKey).(string)
		
		authorized := false
		for _, role := range allowedRoles {
			if userRole == role {
				authorized = true
				break
			}
		}

		if !authorized {
			http.Error(w, "Permission denied", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	}
}