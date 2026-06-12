# User Management REST API 🚀

A robust, premium Go-based RESTful API designed with **GoFiber**, **SQLC** (compiled DB layer), **pgx/v5** (PostgreSQL pool), **Uber Zap** (structured logging), and **go-playground/validator** (input validation).

This service manages users with their dynamic calculated age derived from their date of birth (DOB).

---

## 🛠️ Tech Stack & Architecture

- **Web Framework**: [GoFiber (v2)](https://github.com/gofiber/fiber)
- **Database Engine**: PostgreSQL
- **SQL Code Generator**: [SQLC](https://sqlc.dev/) (Version 2 compile format)
- **Database Driver**: [jackc/pgx/v5](https://github.com/jackc/pgx)
- **Structured Logger**: [Uber Zap](https://github.com/uber-go/zap)
- **Validator**: [go-playground/validator (v10)](https://github.com/go-playground/validator)

### Layered Architecture Design

The service implements a clean, decoupling layered architecture:
1. `cmd/server/main.go`: Application bootstrap (config, logging, pool, migrations, routes registration, and graceful shutdown handling).
2. `internal/handler/`: REST layer. Responsible for parsing JSON bodies, validating constraints (including custom past-date validation), calling the service, mapping internal errors to HTTP statuses, and return formatting.
3. `internal/service/`: Business domain logic, including dynamic age calculations and database model conversion to DTO formats.
4. `internal/repository/`: Data access layer abstraction. Exposes database transactions and queries powered by SQLC.
5. `internal/middleware/`: Global HTTP filters (Panic recovery, custom Request UUID injection, and structured request duration logging).
6. `db/migrations/`: Database schema evolutions automatically executed programmatically on start.

---

## ⚙️ Configuration

The service can be configured via environment variables:

| Environment Variable | Description | Default |
| --- | --- | --- |
| `PORT` | The port on which the web server listens | `3000` |
| `DATABASE_URL` | PostgreSQL connection URL string | `postgres://shantanukale@localhost:5432/user_management?sslmode=disable` |

---

## 🏃 Getting Started (Local Development)

### Prerequisites
- **Go** (1.26 or higher)
- **PostgreSQL** (running locally)

### Step 1: Initialize Database
Make sure PostgreSQL is running, and create a database named `user_management`:
```bash
createdb user_management
```

### Step 2: Run the Server
Simply execute the Go package. Database migrations will automatically run on startup!
```bash
go run cmd/server/main.go
```

---

## 🐳 Docker Deployment (One-Click Setup)

Docker and Docker Compose config files are included for streamlined multi-container provisioning:

To launch both PostgreSQL and the REST application in isolated containers:
```bash
docker-compose up --build
```
This sets up:
- PostgreSQL container listening on host port **5433** (to prevent conflicts with any local DB on 5432).
- Go API container listening on host port **3000** (starts automatically after Postgres passes its health check).

---

## 🧪 Running Unit Tests

To execute unit tests for the dynamic age calculation logic:
```bash
go test ./internal/service -v
```

---

## 🔄 API Documentation

All responses include the custom `X-Request-ID` tracking header.

### 1. Create User
- **Method**: `POST`
- **Path**: `/users`
- **Request Body**:
  ```json
  {
    "name": "Alice",
    "dob": "1990-05-10"
  }
  ```
- **Success Response** (`201 Created`):
  ```json
  {
    "id": 1,
    "name": "Alice",
    "dob": "1990-05-10"
  }
  ```

---

### 2. Get User by ID
- **Method**: `GET`
- **Path**: `/users/:id`
- **Success Response** (`200 OK`):
  ```json
  {
    "id": 1,
    "name": "Alice",
    "dob": "1990-05-10",
    "age": 36
  }
  ```
- **Error Response** (`404 Not Found`):
  ```json
  {
    "error": "User not found"
  }
  ```

---

### 3. Update User
- **Method**: `PUT`
- **Path**: `/users/:id`
- **Request Body**:
  ```json
  {
    "name": "Alice Updated",
    "dob": "1991-03-15"
  }
  ```
- **Success Response** (`200 OK`):
  ```json
  {
    "id": 1,
    "name": "Alice Updated",
    "dob": "1991-03-15"
  }
  ```

---

### 4. Delete User
- **Method**: `DELETE`
- **Path**: `/users/:id`
- **Success Response** (`204 No Content`): No response body.

---

### 5. List All Users (with optional Pagination)
- **Method**: `GET`
- **Path**: `/users`
- **Query Parameters**:
  - `page` (optional): Page number (e.g. `1`)
  - `limit` (optional): Items per page (e.g. `10`)
- **Success Response** (`200 OK`):
  ```json
  [
    {
      "id": 1,
      "name": "Alice Updated",
      "dob": "1991-03-15",
      "age": 35
    }
  ]
  ```
