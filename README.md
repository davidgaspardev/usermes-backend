# UserMes - Manufacturing Execution System

[![CI](https://github.com/YOUR_USERNAME/usermes-backend/actions/workflows/ci.yml/badge.svg)](https://github.com/YOUR_USERNAME/usermes-backend/actions/workflows/ci.yml)
[![Coverage](https://img.shields.io/badge/coverage-92.3%25-green.svg)](https://github.com/YOUR_USERNAME/usermes-backend/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/YOUR_USERNAME/usermes-backend)](https://goreportcard.com/report/github.com/YOUR_USERNAME/usermes-backend)

**UserMes** is a modern Manufacturing Execution System (MES) built with Go, designed to manage production resources, plants, and operations in industrial environments.

## 🏭 What is UserMes?

UserMes is a comprehensive MES platform that enables manufacturers to:

- **Manage Plants** - Create and manage multiple manufacturing facilities
- **Track Resources** - Monitor machines, equipment, and production lines
- **Control Access** - Secure user authentication and authorization
- **Multi-tenant** - Support multiple plants with data isolation
- **Real-time Operations** - Track production, downtime, and performance

## 🏗️ Architecture

UserMes follows **Hexagonal Architecture** (Ports & Adapters) combined with **Modular Monolith** principles, ensuring:

- ✅ **Clean separation of concerns**
- ✅ **High testability** (business logic isolated from infrastructure)
- ✅ **Flexibility** (easy to swap adapters without touching core logic)
- ✅ **Scalability** (modules can evolve into microservices if needed)

### 📐 Architecture Layers

```
┌─────────────────────────────────────────────────────────┐
│              ADAPTERS (Infrastructure)                   │
│  HTTP Handlers, Repositories, JWT, External Services    │
└───────────────────┬─────────────────────────────────────┘
                    │ depends on (interfaces)
┌───────────────────▼─────────────────────────────────────┐
│                   PORTS (Interfaces)                     │
│       Input Ports (Use Cases) + Output Ports            │
└───────────────────┬─────────────────────────────────────┘
                    │ implemented by
┌───────────────────▼─────────────────────────────────────┐
│            APPLICATION (Use Cases)                       │
│         Business Logic Orchestration                     │
└───────────────────┬─────────────────────────────────────┘
                    │ uses
┌───────────────────▼─────────────────────────────────────┐
│              DOMAIN (Core Business)                      │
│    Entities, Value Objects, Domain Logic                │
└─────────────────────────────────────────────────────────┘
```

## 📦 Module Structure

UserMes is organized into **bounded contexts** (modules) following Domain-Driven Design:

```
internal/modules/
│
├── iam/                          # Identity & Access Management
│   └── user/                     # User authentication & management
│       ├── domain/               # User entity, value objects (Email, Password)
│       ├── application/          # Use cases (Register, Login, UpdateUser)
│       └── infrastructure/       # HTTP handlers, JWT, repositories
│
├── organization/                 # Organizational Structure
│   └── plant/                    # Manufacturing plants/facilities
│       ├── domain/               # Plant entity, business rules
│       ├── application/          # Use cases (Create, Update, Deactivate)
│       └── infrastructure/       # HTTP handlers, repositories
│
└── production/                   # Production Execution (Core MES)
    └── resource/                 # Production resources (machines, equipment)
        ├── domain/               # Resource entity, stop factor logic
        ├── application/          # Use cases (Create, Update, Track)
        └── infrastructure/       # HTTP handlers, repositories
```

## 🎯 Golden Rules

### 1. **Weak References Between Modules**

Modules communicate only through **IDs (UUIDs)**, never through direct object references:

```go
// ✅ CORRECT - Weak reference
type Resource struct {
    id        uuid.UUID
    plantCode string      // Reference to Plant by CODE
    code      string
}

// ❌ WRONG - Strong coupling
type Resource struct {
    id    uuid.UUID
    plant *entity.Plant   // Direct reference creates coupling!
}
```

**Why?** This ensures:
- Modules remain independent
- Each module can be developed/tested/deployed separately
- Future migration to microservices is straightforward
- Clear boundaries between domains

### 2. **Ports are Interfaces, Adapters are Structs**

```go
// Port (Interface) - Defines WHAT the system does
type PlantService interface {
    Create(ctx context.Context, ...) (*entity.Plant, error)
    GetByCode(ctx context.Context, code string) (*entity.Plant, error)
}

// Use Case (Struct) - Implements HOW it's done
type PlantServiceImpl struct {
    repository output.PlantRepository  // Depends on interface!
}

// Adapter (Struct) - Implements infrastructure details
type PlantHandler struct {
    service input.PlantService  // Depends on interface!
}
```

### 3. **Test What Matters**

We **only test business logic** (domain + application layers):

```
✅ TESTED (Business Logic)
├── domain/entity/        # Entity behavior, invariants
├── domain/valueobject/   # Value object validation
└── application/usecase/  # Use case orchestration

🚫 NOT TESTED (Infrastructure)
└── infrastructure/       # HTTP handlers, databases, external APIs
```

**Current Coverage: 92.3%** of business logic ✅

### 4. **Multi-tenancy by Design**

Every production resource is scoped to a **Plant**:

- Resources belong to a Plant (identified by `plantCode`)
- Users can access multiple Plants
- Data isolation is enforced at the application level
- API routes reflect this hierarchy:

```
/v1/plants/:plant_code/production/resources
```

## 🚀 Quick Start

### Prerequisites
- Go 1.22.1+
- (Optional) PostgreSQL for production

### Installation & Run

```bash
# Clone repository
git clone https://github.com/YOUR_USERNAME/usermes-backend.git
cd usermes-backend

# Install dependencies
go mod download

# Run application
go run cmd/main.go

# Server starts at http://localhost:3001
```

### Configuration

Environment variables:

```bash
PORT=3001                    # Server port (default: 3001)
JWT_SECRET=your-secret-key   # JWT signing secret
TOKEN_DURATION=24h           # Token expiration time
```

## 📚 API Documentation

### Base URL Structure

```
/v1/users/*                              # IAM - User management
/v1/plants/*                             # Organization - Plant management
/v1/plants/:plant_code/production/*      # Production - Scoped to plant
```

### Example: Create a Plant

```bash
# 1. Register user
curl -X POST http://localhost:3001/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "owner@company.com",
    "password": "SecurePass123!",
    "username": "owner",
    "name": "Plant Owner"
  }'

# 2. Login
curl -X POST http://localhost:3001/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "owner",
    "password": "SecurePass123!"
  }'

# Response: { "token": "eyJhbGc..." }

# 3. Create plant
curl -X POST http://localhost:3001/v1/plants \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "code": "SP01",
    "name": "São Paulo Plant",
    "latitude": -23.5505,
    "longitude": -46.6333
  }'

# 4. Create resource in plant
curl -X POST http://localhost:3001/v1/plants/SP01/production/resources \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "code": "MACHINE-001",
    "type": "CNC",
    "stop_factor": 10
  }'
```

## 🧪 Testing

```bash
# Run all business logic tests
go test ./internal/modules/.../domain/... ./internal/modules/.../application/... -v

# Run tests with coverage
go test -coverprofile=coverage.out -covermode=atomic \
  ./internal/modules/iam/user/domain/... \
  ./internal/modules/iam/user/application/... \
  ./internal/modules/organization/plant/domain/... \
  ./internal/modules/organization/plant/application/... \
  ./internal/modules/production/resource/domain/... \
  ./internal/modules/production/resource/application/...

# View coverage report
go tool cover -html=coverage.out
```

## 📊 Project Status

| Module | Status | Domain Coverage | Application Coverage |
|--------|--------|-----------------|---------------------|
| **iam/user** | ✅ Production Ready | 95.2% | 89.5% |
| **organization/plant** | ✅ Production Ready | 100% | 94.1% |
| **production/resource** | ✅ Production Ready | 98.3% | 92.7% |

## 🗺️ Roadmap

### Phase 1 - Foundation (Current)
- [x] User management (IAM)
- [x] Plant management (Organization)
- [x] Resource management (Production)
- [x] Multi-tenancy support
- [ ] PostgreSQL integration

### Phase 2 - Core MES Features
- [ ] Shift management
- [ ] Production orders
- [ ] Stop/downtime tracking
- [ ] Real-time resource status

### Phase 3 - Advanced Features
- [ ] Quality management
- [ ] Maintenance tracking
- [ ] Inventory management
- [ ] Analytics & reporting (OEE, MTBF, MTTR)

### Phase 4 - Platform
- [ ] WebSocket for real-time updates
- [ ] Event sourcing
- [ ] Multi-language support
- [ ] Mobile API

## 🔒 Security Features

- **JWT Authentication** - Secure token-based auth with 24h expiration
- **Password Hashing** - bcrypt with salt
- **Email Validation** - RFC 5322 compliant
- **Input Sanitization** - All endpoints validate input
- **UUID-based IDs** - Prevents enumeration attacks
- **Multi-tenant Isolation** - Plant-scoped data access

## 🛠️ Development

### Project Structure

```
usermes-backend/
├── cmd/
│   └── main.go                    # Application entry point
├── internal/
│   ├── modules/                   # Business modules (bounded contexts)
│   │   ├── iam/user/
│   │   ├── organization/plant/
│   │   └── production/resource/
│   └── shared/                    # Shared infrastructure
│       └── infrastructure/
│           ├── http/server/       # HTTP server setup
│           └── security/          # JWT implementation
├── go.mod
└── README.md
```

### Adding a New Module

1. Create module structure:
```bash
mkdir -p internal/modules/your-domain/your-module/{domain/{entity,valueobject,errors},application/{port/{input,output},usecase},infrastructure/{adapter/{input/http,output/persistence},dto}}
```

2. Implement in order:
   - Domain entities & value objects
   - Domain tests
   - Application ports (interfaces)
   - Application use cases
   - Application tests
   - Infrastructure adapters (HTTP, repositories)

3. **Remember the golden rule**: Use weak references (UUIDs) to other modules!

## 📖 Documentation

- [User Module](./internal/modules/iam/user/README.md)
- [Plant Module](./internal/modules/organization/plant/README.md)
- [Resource Module](./internal/modules/production/resource/README.md)

## 🤝 Contributing

1. Follow the architecture principles
2. Write tests for domain and application layers
3. Use weak references between modules
4. Keep coverage above 90%
5. Follow Go best practices

## 📄 License

MIT License

---

**Built with ❤️ for modern manufacturing**