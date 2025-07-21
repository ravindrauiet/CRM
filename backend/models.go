package main

type User struct {
    ID          int    `json:"id"`
    Username    string `json:"username"`
    PasswordHash string `json:"-"`
    Designation string `json:"designation"`
    IsAdmin     bool   `json:"is_admin"`
}

type Task struct {
    ID          int    `json:"id"`
    JobID       string `json:"job_id"`
    Description string `json:"description"`
    Priority    string `json:"priority"`
    Deadline    string `json:"deadline"`
}

type TaskAssignment struct {
    ID     int `json:"id"`
    TaskID int `json:"task_id"`
    UserID int `json:"user_id"`
}

type TaskUpdate struct {
    ID        int    `json:"id"`
    TaskID    int    `json:"task_id"`
    UserID    int    `json:"user_id"`
    Status    string `json:"status"`
    Comment   string `json:"comment"`
    CreatedAt string `json:"created_at"`
} 