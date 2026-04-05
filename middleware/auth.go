package middleware

import (
  "context"
  "net/http"
  "strings"

  "github.com/o1egl/paseto"
)

type contextKey string

const (
  UserIDKey contextKey = "user_id"
  RoleKey   contextKey = "role"
)

func RequireAuth(pasetoKey string, next http.HandlerFunc) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    
    authHeader := r.Header.Get("Authorization")
    if authHeader == "" {
      http.Error(w, "Missing authorization header", http.StatusUnauthorized)
      return
    }

    
    parts := strings.Split(authHeader, " ")
    if len(parts) != 2 || parts[0] != "Bearer" {
      http.Error(w, "Invalid authorization format. Expected: Bearer <token>", http.StatusUnauthorized)
      return
    }
    tokenString := parts[1]

    
    pasetoMaker := paseto.NewV2()
    var claims map[string]interface{}
    
    err := pasetoMaker.Decrypt(tokenString, []byte(pasetoKey), &claims, nil)
    if err != nil {
      http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
      return
    }

    
    userIDFloat, okID := claims["user_id"].(float64)
    role, okRole := claims["role"].(string)

    if !okID || !okRole {
      http.Error(w, "Corrupted token payload", http.StatusUnauthorized)
      return
    }


    ctx := context.WithValue(r.Context(), UserIDKey, int(userIDFloat))
    ctx = context.WithValue(ctx, RoleKey, role)

    
    next.ServeHTTP(w, r.WithContext(ctx))
  }
}