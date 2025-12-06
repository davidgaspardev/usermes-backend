# User Module

API endpoints for user management including registration, authentication, and profile operations.

## Base URL
All endpoints are prefixed with `/api/users`

---

## 🔓 Public Endpoints (No Authentication)

### Register User
Create a new user account.

**POST** `/api/users/register`

```bash
curl -X POST http://localhost:3001/api/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@example.com",
    "password": "StrongPassword123!",
    "username": "johndoe",
    "name": "John Doe"
  }'
```

**Response (201 Created):**
```json
{
  "message": "User registered successfully",
  "data": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "email": "john.doe@example.com",
    "username": "johndoe",
    "name": "John Doe",
    "is_active": true,
    "created_at": "2023-12-01T10:00:00Z",
    "updated_at": "2023-12-01T10:00:00Z",
    "last_login_at": null
  }
}
```

### Login
Authenticate user and receive access token.

**POST** `/api/users/login`

```bash
curl -X POST http://localhost:3001/api/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "password": "StrongPassword123!"
  }'
```

**Response (200 OK):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "email": "john.doe@example.com",
    "username": "johndoe",
    "name": "John Doe",
    "is_active": true,
    "created_at": "2023-12-01T10:00:00Z",
    "updated_at": "2023-12-01T10:00:00Z",
    "last_login_at": "2023-12-01T10:30:00Z"
  }
}
```

---

## 🔒 Protected Endpoints (Authentication Required)

> **Note:** Include the JWT token in the Authorization header: `Authorization: Bearer <token>`

### Get Current User
Get authenticated user's profile.

**GET** `/api/users/me`

```bash
curl -X GET http://localhost:3001/api/users/me \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### Get User by ID
Get user profile by user ID.

**GET** `/api/users/:id`

```bash
curl -X GET http://localhost:3001/api/users/123e4567-e89b-12d3-a456-426614174000 \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**Response (200 OK):**
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "email": "john.doe@example.com",
  "username": "johndoe",
  "name": "John Doe",
  "is_active": true,
  "created_at": "2023-12-01T10:00:00Z",
  "updated_at": "2023-12-01T10:00:00Z",
  "last_login_at": "2023-12-01T10:30:00Z"
}
```

### Update User
Update user profile information.

**PUT** `/api/users/:id`

```bash
curl -X PUT http://localhost:3001/api/users/123e4567-e89b-12d3-a456-426614174000 \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Smith"
  }'
```

**Response (200 OK):**
```json
{
  "message": "User updated successfully",
  "data": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "email": "john.doe@example.com",
    "username": "johndoe",
    "name": "John Smith",
    "is_active": true,
    "created_at": "2023-12-01T10:00:00Z",
    "updated_at": "2023-12-01T11:00:00Z",
    "last_login_at": "2023-12-01T10:30:00Z"
  }
}
```

### Change Password
Change user password.

**POST** `/api/users/:id/change-password`

```bash
curl -X POST http://localhost:3001/api/users/123e4567-e89b-12d3-a456-426614174000/change-password \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -d '{
    "old_password": "StrongPassword123!",
    "new_password": "NewStrongPassword456!"
  }'
```

**Response (200 OK):**
```json
{
  "message": "Password changed successfully"
}
```

### Deactivate User
Deactivate user account.

**POST** `/api/users/:id/deactivate`

```bash
curl -X POST http://localhost:3001/api/users/123e4567-e89b-12d3-a456-426614174000/deactivate \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**Response (200 OK):**
```json
{
  "message": "User deactivated successfully"
}
```

### Activate User
Reactivate user account.

**POST** `/api/users/:id/activate`

```bash
curl -X POST http://localhost:3001/api/users/123e4567-e89b-12d3-a456-426614174000/activate \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**Response (200 OK):**
```json
{
  "message": "User activated successfully"
}
```

---

## ❌ Error Responses

All endpoints may return error responses in the following format:

```json
{
  "error": "validation_error",
  "message": "Email is required"
}
```

### Common Error Codes:
- `bad_request` - Invalid request format or parameters
- `validation_error` - Request validation failed
- `unauthorized` - Authentication required or invalid credentials
- `forbidden` - User account is inactive
- `not_found` - User not found
- `conflict` - Email or username already exists
- `invalid_email` - Invalid email format
- `invalid_password` - Invalid password format
- `invalid_username` - Invalid username format
- `invalid_name` - Invalid name format
- `internal_error` - Server error

---

## 📋 Request Validation Rules

### Registration:
- **email**: Required, valid email format, max 255 characters
- **password**: Required, min 8 characters, must contain letters and numbers
- **username**: Required, min 3 characters, max 30 characters, alphanumeric with underscores and hyphens only
- **name**: Required, min 2 characters, max 100 characters

### Login:
- **username**: Required, min 3 characters, max 30 characters
- **password**: Required

### Update User:
- **name**: Required, min 2 characters, max 100 characters

### Change Password:
- **old_password**: Required
- **new_password**: Required, min 8 characters, must contain letters and numbers