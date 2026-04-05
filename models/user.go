package models
import ("time")

type User struct {
	Id int `json:"id" db:"id"`
	Username string `json:"user_name" db:"user_name"`
	PasswordHash string `json:"-" db:"password_hash"`
	Role UserRole `json:"role" db:"role"`
	IsActive bool `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"creates_at" db:"created_at"`
}

type UserRole string

const (
    RoleAdmin  UserRole = "admin"
    RoleAnalyst UserRole = "analyst"
    RoleViewer UserRole= "viewer"            
)

