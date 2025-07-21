# MayDiv CRM - Task Management System

A complete CRM (Customer Relationship Management) system built with **Next.js**, **Go**, and **MySQL** for managing tasks, users, and team collaboration.

## 🚀 Features

### ✅ Authentication & Role Management
- **Admin** and **Employee** role-based access
- Simple username/password login with session management
- Role-based dashboard redirects

### ✅ Admin Dashboard
- **Task Overview**: View all tasks across all users
- **Task Statistics**: Real-time stats (Total, Assigned, In Progress, Completed)
- **Task Creation**: Create tasks and assign to multiple employees
- **Team Management**: View employee performance and completion rates
- **User Management**: Add, edit, remove employees with designations

### ✅ Employee Dashboard
- **Personal Tasks**: View only assigned tasks
- **Status Updates**: Mark tasks as "In Progress" or "Completed"
- **Progress Tracking**: View task history and updates

### ✅ Task Management
- **Task Fields**: Job ID, Description, Priority (High/Medium/Low), Deadline
- **Multi-user Assignment**: Assign tasks to one or multiple employees
- **Status Tracking**: Assigned → In Progress → Completed
- **Task Timeline**: Track all updates and status changes
- **Filtering**: Filter tasks by status, employee, or date

### ✅ User & Designation Management
- **Add/Edit/Remove** employees
- **Role Assignment**: Admin vs Employee
- **Designations**: Set job titles (Developer, Designer, QA, etc.)
- **Performance Tracking**: Monitor completion rates per employee

---

## 🛠️ Tech Stack

- **Frontend**: Next.js 15 + Tailwind CSS
- **Backend**: Go (Golang) with REST API
- **Database**: MySQL
- **Session Management**: Cookie-based sessions (no JWT)
- **CORS**: Configured for local development

---

## 📦 Project Structure

```
MayDiv CRM/
├── frontend/                 # Next.js Frontend
│   ├── src/app/
│   │   ├── login/           # Login page
│   │   ├── dashboard/       # Admin & Employee dashboards
│   │   ├── tasks/           # Task management page
│   │   ├── users/           # User management (admin only)
│   │   └── components/      # Shared components
│   └── package.json
├── backend/                 # Go Backend
│   ├── main.go             # Server entry point
│   ├── handlers.go         # API endpoints
│   ├── models.go           # Data models
│   ├── db.go              # Database connection
│   ├── session.go         # Session management
│   └── go.mod
├── schema.sql              # Database schema
└── README.md
```

---

## 🚀 Quick Setup

### 1. Prerequisites
- **Node.js** (v18+)
- **Go** (v1.19+)
- **MySQL** (v8.0+)
- **MySQL Workbench** or **phpMyAdmin** (optional)

### 2. Database Setup
1. Create a new MySQL database: `maydiv_crm`
2. Import the provided `schema.sql` file:
   ```sql
   USE maydiv_crm;
   SOURCE /path/to/schema.sql;
   ```
3. Insert sample data (optional):
   ```sql
   -- Add admin user
   INSERT INTO users (username, password_hash, designation, is_admin) 
   VALUES ('admin', '123456', 'Administrator', TRUE);
   
   -- Add employee user
   INSERT INTO users (username, password_hash, designation, is_admin) 
   VALUES ('ravindra', '123456', 'Developer', FALSE);
   ```

### 3. Backend Setup
1. Navigate to the backend directory:
   ```bash
   cd backend
   ```
2. Install Go dependencies:
   ```bash
   go mod tidy
   go get github.com/go-sql-driver/mysql
   go get github.com/gorilla/sessions
   go get github.com/joho/godotenv
   go get golang.org/x/crypto/bcrypt
   ```
3. Create `.env` file:
   ```env
   DB_USER=root
   DB_PASS=your_mysql_password
   DB_HOST=localhost
   DB_PORT=3306
   DB_NAME=maydiv_crm
   SESSION_KEY=supersecretkey
   ```
4. Start the backend server:
   ```bash
   go run .
   ```
   Server will start on `http://localhost:8080`

### 4. Frontend Setup
1. Navigate to the frontend directory:
   ```bash
   cd frontend
   ```
2. Install dependencies:
   ```bash
   npm install
   ```
3. Create `.env.local` file:
   ```env
   NEXT_PUBLIC_API_URL=http://localhost:8080
   ```
4. Start the development server:
   ```bash
   npm run dev
   ```
   Frontend will be available at `http://localhost:3000`

---

## 🔐 Default Login Credentials

### Admin Access
- **Username**: `admin`
- **Password**: `123456`

### Employee Access
- **Username**: `ravindra`
- **Password**: `123456`

---

## 🎯 Usage Guide

### For Admins
1. **Login** with admin credentials
2. **Create Tasks**: Use "Create New Task" button on dashboard
3. **Assign Tasks**: Select multiple employees during task creation
4. **Monitor Progress**: View real-time statistics and employee performance
5. **Manage Users**: Add/edit/remove team members with designations

### For Employees
1. **Login** with employee credentials
2. **View Tasks**: See only tasks assigned to you
3. **Update Status**: Mark tasks as "In Progress" or "Completed"
4. **Track Progress**: View your task history and deadlines

### Task Workflow
```
Admin Creates Task → Assigns to Employee(s) → Employee Marks "In Progress" → Employee Completes Task
```

---

## 🔧 API Endpoints

### Authentication
- `POST /api/login` - User login
- `POST /api/logout` - User logout

### Tasks
- `GET /api/tasks` - Get all tasks (admin only)
- `POST /api/tasks` - Create new task (admin only)
- `GET /api/mytasks` - Get user's assigned tasks
- `POST /api/tasks/{id}/status` - Update task status

### Users
- `GET /api/users` - Get all users (admin only)
- `POST /api/users` - Create new user (admin only)

---

## 🎨 Customization

### Adding New User Roles
1. Update `users` table schema to add new role column
2. Modify authentication logic in `handlers.go`
3. Add role-based UI components in frontend

### Adding Task Categories
1. Update `tasks` table to include category field
2. Modify task creation form
3. Add category filtering options

### Email Notifications
1. Integrate email service (SendGrid, SES, etc.)
2. Add notification triggers for task assignments/updates
3. Create email templates

---

## 🚧 Production Deployment

### Security Enhancements (Recommended)
1. **Password Hashing**: Replace plaintext passwords with bcrypt
2. **HTTPS**: Enable SSL/TLS certificates
3. **Environment Variables**: Secure credential management
4. **Input Validation**: Add comprehensive validation
5. **Rate Limiting**: Prevent API abuse

### Database Optimization
1. **Indexes**: Add indexes on frequently queried columns
2. **Connection Pooling**: Optimize database connections
3. **Backup Strategy**: Implement automated backups

---

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 📞 Support

For support and questions:
- Create an issue on GitHub
- Email: support@maydiv.com

---

**Built with ❤️ by MayDiv Team** 