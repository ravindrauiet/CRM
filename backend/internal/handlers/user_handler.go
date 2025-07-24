package handlers

import (
	"encoding/json"
	"net/http"
	"maydiv-crm/internal/models"
	"maydiv-crm/internal/repository"
	"github.com/gorilla/sessions"
)

// UserHandler handles user-related requests
type UserHandler struct {
	userRepo *repository.UserRepository
	sessionStore *sessions.CookieStore
}

// NewUserHandler creates a new user handler
func NewUserHandler(userRepo *repository.UserRepository, sessionStore *sessions.CookieStore) *UserHandler {
	return &UserHandler{
		userRepo: userRepo,
		sessionStore: sessionStore,
	}
}

// HandleUsers handles user CRUD operations
func (h *UserHandler) HandleUsers(w http.ResponseWriter, r *http.Request) {
	// Check if user is admin or subadmin
	if !h.isAdminOrSubadmin(r) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	
	switch r.Method {
	case http.MethodGet:
		h.getUsers(w, r)
	case http.MethodPost:
		h.createUser(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// getUsers retrieves all users
func (h *UserHandler) getUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userRepo.GetAll()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	
	writeJSON(w, users)
}

// createUser creates a new user
func (h *UserHandler) createUser(w http.ResponseWriter, r *http.Request) {
	var userCreate struct {
		Username    string `json:"username"`
		Password    string `json:"password"`
		Designation string `json:"designation"`
		IsAdmin     bool   `json:"is_admin"`
		Role        string `json:"role"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&userCreate); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Set default role if not provided
	if userCreate.Role == "" {
		userCreate.Role = "stage1_employee"
	}
	
	// Create user with hashed password
	user := models.UserCreate{
		Username:     userCreate.Username,
		PasswordHash: userCreate.Password, // For now, store plain text (as per current system)
		Designation:  userCreate.Designation,
		IsAdmin:      userCreate.IsAdmin,
		Role:         userCreate.Role,
	}
	
	if err := h.userRepo.Create(&user); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	
	writeJSON(w, map[string]interface{}{"success": true})
}

// isAdmin checks if the current user is an admin
func (h *UserHandler) isAdmin(r *http.Request) bool {
	session, _ := h.sessionStore.Get(r, "session")
	isAdmin, ok := session.Values["is_admin"]
	if !ok {
		return false
	}
	
	adminStatus, ok := isAdmin.(bool)
	return ok && adminStatus
}

// isSubadmin checks if the current user is a subadmin
func (h *UserHandler) isSubadmin(r *http.Request) bool {
	session, _ := h.sessionStore.Get(r, "session")
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		return false
	}
	
	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		return false
	}
	
	return user.Role == "subadmin"
}

// isAdminOrSubadmin checks if the current user is an admin or subadmin
func (h *UserHandler) isAdminOrSubadmin(r *http.Request) bool {
	return h.isAdmin(r) || h.isSubadmin(r)
} 