package repository

import (
	"database/sql"
	"maydiv-crm/internal/models"
)

// UserRepository handles user-related database operations
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// GetByUsername retrieves a user by username
func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRow(
		"SELECT id, username, password_hash, designation, is_admin FROM users WHERE username = ?",
		username,
	).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Designation, &user.IsAdmin)
	
	if err != nil {
		return nil, err
	}
	
	return user, nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(id int) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRow(
		"SELECT id, username, password_hash, designation, is_admin FROM users WHERE id = ?",
		id,
	).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Designation, &user.IsAdmin)
	
	if err != nil {
		return nil, err
	}
	
	return user, nil
}

// GetAll retrieves all users
func (r *UserRepository) GetAll() ([]models.UserResponse, error) {
	rows, err := r.db.Query("SELECT id, username, designation, is_admin FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var users []models.UserResponse
	for rows.Next() {
		var user models.UserResponse
		if err := rows.Scan(&user.ID, &user.Username, &user.Designation, &user.IsAdmin); err != nil {
			continue
		}
		users = append(users, user)
	}
	
	return users, nil
}

// Create creates a new user
func (r *UserRepository) Create(user *models.UserCreate) error {
	_, err := r.db.Exec(
		"INSERT INTO users (username, password_hash, designation, is_admin) VALUES (?, ?, ?, ?)",
		user.Username, user.PasswordHash, user.Designation, user.IsAdmin,
	)
	return err
}

// Update updates an existing user
func (r *UserRepository) Update(id int, user *models.UserCreate) error {
	_, err := r.db.Exec(
		"UPDATE users SET username = ?, password_hash = ?, designation = ?, is_admin = ? WHERE id = ?",
		user.Username, user.PasswordHash, user.Designation, user.IsAdmin, id,
	)
	return err
}

// Delete deletes a user by ID
func (r *UserRepository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM users WHERE id = ?", id)
	return err
} 