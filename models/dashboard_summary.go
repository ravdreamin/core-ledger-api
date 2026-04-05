package models

type DashboardSummary struct {
    TotalIncome   int64 `json:"total_income"`
    TotalExpenses int64 `json:"total_expenses"`
    NetBalance    int64 `json:"net_balance"`
}