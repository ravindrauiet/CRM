package main

import (
    "database/sql"
    "fmt"
    "log"
    "os"
    _ "github.com/go-sql-driver/mysql"
)

func InitDB() *sql.DB {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASS"),
        os.Getenv("DB_HOST"),
        os.Getenv("DB_PORT"),
        os.Getenv("DB_NAME"),
    )
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        panic(err)
    }
    if err := db.Ping(); err != nil {
        panic(err)
    }
    
    // Check if tables exist, if not create them
    ensureTablesExist(db)
    
    return db
}

func ensureTablesExist(db *sql.DB) {
    // Check if users table exists
    var tableName string
    err := db.QueryRow("SHOW TABLES LIKE 'users'").Scan(&tableName)
    if err != nil {
        log.Println("Tables don't exist, creating them...")
        createTables(db)
    } else {
        log.Println("Database tables already exist")
    }
}

func createTables(db *sql.DB) {
    schema := `
    -- Users table
    CREATE TABLE IF NOT EXISTS users (
        id INT AUTO_INCREMENT PRIMARY KEY,
        username VARCHAR(50) UNIQUE NOT NULL,
        password_hash VARCHAR(255) NOT NULL,
        designation VARCHAR(100) NOT NULL,
        is_admin BOOLEAN DEFAULT FALSE,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

    -- Tasks table
    CREATE TABLE IF NOT EXISTS tasks (
        id INT AUTO_INCREMENT PRIMARY KEY,
        job_id VARCHAR(50) NOT NULL,
        description TEXT NOT NULL,
        priority ENUM('Low', 'Medium', 'High', 'Critical') DEFAULT 'Medium',
        deadline DATE NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

    -- Task assignments table (many-to-many relationship)
    CREATE TABLE IF NOT EXISTS task_assignments (
        id INT AUTO_INCREMENT PRIMARY KEY,
        task_id INT NOT NULL,
        user_id INT NOT NULL,
        assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
        UNIQUE KEY unique_task_user (task_id, user_id)
    );

    -- Task updates table (for status updates and comments)
    CREATE TABLE IF NOT EXISTS task_updates (
        id INT AUTO_INCREMENT PRIMARY KEY,
        task_id INT NOT NULL,
        user_id INT NOT NULL,
        status ENUM('Assigned', 'In Progress', 'Completed', 'On Hold', 'Cancelled') DEFAULT 'Assigned',
        comment TEXT,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    );
    `
    
    _, err := db.Exec(schema)
    if err != nil {
        log.Printf("Error creating tables: %v", err)
        return
    }
    
    // Insert sample data if tables are empty
    var count int
    err = db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
    if err != nil || count == 0 {
        insertSampleData(db)
    }
}

func insertSampleData(db *sql.DB) {
    log.Println("Inserting sample data...")
    
    // Insert sample users
    users := []struct {
        username, password, designation string
        isAdmin                        bool
    }{
        {"admin", "admin123", "System Administrator", true},
        {"john.doe", "password123", "Software Developer", false},
        {"jane.smith", "password123", "Project Manager", false},
        {"mike.wilson", "password123", "QA Engineer", false},
    }
    
    for _, user := range users {
        _, err := db.Exec("INSERT INTO users (username, password_hash, designation, is_admin) VALUES (?, ?, ?, ?)",
            user.username, user.password, user.designation, user.isAdmin)
        if err != nil {
            log.Printf("Error inserting user %s: %v", user.username, err)
        }
    }
    
    // Insert sample tasks
    tasks := []struct {
        jobID, description, priority, deadline string
    }{
        {"TASK-001", "Implement user authentication system", "High", "2024-02-15"},
        {"TASK-002", "Design database schema for CRM", "Medium", "2024-02-20"},
        {"TASK-003", "Create responsive dashboard UI", "High", "2024-02-25"},
        {"TASK-004", "Write API documentation", "Low", "2024-03-01"},
    }
    
    for _, task := range tasks {
        result, err := db.Exec("INSERT INTO tasks (job_id, description, priority, deadline) VALUES (?, ?, ?, ?)",
            task.jobID, task.description, task.priority, task.deadline)
        if err != nil {
            log.Printf("Error inserting task %s: %v", task.jobID, err)
            continue
        }
        
        taskID, _ := result.LastInsertId()
        
        // Assign tasks to users
        if task.jobID == "TASK-001" {
            db.Exec("INSERT INTO task_assignments (task_id, user_id) VALUES (?, ?)", taskID, 2) // john.doe
            db.Exec("INSERT INTO task_assignments (task_id, user_id) VALUES (?, ?)", taskID, 3) // jane.smith
        } else if task.jobID == "TASK-002" {
            db.Exec("INSERT INTO task_assignments (task_id, user_id) VALUES (?, ?)", taskID, 2) // john.doe
        } else if task.jobID == "TASK-003" {
            db.Exec("INSERT INTO task_assignments (task_id, user_id) VALUES (?, ?)", taskID, 3) // jane.smith
        } else if task.jobID == "TASK-004" {
            db.Exec("INSERT INTO task_assignments (task_id, user_id) VALUES (?, ?)", taskID, 4) // mike.wilson
        }
    }
    
    log.Println("Sample data inserted successfully")
} 