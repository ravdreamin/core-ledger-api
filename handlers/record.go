package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ravdreamin/core-ledger-api/middleware"
)


type CreateRecordRequest struct {
	UserID      int     `json:"user_id"`
	CategoryID  int     `json:"category_id"`
	Amount      float64 `json:"amount"`
	Type        string  `json:"type"`
	Description string  `json:"description"`
}

func CreateRecordHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		
		var req CreateRecordRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		
		realUserID := r.Context().Value(middleware.UserIDKey).(int)
		req.UserID = realUserID

		
		_, err := pool.Exec(r.Context(),
			"INSERT INTO records (user_id, category_id, amount, type, description, date) VALUES ($1, $2, $3, $4, $5, NOW())",
			req.UserID, req.CategoryID, req.Amount, req.Type, req.Description,
		)

		if err != nil {
			http.Error(w, "Failed to insert record into database", http.StatusInternalServerError)
			return
		}

		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated) 
		
		err = json.NewEncoder(w).Encode(map[string]string{
			"message": "Record created successfully",
		})
		if err != nil {
			
			return
		}
	}
}

type RecordResponse struct {
	ID          int     `json:"id"`
	Amount      float64 `json:"amount"`
	Type        string  `json:"type"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
}

func GetRecordsHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(middleware.UserIDKey).(int)

		query := `
			SELECT r.id, r.amount, r.type, c.name as category, r.description
			FROM records r
			JOIN categories c ON r.category_id = c.id
			WHERE r.user_id = $1
			ORDER BY r.date DESC
		`

		rows, err := pool.Query(r.Context(), query, userID)
		if err != nil {
			http.Error(w, "Failed to fetch records", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		records := make([]RecordResponse, 0)
		for rows.Next() {
			var rec RecordResponse
			if err := rows.Scan(&rec.ID, &rec.Amount, &rec.Type, &rec.Category, &rec.Description); err != nil {
				continue
			}
			records = append(records, rec)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(records)
	}
}