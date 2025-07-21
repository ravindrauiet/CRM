# MayDiv CRM Backend

A clean, organized Go backend for the MayDiv CRM system with proper separation of concerns and easy-to-understand structure.

## Project Structure

```
backend/
├── cmd/
│   ├── server/
│   │   └── main.go          # Main application entry point
│   └── setup/
│       └── main.go          # Database setup utility
├── internal/
│   ├── database/
│   │   ├── connection.go    # Database connection management
│   │   └── migrations.go    # Database migrations and seeding
│   ├── handlers/
│   │   ├── auth_handler.go  # Authentication endpoints
│   │   ├── user_handler.go  # User management endpoints
│   │   └── task_handler.go  # Task management endpoints
│   ├── models/
│   │   ├── user.go          # User-related data structures
│   │   └── task.go          # Task-related data structures
│   ├── repository/
│   │   ├── user_repository.go # User data access layer
│   │   └── task_repository.go # Task data access layer
│   └── services/
│       ├── auth_service.go  # Authentication business logic
│       └── errors.go        # Service layer error definitions
├── schema.sql               # Database schema
├── go.mod                   # Go module file
└── README.md               # This file
```

## Architecture

This backend follows a clean architecture pattern with clear separation of concerns:

### Layers

1. **Handlers** (`internal/handlers/`) - HTTP request/response handling
2. **Services** (`internal/services/`) - Business logic
3. **Repository** (`internal/repository/`) - Data access layer
4. **Models** (`internal/models/`) - Data structures
5. **Database** (`internal/database/`) - Database connection and migrations

### Key Features

- **Clean Architecture**: Clear separation between layers
- **Repository Pattern**: Abstracted data access
- **Service Layer**: Business logic isolation
- **Proper Error Handling**: Structured error responses
- **Session-based Authentication**: Secure user sessions
- **CORS Support**: Frontend integration ready
- **Auto Migration**: Database setup automation

## Setup Instructions

### 1. Environment Configuration

Create a `.env` file in the backend directory:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASS=your_password
DB_NAME=maydiv_crm

# Session Configuration
SESSION_KEY=your-secret-session-key-change-this-in-production
```

### 2. Database Setup

First, create the MySQL database:

```sql
CREATE DATABASE maydiv_crm;
```

### 3. Run the Application

```bash
# Navigate to backend directory
cd backend

# Run the server
go run cmd/server/main.go
```

The server will automatically:
- Connect to the database
- Run migrations to create tables
- Seed initial data if tables are empty
- Start the HTTP server on port 8080

## API Endpoints

### Authentication
- `POST /api/login` - User login
- `POST /api/logout` - User logout

### Users (Admin Only)
- `GET /api/users` - Get all users
- `POST /api/users` - Create new user

### Tasks
- `GET /api/tasks` - Get all tasks (Admin only)
- `POST /api/tasks` - Create new task (Admin only)
- `GET /api/mytasks` - Get tasks assigned to current user
- `POST /api/tasks/{id}/status` - Update task status

## Sample Data

The system comes with pre-seeded data:

### Users
- **admin** (password: admin123) - System Administrator
- **john.doe** (password: password123) - Software Developer
- **jane.smith** (password: password123) - Project Manager
- **mike.wilson** (password: password123) - QA Engineer

### Sample Tasks
- TASK-001: Implement user authentication system (High priority)
- TASK-002: Design database schema for CRM (Medium priority)
- TASK-003: Create responsive dashboard UI (High priority)
- TASK-004: Write API documentation (Low priority)

## Development

### Adding New Features

1. **Models**: Define data structures in `internal/models/`
2. **Repository**: Add data access methods in `internal/repository/`
3. **Services**: Implement business logic in `internal/services/`
4. **Handlers**: Create HTTP endpoints in `internal/handlers/`
5. **Routes**: Register new routes in `cmd/server/main.go`

### Code Organization Principles

- **Single Responsibility**: Each file has one clear purpose
- **Dependency Injection**: Services receive dependencies as parameters
- **Error Handling**: Proper error propagation through layers
- **Type Safety**: Strong typing with Go structs
- **Documentation**: Clear comments and README files

## Troubleshooting

### Common Issues

1. **Database Connection Error**: Check your `.env` file and MySQL server
2. **Port Already in Use**: Change the port in `cmd/server/main.go`
3. **CORS Issues**: Verify frontend URL in CORS middleware
4. **Session Issues**: Check SESSION_KEY in `.env` file

### Logs

The application provides detailed logging for:
- Database connections
- Authentication attempts
- API requests
- Error conditions

## Security Notes

- Passwords are stored as plaintext (as requested for MVP)
- Session keys should be changed in production
- CORS is configured for localhost:3000 (frontend)
- Admin-only endpoints are properly protected

## Next Steps

For production deployment:
1. Use bcrypt for password hashing
2. Implement proper session management
3. Add input validation
4. Set up HTTPS
5. Configure proper CORS for production domains
6. Add rate limiting
7. Implement logging and monitoring 