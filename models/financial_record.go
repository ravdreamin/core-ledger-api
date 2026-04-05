package models
import (
	"time"

)

type FinancialRecord struct {
	Id int `json:"id" db:"id"`
	UserID int `json:"user_id" db:"user_id"`
	Amount int64 `json:"amount" db:"amount"`
	Type  FinancialType `json:"type" db:"type"`
	Category string `json:"category" db:"category"`
	Date time.Time `json:"created_at" db:"created_at"`
	Description string `json:"description" db:"description"`
}

type FinancialType string 
type FinancialCategory string

const(
	TypeExpense FinancialType = "expense"
	TypeIncome FinancialType = "income"
)




