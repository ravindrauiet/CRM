package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"maydiv-crm/internal/database"
	"maydiv-crm/internal/handlers"
	"maydiv-crm/internal/repository"
	"maydiv-crm/internal/services"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	// Try multiple possible paths for .env file
	envPaths := []string{
		".env",                    // Current directory
		"../.env",                 // Parent directory
		"../../.env",              // Two levels up
		"../../../.env",           // Three levels up
		"../../../../.env",        // Four levels up
	}
	
	var envLoaded bool
	for _, path := range envPaths {
		if err := godotenv.Load(path); err == nil {
			log.Printf("Loaded .env file from: %s", path)
			envLoaded = true
			break
		}
	}
	
	if !envLoaded {
		log.Println("No .env file found, using system environment variables")
	}
	
	// Initialize database
	db, err := database.NewConnection()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()
	
	// Run migrations and seed data
	if err := db.Migrate(); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}
	
	if err := db.Seed(); err != nil {
		log.Fatal("Failed to seed database:", err)
	}
	
	// Initialize repositories
	userRepo := repository.NewUserRepository(db.DB)
	taskRepo := repository.NewTaskRepository(db.DB)
	pipelineRepo := repository.NewPipelineRepository(db.DB)
	
	// Initialize services
	authService := services.NewAuthService(userRepo)
	notificationService := services.NewNotificationService(db.DB)
	
	// Initialize session store
	sessionKey := os.Getenv("SESSION_KEY")
	if sessionKey == "" {
		sessionKey = "default-session-key-change-in-production"
		log.Println("Warning: Using default session key. Set SESSION_KEY in .env for production.")
	}
	sessionStore := sessions.NewCookieStore([]byte(sessionKey))
	
	// Configure session store for better security and compatibility
	sessionStore.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
	}
	
	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService, sessionStore)
	userHandler := handlers.NewUserHandler(userRepo, sessionStore)
	taskHandler := handlers.NewTaskHandler(taskRepo, sessionStore)
	pipelineHandler := handlers.NewPipelineHandler(pipelineRepo, userRepo, sessionStore, notificationService)
	
	// Setup routes
	mux := http.NewServeMux()
	
	// Auth routes
	mux.HandleFunc("/api/login", authHandler.Login)
	mux.HandleFunc("/api/logout", authHandler.Logout)
	
	// Session check endpoint
	mux.HandleFunc("/api/session", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		
		session, err := sessionStore.Get(r, "session")
		if err != nil {
			http.Error(w, "Session error", http.StatusUnauthorized)
			return
		}
		
		userID, ok := session.Values["user_id"].(int)
		if !ok {
			http.Error(w, "Not authenticated", http.StatusUnauthorized)
			return
		}
		
		// Get user details
		user, err := userRepo.GetByID(userID)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		
		response := map[string]interface{}{
			"authenticated": true,
			"user_id": userID,
			"username": user.Username,
			"is_admin": user.IsAdmin,
			"role": user.Role,
		}
		
		json.NewEncoder(w).Encode(response)
	})
	
	// User routes
	mux.HandleFunc("/api/users", userHandler.HandleUsers)
	
	// Legacy Task routes (keeping for backward compatibility)
	mux.HandleFunc("/api/tasks", taskHandler.HandleTasks)
	mux.HandleFunc("/api/mytasks", taskHandler.HandleMyTasks)
	mux.HandleFunc("/api/tasks/", taskHandler.HandleTaskStatus)
	
	// New Pipeline routes
	mux.HandleFunc("/api/pipeline/jobs", pipelineHandler.HandleJobs)
	mux.HandleFunc("/api/pipeline/myjobs", pipelineHandler.HandleMyJobs)
	mux.HandleFunc("/api/debug", pipelineHandler.HandleDebug)
	
	// Test endpoint
	mux.HandleFunc("/api/test", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "Test endpoint working"}`))
	})
	
	// Debug endpoint to check stage1 data
	mux.HandleFunc("/api/debug/stage1", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		
		// Get all stage1 data
		rows, err := db.DB.Query("SELECT * FROM stage1_data")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		
		var results []map[string]interface{}
		columns, _ := rows.Columns()
		count := len(columns)
		values := make([]interface{}, count)
		valuePtrs := make([]interface{}, count)
		
		for i := range columns {
			valuePtrs[i] = &values[i]
		}
		
		for rows.Next() {
			err := rows.Scan(valuePtrs...)
			if err != nil {
				continue
			}
			
			row := make(map[string]interface{})
			for i, col := range columns {
				val := values[i]
				row[col] = val
			}
			results = append(results, row)
		}
		
		json.NewEncoder(w).Encode(map[string]interface{}{
			"stage1_data": results,
			"count": len(results),
		})
	})
	
		// Debug endpoint to check stage2 data
	mux.HandleFunc("/api/debug/stage2", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Get all stage2 data
		rows, err := db.DB.Query("SELECT * FROM stage2_data")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var results []map[string]interface{}
		columns, _ := rows.Columns()
		count := len(columns)
		values := make([]interface{}, count)
		valuePtrs := make([]interface{}, count)

		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		for rows.Next() {
			err := rows.Scan(valuePtrs...)
			if err != nil {
				continue
			}

			row := make(map[string]interface{})
			for i, col := range columns {
				val := values[i]
				row[col] = val
			}
			results = append(results, row)
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"stage2_data": results,
			"count": len(results),
		})
	})

	// Debug endpoint to check stage3 data
	mux.HandleFunc("/api/debug/stage3", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Get all stage3 data
		rows, err := db.DB.Query("SELECT * FROM stage3_data")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var results []map[string]interface{}
		columns, _ := rows.Columns()
		count := len(columns)
		values := make([]interface{}, count)
		valuePtrs := make([]interface{}, count)

		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		for rows.Next() {
			err := rows.Scan(valuePtrs...)
			if err != nil {
				continue
			}

			row := make(map[string]interface{})
			for i, col := range columns {
				val := values[i]
				row[col] = val
			}
			results = append(results, row)
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"stage3_data": results,
			"count": len(results),
		})
	})

	// Debug endpoint to check stage4 data
	mux.HandleFunc("/api/debug/stage4", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Get all stage4 data
		rows, err := db.DB.Query("SELECT * FROM stage4_data")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var results []map[string]interface{}
		columns, _ := rows.Columns()
		count := len(columns)
		values := make([]interface{}, count)
		valuePtrs := make([]interface{}, count)

		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		for rows.Next() {
			err := rows.Scan(valuePtrs...)
			if err != nil {
				continue
			}

			row := make(map[string]interface{})
			for i, col := range columns {
				val := values[i]
				row[col] = val
			}
			results = append(results, row)
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"stage4_data": results,
			"count": len(results),
		})
	})
	
	// Debug endpoint to check job status
	mux.HandleFunc("/api/debug/jobs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		
		// Get all pipeline jobs with their current stage
		rows, err := db.DB.Query(`
			SELECT pj.id, pj.job_no, pj.current_stage, pj.status, 
			       s1.consignee, s1.commodity,
			       s2.hsn_code, s2.filing_requirement
			FROM pipeline_jobs pj
			LEFT JOIN stage1_data s1 ON pj.id = s1.job_id
			LEFT JOIN stage2_data s2 ON pj.id = s2.job_id
			ORDER BY pj.id
		`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		
		var results []map[string]interface{}
		for rows.Next() {
			var id int
			var jobNo, currentStage, status, consignee, commodity, hsnCode, filingRequirement sql.NullString
			
			err := rows.Scan(&id, &jobNo, &currentStage, &status, &consignee, &commodity, &hsnCode, &filingRequirement)
			if err != nil {
				continue
			}
			
			row := map[string]interface{}{
				"id": id,
				"job_no": jobNo.String,
				"current_stage": currentStage.String,
				"status": status.String,
				"consignee": consignee.String,
				"commodity": commodity.String,
				"hsn_code": hsnCode.String,
				"filing_requirement": filingRequirement.String,
			}
			results = append(results, row)
		}
		
		json.NewEncoder(w).Encode(map[string]interface{}{
			"jobs": results,
			"count": len(results),
		})
	})
	
	// Test email endpoint
	mux.HandleFunc("/api/test-email", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		
		w.Header().Set("Content-Type", "application/json")
		
		// Test email connection
		if err := notificationService.EmailService.TestEmailConnection(); err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			})
			return
		}
		
		// Send a test email
		adminEmail := os.Getenv("ADMIN_EMAIL")
		if adminEmail == "" {
			adminEmail = "admin@maydiv.com"
		}
		
		emailData := services.StageCompletionEmail{
			JobNo:       "TEST-001",
			JobTitle:    "Test Job",
			Stage:       "stage2",
			StageName:   "Stage 2 - Customs & Documentation",
			CompletedBy: "Test User",
			CompletedAt: "2024-01-20 10:30:00",
			NextStage:   "Stage 3 - Clearance & Logistics",
			AdminEmail:  adminEmail,
		}
		
		if err := notificationService.EmailService.SendStageCompletionEmail(emailData); err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			})
			return
		}
		
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "Test email sent successfully",
		})
	})
	
	// Handle all pipeline job routes with ID
	mux.HandleFunc("/api/pipeline/jobs/", func(w http.ResponseWriter, r *http.Request) {
		// Route to specific job handlers based on path
		path := r.URL.Path
		fmt.Printf("Pipeline job route accessed: %s\n", path)
		
		if strings.Contains(path, "/stage2") {
			pipelineHandler.HandleStage2Update(w, r)
		} else if strings.Contains(path, "/stage3") {
			pipelineHandler.HandleStage3Update(w, r)
		} else if strings.Contains(path, "/stage4") {
			pipelineHandler.HandleStage4Update(w, r)
		} else {
			// Default to job details
			pipelineHandler.HandleJobByID(w, r)
		}
	})
	
	// CORS middleware
	handler := withCORS(mux)
	
	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

// CORS middleware
func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
} 