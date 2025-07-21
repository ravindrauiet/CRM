package main

import (
    "database/sql"
    "encoding/json"
    "net/http"
    "golang.org/x/crypto/bcrypt"
    "github.com/gorilla/sessions"
    "log"
    "strings"
)

// Helper: write JSON
func writeJSON(w http.ResponseWriter, v interface{}) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(v)
}

// Helper: hash password
func hashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

// Helper: check password
func checkPassword(hash, password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

// Register all routes
func RegisterRoutes(mux *http.ServeMux, db *sql.DB, store *sessions.CookieStore) {
    // Helper function to check if user is admin
    isAdmin := func(r *http.Request) (bool, int) {
        session, _ := store.Get(r, "session")
        userID, ok := session.Values["user_id"]
        isAdminVal, adminOk := session.Values["is_admin"]
        if !ok || !adminOk {
            return false, 0
        }
        var uid int
        switch v := userID.(type) {
        case int:
            uid = v
        case int64:
            uid = int(v)
        case float64:
            uid = int(v)
        default:
            return false, 0
        }
        return isAdminVal.(bool), uid
    }

    // Helper function to get current user ID
    getUserID := func(r *http.Request) (int, bool) {
        session, _ := store.Get(r, "session")
        userID, ok := session.Values["user_id"]
        if !ok {
            return 0, false
        }
        var uid int
        switch v := userID.(type) {
        case int:
            uid = v
        case int64:
            uid = int(v)
        case float64:
            uid = int(v)
        default:
            return 0, false
        }
        return uid, true
    }

    mux.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            w.WriteHeader(http.StatusMethodNotAllowed)
            return
        }
        var creds struct {
            Username string `json:"username"`
            Password string `json:"password"`
        }
        if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
            w.WriteHeader(http.StatusBadRequest)
            return
        }
        log.Println("Login attempt for username:", creds.Username)
        var user User
        err := db.QueryRow("SELECT id, username, password_hash, designation, is_admin FROM users WHERE username = ?", creds.Username).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Designation, &user.IsAdmin)
        if err != nil {
            log.Println("User not found or DB error:", err)
            w.WriteHeader(http.StatusUnauthorized)
            return
        }
        log.Println("DB returned user:", user.Username, "hash:", user.PasswordHash)
        log.Println("Password provided:", creds.Password)
        passOk := user.PasswordHash == creds.Password
        log.Println("Plaintext password check for user", creds.Username, ":", passOk)
        if !passOk {
            w.WriteHeader(http.StatusUnauthorized)
            return
        }
        session, _ := store.Get(r, "session")
        session.Values["user_id"] = user.ID
        session.Values["is_admin"] = user.IsAdmin
        session.Save(r, w)
        writeJSON(w, map[string]interface{}{"success": true, "is_admin": user.IsAdmin})
    })

    mux.HandleFunc("/api/logout", func(w http.ResponseWriter, r *http.Request) {
        session, _ := store.Get(r, "session")
        session.Options.MaxAge = -1
        session.Save(r, w)
        writeJSON(w, map[string]interface{}{"success": true})
    })

    // Get tasks for employee
    mux.HandleFunc("/api/mytasks", func(w http.ResponseWriter, r *http.Request) {
        uid, ok := getUserID(r)
        if !ok {
            w.WriteHeader(http.StatusUnauthorized)
            return
        }

        rows, err := db.Query(`
            SELECT t.id, t.job_id, t.description, t.priority, t.deadline, 
                   COALESCE(tu.status, 'Assigned') as status
            FROM tasks t
            JOIN task_assignments ta ON t.id = ta.task_id
            LEFT JOIN task_updates tu ON t.id = tu.task_id AND tu.user_id = ? 
                AND tu.id = (SELECT MAX(id) FROM task_updates WHERE task_id = t.id AND user_id = ?)
            WHERE ta.user_id = ?`, uid, uid, uid)
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            return
        }
        defer rows.Close()
        var tasks []map[string]interface{}
        for rows.Next() {
            var id int
            var jobID, description, priority, deadline, status string
            if err := rows.Scan(&id, &jobID, &description, &priority, &deadline, &status); err != nil {
                continue
            }
            tasks = append(tasks, map[string]interface{}{
                "id": id, "job_id": jobID, "description": description, 
                "priority": priority, "deadline": deadline, "status": status,
            })
        }
        writeJSON(w, tasks)
    })

    // Get all tasks (admin only)
    mux.HandleFunc("/api/tasks", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodGet {
            adminStatus, _ := isAdmin(r)
            if !adminStatus {
                w.WriteHeader(http.StatusForbidden)
                return
            }

            log.Println("Executing tasks query...")
            rows, err := db.Query(`
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
                GROUP BY t.id, t.job_id, t.description, t.priority, t.deadline, latest_status.status`)
            if err != nil {
                log.Println("Database error in /api/tasks:", err)
                w.WriteHeader(http.StatusInternalServerError)
                return
            }
            defer rows.Close()
            var tasks []map[string]interface{}
            for rows.Next() {
                var id int
                var jobID, description, priority, deadline, assignedTo, status sql.NullString
                if err := rows.Scan(&id, &jobID, &description, &priority, &deadline, &assignedTo, &status); err != nil {
                    continue
                }
                var assignedList []string
                if assignedTo.Valid && assignedTo.String != "" {
                    assignedList = strings.Split(assignedTo.String, ",")
                }
                tasks = append(tasks, map[string]interface{}{
                    "id": id, "job_id": jobID, "description": description,
                    "priority": priority, "deadline": deadline, 
                    "assigned_to": assignedList, "status": status.String,
                })
            }
            writeJSON(w, tasks)
        } else if r.Method == http.MethodPost {
            // Create new task (admin only)
            adminStatus, _ := isAdmin(r)
            if !adminStatus {
                w.WriteHeader(http.StatusForbidden)
                return
            }

            var task struct {
                JobID       string `json:"job_id"`
                Description string `json:"description"`
                Priority    string `json:"priority"`
                Deadline    string `json:"deadline"`
                AssignedTo  []int  `json:"assigned_to"`
            }
            if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
                w.WriteHeader(http.StatusBadRequest)
                return
            }

            result, err := db.Exec("INSERT INTO tasks (job_id, description, priority, deadline) VALUES (?, ?, ?, ?)",
                task.JobID, task.Description, task.Priority, task.Deadline)
            if err != nil {
                w.WriteHeader(http.StatusInternalServerError)
                return
            }

            taskID, _ := result.LastInsertId()
            for _, userID := range task.AssignedTo {
                db.Exec("INSERT INTO task_assignments (task_id, user_id) VALUES (?, ?)", taskID, userID)
            }

            writeJSON(w, map[string]interface{}{"success": true, "task_id": taskID})
        }
    })

    // Update task status
    mux.HandleFunc("/api/tasks/", func(w http.ResponseWriter, r *http.Request) {
        path := strings.TrimPrefix(r.URL.Path, "/api/tasks/")
        parts := strings.Split(path, "/")
        if len(parts) < 2 || parts[1] != "status" {
            w.WriteHeader(http.StatusNotFound)
            return
        }
        
        taskID := parts[0]
        uid, ok := getUserID(r)
        if !ok {
            w.WriteHeader(http.StatusUnauthorized)
            return
        }

        if r.Method == http.MethodPost {
            var update struct {
                Status  string `json:"status"`
                Comment string `json:"comment"`
            }
            if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
                w.WriteHeader(http.StatusBadRequest)
                return
            }

            _, err := db.Exec("INSERT INTO task_updates (task_id, user_id, status, comment) VALUES (?, ?, ?, ?)",
                taskID, uid, update.Status, update.Comment)
            if err != nil {
                w.WriteHeader(http.StatusInternalServerError)
                return
            }

            writeJSON(w, map[string]interface{}{"success": true})
        }
    })

    // User management (admin only)
    mux.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodGet {
            adminStatus, _ := isAdmin(r)
            if !adminStatus {
                w.WriteHeader(http.StatusForbidden)
                return
            }

            log.Println("Executing users query...")
            rows, err := db.Query("SELECT id, username, designation, is_admin FROM users")
            if err != nil {
                log.Println("Database error in /api/users:", err)
                w.WriteHeader(http.StatusInternalServerError)
                return
            }
            defer rows.Close()
            var users []User
            for rows.Next() {
                var user User
                if err := rows.Scan(&user.ID, &user.Username, &user.Designation, &user.IsAdmin); err != nil {
                    continue
                }
                users = append(users, user)
            }
            writeJSON(w, users)
        } else if r.Method == http.MethodPost {
            // Create user (admin only)
            adminStatus, _ := isAdmin(r)
            if !adminStatus {
                w.WriteHeader(http.StatusForbidden)
                return
            }

            var user User
            if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
                w.WriteHeader(http.StatusBadRequest)
                return
            }

            _, err := db.Exec("INSERT INTO users (username, password_hash, designation, is_admin) VALUES (?, ?, ?, ?)",
                user.Username, user.PasswordHash, user.Designation, user.IsAdmin)
            if err != nil {
                w.WriteHeader(http.StatusInternalServerError)
                return
            }

            writeJSON(w, map[string]interface{}{"success": true})
        }
    })
} 