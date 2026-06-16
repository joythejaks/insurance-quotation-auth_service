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

## 🏃 Menjalankan Aplikasi

1. **Instalasi Dependency**:

   ```bash
   go mod tidy
   ```

2. **Menjalankan Migrasi Database**:

   ```bash
   make migrate-up
   ```

3. **Generate Swagger Docs**:

   ```bash
   swag init -g cmd/api/main.go
   ```

4. **Jalankan Aplikasi**:
   ```bash
   go run cmd/api/main.go
   ```

## 📖 Dokumentasi API

Setelah aplikasi berjalan, buka browser dan akses:
`http://localhost:8080/swagger/index.html`

## 📂 Struktur Proyek

```text
├── cmd/api             # Entry point aplikasi
├── docs/               # Dokumentasi Swagger yang di-generate
├── internal/
│   ├── config/         # Konfigurasi aplikasi & env
│   ├── dto/            # Data Transfer Objects (Request/Response)
│   ├── handler/        # Layer Controller / Entry point HTTP
│   ├── middleware/     # Middleware (Auth, Log, Error)
│   ├── model/          # Entity / Skema Database
│   ├── repository/     # Layer akses Database
│   ├── service/        # Layer Logika Bisnis
│   ├── utils/          # Helper (JWT, Logger, Response)
│   └── router/         # Definisi rute API
└── migrations/         # File SQL migrasi database
```
