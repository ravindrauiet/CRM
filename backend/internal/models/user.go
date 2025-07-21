package models

import "time"

// User represents a user in the system
type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"` // Don't expose password in JSON
	Designation  string    `json:"designation"`
	IsAdmin      bool      `json:"is_admin"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
}

// UserCreate represents the data needed to create a new user
type UserCreate struct {
	Username     string `json:"username" validate:"required"`
	PasswordHash string `json:"password_hash" validate:"required"`
	Designation  string `json:"designation" validate:"required"`
	IsAdmin      bool   `json:"is_admin"`
}

// UserLogin represents login credentials
type UserLogin struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// UserResponse represents the user data sent to frontend (without password)
type UserResponse struct {
	ID          int    `json:"id"`
	Username    string `json:"username"`
	Designation string `json:"designation"`
	IsAdmin     bool   `json:"is_admin"`
} 