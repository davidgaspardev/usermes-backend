# Quick Start Guide - UserMes Backend

Quick guide to get started with UserMes Backend.

## 🚀 Quick Installation

### Prerequisites

- Go 1.22.1 or higher
- curl (for testing)

### Steps

1. **Clone the repository**
```bash
git clone <your-repository>
cd usermes-backend
```

2. **Install dependencies**
```bash
go mod download
```

3. **Run the application**
```bash
go run cmd/main.go
```

The API will be available at: `http://localhost:3000`

## 📚 Available Modules

The application includes two modules:

1. **User Module** - User authentication and management
   - Base path: `/api/users`
   - Features: Register, login, profile management, password change

2. **Resource Module** - Resource management (NEW!)
   - Base path: `/api/resources`
   - Features: Full CRUD operations for resources

## 🧪 Testing the APIs

### User Module Examples

**Register a user:**
```bash
curl -X POST http://localhost:3000/api/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "Password123",
    "name": "Test User"
  }'
```

**Login:**
```bash
curl -X POST http://localhost:3000/api/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "Password123"
  }'
```

### Resource Module Examples

**Create a resource:**
```bash
curl -X POST http://localhost:3000/api/resources \
  -H "Content-Type: application/json" \
  -d '{
    "code": "RES001",
    "shift_id": "SHIFT123",
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
curl http://localhost:3000/api/resources/code/RES001
```

**Update a resource:**
```bash
curl -X PUT http://localhost:3000/api/resources/{id} \
  -H "Content-Type: application/json" \
  -d '{
    "code": "RES001-UPDATED",
    "shift_id": "SHIFT456",
    "type": "OPERATOR",
    "stop_factor": 10
  }'
```

**Delete a resource:**
```bash
curl -X DELETE http://localhost:3000/api/resources/{id}
```

## 📖 Documentation

- **Complete API Documentation**: See [API_ENDPOINTS.md](API_ENDPOINTS.md)
- **Resource Module Details**: See [internal/modules/resource/README.md](internal/modules/resource/README.md)
- **Architecture**: See [README.md](README.md)

## ✅ Quick Test

Open another terminal and run:

```bash
# Health check
curl http://localhost:3000/health
```

Should return:
```json
{"status":"ok","timestamp":1705318800}
```

## 📝 Your First User

### 1. Register a user

```bash
curl -X POST http://localhost:3000/api/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "your@email.com",
    "password": "YourPassword123",
    "name": "Your Name"
  }'
```

**Expected response:**
```json
{
  "message": "User registered successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "your@email.com",
    "name": "Your Name",
    "is_active": true,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

### 2. Login

```bash
curl -X POST http://localhost:3000/api/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "your@email.com",
    "password": "YourPassword123"
  }'
```

**Expected response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "your@email.com",
    "name": "Your Name",
    "is_active": true,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z",
    "last_login_at": "2024-01-15T11:00:00Z"
  }
}
```

**⚠️ IMPORTANT:** Copy the returned token! You'll need it for the next requests.

### 3. Get your profile

Replace `<YOUR-TOKEN>` with the token you received from login:

```bash
curl -X GET http://localhost:3000/api/users/me \
  -H "Authorization: Bearer <YOUR-TOKEN>"
```

**Expected response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "your@email.com",
  "name": "Your Name",
  "is_active": true,
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z",
  "last_login_at": "2024-01-15T11:00:00Z"
}
```

## 🎯 Automated Test

Run all tests at once:

```bash
# Terminal 1: Start the server
go run cmd/main.go

# Terminal 2: Run the tests
./simple_test.sh
```

## 🛠️ Using Makefile

```bash
# Run the application
make run

# Build
make build

# Run tests
make test

# See all available commands
make help
```

## 📁 Project Structure

```
usermes-backend/
├── cmd/                          # Entry point
│   └── main.go
├── internal/
│   ├── modules/                  # Monolith modules
│   │   └── user/                 # User module
│   │       ├── domain/           # Entities and business rules
│   │       ├── application/      # Use cases
│   │       └── infrastructure/   # HTTP, Database, etc
│   └── shared/                   # Shared code
│       └── infrastructure/
└── pkg/                          # Public packages
```

## 📚 Available Endpoints

### Public (no authentication required)

- `POST /api/users/register` - Register new user
- `POST /api/users/login` - Login
- `GET /health` - Health check

### Protected (requires token)

- `GET /api/users/me` - Get logged user profile
- `GET /api/users/:id` - Get user by ID
- `PUT /api/users/:id` - Update user name
- `POST /api/users/:id/change-password` - Change password
- `POST /api/users/:id/deactivate` - Deactivate account
- `POST /api/users/:id/activate` - Activate account

## 🔐 Password Rules

- Minimum 8 characters
- Maximum 72 characters
- Must contain at least one letter
- Must contain at least one number

## 💡 Tips

### Save token in variable (Linux/Mac)

```bash
# Login and save token
TOKEN=$(curl -s -X POST http://localhost:3000/api/users/login \
  -H "Content-Type: application/json" \
  -d '{"email":"your@email.com","password":"YourPassword123"}' \
  | grep -o '"token":"[^"]*"' | sed 's/"token":"//; s/"$//')

# Use the token
curl -X GET http://localhost:3000/api/users/me \
  -H "Authorization: Bearer $TOKEN"
```

### Save token in variable (Windows PowerShell)

```powershell
# Login and save token
$response = Invoke-RestMethod -Uri "http://localhost:3000/api/users/login" `
  -Method POST `
  -ContentType "application/json" `
  -Body '{"email":"your@email.com","password":"YourPassword123"}'

$token = $response.token

# Use the token
Invoke-RestMethod -Uri "http://localhost:3000/api/users/me" `
  -Method GET `
  -Headers @{Authorization="Bearer $token"}
```

## 🐛 Troubleshooting

### Port 3000 already in use

```bash
# Find process using the port
lsof -i :3000

# Kill the process
kill -9 <PID>
```

### Error "Authorization header is required"

Make sure to include the authorization header:
```bash
-H "Authorization: Bearer <your-token>"
```

### Error "Invalid or expired token"

The token expires after 24 hours. Login again to get a new token.

### Error "email already exists"

Use a different email or login with the existing email.

## 🔄 Reset Data

Since it's using in-memory repository, just restart the server:

```bash
# Press Ctrl+C to stop the server
# Run again
go run cmd/main.go
```

All data will be lost and you can start from scratch.

## 📖 Next Steps

1. Read the [Complete Documentation](README.md)
2. Understand the [Architecture](ARCHITECTURE.md)
3. See more [Usage Examples](EXAMPLES.md)
4. Contribute to the project!

## 🆘 Need Help?

- Open an issue on GitHub
- Check the complete documentation
- Review the examples

## 📝 Important Notes

⚠️ **Development Environment:**
- Data stored in memory (not persisted)
- JWT secret is hardcoded (change in production)
- CORS configured to accept any origin

⚠️ **Before Going to Production:**
- Configure environment variables
- Use real database (PostgreSQL, MySQL)
- Configure HTTPS
- Implement rate limiting
- Configure proper logging
- Use secure and unique JWT secret

## ✨ Quick Examples

### Complete flow in one script

```bash
#!/bin/bash

BASE_URL="http://localhost:3000"

# 1. Register
curl -X POST $BASE_URL/api/users/register \
  -H "Content-Type: application/json" \
  -d '{"email":"demo@test.com","password":"Demo123","name":"Demo User"}'

# 2. Login and save token
TOKEN=$(curl -s -X POST $BASE_URL/api/users/login \
  -H "Content-Type: application/json" \
  -d '{"email":"demo@test.com","password":"Demo123"}' \
  | grep -o '"token":"[^"]*"' | sed 's/"token":"//; s/"$//')

# 3. View profile
curl -X GET $BASE_URL/api/users/me \
  -H "Authorization: Bearer $TOKEN"
```

Save this as `test.sh`, give execution permission (`chmod +x test.sh`) and run (`./test.sh`).

---

**Ready to start!** 🚀

If everything worked, you're ready to develop new features or add new modules to the monolith.