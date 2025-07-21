package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"maydiv-crm/internal/models"
	"maydiv-crm/internal/repository"
	"github.com/gorilla/sessions"
)

// TaskHandler handles task-related requests
type TaskHandler struct {
	taskRepo *repository.TaskRepository
	sessionStore *sessions.CookieStore
}

// NewTaskHandler creates a new task handler
func NewTaskHandler(taskRepo *repository.TaskRepository, sessionStore *sessions.CookieStore) *TaskHandler {
	return &TaskHandler{
		taskRepo: taskRepo,
		sessionStore: sessionStore,
	}
}

// HandleTasks handles task CRUD operations
func (h *TaskHandler) HandleTasks(w http.ResponseWriter, r *http.Request) {
	// Check if user is admin for task management
	if !h.isAdmin(r) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	
	switch r.Method {
	case http.MethodGet:
		h.getTasks(w, r)
	case http.MethodPost:
		h.createTask(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// HandleMyTasks handles getting tasks assigned to the current user
func (h *TaskHandler) HandleMyTasks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	userID := h.getUserID(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	
	tasks, err := h.taskRepo.GetByUserID(userID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	
	writeJSON(w, tasks)
}

// HandleTaskStatus handles task status updates
func (h *TaskHandler) HandleTaskStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// Parse task ID from URL
	path := strings.TrimPrefix(r.URL.Path, "/api/tasks/")
	parts := strings.Split(path, "/")
	if len(parts) < 2 || parts[1] != "status" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	
	taskID, err := strconv.Atoi(parts[0])
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}
	
	userID := h.getUserID(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	
	var update models.TaskUpdateCreate
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	if err := h.taskRepo.UpdateStatus(taskID, userID, &update); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	
	writeJSON(w, map[string]interface{}{"success": true})
}

// getTasks retrieves all tasks
func (h *TaskHandler) getTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.taskRepo.GetAll()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	
	writeJSON(w, tasks)
}

// createTask creates a new task
func (h *TaskHandler) createTask(w http.ResponseWriter, r *http.Request) {
	var task models.TaskCreate
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	taskID, err := h.taskRepo.Create(&task)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	
	writeJSON(w, map[string]interface{}{"success": true, "task_id": taskID})
}

// isAdmin checks if the current user is an admin
func (h *TaskHandler) isAdmin(r *http.Request) bool {
	session, _ := h.sessionStore.Get(r, "session")
	isAdmin, ok := session.Values["is_admin"]
	if !ok {
		return false
	}
	
	adminStatus, ok := isAdmin.(bool)
	return ok && adminStatus
}

// getUserID gets the current user ID from session
func (h *TaskHandler) getUserID(r *http.Request) int {
	session, _ := h.sessionStore.Get(r, "session")
	userID, ok := session.Values["user_id"]
	if !ok {
		return 0
	}
	
	switch v := userID.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	default:
		return 0
	}
} 