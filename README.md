# Insurance Auth Service

Microservice for handling Authentication and Authorization within the Insurance Quotation ecosystem. This service is built using Go following Clean Architecture principles and is production-ready.

## 🚀 Key Features

- **Authentication**: User registration, Login (JWT), Logout, and Token Refresh.
- **Authorization**:
  - **RBAC (Role-Based Access Control)**: Restricting access based on roles (ADMIN, USER).
  - **ACL (Access Control List)**: Granular access control at the permission level (e.g., `manage_users`, `view_all_quotations`).
- **User Management**: Complete User CRUD with Search (by Name/Email) and Pagination features.
- **Standardized Response**: Consistent JSON format for all endpoints.
- **Global Error Handling**: Centralized middleware to handle application errors.
- **Structured Logging**: High-performance structured logging (JSON) using Uber Zap.
- **API Documentation**: Automatic Swagger UI integration.

## 🛠 Tech Stack

- **Language**: Go (Golang)
- **Framework**: Gin Gonic
- **ORM**: GORM
- **Database**: PostgreSQL
- **Security**: JWT (v5), BCrypt for password hashing
- **Logging**: Uber Zap
- **Docs**: Swaggo (Swagger)

## 📋 Prerequisites

- Go 1.21 or higher
- PostgreSQL
- Migrate (Golang Migrate) tool for database schema management

## ⚙️ Environment Configuration

Create a `.env` file in the root directory:

```env
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=secret
DB_NAME=insurance_auth_db
DB_PORT=5432
APP_PORT=8080
JWT_SECRET=your_super_secret_key
```

## 🏃 Running the Application

1. **Install Dependencies**:

```bash
   go mod tidy
```

2. **Run Database Migrations**:

```bash
   make migrate-up
```

3. **Generate Swagger Docs**:

```bash
   swag init -g cmd/api/main.go
```

4. **Start the Application**:

```bash
   go run cmd/api/main.go
```

## 📖 API Documentation

Once the application is running, open your browser and navigate to:
`http://localhost:8080/swagger/index.html`

## 📂 Project Structure

```text
├── cmd/api             # Application entry point
├── docs/               # Generated Swagger documentation
├── internal/
│   ├── config/         # Application & environment configuration
│   ├── dto/            # Data Transfer Objects (Request/Response)
│   ├── handler/        # Controller layer / HTTP entry point
│   ├── middleware/     # Middleware (Auth, Log, Error)
│   ├── model/          # Entity / Database schema
│   ├── repository/     # Database access layer
│   ├── service/        # Business logic layer
│   ├── utils/          # Helpers (JWT, Logger, Response)
│   └── router/         # API route definitions
└── migrations/         # Database migration SQL files
```
