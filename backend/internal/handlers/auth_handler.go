package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"maydiv-crm/internal/models"
	"maydiv-crm/internal/services"
	"github.com/gorilla/sessions"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	authService *services.AuthService
	sessionStore *sessions.CookieStore
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *services.AuthService, sessionStore *sessions.CookieStore) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		sessionStore: sessionStore,
	}
}

// Login handles user login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var credentials models.UserLogin
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	log.Println("Login attempt for username:", credentials.Username)
	
	user, err := h.authService.Authenticate(&credentials)
	if err != nil {
		log.Println("Authentication failed for user", credentials.Username, ":", err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	
	// Create session
	session, _ := h.sessionStore.Get(r, "session")
	session.Values["user_id"] = user.ID
	session.Values["is_admin"] = user.IsAdmin
	session.Save(r, w)
	
	log.Println("User", credentials.Username, "logged in successfully")
	
	response := map[string]interface{}{
		"success": true,
		"is_admin": user.IsAdmin,
	}
	
	writeJSON(w, response)
}

// Logout handles user logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := h.sessionStore.Get(r, "session")
	session.Options.MaxAge = -1
	session.Save(r, w)
	
	writeJSON(w, map[string]interface{}{"success": true})
}

// Helper function to write JSON responses
func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
} 