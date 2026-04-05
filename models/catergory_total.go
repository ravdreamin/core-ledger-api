package models

type CategoryTotal struct {
    Category    string `json:"category"`
    TotalAmount int64  `json:"total_amount"`
    Type        string `json:"type"`
}