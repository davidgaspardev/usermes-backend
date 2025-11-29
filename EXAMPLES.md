# API Usage Examples

This document contains practical examples of how to use the UserMes API.

## 📋 Table of Contents

- [Configuration](#configuration)
- [Authentication](#authentication)
- [User Operations](#user-operations)
- [Examples with cURL](#examples-with-curl)
- [Examples with JavaScript/Fetch](#examples-with-javascriptfetch)
- [Examples with Go](#examples-with-go)
- [Status Codes](#status-codes)
- [Error Handling](#error-handling)

## Configuration

Base URL: `http://localhost:3000`

All API routes start with `/api`

## Authentication

The API uses JWT (JSON Web Tokens) for authentication. After logging in, you'll receive a token that must be included in the `Authorization` header of all protected requests.

Format: `Authorization: Bearer <your-token>`

## User Operations

### 1. Register New User

**Endpoint:** `POST /api/users/register`

**Headers:**
```
Content-Type: application/json
```

**Body:**
```json
{
  "email": "john@example.com",
  "password": "SecurePass123",
  "name": "John Silva"
}
```

**Success Response (201):**
```json
{
  "message": "User registered successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "john@example.com",
    "name": "John Silva",
    "is_active": true,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

### 2. Login

**Endpoint:** `POST /api/users/login`

**Headers:**
```
Content-Type: application/json
```

**Body:**
```json
{
  "email": "john@example.com",
  "password": "SecurePass123"
}
```

**Success Response (200):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "john@example.com",
    "name": "John Silva",
    "is_active": true,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z",
    "last_login_at": "2024-01-15T11:00:00Z"
  }
}
```

### 3. Get Current Profile

**Endpoint:** `GET /api/users/me`

**Headers:**
```
Authorization: Bearer <your-token>
```

**Success Response (200):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "john@example.com",
  "name": "John Silva",
  "is_active": true,
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z",
  "last_login_at": "2024-01-15T11:00:00Z"
}
```

### 4. Get User by ID

**Endpoint:** `GET /api/users/:id`

**Headers:**
```
Authorization: Bearer <your-token>
```

**Success Response (200):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "john@example.com",
  "name": "John Silva",
  "is_active": true,
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z",
  "last_login_at": "2024-01-15T11:00:00Z"
}
```

### 5. Update User

**Endpoint:** `PUT /api/users/:id`

**Headers:**
```
Content-Type: application/json
Authorization: Bearer <your-token>
```

**Body:**
```json
{
  "name": "John Peter Silva"
}
```

**Success Response (200):**
```json
{
  "message": "User updated successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "john@example.com",
    "name": "John Peter Silva",
    "is_active": true,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T12:00:00Z",
    "last_login_at": "2024-01-15T11:00:00Z"
  }
}
```

### 6. Change Password

**Endpoint:** `POST /api/users/:id/change-password`

**Headers:**
```
Content-Type: application/json
Authorization: Bearer <your-token>
```

**Body:**
```json
{
  "old_password": "SecurePass123",
  "new_password": "NewSecurePass456"
}
```

**Success Response (200):**
```json
{
  "message": "Password changed successfully"
}
```

### 7. Deactivate Account

**Endpoint:** `POST /api/users/:id/deactivate`

**Headers:**
```
Authorization: Bearer <your-token>
```

**Success Response (200):**
```json
{
  "message": "User deactivated successfully"
}
```

### 8. Activate Account

**Endpoint:** `POST /api/users/:id/activate`

**Headers:**
```
Authorization: Bearer <your-token>
```

**Success Response (200):**
```json
{
  "message": "User activated successfully"
}
```

### 9. Health Check

**Endpoint:** `GET /health`

**Success Response (200):**
```json
{
  "status": "ok",
  "timestamp": 1705318800
}
```

## Examples with cURL

### Register User

```bash
curl -X POST http://localhost:3000/api/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "maria@example.com",
    "password": "SenhaForte123",
    "name": "Maria Santos"
  }'
```

### Login

```bash
curl -X POST http://localhost:3000/api/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "maria@example.com",
    "password": "SenhaForte123"
  }'
```

### Save token in variable (Linux/Mac)

```bash
TOKEN=$(curl -s -X POST http://localhost:3000/api/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "maria@example.com",
    "password": "SenhaForte123"
  }' | jq -r '.token')

echo $TOKEN
```

### Get Profile (using token)

```bash
curl -X GET http://localhost:3000/api/users/me \
  -H "Authorization: Bearer $TOKEN"
```

### Update Name

```bash
curl -X PUT http://localhost:3000/api/users/550e8400-e29b-41d4-a716-446655440000 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "Maria Santos Silva"
  }'
```

### Change Password

```bash
curl -X POST http://localhost:3000/api/users/550e8400-e29b-41d4-a716-446655440000/change-password \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "old_password": "SenhaForte123",
    "new_password": "NovaSenhaForte456"
  }'
```

## Examples with JavaScript/Fetch

### Register User

```javascript
async function registerUser(email, password, name) {
  try {
    const response = await fetch('http://localhost:3000/api/users/register', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ email, password, name }),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.message);
    }

    const data = await response.json();
    console.log('User registered:', data);
    return data;
  } catch (error) {
    console.error('Error registering user:', error);
    throw error;
  }
}

// Usage
registerUser('pedro@example.com', 'Senha123', 'Pedro Costa');
```

### Login

```javascript
async function login(email, password) {
  try {
    const response = await fetch('http://localhost:3000/api/users/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ email, password }),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.message);
    }

    const data = await response.json();
    
    // Save token to localStorage
    localStorage.setItem('token', data.token);
    localStorage.setItem('user', JSON.stringify(data.user));
    
    console.log('Login successful');
    return data;
  } catch (error) {
    console.error('Error logging in:', error);
    throw error;
  }
}

// Usage
login('pedro@example.com', 'Senha123');
```

### Get Profile

```javascript
async function getProfile() {
  try {
    const token = localStorage.getItem('token');
    
    if (!token) {
      throw new Error('Token not found. Please login first.');
    }

    const response = await fetch('http://localhost:3000/api/users/me', {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.message);
    }

    const user = await response.json();
    console.log('User profile:', user);
    return user;
  } catch (error) {
    console.error('Error getting profile:', error);
    throw error;
  }
}

// Usage
getProfile();
```

### Update Name

```javascript
async function updateUserName(userId, newName) {
  try {
    const token = localStorage.getItem('token');
    
    const response = await fetch(`http://localhost:3000/api/users/${userId}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`,
      },
      body: JSON.stringify({ name: newName }),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.message);
    }

    const data = await response.json();
    console.log('Name updated:', data);
    return data;
  } catch (error) {
    console.error('Error updating name:', error);
    throw error;
  }
}

// Usage
updateUserName('550e8400-e29b-41d4-a716-446655440000', 'Novo Nome');
```

### Complete API Client

```javascript
class UserMesAPI {
  constructor(baseURL = 'http://localhost:3000') {
    this.baseURL = baseURL;
    this.token = localStorage.getItem('token');
  }

  setToken(token) {
    this.token = token;
    localStorage.setItem('token', token);
  }

  clearToken() {
    this.token = null;
    localStorage.removeItem('token');
  }

  async request(endpoint, options = {}) {
    const url = `${this.baseURL}${endpoint}`;
    const headers = {
      'Content-Type': 'application/json',
      ...options.headers,
    };

    if (this.token && !options.skipAuth) {
      headers['Authorization'] = `Bearer ${this.token}`;
    }

    const config = {
      ...options,
      headers,
    };

    try {
      const response = await fetch(url, config);
      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.message || 'Request error');
      }

      return data;
    } catch (error) {
      console.error('API Error:', error);
      throw error;
    }
  }

  async register(email, password, name) {
    const data = await this.request('/api/users/register', {
      method: 'POST',
      body: JSON.stringify({ email, password, name }),
      skipAuth: true,
    });
    return data;
  }

  async login(email, password) {
    const data = await this.request('/api/users/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
      skipAuth: true,
    });
    
    this.setToken(data.token);
    return data;
  }

  async logout() {
    this.clearToken();
  }

  async getProfile() {
    return await this.request('/api/users/me');
  }

  async getUserById(userId) {
    return await this.request(`/api/users/${userId}`);
  }

  async updateUser(userId, name) {
    return await this.request(`/api/users/${userId}`, {
      method: 'PUT',
      body: JSON.stringify({ name }),
    });
  }

  async changePassword(userId, oldPassword, newPassword) {
    return await this.request(`/api/users/${userId}/change-password`, {
      method: 'POST',
      body: JSON.stringify({ old_password: oldPassword, new_password: newPassword }),
    });
  }

  async deactivateUser(userId) {
    return await this.request(`/api/users/${userId}/deactivate`, {
      method: 'POST',
    });
  }

  async activateUser(userId) {
    return await this.request(`/api/users/${userId}/activate`, {
      method: 'POST',
    });
  }

  async healthCheck() {
    return await this.request('/health', { skipAuth: true });
  }
}

// Usage
const api = new UserMesAPI();

// Register
await api.register('ana@example.com', 'Senha123', 'Ana Silva');

// Login
const loginData = await api.login('ana@example.com', 'Senha123');

// Get profile
const profile = await api.getProfile();

// Update name
await api.updateUser(profile.id, 'Ana Paula Silva');

// Logout
api.logout();
```

## Examples with Go

### API Client

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
)

const baseURL = "http://localhost:3000"

type Client struct {
    httpClient *http.Client
    token      string
}

func NewClient() *Client {
    return &Client{
        httpClient: &http.Client{},
    }
}

func (c *Client) SetToken(token string) {
    c.token = token
}

func (c *Client) request(method, endpoint string, body interface{}, needsAuth bool) ([]byte, error) {
    var reqBody io.Reader
    if body != nil {
        jsonData, err := json.Marshal(body)
        if err != nil {
            return nil, err
        }
        reqBody = bytes.NewBuffer(jsonData)
    }

    req, err := http.NewRequest(method, baseURL+endpoint, reqBody)
    if err != nil {
        return nil, err
    }

    req.Header.Set("Content-Type", "application/json")
    if needsAuth && c.token != "" {
        req.Header.Set("Authorization", "Bearer "+c.token)
    }

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    respBody, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    if resp.StatusCode >= 400 {
        return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
    }

    return respBody, nil
}

type RegisterRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
    Name     string `json:"name"`
}

type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type LoginResponse struct {
    Token string      `json:"token"`
    User  UserResponse `json:"user"`
}

type UserResponse struct {
    ID        string `json:"id"`
    Email     string `json:"email"`
    Name      string `json:"name"`
    IsActive  bool   `json:"is_active"`
    CreatedAt string `json:"created_at"`
    UpdatedAt string `json:"updated_at"`
}

func (c *Client) Register(email, password, name string) (*UserResponse, error) {
    req := RegisterRequest{
        Email:    email,
        Password: password,
        Name:     name,
    }

    respBody, err := c.request("POST", "/api/users/register", req, false)
    if err != nil {
        return nil, err
    }

    var response struct {
        Message string       `json:"message"`
        Data    UserResponse `json:"data"`
    }

    if err := json.Unmarshal(respBody, &response); err != nil {
        return nil, err
    }

    return &response.Data, nil
}

func (c *Client) Login(email, password string) (*LoginResponse, error) {
    req := LoginRequest{
        Email:    email,
        Password: password,
    }

    respBody, err := c.request("POST", "/api/users/login", req, false)
    if err != nil {
        return nil, err
    }

    var response LoginResponse
    if err := json.Unmarshal(respBody, &response); err != nil {
        return nil, err
    }

    c.SetToken(response.Token)
    return &response, nil
}

func (c *Client) GetProfile() (*UserResponse, error) {
    respBody, err := c.request("GET", "/api/users/me", nil, true)
    if err != nil {
        return nil, err
    }

    var user UserResponse
    if err := json.Unmarshal(respBody, &user); err != nil {
        return nil, err
    }

    return &user, nil
}

func main() {
    client := NewClient()

    // Register
    user, err := client.Register("carlos@example.com", "Senha123", "Carlos Oliveira")
    if err != nil {
        fmt.Printf("Error registering: %v\n", err)
        return
    }
    fmt.Printf("User registered: %+v\n", user)

    // Login
    loginResp, err := client.Login("carlos@example.com", "Senha123")
    if err != nil {
        fmt.Printf("Error logging in: %v\n", err)
        return
    }
    fmt.Printf("Login successful. Token: %s\n", loginResp.Token)

    // Get profile
    profile, err := client.GetProfile()
    if err != nil {
        fmt.Printf("Error getting profile: %v\n", err)
        return
    }
    fmt.Printf("Profile: %+v\n", profile)
}
```

## Status Codes

| Code | Meaning | Description |
|------|---------|-------------|
| 200 | OK | Request successful |
| 201 | Created | Resource created successfully |
| 400 | Bad Request | Invalid or malformed request |
| 401 | Unauthorized | Not authenticated or invalid token |
| 403 | Forbidden | No permission to access the resource |
| 404 | Not Found | Resource not found |
| 409 | Conflict | Conflict (e.g., email already exists) |
| 500 | Internal Server Error | Internal server error |

## Error Handling

All error responses follow the same format:

```json
{
  "error": "error_code",
  "message": "Human-readable error description"
}
```

### Error Examples

#### Email already exists (409)
```json
{
  "error": "conflict",
  "message": "email already exists"
}
```

#### Invalid credentials (401)
```json
{
  "error": "unauthorized",
  "message": "Invalid email or password"
}
```

#### Invalid token (401)
```json
{
  "error": "unauthorized",
  "message": "Invalid or expired token"
}
```

#### Password too short (400)
```json
{
  "error": "invalid_password",
  "message": "password must be at least 8 characters long"
}
```

#### User not found (404)
```json
{
  "error": "not_found",
  "message": "user not found"
}
```

#### Field validation (400)
```json
{
  "error": "validation_error",
  "message": "email is required"
}
```

## Security Tips

1. **Never expose JWT token** in logs or console in production
2. **Use HTTPS** in production to protect token in transit
3. **Store token securely** (httpOnly cookies or memory, not localStorage in production)
4. **Implement refresh tokens** for long-duration sessions
5. **Validate all inputs** on client-side before sending
6. **Implement rate limiting** to prevent brute force attacks
7. **Use strong passwords** with at least 8 characters, letters and numbers
8. **Implement proper logout** by clearing tokens and sessions

## Complete Authentication Flow

```javascript
// 1. User registers
const registerData = await api.register('user@example.com', 'Pass123', 'User Name');

// 2. User logs in
const loginData = await api.login('user@example.com', 'Pass123');
// Token is saved automatically

// 3. Make authenticated requests
const profile = await api.getProfile();
const userData = await api.getUserById(profile.id);
await api.updateUser(profile.id, 'New Name');

// 4. Logout
api.logout();
// Token is removed
```

## Tests with Postman

1. Import the collection in Postman
2. Configure the `baseUrl` variable as `http://localhost:3000`
3. After login, save the returned token in an environment variable
4. Use `{{token}}` in authorization headers of other requests

## Troubleshooting

### Error "Authorization header is required"
- Check if you're including the `Authorization: Bearer <token>` header
- Confirm that the token is being sent correctly

### Error "Invalid or expired token"
- Login again to get a new token
- Tokens expire after 24 hours

### Error "email already exists"
- Use a different email for registration
- Or login with the existing email

### CORS Error
- Make sure the server is running
- Check if you're making requests to the correct URL
- In development, CORS is configured to accept all origins