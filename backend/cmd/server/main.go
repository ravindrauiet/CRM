package main

import (
	"log"
	"net/http"
	"os"
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
	
	// Initialize services
	authService := services.NewAuthService(userRepo)
	
	// Initialize session store
	sessionStore := sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
	
	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService, sessionStore)
	userHandler := handlers.NewUserHandler(userRepo, sessionStore)
	taskHandler := handlers.NewTaskHandler(taskRepo, sessionStore)
	
	// Setup routes
	mux := http.NewServeMux()
	
	// Auth routes
	mux.HandleFunc("/api/login", authHandler.Login)
	mux.HandleFunc("/api/logout", authHandler.Logout)
	
	// User routes
	mux.HandleFunc("/api/users", userHandler.HandleUsers)
	
	// Task routes
	mux.HandleFunc("/api/tasks", taskHandler.HandleTasks)
	mux.HandleFunc("/api/mytasks", taskHandler.HandleMyTasks)
	mux.HandleFunc("/api/tasks/", taskHandler.HandleTaskStatus)
	
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