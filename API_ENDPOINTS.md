# API Endpoints Documentation

Complete API documentation for UserMes Backend.

## Base URL

```
http://localhost:3000
```

## Table of Contents

- [Health Check](#health-check)
- [User Module](#user-module)
- [Resource Module](#resource-module)

---

## Health Check

### Get Health Status

```http
GET /health
```

**Response:** `200 OK`

```json
{
  "status": "healthy"
}
```

---

## User Module

Base path: `/api/users`

### 1. Register User

Create a new user account.

```http
POST /api/users/register
Content-Type: application/json
```

**Request Body:**

```json
{
  "email": "user@example.com",
  "password": "SecurePassword123",
  "name": "John Doe"
}
```

**Response:** `201 Created`

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "user@example.com",
  "name": "John Doe",
  "is_active": true,
  "created_at": "2024-01-15T10:00:00Z",
  "updated_at": "2024-01-15T10:00:00Z",
  "last_login_at": null
}
```

**Validation Rules:**
- Email: Valid email format
- Password: Minimum 8 characters, at least one uppercase, one lowercase, and one number
- Name: Minimum 2 characters

**Error Responses:**
- `400 Bad Request`: Invalid email, password, or name
- `409 Conflict`: Email already exists

---

### 2. Login

Authenticate and receive access token.

```http
POST /api/users/login
Content-Type: application/json
```

**Request Body:**

```json
{
  "email": "user@example.com",
  "password": "SecurePassword123"
}
```

**Response:** `200 OK`

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "name": "John Doe",
    "is_active": true,
    "created_at": "2024-01-15T10:00:00Z",
    "updated_at": "2024-01-15T10:00:00Z",
    "last_login_at": "2024-01-15T11:30:00Z"
  }
}
```

**Error Responses:**
- `400 Bad Request`: Missing credentials
- `401 Unauthorized`: Invalid credentials
- `403 Forbidden`: User is inactive

---

### 3. Get Current User (Protected)

Get authenticated user's profile.

```http
GET /api/users/me
Authorization: Bearer <token>
```

**Response:** `200 OK`

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "user@example.com",
  "name": "John Doe",
  "is_active": true,
  "created_at": "2024-01-15T10:00:00Z",
  "updated_at": "2024-01-15T10:00:00Z",
  "last_login_at": "2024-01-15T11:30:00Z"
}
```

**Error Responses:**
- `401 Unauthorized`: Missing or invalid token

---

### 4. Get User by ID (Protected)

Retrieve a specific user by ID.

```http
GET /api/users/:id
Authorization: Bearer <token>
```

**Response:** `200 OK`

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "user@example.com",
  "name": "John Doe",
  "is_active": true,
  "created_at": "2024-01-15T10:00:00Z",
  "updated_at": "2024-01-15T10:00:00Z",
  "last_login_at": "2024-01-15T11:30:00Z"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid UUID format
- `401 Unauthorized`: Missing or invalid token
- `404 Not Found`: User not found

---

### 5. Update User (Protected)

Update user's name.

```http
PUT /api/users/:id
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**

```json
{
  "name": "Jane Doe"
}
```

**Response:** `200 OK`

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "user@example.com",
  "name": "Jane Doe",
  "is_active": true,
  "created_at": "2024-01-15T10:00:00Z",
  "updated_at": "2024-01-15T12:00:00Z",
  "last_login_at": "2024-01-15T11:30:00Z"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid name or UUID
- `401 Unauthorized`: Missing or invalid token
- `404 Not Found`: User not found

---

### 6. Change Password (Protected)

Change user's password.

```http
POST /api/users/:id/change-password
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**

```json
{
  "old_password": "SecurePassword123",
  "new_password": "NewSecurePassword456"
}
```

**Response:** `200 OK`

```json
{
  "message": "Password changed successfully"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid password format
- `401 Unauthorized`: Missing or invalid token / wrong old password
- `404 Not Found`: User not found

---

### 7. Deactivate User (Protected)

Deactivate a user account.

```http
POST /api/users/:id/deactivate
Authorization: Bearer <token>
```

**Response:** `200 OK`

```json
{
  "message": "User deactivated successfully"
}
```

**Error Responses:**
- `401 Unauthorized`: Missing or invalid token
- `404 Not Found`: User not found

---

### 8. Activate User (Protected)

Reactivate a deactivated user account.

```http
POST /api/users/:id/activate
Authorization: Bearer <token>
```

**Response:** `200 OK`

```json
{
  "message": "User activated successfully"
}
```

**Error Responses:**
- `401 Unauthorized`: Missing or invalid token
- `404 Not Found`: User not found

---

## Resource Module

Base path: `/api/resources`

### 1. Create Resource

Create a new resource.

```http
POST /api/resources
Content-Type: application/json
```

**Request Body:**

```json
{
  "code": "RES001",
  "shift_id": "SHIFT123",
  "type": "MACHINE",
  "stop_factor": 5
}
```

**Note:** `shift_id` is optional and can be omitted or set to `null`.

**Response:** `201 Created`

```json
{
  "id": "660e8400-e29b-41d4-a716-446655440000",
  "code": "RES001",
  "shift_id": "SHIFT123",
  "type": "MACHINE",
  "stop_factor": 5,
  "created_at": "2024-01-15T10:00:00Z",
  "updated_at": "2024-01-15T10:00:00Z"
}
```

**Validation Rules:**
- Code: 2-50 characters, unique, required
- Type: 2-50 characters, required
- Stop Factor: Non-negative integer (>= 0), required
- Shift ID: Optional

**Error Responses:**
- `400 Bad Request`: Invalid code, type, or stop factor
- `409 Conflict`: Code already exists

---

### 2. Get All Resources

Retrieve all resources with pagination.

```http
GET /api/resources?limit=10&offset=0
```

**Query Parameters:**
- `limit` (optional): Number of results to return (default: 10)
- `offset` (optional): Number of results to skip (default: 0)

**Response:** `200 OK`

```json
{
  "resources": [
    {
      "id": "660e8400-e29b-41d4-a716-446655440000",
      "code": "RES001",
      "shift_id": "SHIFT123",
      "type": "MACHINE",
      "stop_factor": 5,
      "created_at": "2024-01-15T10:00:00Z",
      "updated_at": "2024-01-15T10:00:00Z"
    },
    {
      "id": "770e8400-e29b-41d4-a716-446655440000",
      "code": "RES002",
      "shift_id": null,
      "type": "OPERATOR",
      "stop_factor": 0,
      "created_at": "2024-01-15T11:00:00Z",
      "updated_at": "2024-01-15T11:00:00Z"
    }
  ],
  "total": 2,
  "limit": 10,
  "offset": 0
}
```

---

### 3. Get Resource by ID

Retrieve a specific resource by ID.

```http
GET /api/resources/:id
```

**Response:** `200 OK`

```json
{
  "id": "660e8400-e29b-41d4-a716-446655440000",
  "code": "RES001",
  "shift_id": "SHIFT123",
  "type": "MACHINE",
  "stop_factor": 5,
  "created_at": "2024-01-15T10:00:00Z",
  "updated_at": "2024-01-15T10:00:00Z"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid UUID format
- `404 Not Found`: Resource not found

---

### 4. Get Resource by Code

Retrieve a resource by its unique code.

```http
GET /api/resources/code/:code
```

**Example:**

```http
GET /api/resources/code/RES001
```

**Response:** `200 OK`

```json
{
  "id": "660e8400-e29b-41d4-a716-446655440000",
  "code": "RES001",
  "shift_id": "SHIFT123",
  "type": "MACHINE",
  "stop_factor": 5,
  "created_at": "2024-01-15T10:00:00Z",
  "updated_at": "2024-01-15T10:00:00Z"
}
```

**Error Responses:**
- `400 Bad Request`: Code is required
- `404 Not Found`: Resource not found

---

### 5. Get Resources by Type

Retrieve all resources of a specific type.

```http
GET /api/resources/type/:type?limit=10&offset=0
```

**Example:**

```http
GET /api/resources/type/MACHINE?limit=20&offset=0
```

**Response:** `200 OK`

```json
{
  "resources": [
    {
      "id": "660e8400-e29b-41d4-a716-446655440000",
      "code": "RES001",
      "shift_id": "SHIFT123",
      "type": "MACHINE",
      "stop_factor": 5,
      "created_at": "2024-01-15T10:00:00Z",
      "updated_at": "2024-01-15T10:00:00Z"
    }
  ],
  "total": 1,
  "limit": 20,
  "offset": 0
}
```

**Error Responses:**
- `400 Bad Request`: Type is required

---

### 6. Get Resources by Shift ID

Retrieve all resources assigned to a specific shift.

```http
GET /api/resources/shift/:shiftId?limit=10&offset=0
```

**Example:**

```http
GET /api/resources/shift/SHIFT123?limit=10&offset=0
```

**Response:** `200 OK`

```json
{
  "resources": [
    {
      "id": "660e8400-e29b-41d4-a716-446655440000",
      "code": "RES001",
      "shift_id": "SHIFT123",
      "type": "MACHINE",
      "stop_factor": 5,
      "created_at": "2024-01-15T10:00:00Z",
      "updated_at": "2024-01-15T10:00:00Z"
    }
  ],
  "total": 1,
  "limit": 10,
  "offset": 0
}
```

**Error Responses:**
- `400 Bad Request`: Shift ID is required

---

### 7. Update Resource

Update an existing resource.

```http
PUT /api/resources/:id
Content-Type: application/json
```

**Request Body:**

```json
{
  "code": "RES001-UPDATED",
  "shift_id": "SHIFT456",
  "type": "OPERATOR",
  "stop_factor": 10
}
```

**Response:** `200 OK`

```json
{
  "id": "660e8400-e29b-41d4-a716-446655440000",
  "code": "RES001-UPDATED",
  "shift_id": "SHIFT456",
  "type": "OPERATOR",
  "stop_factor": 10,
  "created_at": "2024-01-15T10:00:00Z",
  "updated_at": "2024-01-15T13:00:00Z"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid code, type, stop factor, or UUID
- `404 Not Found`: Resource not found
- `409 Conflict`: Code already exists (if changed)

---

### 8. Delete Resource

Delete a resource by ID.

```http
DELETE /api/resources/:id
```

**Response:** `200 OK`

```json
{
  "message": "Resource deleted successfully",
  "data": null
}
```

**Error Responses:**
- `400 Bad Request`: Invalid UUID format
- `404 Not Found`: Resource not found

---

## Error Response Format

All API errors follow this standard format:

```json
{
  "error": "error_code",
  "message": "Human-readable error message"
}
```

### Common Error Codes

| Status Code | Error Code | Description |
|-------------|------------|-------------|
| 400 | `invalid_request` | Invalid request body or parameters |
| 400 | `invalid_id` | Invalid UUID format |
| 400 | `invalid_email` | Invalid email format |
| 400 | `invalid_password` | Password doesn't meet requirements |
| 400 | `invalid_code` | Invalid resource code |
| 400 | `invalid_type` | Invalid resource type |
| 400 | `invalid_stop_factor` | Stop factor must be non-negative |
| 401 | `unauthorized` | Missing or invalid authentication token |
| 403 | `forbidden` | User account is inactive |
| 404 | `not_found` | Resource not found |
| 404 | `user_not_found` | User not found |
| 404 | `resource_not_found` | Resource not found |
| 409 | `conflict` | Resource already exists |
| 409 | `email_already_exists` | Email already registered |
| 409 | `code_already_exists` | Resource code already exists |
| 500 | `internal_error` | Internal server error |

---

## Authentication

Most user endpoints require authentication using JWT (JSON Web Token).

### How to Authenticate

1. **Login** to receive a token:
   ```http
   POST /api/users/login
   ```

2. **Use the token** in subsequent requests:
   ```http
   GET /api/users/me
   Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
   ```

### Token Expiration

Tokens expire after 24 hours by default. After expiration, you need to login again to receive a new token.

---

## Examples with cURL

### User Examples

**Register a new user:**
```bash
curl -X POST http://localhost:3000/api/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "SecurePass123",
    "name": "John Doe"
  }'
```

**Login:**
```bash
curl -X POST http://localhost:3000/api/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "SecurePass123"
  }'
```

**Get current user (replace TOKEN with actual token):**
```bash
curl http://localhost:3000/api/users/me \
  -H "Authorization: Bearer TOKEN"
```

### Resource Examples

**Create a resource:**
```bash
curl -X POST http://localhost:3000/api/resources \
  -H "Content-Type: application/json" \
  -d '{
    "code": "MACHINE001",
    "shift_id": "MORNING",
    "type": "MACHINE",
    "stop_factor": 5
  }'
```

**Get all resources:**
```bash
curl "http://localhost:3000/api/resources?limit=10&offset=0"
```

**Get resource by code:**
```bash
curl http://localhost:3000/api/resources/code/MACHINE001
```

**Get resources by type:**
```bash
curl "http://localhost:3000/api/resources/type/MACHINE?limit=20"
```

**Update a resource:**
```bash
curl -X PUT http://localhost:3000/api/resources/660e8400-e29b-41d4-a716-446655440000 \
  -H "Content-Type: application/json" \
  -d '{
    "code": "MACHINE001-UPDATED",
    "shift_id": "EVENING",
    "type": "MACHINE",
    "stop_factor": 10
  }'
```

**Delete a resource:**
```bash
curl -X DELETE http://localhost:3000/api/resources/660e8400-e29b-41d4-a716-446655440000
```

---

## Rate Limiting

Currently, there are no rate limits implemented. In production, consider implementing rate limiting to prevent abuse.

---

## CORS

CORS is enabled for all origins in development. In production, configure specific allowed origins.

---

## Testing

### Using Postman

1. Import the endpoints into Postman
2. Create an environment with `base_url = http://localhost:3000`
3. For authenticated requests, add a token variable and use `{{token}}`

### Using HTTPie

```bash
# Register
http POST localhost:3000/api/users/register \
  email=user@example.com password=Password123 name="Test User"

# Login
http POST localhost:3000/api/users/login \
  email=user@example.com password=Password123

# Get current user (with token)
http localhost:3000/api/users/me Authorization:"Bearer TOKEN"
```

---

## Version

API Version: 1.0.0

Last Updated: 2024-01-15