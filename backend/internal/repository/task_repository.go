package repository

import (
	"database/sql"
	"log"
	"maydiv-crm/internal/models"
	"strings"
	"time"
)

// TaskRepository handles task-related database operations
type TaskRepository struct {
	db *sql.DB
}

// NewTaskRepository creates a new task repository
func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

// GetAll retrieves all tasks with assigned users and latest status
func (r *TaskRepository) GetAll() ([]models.TaskResponse, error) {
	query := `
		SELECT t.id, t.job_id, t.description, t.priority, t.deadline,
		       GROUP_CONCAT(DISTINCT u.username) as assigned_to,
		       COALESCE(latest_status.status, 'Assigned') as status
		FROM tasks t
		LEFT JOIN task_assignments ta ON t.id = ta.task_id
		LEFT JOIN users u ON ta.user_id = u.id
		LEFT JOIN (
			SELECT tu1.task_id, tu1.status
			FROM task_updates tu1
			INNER JOIN (
				SELECT task_id, MAX(id) as max_id
				FROM task_updates
				GROUP BY task_id
			) tu2 ON tu1.task_id = tu2.task_id AND tu1.id = tu2.max_id
		) latest_status ON t.id = latest_status.task_id
		GROUP BY t.id, t.job_id, t.description, t.priority, t.deadline, latest_status.status
	`
	
	rows, err := r.db.Query(query)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil, err
	}
	defer rows.Close()
	
	var tasks []models.TaskResponse
	for rows.Next() {
		var task models.TaskResponse
		var assignedTo sql.NullString
		var status sql.NullString
		
		if err := rows.Scan(&task.ID, &task.JobID, &task.Description, &task.Priority, &task.Deadline, &assignedTo, &status); err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		
		// Parse assigned users
		if assignedTo.Valid && assignedTo.String != "" {
			task.AssignedTo = strings.Split(assignedTo.String, ",")
		} else {
			task.AssignedTo = []string{}
		}
		
		// Set status
		if status.Valid {
			task.Status = status.String
		} else {
			task.Status = "Assigned"
		}
		
		tasks = append(tasks, task)
	}
	
	return tasks, nil
}

// GetByUserID retrieves tasks assigned to a specific user
func (r *TaskRepository) GetByUserID(userID int) ([]models.TaskResponse, error) {
	query := `
		SELECT t.id, t.job_id, t.description, t.priority, t.deadline,
		       COALESCE(tu.status, 'Assigned') as status
		FROM tasks t
		JOIN task_assignments ta ON t.id = ta.task_id
		LEFT JOIN task_updates tu ON t.id = tu.task_id AND tu.user_id = ? 
			AND tu.id = (SELECT MAX(id) FROM task_updates WHERE task_id = t.id AND user_id = ?)
		WHERE ta.user_id = ?
	`
	
	rows, err := r.db.Query(query, userID, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var tasks []models.TaskResponse
	for rows.Next() {
		var task models.TaskResponse
		var status sql.NullString
		
		if err := rows.Scan(&task.ID, &task.JobID, &task.Description, &task.Priority, &task.Deadline, &status); err != nil {
			continue
		}
		
		if status.Valid {
			task.Status = status.String
		} else {
			task.Status = "Assigned"
		}
		
		task.AssignedTo = []string{} // Will be populated separately if needed
		tasks = append(tasks, task)
	}
	
	return tasks, nil
}

// Create creates a new task with assignments
func (r *TaskRepository) Create(task *models.TaskCreate) (int64, error) {
	// Parse deadline
	deadline, err := time.Parse("2006-01-02", task.Deadline)
	if err != nil {
		return 0, err
	}
	
	// Insert task
	result, err := r.db.Exec(
		"INSERT INTO tasks (job_id, description, priority, deadline) VALUES (?, ?, ?, ?)",
		task.JobID, task.Description, task.Priority, deadline,
	)
	if err != nil {
		return 0, err
	}
	
	taskID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	
	// Assign task to users
	for _, userID := range task.AssignedTo {
		_, err := r.db.Exec("INSERT INTO task_assignments (task_id, user_id) VALUES (?, ?)", taskID, userID)
		if err != nil {
			// Log error but continue with other assignments
			continue
		}
	}
	
	return taskID, nil
}

// UpdateStatus updates the status of a task for a specific user
func (r *TaskRepository) UpdateStatus(taskID int, userID int, update *models.TaskUpdateCreate) error {
	_, err := r.db.Exec(
		"INSERT INTO task_updates (task_id, user_id, status, comment) VALUES (?, ?, ?, ?)",
		taskID, userID, update.Status, update.Comment,
	)
	return err
}

// GetByID retrieves a task by ID
func (r *TaskRepository) GetByID(id int) (*models.Task, error) {
	task := &models.Task{}
	err := r.db.QueryRow(
		"SELECT id, job_id, description, priority, deadline FROM tasks WHERE id = ?",
		id,
	).Scan(&task.ID, &task.JobID, &task.Description, &task.Priority, &task.Deadline)
	
	if err != nil {
		return nil, err
	}
	
	return task, nil
}

// Delete deletes a task by ID
func (r *TaskRepository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM tasks WHERE id = ?", id)
	return err
} 