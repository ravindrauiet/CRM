package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"maydiv-crm/internal/models"
	"maydiv-crm/internal/repository"

	"github.com/gorilla/sessions"
)

type PipelineHandler struct {
	pipelineRepo *repository.PipelineRepository
	userRepo     *repository.UserRepository
	sessionStore *sessions.CookieStore
}

func NewPipelineHandler(pipelineRepo *repository.PipelineRepository, userRepo *repository.UserRepository, sessionStore *sessions.CookieStore) *PipelineHandler {
	return &PipelineHandler{
		pipelineRepo: pipelineRepo,
		userRepo:     userRepo,
		sessionStore: sessionStore,
	}
}

// HandleJobs handles GET /api/pipeline/jobs (admin only - all jobs) and POST /api/pipeline/jobs (admin only - create job)
func (h *PipelineHandler) HandleJobs(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.getAllJobs(w, r)
	} else if r.Method == http.MethodPost {
		h.createJob(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// HandleMyJobs handles GET /api/pipeline/myjobs - gets jobs assigned to current user based on their role
func (h *PipelineHandler) HandleMyJobs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := h.getUserID(r)
	fmt.Printf("HandleMyJobs - getUserID returned: %d\n", userID)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get user role
	fmt.Printf("Looking up user with ID: %d\n", userID)
	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		fmt.Printf("User not found error: %v\n", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	fmt.Printf("Found user: %s (ID: %d, Role: %s)\n", user.Username, user.ID, user.Role)

	// Admin gets all jobs
	if user.IsAdmin {
		jobs, err := h.pipelineRepo.GetAllJobs()
		if err != nil {
			http.Error(w, "Failed to fetch jobs", http.StatusInternalServerError)
			return
		}
		writeJSON(w, jobs)
		return
	}

	// Get jobs based on user role
	fmt.Printf("Getting jobs for user ID %d with role %s\n", userID, user.Role)
	jobs, err := h.pipelineRepo.GetJobsByUserRole(userID, user.Role)
	if err != nil {
		fmt.Printf("Error getting jobs: %v\n", err)
		http.Error(w, "Failed to fetch jobs", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Found %d jobs for user\n", len(jobs))
	writeJSON(w, jobs)
}

// Debug endpoint to check database state
func (h *PipelineHandler) HandleDebug(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get all users
	users, err := h.userRepo.GetAll()
	if err != nil {
		http.Error(w, "Failed to get users", http.StatusInternalServerError)
		return
	}

	// Get all jobs
	jobs, err := h.pipelineRepo.GetAllJobs()
	if err != nil {
		http.Error(w, "Failed to get jobs", http.StatusInternalServerError)
		return
	}

	debugInfo := map[string]interface{}{
		"users": users,
		"jobs":  jobs,
	}

	writeJSON(w, debugInfo)
}

// HandleJobByID handles GET /api/pipeline/jobs/{id} - get specific job details
func (h *PipelineHandler) HandleJobByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := h.getUserID(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Extract job ID from URL
	pathParts := strings.Split(r.URL.Path, "/")
	fmt.Printf("URL Path: %s\n", r.URL.Path)
	fmt.Printf("Path parts: %v\n", pathParts)
	fmt.Printf("Path parts length: %d\n", len(pathParts))
	
	if len(pathParts) < 5 {
		fmt.Printf("Invalid path length: %d\n", len(pathParts))
		http.Error(w, "Invalid job ID", http.StatusBadRequest)
		return
	}

	fmt.Printf("Job ID from path: %s\n", pathParts[4])
	jobID, err := strconv.Atoi(pathParts[4])
	if err != nil {
		fmt.Printf("Error converting job ID: %v\n", err)
		http.Error(w, "Invalid job ID", http.StatusBadRequest)
		return
	}
	fmt.Printf("Parsed job ID: %d\n", jobID)

	job, err := h.pipelineRepo.GetJobByID(jobID)
	if err != nil {
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	// Check if user has access to this job
	if !h.hasJobAccess(userID, job) {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	writeJSON(w, job)
}

// HandleStage2Update handles PUT /api/pipeline/jobs/{id}/stage2
func (h *PipelineHandler) HandleStage2Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := h.getUserID(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check if user is stage2 employee or admin
	user, err := h.userRepo.GetByID(userID)
	if err != nil || (!user.IsAdmin && user.Role != "stage2_employee") {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	jobID, err := h.extractJobID(r)
	if err != nil {
		http.Error(w, "Invalid job ID", http.StatusBadRequest)
		return
	}

	var req models.Stage2UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err = h.pipelineRepo.UpdateStage2Data(jobID, &req, userID)
	if err != nil {
		http.Error(w, "Failed to update stage 2 data", http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{"message": "Stage 2 data updated successfully"})
}

// HandleStage3Update handles PUT /api/pipeline/jobs/{id}/stage3
func (h *PipelineHandler) HandleStage3Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := h.getUserID(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check if user is stage3 employee or admin
	user, err := h.userRepo.GetByID(userID)
	if err != nil || (!user.IsAdmin && user.Role != "stage3_employee") {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	jobID, err := h.extractJobID(r)
	if err != nil {
		http.Error(w, "Invalid job ID", http.StatusBadRequest)
		return
	}

	var req models.Stage3UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err = h.pipelineRepo.UpdateStage3Data(jobID, &req, userID)
	if err != nil {
		http.Error(w, "Failed to update stage 3 data", http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{"message": "Stage 3 data updated successfully"})
}

// HandleStage4Update handles PUT /api/pipeline/jobs/{id}/stage4
func (h *PipelineHandler) HandleStage4Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := h.getUserID(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check if user is customer or admin
	user, err := h.userRepo.GetByID(userID)
	if err != nil || (!user.IsAdmin && user.Role != "customer") {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	jobID, err := h.extractJobID(r)
	if err != nil {
		http.Error(w, "Invalid job ID", http.StatusBadRequest)
		return
	}

	var req models.Stage4UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err = h.pipelineRepo.UpdateStage4Data(jobID, &req, userID)
	if err != nil {
		http.Error(w, "Failed to update stage 4 data", http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{"message": "Stage 4 data updated successfully"})
}

// Private helper methods

func (h *PipelineHandler) getAllJobs(w http.ResponseWriter, r *http.Request) {
	if !h.isAdmin(r) {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	jobs, err := h.pipelineRepo.GetAllJobs()
	if err != nil {
		http.Error(w, "Failed to fetch jobs", http.StatusInternalServerError)
		return
	}

	writeJSON(w, jobs)
}

func (h *PipelineHandler) createJob(w http.ResponseWriter, r *http.Request) {
	if !h.isAdmin(r) {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	userID := h.getUserID(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

		var req models.Stage1CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Printf("JSON decode error: %v\n", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	// Validate required fields
	if req.JobNo == "" {
		http.Error(w, "Job number is required", http.StatusBadRequest)
		return
	}

	job, err := h.pipelineRepo.CreateJob(&req, userID)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			http.Error(w, "Job number already exists", http.StatusConflict)
			return
		}
		http.Error(w, "Failed to create job", http.StatusInternalServerError)
		return
	}

	writeJSON(w, job)
}

func (h *PipelineHandler) extractJobID(r *http.Request) (int, error) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 5 {
		return 0, fmt.Errorf("invalid path")
	}

	return strconv.Atoi(pathParts[4])
}

func (h *PipelineHandler) hasJobAccess(userID int, job *models.PipelineJobResponse) bool {
	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		return false
	}

	// Admin has access to all jobs
	if user.IsAdmin {
		return true
	}

	// Check role-based access
	switch user.Role {
	case "stage2_employee":
		return job.AssignedToStage2 != nil && *job.AssignedToStage2 == userID
	case "stage3_employee":
		return job.AssignedToStage3 != nil && *job.AssignedToStage3 == userID
	case "customer":
		return job.CustomerID != nil && *job.CustomerID == userID
	default:
		return false
	}
}

func (h *PipelineHandler) isAdmin(r *http.Request) bool {
	session, err := h.sessionStore.Get(r, "session")
	if err != nil {
		return false
	}

	userID, ok := session.Values["user_id"].(int)
	if !ok {
		return false
	}

	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		return false
	}

	return user.IsAdmin
}

func (h *PipelineHandler) getUserID(r *http.Request) int {
	session, err := h.sessionStore.Get(r, "session")
	if err != nil {
		return 0
	}

	userID, ok := session.Values["user_id"].(int)
	if !ok {
		return 0
	}

	return userID
} 