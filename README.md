# RBAC System

Minimal Role-Based Access Control system with a Go/Gin backend and React/Vite frontend.

## Project Structure

```
backend/   — Go API (Gin, JWT, in-memory storage)
frontend/  — React SPA (Vite, plain CSS)
```

## Running

### Backend

```bash
cd backend
go run .
```

API starts on `http://localhost:8080`.

### Frontend

```bash
cd frontend
npm install
npm run dev
```

Opens on `http://localhost:5173`.

## Seed Data

The app boots with:

| Entity     | Name         | Details                              |
|------------|-------------|--------------------------------------|
| User       | `admin`     | password: `admin123`, role: admin    |
| Role       | `admin`     | all permissions                      |
| Role       | `viewer`    | `user:read`, `role:read`             |
| Permission | `user:read` |                                      |
| Permission | `user:write`|                                      |
| Permission | `role:read` |                                      |
| Permission | `role:write`|                                      |

## Usage Examples

### Login
```bash
curl -s -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

### List users (protected — requires `user:read`)
```bash
curl -s http://localhost:8080/api/users \
  -H "Authorization: Bearer <token>"
```

### Create a role
```bash
curl -s -X POST http://localhost:8080/api/roles \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"name":"editor"}'
```

### Assign roles to a user
```bash
curl -s -X PUT http://localhost:8080/api/users/<user_id>/roles \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"role_ids":["<role_id>"]}'
```

### Assign permissions to a role
```bash
curl -s -X PUT http://localhost:8080/api/roles/<role_id>/permissions \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"permission_ids":["<perm_id>"]}'
```

## Routes

| Method | Path                            | Permission    |
|--------|---------------------------------|---------------|
| POST   | `/login`                        | public        |
| GET    | `/api/users`                    | `user:read`   |
| POST   | `/api/users`                    | `user:write`  |
| PUT    | `/api/users/:id/roles`          | `role:write`  |
| GET    | `/api/roles`                    | `role:read`   |
| POST   | `/api/roles`                    | `role:write`  |
| PUT    | `/api/roles/:id/permissions`    | `role:write`  |
| GET    | `/api/permissions`              | `role:read`   |
| POST   | `/api/permissions`              | `role:write`  |
