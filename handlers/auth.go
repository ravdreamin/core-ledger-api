package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/o1egl/paseto"
	"golang.org/x/crypto/bcrypt"
)
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoginHandler(pool *pgxpool.Pool, pasetoKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		var storedHash string
		var storedRole string
		var storedID int

		err := pool.QueryRow(r.Context(),"SELECT id,role, password_hash FROM users WHERE username=$1",req.Username).Scan(&storedID,&storedRole,&storedHash)
		if err != nil {
    if errors.Is(err, pgx.ErrNoRows) {
        http.Error(w, "Invalid username or password", http.StatusUnauthorized)
        return 
    }
    http.Error(w, "Internal server error", http.StatusInternalServerError)
    return
}

err = bcrypt.CompareHashAndPassword([]byte(storedHash),[]byte(req.Password))
if err != nil{
	http.Error(w,"Invalid username or password", http.StatusUnauthorized)
	return
}

claims := map[string]interface{}{
    "user_id": storedID,
    "role":    storedRole,
    "exp":     time.Now().Add(24 * time.Hour), // The token dies in 24 hours
}


pasetoMaker := paseto.NewV2()
token, err := pasetoMaker.Encrypt([]byte(pasetoKey), claims, nil)
if err != nil {
    http.Error(w, "Failed to generate token", http.StatusInternalServerError)
    return
}

w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusOK)

err = json.NewEncoder(w).Encode(map[string]string{
    "token": token,
})
if err != nil {
    log.Printf("Failed to write response: %v", err)
}

	}
}