package main

import (
    "log"
    "net/http"
    "os"
    "github.com/joho/godotenv"
)

func LoadEnv() {
    err := godotenv.Load(".env")
    if err != nil {
        log.Println("No .env file found or error loading .env file")
    }
}

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

func main() {
    LoadEnv()
    db := InitDB()
    defer db.Close()
    sessionStore := InitSessionStore([]byte(os.Getenv("SESSION_KEY")))

    mux := http.NewServeMux()
    RegisterRoutes(mux, db, sessionStore)

    log.Println("Server started at :8080")
    log.Fatal(http.ListenAndServe(":8080", withCORS(mux)))
} 