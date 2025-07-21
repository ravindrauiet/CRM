package database

import (
	"log"
)

// Migrate runs database migrations
func (db *DB) Migrate() error {
	log.Println("Running database migrations...")
	
	// Create tables using your existing schema - execute each statement separately
	statements := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INT AUTO_INCREMENT PRIMARY KEY,
			username VARCHAR(50) NOT NULL UNIQUE,
			password_hash VARCHAR(255) NOT NULL,
			designation VARCHAR(50) NOT NULL,
			is_admin BOOLEAN DEFAULT FALSE
		)`,
		
		`CREATE TABLE IF NOT EXISTS tasks (
			id INT AUTO_INCREMENT PRIMARY KEY,
			job_id VARCHAR(50) NOT NULL,
			description TEXT NOT NULL,
			priority ENUM('Low', 'Medium', 'High') NOT NULL,
			deadline DATE NOT NULL
		)`,
		
		`CREATE TABLE IF NOT EXISTS task_assignments (
			id INT AUTO_INCREMENT PRIMARY KEY,
			task_id INT NOT NULL,
			user_id INT NOT NULL,
			FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
		
		`CREATE TABLE IF NOT EXISTS task_updates (
			id INT AUTO_INCREMENT PRIMARY KEY,
			task_id INT NOT NULL,
			user_id INT NOT NULL,
			status ENUM('Assigned', 'In Progress', 'Completed') NOT NULL,
			comment TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
	}
	
	for _, statement := range statements {
		_, err := db.Exec(statement)
		if err != nil {
			return err
		}
	}
	
	log.Println("Database migrations completed successfully")
	return nil
}

// Seed inserts initial data if tables are empty
func (db *DB) Seed() error {
	// Check if users table is empty
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return err
	}
	
	if count > 0 {
		log.Println("Database already has data, skipping seed")
		return nil
	}
	
	log.Println("Seeding database with initial data...")
	
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
	
	log.Println("Database seeding completed successfully")
	return nil
} 