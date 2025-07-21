package main

import (
    "database/sql"
    "fmt"
    "log"
    "os"
    "path/filepath"
    "github.com/joho/godotenv"
    _ "github.com/go-sql-driver/mysql"
)

func main() {
    // Load environment variables from parent directory
    err := godotenv.Load("../.env")
    if err != nil {
        log.Println("No .env file found, using system environment variables")
    }

    // Connect to database
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASS"),
        os.Getenv("DB_HOST"),
        os.Getenv("DB_PORT"),
        os.Getenv("DB_NAME"),
    )
    
    log.Println("Connecting to database...")
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        log.Fatal("Error opening database:", err)
    }
    defer db.Close()

    if err := db.Ping(); err != nil {
        log.Fatal("Error connecting to database:", err)
    }
    log.Println("Successfully connected to database")

    // Read and execute schema
    schemaPath := filepath.Join("..", "schema.sql")
    schema, err := os.ReadFile(schemaPath)
    if err != nil {
        log.Fatal("Error reading schema.sql:", err)
    }

    log.Println("Creating tables and inserting sample data...")
    _, err = db.Exec(string(schema))
    if err != nil {
        log.Fatal("Error executing schema:", err)
    }

    log.Println("Database setup completed successfully!")
    log.Println("Sample users created:")
    log.Println("- admin (password: admin123) - System Administrator")
    log.Println("- john.doe (password: password123) - Software Developer")
    log.Println("- jane.smith (password: password123) - Project Manager")
    log.Println("- mike.wilson (password: password123) - QA Engineer")
} 