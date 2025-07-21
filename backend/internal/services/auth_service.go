package services

import (
	"maydiv-crm/internal/models"
	"maydiv-crm/internal/repository"
)

// AuthService handles authentication logic
type AuthService struct {
	userRepo *repository.UserRepository
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

// Authenticate validates user credentials
func (s *AuthService) Authenticate(credentials *models.UserLogin) (*models.User, error) {
	user, err := s.userRepo.GetByUsername(credentials.Username)
	if err != nil {
		return nil, err
	}
	
	// For now, using plaintext password comparison as requested
	if user.PasswordHash != credentials.Password {
		return nil, ErrInvalidCredentials
	}
	
	return user, nil
}

// GetUserByID retrieves a user by ID
func (s *AuthService) GetUserByID(id int) (*models.User, error) {
	return s.userRepo.GetByID(id)
}

// IsAdmin checks if a user is an admin
func (s *AuthService) IsAdmin(userID int) (bool, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return false, err
	}
	
	return user.IsAdmin, nil
} 