# ERP System - Core Identity & RBAC Foundation

This project is a Go-based ERP system currently in its foundational phase, focusing on secure identity management and Role-Based Access Control (RBAC).

## Current Project State

The project currently implements a robust authentication and authorization layer:
- **Identity Management**: User registration with Bcrypt password hashing and login with email verification checks.
- **RBAC Architecture**: A 5-table PostgreSQL schema supporting many-to-many relationships between Users, Roles, and Permissions.
- **Security**: Stateless JWT-based authentication with roles embedded in token claims.
- **Middleware**: Custom Go middleware for JWT verification (`AuthMiddleware`) and granular role-based access control (`RequireRole`).

### Defined Roles (Baseline)
1. **Administrator**: Full system access.

---

## Getting Started

### 1. Prerequisites
- Docker & Docker Compose
- Go 1.22+
- PostgreSQL client (`psql`)

### 2. Environment Setup
Create a `.env` file in the root directory (or update the existing one) with the following variables:
```env
POSTGRES_USER=postgres
POSTGRES_PASSWORD=password
POSTGRES_DB=erp_db
DB_HOST=localhost
DB_PORT=5432
REDIS_PORT=6379
JWT_SECRET=your_super_secret_key_change_me
```

### 3. Infrastructure & Database
Start the database services and apply the initial schema:
```bash
# Start Postgres and Redis
docker-compose up -d

# Apply the RBAC schema and seed data
psql -h localhost -U postgres -d erp_db -f 001_initial_schema.sql
```

### 4. Running the Server
```bash
go run cmd/server/main.go
```
The server will start on `http://localhost:8080`.

---

## API Reference

### Public Endpoints

#### Health Check
`GET /health`
- **Description**: Returns the system status.
- **Auth**: None

#### User Registration
`POST /api/v1/auth/register`
- **Description**: Creates a new user account.
- **Body**: `{"email": "user@example.com", "password": "password123", "department": "Engineering"}`
- **Note**: New users have **no roles** by default and require manual verification in the DB for now.

#### User Login
`POST /api/v1/auth/login`
- **Description**: Authenticates user and returns a JWT.
- **Body**: `{"email": "user@example.com", "password": "password123"}`
- **Returns**: `{"token": "eyJhbG...", "message": "Login successful"}`

### Protected Endpoints (Requires `Authorization: Bearer <token>`)

#### User Profile
`GET /api/v1/users/profile`
- **Auth**: `AuthMiddleware`
- **Description**: Returns the claims stored in the user's JWT (ID and Roles).

#### Admin Dashboard
`GET /api/v1/admin/dashboard`
- **Auth**: `AuthMiddleware` + `RequireRole("Administrator")`
- **Description**: Access restricted to users with the Administrator role.

---

## Manual Testing Workflow

1. **Register** a user via the registration endpoint.
2. **Verify the user** manually in the database:
   ```sql
   UPDATE users SET is_verified = true WHERE email = 'user@example.com';
   ```
3. **Login** to receive your JWT.
4. **Assign a Role** (Optional) to test RBAC:
   ```sql
   -- Give Admin role (ID 1) to your user
   INSERT INTO user_roles (user_id, role_id) 
   SELECT id, 1 FROM users WHERE email = 'user@example.com';
   ```
5. **Access Protected Routes** using the token in the `Authorization` header.
