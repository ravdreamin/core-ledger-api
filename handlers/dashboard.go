package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ravdreamin/core-ledger-api/middleware"
)

type DashboardSummary struct {
	TotalIncome  float64 `json:"total_income"`
	TotalExpense float64 `json:"total_expense"`
	Balance      float64 `json:"balance"`
}

func DashboardSummaryHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Get the authenticated user ID straight from the Middleware Bouncer
		userID := r.Context().Value(middleware.UserIDKey).(int)

		// 2. The Optimized SQL Query
		query := `
			SELECT 
				COALESCE(SUM(amount) FILTER (WHERE type = 'income'), 0) AS total_income,
				COALESCE(SUM(amount) FILTER (WHERE type = 'expense'), 0) AS total_expense
			FROM records 
			WHERE user_id = $1
		`

		var summary DashboardSummary

		// 3. Ask Postgres to do the math
		err := pool.QueryRow(r.Context(), query, userID).Scan(&summary.TotalIncome, &summary.TotalExpense)
		if err != nil {
			http.Error(w, "Failed to calculate dashboard summary", http.StatusInternalServerError)
			return
		}

		// 4. Calculate the final balance in Go's memory
		summary.Balance = summary.TotalIncome - summary.TotalExpense

		// 5. Fire the JSON back to the frontend
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(summary)
	}
}

type CategoryReport struct {
	CategoryName string  `json:"category_name"`
	TotalAmount  float64 `json:"total_amount"`
}

func CategoryReportHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(middleware.UserIDKey).(int)

		query := `
			SELECT c.name, COALESCE(SUM(r.amount), 0)
			FROM records r
			JOIN categories c ON r.category_id = c.id
			WHERE r.user_id = $1
			GROUP BY c.name
		`

		rows, err := pool.Query(r.Context(), query, userID)
		if err != nil {
			http.Error(w, "Failed to query category report", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		reports := make([]CategoryReport, 0)
		for rows.Next() {
			var report CategoryReport
			if err := rows.Scan(&report.CategoryName, &report.TotalAmount); err != nil {
				http.Error(w, "Failed to scan row", http.StatusInternalServerError)
				return
			}
			reports = append(reports, report)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(reports)
	}
}