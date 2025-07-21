package models

import "time"

// Task represents a task in the system
type Task struct {
	ID          int       `json:"id"`
	JobID       string    `json:"job_id"`
	Description string    `json:"description"`
	Priority    string    `json:"priority"`
	Deadline    time.Time `json:"deadline"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
}

// TaskCreate represents the data needed to create a new task
type TaskCreate struct {
	JobID       string `json:"job_id" validate:"required"`
	Description string `json:"description" validate:"required"`
	Priority    string `json:"priority" validate:"required,oneof=Low Medium High"`
	Deadline    string `json:"deadline" validate:"required"`
	AssignedTo  []int  `json:"assigned_to"`
}

// TaskResponse represents task data with additional information for frontend
type TaskResponse struct {
	ID          int       `json:"id"`
	JobID       string    `json:"job_id"`
	Description string    `json:"description"`
	Priority    string    `json:"priority"`
	Deadline    time.Time `json:"deadline"`
	AssignedTo  []string  `json:"assigned_to"`
	Status      string    `json:"status"`
}

// TaskAssignment represents the assignment of a task to a user
type TaskAssignment struct {
	ID        int       `json:"id"`
	TaskID    int       `json:"task_id"`
	UserID    int       `json:"user_id"`
	AssignedAt time.Time `json:"assigned_at"`
}

// TaskUpdate represents a status update or comment on a task
type TaskUpdate struct {
	ID        int       `json:"id"`
	TaskID    int       `json:"task_id"`
	UserID    int       `json:"user_id"`
	Status    string    `json:"status"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}

// TaskUpdateCreate represents the data needed to create a task update
type TaskUpdateCreate struct {
	Status  string `json:"status" validate:"required,oneof=Assigned In Progress Completed"`
	Comment string `json:"comment"`
} 