-- Database schema for MayDiv CRM

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

-- Insert sample data
INSERT INTO users (username, password_hash, designation, is_admin) VALUES
('admin', 'admin123', 'System Administrator', TRUE),
('john.doe', 'password123', 'Software Developer', FALSE),
('jane.smith', 'password123', 'Project Manager', FALSE),
('mike.wilson', 'password123', 'QA Engineer', FALSE);

INSERT INTO tasks (job_id, description, priority, deadline) VALUES
('TASK-001', 'Implement user authentication system', 'High', '2024-02-15'),
('TASK-002', 'Design database schema for CRM', 'Medium', '2024-02-20'),
('TASK-003', 'Create responsive dashboard UI', 'High', '2024-02-25'),
('TASK-004', 'Write API documentation', 'Low', '2024-03-01');

INSERT INTO task_assignments (task_id, user_id) VALUES
(1, 2), -- TASK-001 assigned to john.doe
(1, 3), -- TASK-001 also assigned to jane.smith
(2, 2), -- TASK-002 assigned to john.doe
(3, 3), -- TASK-003 assigned to jane.smith
(4, 4); -- TASK-004 assigned to mike.wilson

INSERT INTO task_updates (task_id, user_id, status, comment) VALUES
(1, 2, 'In Progress', 'Started working on authentication module'),
(1, 3, 'In Progress', 'Reviewing the authentication flow'),
(2, 2, 'Completed', 'Database schema design completed'),
(3, 3, 'In Progress', 'Working on dashboard components'); 